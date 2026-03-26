//go:build ignore

// Check if our calculation parameters match SF
package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

const nJD = 2450800.900009
const lat, lon = 30.9, 121.15

func main() {
	sweph.Init("/home/ecs-user/SolarSage/ephe")

	// SF Meta values:
	// DeltaT = +62s; ET = 9:37:02 am Dec 18 1997; JDE = 2450800.900729
	// ST(0°) = 15:24:10; LST = 23:28:46; Ob = 23°26'22''

	jde := 2450800.900729

	// Check DeltaT
	deltaT := sweph.DeltaT(nJD)
	fmt.Printf("DeltaT (Swiss Eph at nJD): %.1f s (SF says +62s)\n", deltaT*86400)

	deltaT1997 := sweph.DeltaT(sweph.JulDay(1997, 12, 18, 9.5, true))
	fmt.Printf("DeltaT (Swiss Eph at 1997-12-18): %.1f s\n", deltaT1997*86400)

	// Check Obliquity (黄道倾角)
	ob, _ := sweph.Obliquity(nJD)
	fmt.Printf("\nObliquity (nJD): %.8f° = 23°%.0f'%.1f''\n",
		ob, (ob-23)*60, ((ob-23)*60-26)*60)

	// SF says Ob = 23°26'22'' = 23.439444°
	sfOb := 23.0 + 26.0/60 + 22.0/3600
	fmt.Printf("Obliquity (SF Meta): %.8f° = 23°26'22''\n", sfOb)
	fmt.Printf("Diff: %.6f'' (arcsec)\n", (ob-sfOb)*3600)

	obJDE, _ := sweph.Obliquity(jde)
	fmt.Printf("Obliquity (JDE): %.8f°\n", obJDE)

	// Check ARMC (Right Ascension of MC)
	// SF: LST = 23:28:46 = 23.4794 hours = 352.191°
	sfLST := (23.0 + 28.0/60 + 46.0/3600) * 15 // hours to degrees
	fmt.Printf("\nLST (SF Meta): %.6f°\n", sfLST)

	// Our ARMC from houses
	result, _ := sweph.Houses(nJD, lat, lon, 'P')
	fmt.Printf("ARMC (our, nJD): %.6f°\n", result.ARMC)
	fmt.Printf("Diff from SF LST: %.6f° (%.2f')\n", result.ARMC-sfLST, (result.ARMC-sfLST)*60)

	// Check with JDE
	resultJDE, _ := sweph.Houses(jde, lat, lon, 'P')
	fmt.Printf("ARMC (JDE): %.6f°\n", resultJDE.ARMC)
	fmt.Printf("Diff from SF LST: %.6f°\n", resultJDE.ARMC-sfLST)

	// ASC comparison
	fmt.Println("\n=== ASC Comparison ===")
	fmt.Printf("ASC (nJD): %.6f°\n", result.ASC)
	fmt.Printf("ASC (JDE): %.6f°\n", resultJDE.ASC)

	sfASC := 96.529167 // 06°Cancer31'45''
	fmt.Printf("ASC (SF): %.6f°\n", sfASC)
	fmt.Printf("Diff (nJD): %.6f° (%.2f')\n", result.ASC-sfASC, (result.ASC-sfASC)*60)
	fmt.Printf("Diff (JDE): %.6f° (%.2f')\n", resultJDE.ASC-sfASC, (resultJDE.ASC-sfASC)*60)
}
