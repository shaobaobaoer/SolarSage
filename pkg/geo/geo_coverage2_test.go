package geo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimezoneFromCoords_OceanFallback(t *testing.T) {
	// Point in the middle of the Pacific Ocean - may not have a timezone
	// The fallback logic uses longitude-based estimation
	tz := TimezoneFromCoords(0, -170)
	if tz == "" {
		t.Error("should return a timezone even for ocean coordinates")
	}
}

func TestTimezoneFromCoords_UTCZone(t *testing.T) {
	// Greenwich area should be UTC or Europe/London
	tz := TimezoneFromCoords(51.48, 0.0)
	if tz == "" {
		t.Error("should return a timezone for Greenwich")
	}
}

func TestGeocodeNominatim_InvalidLat(t *testing.T) {
	results := []nominatimResult{
		{Lat: "not-a-number", Lon: "116.4074", DisplayName: "Bad"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(results)
	}))
	defer ts.Close()

	origClient := httpClient
	transport := &http.Transport{}
	transport.RegisterProtocol("https", &rewriteTransport{ts.URL, ts.Client().Transport})
	httpClient = &http.Client{Transport: transport}
	defer func() { httpClient = origClient }()

	_, err := Geocode("Test")
	if err == nil {
		t.Error("Expected error for invalid latitude")
	}
}

func TestGeocodeNominatim_InvalidLon(t *testing.T) {
	results := []nominatimResult{
		{Lat: "39.9042", Lon: "not-a-number", DisplayName: "Bad"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(results)
	}))
	defer ts.Close()

	origClient := httpClient
	transport := &http.Transport{}
	transport.RegisterProtocol("https", &rewriteTransport{ts.URL, ts.Client().Transport})
	httpClient = &http.Client{Transport: transport}
	defer func() { httpClient = origClient }()

	_, err := Geocode("Test")
	if err == nil {
		t.Error("Expected error for invalid longitude")
	}
}
