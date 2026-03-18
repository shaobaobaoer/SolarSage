package primary

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath := filepath.Join("..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	code := m.Run()
	sweph.Close()
	os.Exit(code)
}

// natalJD: 1990-01-01 12:00 UTC
var natalJD = sweph.JulDay(1990, 1, 1, 12.0, true)

const (
	testLat = 51.5 // London
	testLon = 0.0
)

func TestCalcPrimaryDirections_Basic(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
		Planets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMars,
		},
		Aspects: []models.AspectType{
			models.AspectConjunction, models.AspectOpposition, models.AspectSquare,
		},
		Key:    KeyNaibod,
		MaxAge: 80,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections: %v", err)
	}
	if len(result.Directions) == 0 {
		t.Fatal("expected at least one direction, got 0")
	}

	// All ages should be positive and within range
	for _, d := range result.Directions {
		if d.AgeExact <= 0 || d.AgeExact > 80 {
			t.Errorf("direction %s->%s age=%f out of range (0,80]",
				d.Promissor, d.Significator, d.AgeExact)
		}
		if d.Arc <= 0 {
			t.Errorf("direction %s->%s arc=%f should be positive",
				d.Promissor, d.Significator, d.Arc)
		}
	}
}

func TestCalcPrimaryDirections_BothDirectionTypes(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
		Planets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMars,
			models.PlanetJupiter, models.PlanetSaturn,
		},
		MaxAge: 100,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections: %v", err)
	}

	hasDirect := false
	hasConverse := false
	for _, d := range result.Directions {
		if d.Direction == DirectDirect {
			hasDirect = true
		}
		if d.Direction == DirectConverse {
			hasConverse = true
		}
	}
	if !hasDirect {
		t.Error("expected at least one DIRECT direction")
	}
	if !hasConverse {
		t.Error("expected at least one CONVERSE direction")
	}
}

func TestCalcPrimaryDirections_IncludesAngles(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
		Planets: []models.PlanetID{models.PlanetSun},
		Aspects: []models.AspectType{models.AspectConjunction},
		MaxAge:  100,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections: %v", err)
	}

	hasMC := false
	hasASC := false
	for _, d := range result.Directions {
		if d.Significator == "MC" || d.Promissor == "MC" {
			hasMC = true
		}
		if d.Significator == "ASC" || d.Promissor == "ASC" {
			hasASC = true
		}
	}
	if !hasMC {
		t.Error("expected MC to appear as significator or promissor")
	}
	if !hasASC {
		t.Error("expected ASC to appear as significator or promissor")
	}
}

func TestNaibodKeyConversion(t *testing.T) {
	// 1 degree of arc should be approximately 1.0146 years with Naibod key
	// arcToYears(1.0, KeyNaibod, 0) = 1.0 / 0.98556 ≈ 1.01465
	years := arcToYears(1.0, KeyNaibod, 0)
	expected := 1.0 / NaibodRate // ≈ 1.01465
	if math.Abs(years-expected) > 0.0001 {
		t.Errorf("1 degree Naibod = %f years, expected %f", years, expected)
	}
	// Should be approximately 1.0146
	if years < 1.01 || years > 1.02 {
		t.Errorf("Naibod 1 degree = %f years, expected ~1.0146", years)
	}
}

func TestNaibodRate(t *testing.T) {
	// NaibodRate should be approximately 0.98556 degrees per year
	if math.Abs(NaibodRate-0.98556) > 0.00001 {
		t.Errorf("NaibodRate = %f, expected 0.98556", NaibodRate)
	}
}

func TestArcToYears_SolarArc(t *testing.T) {
	// With solar arc key and rate of 1.0 deg/day, 10 degrees = 10 years
	years := arcToYears(10.0, KeySolarArc, 1.0)
	if math.Abs(years-10.0) > 0.0001 {
		t.Errorf("solar arc 10 deg at 1.0/day = %f years, expected 10", years)
	}

	// Fallback to Naibod when rate is 0
	years = arcToYears(10.0, KeySolarArc, 0)
	expected := 10.0 / NaibodRate
	if math.Abs(years-expected) > 0.0001 {
		t.Errorf("solar arc fallback = %f, expected %f", years, expected)
	}
}

