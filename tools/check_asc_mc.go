//go:build ignore

// Analyze ASC/MC calculation differences
package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

const nJD = 2450800.900009
const lat, lon = 30.9, 121.15

func main() {
	sweph.Init("/home/ecs-user/SolarSage/ephe")

	// SF Meta values (from meta file)
	// Ascendant 06°Cancer31'45'' = 96.529167°
	// MC 21°Pisces29'51'' = 351.497500°
	sfASC := 96.529167
	sfMC := 351.497500

	// Our calculations via chart package
	ourASC, _ := chart.CalcSpecialPointLongitude(models.PointASC, lat, lon, nJD, models.HousePlacidus)
	ourMC, _ := chart.CalcSpecialPointLongitude(models.PointMC, lat, lon, nJD, models.HousePlacidus)

	fmt.Println("=== ASC/MC Comparison ===")
	fmt.Printf("ASC: Our=%.6f° SF=%.6f° Diff=%.6f° (%.1f\")\n",
		ourASC, sfASC, ourASC-sfASC, (ourASC-sfASC)*3600)
	fmt.Printf("MC:Our=%.6f° SF=%.6f° Diff=%.6f° (%.1f\")\n",
		ourMC, sfMC, ourMC-sfMC, (ourMC-sfMC)*3600)

	// Direct call to sweph.Houses
	fmt.Println("\n=== Direct sweph.Houses call ===")
	result, _ := sweph.Houses(nJD, lat, lon, 'P')
	fmt.Printf("ASC: %.6f°\n", result.ASC)
	fmt.Printf("MC:%.6f°\n", result.MC)
	fmt.Printf("ARMC: %.6f°\n", result.ARMC)

	// Compare all house cusps
	fmt.Println("\n=== House Cusps Comparison ===")
	// SF meta house cusps:
	sfCusps := []float64{
		96.529167, // 1: 06°Cancer31'45''
		118.654,// 2: 28°Cancer39'14''
		142.691,// 3: 22°Leo41'27''
		171.497,// 4: 21°Virgo29'51''
		206.129,// 5: 26°Libra07'35''
		242.993,// 6: 02°Sagittarius58'32''
		276.529,// 7: 06°Capricorn31'45''
		298.654,// 8: 28°Capricorn39'14''
		322.691,// 9: 22°Aquarius41'27''
		351.497,// 10: 21°Pisces29'51''
		26.129, // 11: 26°Aries07'35''
		62.993, // 12: 02°Gemini58'32''
	}

	fmt.Printf("%-6s %12s %12s %12s\n", "House", "Our", "SF", "Diff")
	fmt.Println("----------------------------------------")
	for i := 0; i < 12; i++ {
		diff := result.Cusps[i+1] - sfCusps[i]
		fmt.Printf("%-6d %12.6f %12.6f %12.6f\n", i+1, result.Cusps[i+1], sfCusps[i], diff)
	}

	// Impact on SA ASC timing
	fmt.Println("\n=== Impact on Transit Timing ===")
	plutoSpeed := 0.012 // °/day
	ascDiff := ourASC - sfASC
	timeError := ascDiff / plutoSpeed * 24 * 60 // minutes
	fmt.Printf("ASC diff %.6f° / Pluto speed %.3f°/day = %.1f minutes\n",
		ascDiff, plutoSpeed, timeError)
}
