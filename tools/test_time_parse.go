package main

import (
	"fmt"
	"time"
)

func main() {
	// Test time parsing
	testInput := "2026-02-01T00:00:00+08:00"
	
	t, err := time.Parse(time.RFC3339, testInput)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	
	fmt.Printf("Input: %s\n", testInput)
	fmt.Printf("Parsed time: %s\n", t.Format(time.RFC3339))
	fmt.Printf("Location: %v\n", t.Location())
	fmt.Printf("UTC: %s\n", t.UTC().Format(time.RFC3339))
	
	// Extract components in UTC
	utc := t.UTC()
	fmt.Printf("\nUTC components:\n")
	fmt.Printf("  Year: %d\n", utc.Year())
	fmt.Printf("  Month: %d\n", utc.Month())
	fmt.Printf("  Day: %d\n", utc.Day())
	fmt.Printf("  Hour: %d\n", utc.Hour())
	fmt.Printf("  Minute: %d\n", utc.Minute())
	fmt.Printf("  Second: %d\n", utc.Second())
	
	// Calculate hour as float
	hour := float64(utc.Hour()) + float64(utc.Minute())/60.0 + float64(utc.Second())/3600.0
	fmt.Printf("  Hour (float): %.6f\n", hour)
}