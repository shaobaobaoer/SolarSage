package main

import (
	"fmt"
)

func main() {
	// Pluto Sp-Na case
	sfLon := 247.217 // 7.217° Sagittarius
	ourLon := 247.2254 // Our calculation
	
	diff := ourLon - sfLon
	fmt.Printf("SF: %.3f°\n", sfLon)
	fmt.Printf("Our: %.3f°\n", ourLon)
	fmt.Printf("Raw diff: %.3f°\n", diff)
	
	// Normalize
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	fmt.Printf("Normalized diff: %.3f°\n", diff)
	
	// Abs
	if diff < 0 {
		diff = -diff
	}
	fmt.Printf("Abs diff: %.3f°\n", diff)
}