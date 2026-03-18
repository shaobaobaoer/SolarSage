package transit

import (
	"testing"

	"github.com/anthropic/swisseph-mcp/pkg/models"
)

func TestCalcTransitEvents_WithTransitSpecialPoints(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{models.PlanetSun},
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 10,
		TransitPlanets: []models.PlanetID{models.PlanetMoon},
		SpecialPoints: &models.SpecialPointsConfig{
			TransitPoints: []models.SpecialPointID{models.PointASC, models.PointMC},
			NatalPoints:   []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		EventConfig: models.EventConfig{
			IncludeTrNa: true,
			IncludeTrTr: true,
		},
		OrbConfigTransit:      orbs,
		OrbConfigProgressions: orbs,
		OrbConfigSolarArc:     orbs,
		HouseSystem:           models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents with special points error: %v", err)
	}
	_ = events
}

func TestCalcTransitEvents_WithProgressionsSpecialPoints(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{models.PlanetSun},
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		ProgressionsConfig: &models.ProgressionsConfig{
			Enabled: true,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		},
		SpecialPoints: &models.SpecialPointsConfig{
			ProgressionsPoints: []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		EventConfig: models.EventConfig{
			IncludeTrSp: true,
			IncludeSpNa: true,
		},
		OrbConfigTransit:      orbs,
		OrbConfigProgressions: orbs,
		OrbConfigSolarArc:     orbs,
		HouseSystem:           models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents with progressions special points error: %v", err)
	}
	_ = events
}

func TestCalcTransitEvents_WithSolarArcSpecialPoints(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{models.PlanetSun},
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		},
		SpecialPoints: &models.SpecialPointsConfig{
			SolarArcPoints: []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		EventConfig: models.EventConfig{
			IncludeSaNa: true,
		},
		OrbConfigTransit:      orbs,
		OrbConfigProgressions: orbs,
		OrbConfigSolarArc:     orbs,
		HouseSystem:           models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents with solar arc special points error: %v", err)
	}
	_ = events
}

func TestDirectedAngleDiff(t *testing.T) {
	// Test with direct motion (speed > 0)
	got := directedAngleDiff(100, 10, 90, 1.0)
	if got < -180 || got > 180 {
		t.Errorf("directedAngleDiff direct = %f, want in [-180, 180]", got)
	}

	// Test with retrograde motion (speed < 0)
	got2 := directedAngleDiff(100, 10, 90, -0.5)
	if got2 < -180 || got2 > 180 {
		t.Errorf("directedAngleDiff retrograde = %f, want in [-180, 180]", got2)
	}

	// Test conjunction (aspect angle = 0)
	got3 := directedAngleDiff(15, 15, 0, 1.0)
	if got3 != 0 {
		t.Errorf("directedAngleDiff conjunction exact = %f, want 0", got3)
	}
}
