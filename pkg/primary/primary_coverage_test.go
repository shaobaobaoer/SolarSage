package primary

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestAscensionalDifference_Extreme(t *testing.T) {
	// Very high latitude (circumpolar) - tests clamping to [-1, 1]
	ad := ascensionalDifference(80, 85) // extreme declination + latitude
	if ad < -90 || ad > 90 {
		t.Errorf("ascensionalDifference out of range: %f", ad)
	}

	// Negative extreme
	ad2 := ascensionalDifference(-80, 85)
	if ad2 < -90 || ad2 > 90 {
		t.Errorf("ascensionalDifference negative extreme out of range: %f", ad2)
	}
}

func TestSemiArc_Extreme(t *testing.T) {
	// High latitude should still return valid semi-arc
	sa := semiArc(70, 80)
	if sa <= 0 {
		t.Errorf("semiArc should be positive, got %f", sa)
	}
}

func TestSemiArcNocturnal_Extreme(t *testing.T) {
	sa := semiArcNocturnal(70, 80)
	if sa < 0 {
		t.Errorf("semiArcNocturnal should be non-negative, got %f", sa)
	}
}

func TestAspectAngle_Unknown(t *testing.T) {
	angle := aspectAngle("UNKNOWN_ASPECT")
	if angle != 0 {
		t.Errorf("expected 0 for unknown aspect, got %f", angle)
	}
}

func TestCalcPrimaryDirections_WithPtolemyKey(t *testing.T) {
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD:     2451545.0,
		GeoLat:      51.5074,
		GeoLon:      -0.1278,
		Planets:     []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		Aspects:     []models.AspectType{models.AspectConjunction},
		Key:         KeyPtolemy,
		MaxAge:      30,
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections Ptolemy: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCalcPrimaryDirections_HighLatitude(t *testing.T) {
	// Test near-polar latitude to exercise extreme semi-arc paths
	result, err := CalcPrimaryDirections(PrimaryDirectionInput{
		NatalJD:     2451545.0,
		GeoLat:      65.0, // Northern latitude
		GeoLon:      25.0,
		Planets:     []models.PlanetID{models.PlanetSun},
		Aspects:     []models.AspectType{models.AspectConjunction, models.AspectOpposition},
		Key:         KeyNaibod,
		MaxAge:      20,
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcPrimaryDirections high latitude: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
