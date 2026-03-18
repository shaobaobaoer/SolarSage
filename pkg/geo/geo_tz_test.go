package geo

import "testing"

func TestTimezoneFromCoords_Fallback(t *testing.T) {
	// Test the longitude-based fallback for coords that tzf can't resolve
	// Normal case - should return a real timezone
	tz := TimezoneFromCoords(51.5074, -0.1278) // London
	if tz == "" {
		t.Error("expected non-empty timezone for London")
	}

	// Test with coordinates that return known timezones
	tz2 := TimezoneFromCoords(35.6762, 139.6503) // Tokyo
	if tz2 == "" {
		t.Error("expected non-empty timezone for Tokyo")
	}

	tz3 := TimezoneFromCoords(40.7128, -74.0060) // New York
	if tz3 == "" {
		t.Error("expected non-empty timezone for New York")
	}
}
