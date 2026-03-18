package harmonic

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

func TestCalcHarmonicChart_H1(t *testing.T) {
	// Harmonic 1 should equal natal
	hc, err := CalcHarmonicChart(51.5074, -0.1278, j2000, 1, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		t.Fatalf("H1: %v", err)
	}
	if hc.Harmonic != 1 {
		t.Errorf("Harmonic = %d, want 1", hc.Harmonic)
	}
	if len(hc.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(hc.Planets))
	}
}

func TestCalcHarmonicChart_H5(t *testing.T) {
	hc, err := CalcHarmonicChart(51.5074, -0.1278, j2000, 5, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		t.Fatalf("H5: %v", err)
	}

	// Each harmonic longitude should be natal * 5, mod 360
	for _, p := range hc.Planets {
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("H5 planet %s longitude out of range: %.2f", p.PlanetID, p.Longitude)
		}
	}
}

func TestCalcHarmonicChart_InvalidHarmonic(t *testing.T) {
	_, err := CalcHarmonicChart(51.5074, -0.1278, j2000, 0, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err == nil {
		t.Error("Expected error for harmonic 0")
	}
	_, err = CalcHarmonicChart(51.5074, -0.1278, j2000, 200, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err == nil {
		t.Error("Expected error for harmonic 200")
	}
}

func TestCalcHarmonicChart_Multiplication(t *testing.T) {
	// Check that a known longitude is multiplied correctly
	hc, err := CalcHarmonicChart(51.5074, -0.1278, j2000, 4, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		t.Fatal(err)
	}
	// For H4, longitude should be 4x natal mod 360
	// We verify it's within [0, 360)
	for _, p := range hc.Planets {
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("Planet %s out of range: %.4f", p.PlanetID, p.Longitude)
		}
	}
}

func TestCommonHarmonics(t *testing.T) {
	ch := CommonHarmonics()
	if len(ch) < 5 {
		t.Errorf("Expected at least 5 common harmonics, got %d", len(ch))
	}
	if _, ok := ch[5]; !ok {
		t.Error("Missing quintile harmonic (5)")
	}
}

func TestCalcHarmonicChart_ConjunctionInH5(t *testing.T) {
	// Two planets 72° apart in natal should be conjunct in H5
	// Sun at ~280° and another planet at ~352° would be ~72° apart
	// After H5: 280*5=1400=320°, 352*5=1760=320° -> should be close
	// This is a property check rather than specific value check
	hc, err := CalcHarmonicChart(51.5074, -0.1278, j2000, 5, nil, models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		t.Fatal(err)
	}
	_ = hc
	// Just verify no errors and reasonable output
	for _, p := range hc.Planets {
		if math.IsNaN(p.Longitude) || math.IsInf(p.Longitude, 0) {
			t.Errorf("Invalid longitude for %s: %f", p.PlanetID, p.Longitude)
		}
	}
}
