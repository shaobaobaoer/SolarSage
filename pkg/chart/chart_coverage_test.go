package chart

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcDoubleChart_NilSpecialPoints(t *testing.T) {
	lat := 39.9042
	lon := 116.4074
	jd1 := 2448057.5208
	jd2 := 2451545.0
	planets := []models.PlanetID{models.PlanetSun, models.PlanetMoon}
	orbs := models.DefaultOrbConfig()

	inner, outer, _, err := CalcDoubleChart(
		lat, lon, jd1, planets,
		lat, lon, jd2, planets,
		nil, orbs, models.HousePlacidus,
	)
	if err != nil {
		t.Fatalf("CalcDoubleChart error: %v", err)
	}
	if inner == nil || outer == nil {
		t.Fatal("Charts should not be nil")
	}
}

func TestCalcDoubleChart_AllSpecialPoints(t *testing.T) {
	lat := 39.9042
	lon := 116.4074
	jd1 := 2448057.5208
	jd2 := 2451545.0
	planets := []models.PlanetID{models.PlanetSun, models.PlanetMoon}
	orbs := models.DefaultOrbConfig()
	sp := &models.SpecialPointsConfig{
		InnerPoints: []models.SpecialPointID{
			models.PointASC, models.PointMC, models.PointDSC, models.PointIC,
			models.PointVertex, models.PointAntiVertex, models.PointEastPoint,
			models.PointLotFortune, models.PointLotSpirit,
		},
		OuterPoints: []models.SpecialPointID{models.PointASC},
	}

	_, _, _, err := CalcDoubleChart(
		lat, lon, jd1, planets,
		lat, lon, jd2, planets,
		sp, orbs, models.HousePlacidus,
	)
	if err != nil {
		t.Fatalf("CalcDoubleChart with all special points error: %v", err)
	}
}

func TestCalcPlanetLongitude_AllPlanets(t *testing.T) {
	jdUT := 2451545.0
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeTrue,
		models.PlanetNorthNodeMean, models.PlanetSouthNode,
		models.PlanetLilithMean, models.PlanetLilithTrue,
	}
	for _, pid := range planets {
		lon, _, err := CalcPlanetLongitude(pid, jdUT)
		if err != nil {
			t.Errorf("CalcPlanetLongitude(%s) error: %v", pid, err)
			continue
		}
		if lon < 0 || lon >= 360 {
			t.Errorf("CalcPlanetLongitude(%s) = %f, want 0-360", pid, lon)
		}
	}
}

func TestCalcPlanetLongitude_UnknownPlanet(t *testing.T) {
	_, _, err := CalcPlanetLongitude(models.PlanetID("UNKNOWN"), 2451545.0)
	if err == nil {
		t.Error("Expected error for unknown planet")
	}
}

func TestCalcSpecialPointLongitude_Coverage(t *testing.T) {
	jdUT := 2451545.0
	lat := 39.9042
	lon := 116.4074

	points := []models.SpecialPointID{
		models.PointASC, models.PointMC, models.PointDSC, models.PointIC,
		models.PointVertex, models.PointAntiVertex, models.PointEastPoint,
		models.PointLotFortune, models.PointLotSpirit,
	}
	for _, sp := range points {
		result, err := CalcSpecialPointLongitude(sp, lat, lon, jdUT, models.HousePlacidus)
		if err != nil {
			t.Errorf("CalcSpecialPointLongitude(%s) error: %v", sp, err)
			continue
		}
		if result < 0 || result >= 360 {
			t.Errorf("CalcSpecialPointLongitude(%s) = %f, want 0-360", sp, result)
		}
	}

	// Unknown point
	_, err := CalcSpecialPointLongitude(models.SpecialPointID("UNKNOWN"), lat, lon, jdUT, models.HousePlacidus)
	if err == nil {
		t.Error("Expected error for unknown special point")
	}
}

func TestFindHouse_WraparoundCusps(t *testing.T) {
	cusps := []float64{350, 20, 50, 80, 110, 140, 170, 200, 230, 260, 290, 320}
	h := FindHouseForLongitude(355, cusps)
	if h != 1 {
		t.Errorf("FindHouseForLongitude(355, wraparound) = %d, want 1", h)
	}
	h2 := FindHouseForLongitude(10, cusps)
	if h2 != 1 {
		t.Errorf("FindHouseForLongitude(10, wraparound) = %d, want 1", h2)
	}
}

func TestFindHouseForLongitude_NilCusps(t *testing.T) {
	h := FindHouseForLongitude(45, nil)
	if h != 2 {
		t.Errorf("FindHouseForLongitude(45, nil) = %d, want 2", h)
	}
}

func TestIsDayChart_Cases(t *testing.T) {
	// Sun at 0°, ASC at 90° => DSC at 270°
	result := isDayChart(0, 90)
	if !result {
		t.Error("isDayChart(0, 90) should be true")
	}

	// Sun at 180°, ASC at 90° => DSC at 270°
	result2 := isDayChart(180, 90)
	_ = result2
}

func TestCalcNatalFixedHouses_Placidus(t *testing.T) {
	cusps, err := CalcNatalFixedHouses(39.9, 116.4, 2451545.0, models.HousePlacidus)
	if err != nil {
		t.Fatalf("CalcNatalFixedHouses error: %v", err)
	}
	if len(cusps) != 12 {
		t.Errorf("Expected 12 cusps, got %d", len(cusps))
	}
}

func TestWrapAngle_Cases(t *testing.T) {
	tests := []struct {
		a, want float64
	}{
		{0, 0},
		{180, -180},
		{-180, -180},
		{90, 90},
		{-90, -90},
		{270, -90},
		{360, 0},
	}
	for _, tt := range tests {
		got := WrapAngle(tt.a)
		if got != tt.want {
			t.Errorf("WrapAngle(%f) = %f, want %f", tt.a, got, tt.want)
		}
	}
}
