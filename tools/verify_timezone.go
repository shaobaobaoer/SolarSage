package main

import (
	"fmt"
	"time"
)

func main() {
	// Verify timezone conversion
	// If SF shows 2026-02-01 01:17:24 AWST
	// Then UTC should be 2026-01-31 17:17:24
	
	awst, _ := time.LoadLocation("Australia/Perth")
	
	// Create time in AWST
	tAWST := time.Date(2026, 2, 1, 1, 17, 24, 0, awst)
	
	// Convert to UTC
	tUTC := tAWST.UTC()
	
	fmt.Println("Timezone conversion test:")
	fmt.Printf("  AWST time: %s\n", tAWST.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  UTC time:  %s\n", tUTC.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  Offset: %s\n", tAWST.Format("-07:00"))
	
	// Check if AWST is always UTC+8
	fmt.Println("\nAWST offset verification:")
	testDates := []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, awst),
		time.Date(2026, 6, 1, 0, 0, 0, 0, awst),
		time.Date(2026, 12, 1, 0, 0, 0, 0, awst),
	}
	
	for _, t := range testDates {
		_, offset := t.Zone()
		fmt.Printf("  %s: UTC%+d seconds (%+.1f hours)\n", 
			t.Format("2006-01-02"), offset, float64(offset)/3600)
	}
}