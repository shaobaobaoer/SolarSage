//go:build ignore

// Check if we can use SF's obliquity value to match ASC/MC
package main

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

const lat, lon = 30.9, 121.15

// ASC calculation formula (simplified):
// tan(ASC) = cos(ARMC) / (cos(eps) * sin(ARMC) + sin(eps) * tan(lat))
// where eps = obliquity

func calcASC(armc, obliquity, latitude float64) float64 {
	// Simplified formula for Placidus
	// ASC = atan2(cos(ARMC), cos(eps)*sin(ARMC) + sin(eps)*tan(lat))
	cosARMC := math.Cos(armc * math.Pi / 180)
	sinARMC := math.Sin(armc * math.Pi / 180)
	cosEps := math.Cos(obliquity * math.Pi / 180)
	sinEps := math.Sin(obliquity * math.Pi / 180)
	tanLat := math.Tan(latitude * math.Pi / 180)

	ascRad := math.Atan2(cosARMC, cosEps*sinARMC+sinEps*tanLat)
	asc := ascRad * 180 / math.Pi
	if asc < 0 {
		asc += 360
	}
	return asc
}

func main() {
	sweph.Init("/home/ecs-user/SolarSage/ephe")

	nJD := 2450800.900009

	// Our values
	ourOb, _ := sweph.Obliquity(nJD)
	result, _ := sweph.Houses(nJD, lat, lon, 'P')
	ourARMC := result.ARMC
	ourASC := result.ASC

	// SF values
	sfOb := 23.0 + 26.0/60 + 22.0/3600 // 23°26'22''
	sfLST := (23.0 + 28.0/60 + 46.0/3600) * 15
	sfASC := 96.0 + 31.0/60 + 45.0/3600 // 06°Cancer31'45''

	fmt.Println("=== Testing ASC calculation with different obliquity values ===")
	fmt.Printf("Our obliquity: %.8f°\n", ourOb)
	fmt.Printf("SF obliquity:%.8f°\n", sfOb)
	fmt.Printf("Diff: %.6f''\n\n", (ourOb-sfOb)*3600)

	// Calculate ASC with our obliquity
	calcASCOurOb := calcASC(ourARMC, ourOb, lat)
	fmt.Printf("ASC calc with our obliquity: %.6f°\n", calcASCOurOb)
	fmt.Printf("Our actual ASC: %.6f°\n", ourASC)
	fmt.Printf("Diff: %.6f°\n\n", calcASCOurOb-ourASC)

	// Calculate ASC with SF obliquity
	calcASCSFOb := calcASC(sfLST, sfOb, lat)
	fmt.Printf("ASC calc with SF LST + SF obliquity: %.6f°\n", calcASCSFOb)
	fmt.Printf("SF ASC: %.6f°\n", sfASC)
	fmt.Printf("Diff: %.6f°\n\n", calcASCSFOb-sfASC)

	// Calculate ASC with our ARMC but SF obliquity
	calcASCMixed := calcASC(ourARMC, sfOb, lat)
	fmt.Printf("ASC calc with our ARMC + SF obliquity: %.6f°\n", calcASCMixed)
	fmt.Printf("Diff from our ASC: %.6f°\n", calcASCMixed-ourASC)
	fmt.Printf("Diff from SF ASC: %.6f°\n", calcASCMixed-sfASC)

	// The key question: can we make Swiss Eph use SF's obliquity?
	fmt.Println("\n=== Can we override Swiss Eph obliquity? ===")
	fmt.Println("Swiss Ephemeris uses IAU2006 obliquity model")
	fmt.Println("SF might use an older model or fixed value")
	fmt.Println("The difference is ~9 arcseconds, which is significant for precise calculations")

	// Impact analysis
	fmt.Println("\n=== Impact on Transit Timing ===")
	ascDiff := ourASC - sfASC
	plutoSpeed := 0.012 // °/day
	timeError := ascDiff / plutoSpeed * 24 * 60
	fmt.Printf("ASC diff: %.6f° (%.2f\")\n", ascDiff, ascDiff*3600)
	fmt.Printf("Time error for Pluto Opposition ASC: %.1f minutes\n", timeError)

	// This matches the 117 min deviation we saw!
}
