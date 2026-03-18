package antiscia

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcAntiscia(t *testing.T) {
	// 10° Aries (10°) -> antiscion at 20° Virgo (170°)
	a := CalcAntiscia(10)
	expected := 170.0
	if math.Abs(a-expected) > 0.01 {
		t.Errorf("Antiscia(10°) = %.2f, want %.2f", a, expected)
	}
}

func TestCalcAntiscia_Symmetry(t *testing.T) {
	// Antiscia pairs: Aries↔Virgo, Taurus↔Leo, Gemini↔Cancer
	tests := []struct {
		lon      float64
		expected float64
	}{
		{0, 180},     // 0° Aries -> 0° Libra (180° - 0° = 180°)
		{90, 90},     // 0° Cancer -> 0° Cancer (self-symmetric on solstice axis)
		{270, 270},   // 0° Capricorn -> 0° Capricorn (self-symmetric)
		{45, 135},    // 15° Taurus -> 15° Leo
	}
	for _, tt := range tests {
		got := CalcAntiscia(tt.lon)
		diff := math.Abs(got - tt.expected)
		if diff > 180 {
			diff = 360 - diff
		}
		if diff > 0.01 {
			t.Errorf("Antiscia(%.0f°) = %.2f, want %.2f", tt.lon, got, tt.expected)
		}
	}
}

func TestCalcContraAntiscia(t *testing.T) {
	// 10° -> 350°
	ca := CalcContraAntiscia(10)
	if math.Abs(ca-350) > 0.01 {
		t.Errorf("ContraAntiscia(10°) = %.2f, want 350", ca)
	}
}

func TestCalcChartAntiscia(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 10},
		{PlanetID: models.PlanetMoon, Longitude: 170},
	}
	points := CalcChartAntiscia(positions)
	if len(points) != 2 {
		t.Fatalf("Expected 2 points, got %d", len(points))
	}
	// Sun and Moon are antiscia of each other (10° and 170° mirror across solstice)
	if points[0].AntisciaSign == "" {
		t.Error("Missing antiscia sign")
	}
}

func TestFindAntisciaPairs(t *testing.T) {
	// Sun at 10°, Moon at 170°: they are antiscia (mirror = 170° and 10°)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 10},
		{PlanetID: models.PlanetMoon, Longitude: 170},
	}
	pairs := FindAntisciaPairs(positions, 2.0)
	found := false
	for _, p := range pairs {
		if p.Type == "antiscia" {
			found = true
		}
	}
	if !found {
		t.Error("Expected to find antiscia pair between Sun(10°) and Moon(170°)")
	}
}

func TestFindAntisciaPairs_DefaultOrb(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 10},
	}
	pairs := FindAntisciaPairs(positions, 0)
	_ = pairs // should not panic
}
