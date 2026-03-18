package lots

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcStandardLots(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 280},
		{PlanetID: models.PlanetMoon, Longitude: 120},
		{PlanetID: models.PlanetMercury, Longitude: 290},
		{PlanetID: models.PlanetVenus, Longitude: 310},
		{PlanetID: models.PlanetMars, Longitude: 180},
		{PlanetID: models.PlanetJupiter, Longitude: 45},
		{PlanetID: models.PlanetSaturn, Longitude: 70},
	}

	asc := 15.0 // ASC at 15° Aries
	results := CalcStandardLots(positions, asc, true)

	if len(results) == 0 {
		t.Fatal("Expected at least some lots")
	}

	// Check Lot of Fortune exists
	found := false
	for _, r := range results {
		if r.Name == "Lot of Fortune" {
			found = true
			if r.Sign == "" {
				t.Error("Fortune lot has empty sign")
			}
			// Day formula: ASC + Moon - Sun = 15 + 120 - 280 = -145 -> 215° (Scorpio)
			break
		}
	}
	if !found {
		t.Error("Lot of Fortune not found")
	}

	// Check Lot of Spirit exists
	found = false
	for _, r := range results {
		if r.Name == "Lot of Spirit" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Lot of Spirit not found")
	}
}

func TestCalcStandardLots_NightChart(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 280},
		{PlanetID: models.PlanetMoon, Longitude: 120},
		{PlanetID: models.PlanetMercury, Longitude: 290},
		{PlanetID: models.PlanetVenus, Longitude: 310},
		{PlanetID: models.PlanetMars, Longitude: 180},
		{PlanetID: models.PlanetJupiter, Longitude: 45},
		{PlanetID: models.PlanetSaturn, Longitude: 70},
	}

	asc := 15.0
	dayResults := CalcStandardLots(positions, asc, true)
	nightResults := CalcStandardLots(positions, asc, false)

	// Fortune day and night should differ (reversed formula)
	var dayFortune, nightFortune float64
	for _, r := range dayResults {
		if r.Name == "Lot of Fortune" {
			dayFortune = r.Longitude
		}
	}
	for _, r := range nightResults {
		if r.Name == "Lot of Fortune" {
			nightFortune = r.Longitude
		}
	}
	if dayFortune == nightFortune {
		t.Error("Day and night Fortune should differ when Sun != Moon")
	}
}

func TestCalcCustomLot(t *testing.T) {
	// Simple: ASC(0) + bodyA(90) - bodyB(30) = 60
	result := CalcCustomLot(0, 90, 30, true, true)
	if result != 60 {
		t.Errorf("Custom lot = %.2f, want 60", result)
	}

	// Night reversal: ASC(0) + bodyB(30) - bodyA(90) = -60 -> 300
	result = CalcCustomLot(0, 90, 30, false, true)
	if result != 300 {
		t.Errorf("Night custom lot = %.2f, want 300", result)
	}
}
