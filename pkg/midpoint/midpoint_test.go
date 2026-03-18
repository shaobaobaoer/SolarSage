package midpoint

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcMidpoints(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 0, Speed: 1},
		{PlanetID: models.PlanetMoon, Longitude: 60, Speed: 13},
		{PlanetID: models.PlanetMars, Longitude: 120, Speed: 0.5},
	}

	tree := CalcMidpoints(positions, 1.5)

	// 3 planets -> 3 midpoints
	if len(tree.Midpoints) != 3 {
		t.Errorf("Expected 3 midpoints, got %d", len(tree.Midpoints))
	}

	// Check Sun/Moon midpoint = 30°
	found := false
	for _, mp := range tree.Midpoints {
		if (mp.BodyA == "SUN" && mp.BodyB == "MOON") || (mp.BodyA == "MOON" && mp.BodyB == "SUN") {
			found = true
			if math.Abs(mp.Longitude-30) > 0.01 {
				t.Errorf("Sun/Moon midpoint = %.2f, expected 30", mp.Longitude)
			}
		}
	}
	if !found {
		t.Error("Sun/Moon midpoint not found")
	}

	// 90° dial sort should have 3 entries
	if len(tree.SortedDial90) != 3 {
		t.Errorf("Expected 3 dial entries, got %d", len(tree.SortedDial90))
	}

	// Dial values should be in [0, 90) and sorted
	for i := 0; i < len(tree.SortedDial90)-1; i++ {
		if tree.SortedDial90[i].Dial90 > tree.SortedDial90[i+1].Dial90 {
			t.Error("90° dial not sorted")
		}
	}
}

func TestMidpointWraparound(t *testing.T) {
	// Test midpoint wrapping around 0°
	mp := midpoint(350, 10)
	if math.Abs(mp-0) > 0.01 && math.Abs(mp-360) > 0.01 {
		t.Errorf("Midpoint(350, 10) = %.2f, expected ~0", mp)
	}
}

func TestHarmonicDiff(t *testing.T) {
	// Conjunction: 0° difference
	diff := harmonicDiff(100, 100, 1)
	if diff > 0.01 {
		t.Errorf("harmonicDiff(100, 100, 1) = %.2f, expected ~0", diff)
	}

	// Opposition: 180° difference on div=1 axis
	diff = harmonicDiff(0, 180, 1)
	if diff > 0.01 {
		t.Errorf("harmonicDiff(0, 180, 1) = %.2f, expected ~0", diff)
	}

	// Square: 90° difference on div=2 axis
	diff = harmonicDiff(0, 90, 2)
	if diff > 0.01 {
		t.Errorf("harmonicDiff(0, 90, 2) = %.2f, expected ~0", diff)
	}
}

func TestCalcMidpoints_Activations(t *testing.T) {
	// Mars at 30° is exactly on the Sun(0)/Moon(60) midpoint
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 0},
		{PlanetID: models.PlanetMoon, Longitude: 60},
		{PlanetID: models.PlanetMars, Longitude: 30},
	}

	tree := CalcMidpoints(positions, 2.0)
	found := false
	for _, a := range tree.Activations {
		if a.Planet == "MARS" && a.MidpointOf[0] == "SUN" && a.MidpointOf[1] == "MOON" {
			found = true
			if a.Orb > 0.01 {
				t.Errorf("Mars on Sun/Moon midpoint orb = %.4f, expected ~0", a.Orb)
			}
		}
	}
	if !found {
		t.Error("Expected Mars activation on Sun/Moon midpoint")
	}
}

func TestCalcMidpoints_DefaultOrb(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 0},
		{PlanetID: models.PlanetMoon, Longitude: 60},
	}
	// orb <= 0 should default to 1.5
	tree := CalcMidpoints(positions, 0)
	if tree == nil {
		t.Fatal("CalcMidpoints returned nil")
	}
}