func TestCalcPrimaryDirections_Defaults(t *testing.T) {
	// Test with minimal input - all defaults should be applied
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections with defaults: %v", err)
	}
	if len(result.Directions) == 0 {
		t.Fatal("expected directions with default config, got 0")
	}
}

func TestCalcPrimaryDirections_PtolemyKey(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
		Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		Key:     KeyPtolemy,
		MaxAge:  50,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections Ptolemy: %v", err)
	}
	// Ptolemy key should produce same results as Naibod
	for _, d := range result.Directions {
		if d.Key != KeyPtolemy {
			t.Errorf("expected key PTOLEMY, got %s", d.Key)
		}
	}
}

func TestCalcPrimaryDirections_SolarArcKey(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD: natalJD,
		GeoLat:  testLat,
		GeoLon:  testLon,
		Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		Key:     KeySolarArc,
		MaxAge:  50,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections SolarArc: %v", err)
	}
	if len(result.Directions) == 0 {
		t.Fatal("expected directions with SolarArc key, got 0")
	}
	for _, d := range result.Directions {
		if d.Key != KeySolarArc {
			t.Errorf("expected key SOLAR_ARC, got %s", d.Key)
		}
	}
}

func TestEclipticToRA(t *testing.T) {
	// At 0 Aries (lon=0), RA should also be 0
	ra := eclipticToRA(0, 23.44)
	if math.Abs(ra) > 0.001 && math.Abs(ra-360) > 0.001 {
		t.Errorf("RA at 0 Aries = %f, expected 0", ra)
	}

	// At 90 degrees (0 Cancer), RA should be close to 90 but shifted by obliquity
	ra90 := eclipticToRA(90, 23.44)
	if ra90 < 85 || ra90 > 95 {
		t.Errorf("RA at 90 lon = %f, expected near 90", ra90)
	}
}

func TestEclipticToDec(t *testing.T) {
	// At 0 Aries, declination should be 0
	dec := eclipticToDec(0, 23.44)
	if math.Abs(dec) > 0.001 {
		t.Errorf("Dec at 0 Aries = %f, expected 0", dec)
	}

	// At 90 degrees (0 Cancer), declination should equal obliquity
	dec90 := eclipticToDec(90, 23.44)
	if math.Abs(dec90-23.44) > 0.01 {
		t.Errorf("Dec at 90 lon = %f, expected ~23.44", dec90)
	}
}

func TestAscensionalDifference(t *testing.T) {
	// At equator (lat=0), AD should be 0 for any declination
	ad := ascensionalDifference(10, 0)
	if math.Abs(ad) > 0.001 {
		t.Errorf("AD at equator = %f, expected 0", ad)
	}

	// At lat=45, dec=23.44 (summer solstice), AD should be positive
	ad45 := ascensionalDifference(23.44, 45)
	if ad45 <= 0 {
		t.Errorf("AD at lat=45, dec=23.44 = %f, expected positive", ad45)
	}
}

func TestSemiArc(t *testing.T) {
	// At equator, semi-arc should be 90 for any declination
	sa := semiArc(10, 0)
	if math.Abs(sa-90) > 0.001 {
		t.Errorf("semi-arc at equator = %f, expected 90", sa)
	}

	// At northern latitudes with positive declination, diurnal SA > 90
	saN := semiArc(23.44, 51.5)
	if saN <= 90 {
		t.Errorf("semi-arc at lat=51.5, dec=23.44 = %f, expected > 90", saN)
	}
}

func TestAspectAngle(t *testing.T) {
	tests := []struct {
		asp   models.AspectType
		angle float64
	}{
		{models.AspectConjunction, 0},
		{models.AspectOpposition, 180},
		{models.AspectSquare, 90},
		{models.AspectTrine, 120},
		{models.AspectSextile, 60},
	}
	for _, tt := range tests {
		got := aspectAngle(tt.asp)
		if got != tt.angle {
			t.Errorf("aspectAngle(%s) = %f, want %f", tt.asp, got, tt.angle)
		}
	}
}
