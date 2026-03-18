package composite

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcDavisonChart_MidpointJD(t *testing.T) {
	jd1 := 2451545.0 // J2000
	jd2 := jd1 + 365.0

	dc, err := CalcDavisonChart(CompositeInput{
		Person1Lat: 51.5074, Person1Lon: -0.1278, Person1JD: jd1,
		Person2Lat: 40.7128, Person2Lon: -74.006, Person2JD: jd2,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcDavisonChart: %v", err)
	}

	expectedJD := (jd1 + jd2) / 2
	if math.Abs(dc.MidpointJD-expectedJD) > 1e-9 {
		t.Errorf("MidpointJD = %f, want %f", dc.MidpointJD, expectedJD)
	}
}

func TestCalcDavisonChart_MidpointLatLon(t *testing.T) {
	lat1, lon1 := 51.5074, -0.1278
	lat2, lon2 := 40.7128, -74.006

	dc, err := CalcDavisonChart(CompositeInput{
		Person1Lat: lat1, Person1Lon: lon1, Person1JD: j2000,
		Person2Lat: lat2, Person2Lon: lon2, Person2JD: j2000 + 365,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcDavisonChart: %v", err)
	}

	expectedLat := (lat1 + lat2) / 2
	expectedLon := (lon1 + lon2) / 2

	if math.Abs(dc.MidpointLat-expectedLat) > 1e-9 {
		t.Errorf("MidpointLat = %f, want %f", dc.MidpointLat, expectedLat)
	}
	if math.Abs(dc.MidpointLon-expectedLon) > 1e-9 {
		t.Errorf("MidpointLon = %f, want %f", dc.MidpointLon, expectedLon)
	}
}

func TestCalcDavisonChart_ChartPopulated(t *testing.T) {
	dc, err := CalcDavisonChart(CompositeInput{
		Person1Lat: 51.5074, Person1Lon: -0.1278, Person1JD: j2000,
		Person2Lat: 40.7128, Person2Lon: -74.006, Person2JD: j2000 + 365,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcDavisonChart: %v", err)
	}

	if len(dc.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(dc.Planets))
	}
	if len(dc.Houses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(dc.Houses))
	}
	if len(dc.Aspects) == 0 {
		t.Error("Expected at least some aspects")
	}

	// All planet longitudes should be in [0, 360)
	for _, p := range dc.Planets {
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("Planet %s longitude out of range: %.2f", p.PlanetID, p.Longitude)
		}
		if p.Sign == "" {
			t.Errorf("Planet %s has empty sign", p.PlanetID)
		}
	}

	// Angles should be populated
	if dc.Angles.ASC == 0 && dc.Angles.MC == 0 {
		t.Error("Expected non-zero angles")
	}
}

func TestCalcDavisonChart_DefaultPlanets(t *testing.T) {
	// Empty planets list should default to 10 planets
	dc, err := CalcDavisonChart(CompositeInput{
		Person1Lat: 51.5074, Person1Lon: -0.1278, Person1JD: j2000,
		Person2Lat: 40.7128, Person2Lon: -74.006, Person2JD: j2000 + 365,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcDavisonChart with default planets: %v", err)
	}

	if len(dc.Planets) != 10 {
		t.Errorf("Expected 10 default planets, got %d", len(dc.Planets))
	}
}
