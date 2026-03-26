package main

import (
	"fmt"
	"github.com/yourusername/SolarSage/pkg/julian"
)

func main() {
	// Test JD to datetime conversion
	// JD 2461072.166667 should correspond to a specific date/time
	
	testJDs := []float64{
		2461072.166667, // Start of transit period
		2461072.5,      // Noon
		2461073.0,      // Next day
	}
	
	timezones := []string{"UTC", "Australia/Perth", "Asia/Shanghai"}
	
	fmt.Println("=== JD to DateTime Conversion Test ===")
	for _, jd := range testJDs {
		fmt.Printf("\nJD: %.6f\n", jd)
		for _, tz := range timezones {
			dt, err := julian.JDToDateTime(jd, tz)
			if err != nil {
				fmt.Printf("  %s: ERROR %v\n", tz, err)
			} else {
				fmt.Printf("  %s: %s\n", tz, dt)
			}
		}
	}
}