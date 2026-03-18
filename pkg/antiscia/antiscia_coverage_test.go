package antiscia

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestFindAntisciaPairs_ZeroOrb(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 10.0},
		{PlanetID: models.PlanetMoon, Longitude: 170.0}, // antiscia of 10° is 170°
	}
	pairs := FindAntisciaPairs(positions, 0) // should use default orb
	if len(pairs) == 0 {
		t.Error("expected to find antiscia pair for Sun at 10° and Moon at 170°")
	}
}

func TestAngleDiff_LargeAngle(t *testing.T) {
	// Test where diff > 180
	diff := angleDiff(10, 350)
	if diff > 180 {
		t.Errorf("angleDiff(10, 350) should be <= 180, got %f", diff)
	}
	expected := 20.0
	if diff < expected-0.1 || diff > expected+0.1 {
		t.Errorf("angleDiff(10, 350) = %f, want ~%f", diff, expected)
	}
}

func TestAngleDiff_NegativeDiff(t *testing.T) {
	// b > a
	diff := angleDiff(10, 50)
	expected := 40.0
	if diff < expected-0.1 || diff > expected+0.1 {
		t.Errorf("angleDiff(10, 50) = %f, want ~%f", diff, expected)
	}
}

func TestFindAntisciaPairs_ContraAntiscia(t *testing.T) {
	// Contra-antiscia of 10° is 360-10 = 350°
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 10.0},
		{PlanetID: models.PlanetMars, Longitude: 350.0},
	}
	pairs := FindAntisciaPairs(positions, 2.0)
	found := false
	for _, p := range pairs {
		if p.Type == "contra_antiscia" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find contra-antiscia pair")
	}
}
