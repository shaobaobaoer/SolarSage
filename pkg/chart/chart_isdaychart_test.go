package chart

import "testing"

func TestIsDayChart_Exported(t *testing.T) {
	// J2000.0 noon - Sun at ~280° Capricorn, ASC depends on location
	// Just verify it doesn't panic and returns a bool
	result := IsDayChart(2451545.0, 100.0)
	_ = result // we just want to cover the exported wrapper

	// Test with different ASC
	result2 := IsDayChart(2451545.0, 280.0)
	_ = result2
}
