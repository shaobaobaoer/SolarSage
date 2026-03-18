package geo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimezoneFromCoords_MoreCases(t *testing.T) {
	tests := []struct {
		lat, lon float64
		want     string
	}{
		{35.7, 139.7, "Asia/Tokyo"},
		{48.9, 2.35, "Europe/Paris"},
		{-33.9, 151.2, "Australia/Sydney"},
	}
	for _, tt := range tests {
		got := TimezoneFromCoords(tt.lat, tt.lon)
		if got != tt.want {
			t.Errorf("TimezoneFromCoords(%f, %f) = %q, want %q", tt.lat, tt.lon, got, tt.want)
		}
	}
}

func TestGeocode_WhitespaceOnly(t *testing.T) {
	_, err := Geocode("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only location name")
	}
}

func TestGeocodeNominatim_Success(t *testing.T) {
	results := []nominatimResult{
		{Lat: "39.9042", Lon: "116.4074", DisplayName: "Beijing, China"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(results)
	}))
	defer ts.Close()

	// Save and restore the original httpClient
	origClient := httpClient
	httpClient = ts.Client()
	defer func() { httpClient = origClient }()

	// Override the URL by testing geocodeNominatim directly via Geocode
	// We need to also override the URL. Instead, let's test the full path through a custom transport.
	transport := &http.Transport{}
	transport.RegisterProtocol("https", &rewriteTransport{ts.URL, ts.Client().Transport})
	httpClient = &http.Client{Transport: transport}

	loc, err := Geocode("Beijing")
	if err != nil {
		t.Fatalf("Geocode error: %v", err)
	}
	if loc.DisplayName != "Beijing, China" {
		t.Errorf("DisplayName = %q, want Beijing, China", loc.DisplayName)
	}
	if loc.Timezone == "" {
		t.Error("Timezone should not be empty")
	}
}

func TestGeocodeNominatim_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]nominatimResult{})
	}))
	defer ts.Close()

	origClient := httpClient
	transport := &http.Transport{}
	transport.RegisterProtocol("https", &rewriteTransport{ts.URL, ts.Client().Transport})
	httpClient = &http.Client{Transport: transport}
	defer func() { httpClient = origClient }()

	_, err := Geocode("NonexistentPlace12345")
	if err == nil {
		t.Error("Expected error for not found location")
	}
}

func TestGeocodeNominatim_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer ts.Close()

	origClient := httpClient
	transport := &http.Transport{}
	transport.RegisterProtocol("https", &rewriteTransport{ts.URL, ts.Client().Transport})
	httpClient = &http.Client{Transport: transport}
	defer func() { httpClient = origClient }()

	_, err := Geocode("Beijing")
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

// rewriteTransport rewrites HTTPS requests to a local test server
type rewriteTransport struct {
	targetURL string
	inner     http.RoundTripper
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.targetURL[len("http://"):]
	if t.inner == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.inner.RoundTrip(req)
}
