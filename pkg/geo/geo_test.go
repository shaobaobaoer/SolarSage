package geo

import (
	"testing"
)

func TestTimezoneFromCoords(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
		want string
	}{
		{"Beijing", 39.9042, 116.4074, "Asia/Shanghai"},
		{"London", 51.5074, -0.1278, "Europe/London"},
		{"New York", 40.7128, -74.0060, "America/New_York"},
		{"Tokyo", 35.6762, 139.6503, "Asia/Tokyo"},
		{"Sydney", -33.8688, 151.2093, "Australia/Sydney"},
		{"Singapore", 1.3521, 103.8198, "Asia/Singapore"},
		{"Dubai", 25.2048, 55.2708, "Asia/Dubai"},
		{"Berlin", 52.5200, 13.4050, "Europe/Berlin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimezoneFromCoords(tt.lat, tt.lon)
			if got != tt.want {
				t.Errorf("TimezoneFromCoords(%f, %f) = %q, want %q", tt.lat, tt.lon, got, tt.want)
			}
		})
	}
}

func TestGeocode_EmptyName(t *testing.T) {
	_, err := Geocode("")
	if err == nil {
		t.Error("Geocode(\"\") should return error")
	}
	_, err = Geocode("   ")
	if err == nil {
		t.Error("Geocode(\"   \") should return error")
	}
}
