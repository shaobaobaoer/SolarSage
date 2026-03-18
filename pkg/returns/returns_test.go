package returns

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

// J2000.0 = 2000-01-01 12:00 UT
const j2000 = 2451545.0

func TestCalcSolarReturn(t *testing.T) {
	// Calculate solar return for J2000.0, searching from ~1 year later
	rc, err := CalcSolarReturn(ReturnInput{
		NatalJD:     j2000,
		NatalLat:    51.5074,
		NatalLon:    -0.1278,
		SearchJD:    j2000 + 360,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcSolarReturn: %v", err)
	}

	// Return should be within ~365.25 days of natal
	daysDiff := rc.ReturnJD - j2000
	if daysDiff < 365 || daysDiff > 366 {
		t.Errorf("Solar return JD diff = %.2f, expected ~365.25", daysDiff)
	}

	// Sun at return should match natal Sun within 0.01 degree
	natalSunLon, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, j2000)
	diff := math.Abs(normDiff(rc.PlanetLon, natalSunLon))
	if diff > 0.01 {
		t.Errorf("Sun return accuracy: %.4f degrees off (natal=%.4f, return=%.4f)",
			diff, natalSunLon, rc.PlanetLon)
	}

	if rc.ReturnType != "solar" {
		t.Errorf("ReturnType = %s, want solar", rc.ReturnType)
	}
	if rc.Chart == nil {
		t.Error("Chart is nil")
	}
	if rc.Age < 0.99 || rc.Age > 1.01 {
		t.Errorf("Age = %.4f, expected ~1.0", rc.Age)
	}
}

func TestCalcLunarReturn(t *testing.T) {
	rc, err := CalcLunarReturn(ReturnInput{
		NatalJD:     j2000,
		NatalLat:    51.5074,
		NatalLon:    -0.1278,
		SearchJD:    j2000 + 25,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcLunarReturn: %v", err)
	}

	daysDiff := rc.ReturnJD - j2000
	if daysDiff < 25 || daysDiff > 30 {
		t.Errorf("Lunar return JD diff = %.2f, expected 27-28", daysDiff)
	}

	natalMoonLon, _, _ := chart.CalcPlanetLongitude(models.PlanetMoon, j2000)
	diff := math.Abs(normDiff(rc.PlanetLon, natalMoonLon))
	if diff > 0.01 {
		t.Errorf("Moon return accuracy: %.4f degrees off", diff)
	}

	if rc.ReturnType != "lunar" {
		t.Errorf("ReturnType = %s, want lunar", rc.ReturnType)
	}
}

func TestCalcSolarReturnSeries(t *testing.T) {
	series, err := CalcSolarReturnSeries(ReturnInput{
		NatalJD:     j2000,
		NatalLat:    51.5074,
		NatalLon:    -0.1278,
		SearchJD:    j2000 + 360,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, 3)
	if err != nil {
		t.Fatalf("CalcSolarReturnSeries: %v", err)
	}

	if len(series) != 3 {
		t.Fatalf("Expected 3 returns, got %d", len(series))
	}

	// Each return should be ~365.25 days after the previous
	for i := 1; i < len(series); i++ {
		gap := series[i].ReturnJD - series[i-1].ReturnJD
		if gap < 364 || gap > 367 {
			t.Errorf("Gap between return %d and %d = %.2f days", i-1, i, gap)
		}
	}
}

func TestCalcLunarReturnSeries(t *testing.T) {
	series, err := CalcLunarReturnSeries(ReturnInput{
		NatalJD:     j2000,
		NatalLat:    51.5074,
		NatalLon:    -0.1278,
		SearchJD:    j2000 + 25,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, 3)
	if err != nil {
		t.Fatalf("CalcLunarReturnSeries: %v", err)
	}

	if len(series) != 3 {
		t.Fatalf("Expected 3 returns, got %d", len(series))
	}

	for i := 1; i < len(series); i++ {
		gap := series[i].ReturnJD - series[i-1].ReturnJD
		if gap < 26 || gap > 29 {
			t.Errorf("Gap between lunar return %d and %d = %.2f days", i-1, i, gap)
		}
	}
}

func TestCalcPlanetReturn(t *testing.T) {
	// Mercury return (synodic period ~116 days)
	rc, err := CalcPlanetReturn(ReturnInput{
		NatalJD:     j2000,
		NatalLat:    51.5074,
		NatalLon:    -0.1278,
		SearchJD:    j2000 + 100,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, models.PlanetMercury)
	if err != nil {
		t.Fatalf("CalcPlanetReturn (Mercury): %v", err)
	}

	natalLon, _, _ := chart.CalcPlanetLongitude(models.PlanetMercury, j2000)
	diff := math.Abs(normDiff(rc.PlanetLon, natalLon))
	if diff > 0.01 {
		t.Errorf("Mercury return accuracy: %.4f degrees off", diff)
	}
}

func TestCalcPlanetReturn_Venus(t *testing.T) {
	rc, err := CalcPlanetReturn(ReturnInput{
		NatalJD:  j2000, NatalLat: 51.5074, NatalLon: -0.1278,
		SearchJD: j2000 + 500, OrbConfig: models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, models.PlanetVenus)
	if err != nil {
		t.Fatalf("Venus return: %v", err)
	}
	if rc.ReturnType != "planetary" {
		t.Errorf("ReturnType = %s", rc.ReturnType)
	}
}

func TestCalcPlanetReturn_Mars(t *testing.T) {
	rc, err := CalcPlanetReturn(ReturnInput{
		NatalJD:  j2000, NatalLat: 51.5074, NatalLon: -0.1278,
		SearchJD: j2000 + 600, OrbConfig: models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, models.PlanetMars)
	if err != nil {
		t.Fatalf("Mars return: %v", err)
	}
	_ = rc
}

func TestCalcPlanetReturn_Jupiter(t *testing.T) {
	rc, err := CalcPlanetReturn(ReturnInput{
		NatalJD:  j2000, NatalLat: 51.5074, NatalLon: -0.1278,
		SearchJD: j2000 + 4000, OrbConfig: models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, models.PlanetJupiter)
	if err != nil {
		t.Fatalf("Jupiter return: %v", err)
	}
	_ = rc
}

func TestCalcPlanetReturn_Saturn(t *testing.T) {
	rc, err := CalcPlanetReturn(ReturnInput{
		NatalJD:  j2000, NatalLat: 51.5074, NatalLon: -0.1278,
		SearchJD: j2000 + 10000, OrbConfig: models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, models.PlanetSaturn)
	if err != nil {
		t.Fatalf("Saturn return: %v", err)
	}
	_ = rc
}

func TestCalcReturnSeriesZeroCount(t *testing.T) {
	_, err := CalcSolarReturnSeries(ReturnInput{
		NatalJD:  j2000,
		NatalLat: 51.5074, NatalLon: -0.1278,
		SearchJD:    j2000 + 360,
		OrbConfig:   models.DefaultOrbConfig(),
		HouseSystem: models.HousePlacidus,
	}, 0)
	if err == nil {
		t.Error("Expected error for count=0")
	}
}
