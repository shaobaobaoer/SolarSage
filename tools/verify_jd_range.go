package main

import (
	"fmt"
	"time"
	"github.com/yourusername/SolarSage/pkg/julian"
)

func main() {
	// Verify JD range from compare tool
	startJD := 2461072.166667
	endJD := 2461438.166655
	
	tz := "Australia/Perth"
	
	fmt.Println("JD Range Verification:")
	fmt.Printf("Start JD: %.6f\n", startJD)
	fmt.Printf("End JD:   %.6f\n", endJD)
	fmt.Println()
	
	startDT, _ := julian.JDToDateTime(startJD, tz)
	endDT, _ := julian.JDToDateTime(endJD, tz)
	
	fmt.Printf("Start (AWST): %s\n", startDT)
	fmt.Printf("End (AWST):   %s\n", endDT)
	
	// Also check UTC
	startUTC, _ := julian.JDToDateTime(startJD, "UTC")
	endUTC, _ := julian.JDToDateTime(endJD, "UTC")
	
	fmt.Println()
	fmt.Printf("Start (UTC):  %s\n", startUTC)
	fmt.Printf("End (UTC):    %s\n", endUTC)
	
	// Verify specific event time
	// SF: 2026-02-01 01:17:24 AWST
	// What JD should this be?
	testDate := "2026-02-01T01:17:24+08:00"
	jdResult, _ := julian.DateTimeToJD(testDate, 1) // Gregorian
	fmt.Println()
	fmt.Printf("SF Event '2026-02-01 01:17:24 AWST' -> JD: %.6f (UT)\n", jdResult.JDUT)
	
	// Convert back
	backToDT, _ := julian.JDToDateTime(jdResult.JDUT, tz)
	fmt.Printf("Convert back to AWST: %s\n", backToDT)
}