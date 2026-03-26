//go:build ignore

// Compare our Natal positions with SF Meta (precise to arcsec)
package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

const nJD = 2450800.900009
const lat, lon = 30.9, 121.15

// SF Meta positions (precise to arcsec from meta file)
// Format: degrees, minutes, seconds
type dms struct {
	d, m, s int
}

func (d dms) toDecimal() float64 {
	return float64(d.d) + float64(d.m)/60 + float64(d.s)/3600
}

// Sign offsets
var signOffset = map[string]float64{
	"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
	"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
	"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
}

var sfMeta = []struct {
	name string
	pid  models.PlanetID
	dms  dms
	sign string
}{
	{"Sun", models.PlanetSun, dms{26, 29, 59}, "Sagittarius"},
	{"Moon", models.PlanetMoon, dms{18, 6, 58}, "Leo"},
	{"Mercury", models.PlanetMercury, dms{23, 56, 0}, "Sagittarius"},
	{"Venus", models.PlanetVenus, dms{2, 33, 9}, "Aquarius"},
	{"Mars", models.PlanetMars, dms{0, 5, 50}, "Aquarius"},
	{"Jupiter", models.PlanetJupiter, dms{19, 32, 52}, "Aquarius"},
	{"Saturn", models.PlanetSaturn, dms{13, 32, 15}, "Aries"},
	{"Uranus", models.PlanetUranus, dms{6, 25, 29}, "Aquarius"},
	{"Neptune", models.PlanetNeptune, dms{28, 28, 3}, "Capricorn"},
	{"Pluto", models.PlanetPluto, dms{6, 16, 27}, "Sagittarius"},
	{"Chiron", models.PlanetChiron, dms{14, 1, 21}, "Scorpio"},
	{"NorthNode", models.PlanetNorthNodeMean, dms{14, 26, 45}, "Virgo"},
}

func main() {
	sweph.Init("/home/ecs-user/SolarSage/ephe")

	fmt.Printf("%-12s %12s %12s %12s %10s\n", "Planet", "Our Calc", "SF Meta", "Diff", "Arcsec")
	fmt.Println("-------------------------------------------------------------")

	for _, sf := range sfMeta {
		ourLon, _, _ := chart.CalcPlanetLongitude(sf.pid, nJD)
		sfLon := signOffset[sf.sign] + sf.dms.toDecimal()
		diff := ourLon - sfLon
		arcsec := diff * 3600

		status := "✓"
		if arcsec > 1 || arcsec < -1 {
			status = "⚠"
		}

		fmt.Printf("%-12s %12.6f %12.6f %12.6f %+.1f\" %s\n",
			sf.name, ourLon, sfLon, diff, arcsec, status)
	}

	// Check ASC/MC
	fmt.Println("\n--- Angles ---")
	asc, _ := chart.CalcSpecialPointLongitude(models.PointASC, lat, lon, nJD, models.HousePlacidus)
	mc, _ := chart.CalcSpecialPointLongitude(models.PointMC, lat, lon, nJD, models.HousePlacidus)

	sfASC := signOffset["Cancer"] + dms{6, 31, 45}.toDecimal()
	sfMC := signOffset["Pisces"] + dms{21, 29, 51}.toDecimal()

	fmt.Printf("%-12s %12.6f %12.6f %12.6f %+.1f\"\n",
		"ASC", asc, sfASC, asc-sfASC, (asc-sfASC)*3600)
	fmt.Printf("%-12s %12.6f %12.6f %12.6f %+.1f\"\n",
		"MC", mc, sfMC, mc-sfMC, (mc-sfMC)*3600)
}
