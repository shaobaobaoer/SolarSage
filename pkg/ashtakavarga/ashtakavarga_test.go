package ashtakavarga

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

// makeSyntheticPositions creates sidereal positions with known longitudes.
func makeSyntheticPositions(lons map[models.PlanetID]float64) []vedic.SiderealPosition {
	var positions []vedic.SiderealPosition
	for id, lon := range lons {
		positions = append(positions, vedic.SiderealPosition{
			PlanetPosition: models.PlanetPosition{PlanetID: id},
			SiderealLon:    lon,
		})
	}
	return positions
}

func TestCalcAshtakavarga_SevenTables(t *testing.T) {
	positions := makeSyntheticPositions(map[models.PlanetID]float64{
		models.PlanetSun:     10,  // Aries
		models.PlanetMoon:    70,  // Gemini
		models.PlanetMars:    130, // Leo
		models.PlanetMercury: 190, // Libra
		models.PlanetJupiter: 250, // Sagittarius
		models.PlanetVenus:   310, // Aquarius
		models.PlanetSaturn:  40,  // Taurus
	})
	ascLon := 0.0 // Aries rising

	result := CalcAshtakavarga(positions, ascLon)

	if len(result.PlanetTables) != 7 {
		t.Fatalf("Expected 7 planet tables, got %d", len(result.PlanetTables))
	}
}

func TestCalcAshtakavarga_BinduTotals(t *testing.T) {
	positions := makeSyntheticPositions(map[models.PlanetID]float64{
		models.PlanetSun:     10,
		models.PlanetMoon:    70,
		models.PlanetMars:    130,
		models.PlanetMercury: 190,
		models.PlanetJupiter: 250,
		models.PlanetVenus:   310,
		models.PlanetSaturn:  40,
	})

	result := CalcAshtakavarga(positions, 0.0)

	for _, table := range result.PlanetTables {
		// Verify Total equals sum of Bindus
		sum := 0
		for _, b := range table.Bindus {
			sum += b
		}
		if sum != table.Total {
			t.Errorf("%s: Total=%d but sum of bindus=%d", table.Planet, table.Total, sum)
		}

		// Each planet's table should have a reasonable number of bindus.
		// With 8 contributors, each having up to 11 benefic houses, the
		// theoretical max can exceed 48 for some planets.
		if table.Total < 1 || table.Total > 88 {
			t.Errorf("%s: Total=%d, out of expected range [1,88]", table.Planet, table.Total)
		}
	}
}

func TestCalcAshtakavarga_SAVIsSumOfTables(t *testing.T) {
	positions := makeSyntheticPositions(map[models.PlanetID]float64{
		models.PlanetSun:     0,
		models.PlanetMoon:    30,
		models.PlanetMars:    60,
		models.PlanetMercury: 90,
		models.PlanetJupiter: 120,
		models.PlanetVenus:   150,
		models.PlanetSaturn:  180,
	})

	result := CalcAshtakavarga(positions, 0.0)

	// SAV[i] should equal sum of all planet table bindus[i]
	for i := 0; i < 12; i++ {
		sum := 0
		for _, table := range result.PlanetTables {
			sum += table.Bindus[i]
		}
		if sum != result.SAV[i] {
			t.Errorf("SAV[%d] = %d, but sum of planet bindus = %d", i, result.SAV[i], sum)
		}
	}

	// SAVTotal should equal sum of SAV
	savSum := 0
	for _, v := range result.SAV {
		savSum += v
	}
	if savSum != result.SAVTotal {
		t.Errorf("SAVTotal = %d, but sum of SAV = %d", result.SAVTotal, savSum)
	}
}

func TestCalcAshtakavarga_SAVTotalRange(t *testing.T) {
	positions := makeSyntheticPositions(map[models.PlanetID]float64{
		models.PlanetSun:     15,
		models.PlanetMoon:    85,
		models.PlanetMars:    145,
		models.PlanetMercury: 205,
		models.PlanetJupiter: 265,
		models.PlanetVenus:   325,
		models.PlanetSaturn:  55,
	})

	result := CalcAshtakavarga(positions, 120.0)

	// The SAV total should be reasonable for any chart configuration.
	if result.SAVTotal < 100 || result.SAVTotal > 500 {
		t.Errorf("SAVTotal = %d, out of expected range [100,500]", result.SAVTotal)
	}
}

func TestCalcAshtakavarga_WithRealChart(t *testing.T) {
	// Use a real sidereal chart from vedic package.
	sc, err := vedic.CalcSiderealChart(51.5074, -0.1278, 2451545.0, vedic.AyanamsaLahiri)
	if err != nil {
		t.Fatalf("CalcSiderealChart: %v", err)
	}

	ascLon := vedic.TropicalToSidereal(sc.Angles.ASC, sc.AyanamsaValue)
	result := CalcAshtakavarga(sc.Planets, ascLon)

	if len(result.PlanetTables) != 7 {
		t.Fatalf("Expected 7 planet tables, got %d", len(result.PlanetTables))
	}

	// Every planet table total must be > 0
	for _, table := range result.PlanetTables {
		if table.Total <= 0 {
			t.Errorf("%s: Total=%d, expected > 0", table.Planet, table.Total)
		}
	}

	if result.SAVTotal <= 0 {
		t.Error("SAVTotal should be > 0 for a real chart")
	}
}

func TestCalcAshtakavarga_AllPlanetsPresent(t *testing.T) {
	positions := makeSyntheticPositions(map[models.PlanetID]float64{
		models.PlanetSun:     10,
		models.PlanetMoon:    70,
		models.PlanetMars:    130,
		models.PlanetMercury: 190,
		models.PlanetJupiter: 250,
		models.PlanetVenus:   310,
		models.PlanetSaturn:  40,
	})

	result := CalcAshtakavarga(positions, 0.0)

	seen := make(map[models.PlanetID]bool)
	for _, table := range result.PlanetTables {
		seen[table.Planet] = true
	}
	for _, p := range traditionalPlanets {
		if !seen[p] {
			t.Errorf("Missing planet table for %s", p)
		}
	}
}

func TestSignIndexFromLon(t *testing.T) {
	tests := []struct {
		lon  float64
		want int
	}{
		{0, 0},     // Aries
		{29.9, 0},  // Aries
		{30.0, 1},  // Taurus
		{359.9, 11}, // Pisces
	}
	for _, tt := range tests {
		got := signIndexFromLon(tt.lon)
		if got != tt.want {
			t.Errorf("signIndexFromLon(%.1f) = %d, want %d", tt.lon, got, tt.want)
		}
	}
}
