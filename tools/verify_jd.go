package main

import (
	"fmt"
	"github.com/yourusername/SolarSage/pkg/julian"
)

func main() {
	// Test case 1 JD from Solar Fire
	// JDE = 2450800.900729, DeltaT = +62s
	// JD_UT = JDE - DeltaT/86400
	
	jde := 2450800.900729
	deltaT := 62.0 // seconds
	jdUT := jde - deltaT/86400.0
	
	fmt.Printf("Solar Fire metadata:\n")
	fmt.Printf("  JDE (Ephemeris Time): %.6f\n", jde)
	fmt.Printf("  DeltaT: %.0f seconds\n", deltaT)
	fmt.Printf("  JD_UT (Universal Time): %.6f\n", jdUT)
	fmt.Printf("  JD_UT (from test): 2450800.900009\n")
	fmt.Printf("  Difference: %.6f\n", jdUT - 2450800.900009)
	
	// Convert to datetime in different timezones
	timezones := []string{"UTC", "Australia/Perth", "Asia/Shanghai"}
	
	fmt.Printf("\nJD %.6f in different timezones:\n", jdUT)
	for _, tz := range timezones {
		dt, err := julian.JDToDateTime(jdUT, tz)
		if err != nil {
			fmt.Printf("  %s: ERROR %v\n", tz, err)
		} else {
			fmt.Printf("  %s: %s\n", tz, dt)
		}
	}
}