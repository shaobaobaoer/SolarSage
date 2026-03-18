package chart_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func init() {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
}

func Example() {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars,
	}

	info, err := chart.CalcSingleChart(
		51.5074, -0.1278, 2451545.0, // London, J2000.0
		planets, models.DefaultOrbConfig(), models.HousePlacidus,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, p := range info.Planets {
		fmt.Printf("%s: %s %s\n", p.PlanetID, p.Sign, models.FormatLonDMS(p.Longitude))
	}
	// Output is dynamic based on ephemeris
}
