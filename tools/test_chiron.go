package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	if _, err := os.Stat(ephePath); err != nil {
		ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
	}
	fmt.Printf("Ephe path: %s\n", ephePath)
	sweph.Init(ephePath)
	defer sweph.Close()

	// Test Chiron position at 2026-03-01 UTC
	jd := sweph.JulDay(2026, 3, 1, 0, true)
	fmt.Printf("JD for 2026-03-01: %.6f\n", jd)

	// Direct sweph call
	fmt.Println("\n=== Direct sweph.CalcUT ===")
	res, err := sweph.CalcUT(jd, sweph.SE_CHIRON)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("Chiron (direct): lon=%.4f° lat=%.4f° speed=%.4f°/d\n",
			res.Longitude, res.Latitude, res.SpeedLong)
	}

	// Via chart.CalcPlanetLongitude
	fmt.Println("\n=== chart.CalcPlanetLongitude ===")
	lon, speed, err := chart.CalcPlanetLongitude(models.PlanetChiron, jd)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("Chiron (chart): lon=%.4f° speed=%.6f°/d\n", lon, speed)
	}

	// Also check NorthNode Mean
	lon2, speed2, err := chart.CalcPlanetLongitude(models.PlanetNorthNodeMean, jd)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("NorthNode Mean: lon=%.4f° speed=%.6f°/d\n", lon2, speed2)
	}

	// All planets
	fmt.Println("\n=== All planets at 2026-03-01 ===")
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
	}
	for _, p := range planets {
		l, s, e := chart.CalcPlanetLongitude(p, jd)
		if e != nil {
			fmt.Printf("  %-20s ERROR: %v\n", p, e)
		} else {
			fmt.Printf("  %-20s lon=%10.4f° (%s)  speed=%+.6f\n", p, l, models.FormatLonDMS(l), s)
		}
	}
}
