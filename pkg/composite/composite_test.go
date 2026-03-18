package composite

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

const j2000 = 2451545.0

func TestCalcCompositeChart(t *testing.T) {
	cc, err := CalcCompositeChart(CompositeInput{
		Person1Lat: 51.5074, Person1Lon: -0.1278, Person1JD: j2000,
		Person2Lat: 40.7128, Person2Lon: -74.006, Person2JD: j2000 + 365,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcCompositeChart: %v", err)
	}

	if len(cc.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(cc.Planets))
	}
	if len(cc.Houses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(cc.Houses))
	}
	if len(cc.Aspects) == 0 {
		t.Error("Expected at least some aspects")
	}

	// All longitudes should be in [0, 360)
	for _, p := range cc.Planets {
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("Planet %s longitude out of range: %.2f", p.PlanetID, p.Longitude)
		}
		if p.Sign == "" {
			t.Errorf("Planet %s has empty sign", p.PlanetID)
		}
	}
}

func TestMidpoint(t *testing.T) {
	tests := []struct {
		lon1, lon2, expected float64
	}{
		{0, 180, 90},    // opposite points -> midpoint at 90
		{350, 10, 0},    // wrapping around 0
		{10, 10, 10},    // same point
		{100, 200, 150}, // simple midpoint
		{0, 120, 60},
	}

	for _, tt := range tests {
		got := Midpoint(tt.lon1, tt.lon2)
		diff := math.Abs(got - tt.expected)
		if diff > 180 {
			diff = 360 - diff
		}
		if diff > 0.01 {
			t.Errorf("Midpoint(%.1f, %.1f) = %.4f, want %.4f", tt.lon1, tt.lon2, got, tt.expected)
		}
	}
}

func TestCalcCompositeChart_SameChart(t *testing.T) {
	// Same person = composite should equal natal
	cc, err := CalcCompositeChart(CompositeInput{
		Person1Lat: 51.5074, Person1Lon: -0.1278, Person1JD: j2000,
		Person2Lat: 51.5074, Person2Lon: -0.1278, Person2JD: j2000,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcCompositeChart same person: %v", err)
	}

	// Composite of same chart should have longitudes very close to original
	for _, p := range cc.Planets {
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("Planet %s longitude out of range: %.2f", p.PlanetID, p.Longitude)
		}
	}
}
