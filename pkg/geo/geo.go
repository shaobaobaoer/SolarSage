package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/anthropic/swisseph-mcp/pkg/models"
)

// builtinLocations provides common locations without network access
var builtinLocations = map[string]models.GeoLocation{
	"北京":        {Latitude: 39.9042, Longitude: 116.4074, Timezone: "Asia/Shanghai", DisplayName: "北京, 中国"},
	"beijing":   {Latitude: 39.9042, Longitude: 116.4074, Timezone: "Asia/Shanghai", DisplayName: "Beijing, China"},
	"上海":        {Latitude: 31.2304, Longitude: 121.4737, Timezone: "Asia/Shanghai", DisplayName: "上海, 中国"},
	"shanghai":  {Latitude: 31.2304, Longitude: 121.4737, Timezone: "Asia/Shanghai", DisplayName: "Shanghai, China"},
	"广州":        {Latitude: 23.1291, Longitude: 113.2644, Timezone: "Asia/Shanghai", DisplayName: "广州, 中国"},
	"深圳":        {Latitude: 22.5431, Longitude: 114.0579, Timezone: "Asia/Shanghai", DisplayName: "深圳, 中国"},
	"香港":        {Latitude: 22.3193, Longitude: 114.1694, Timezone: "Asia/Hong_Kong", DisplayName: "香港"},
	"台北":        {Latitude: 25.0330, Longitude: 121.5654, Timezone: "Asia/Taipei", DisplayName: "台北, 台湾"},
	"东京":        {Latitude: 35.6762, Longitude: 139.6503, Timezone: "Asia/Tokyo", DisplayName: "东京, 日本"},
	"tokyo":     {Latitude: 35.6762, Longitude: 139.6503, Timezone: "Asia/Tokyo", DisplayName: "Tokyo, Japan"},
	"london":    {Latitude: 51.5074, Longitude: -0.1278, Timezone: "Europe/London", DisplayName: "London, UK"},
	"伦敦":        {Latitude: 51.5074, Longitude: -0.1278, Timezone: "Europe/London", DisplayName: "伦敦, 英国"},
	"new york":  {Latitude: 40.7128, Longitude: -74.0060, Timezone: "America/New_York", DisplayName: "New York, USA"},
	"纽约":        {Latitude: 40.7128, Longitude: -74.0060, Timezone: "America/New_York", DisplayName: "纽约, 美国"},
	"paris":     {Latitude: 48.8566, Longitude: 2.3522, Timezone: "Europe/Paris", DisplayName: "Paris, France"},
	"巴黎":        {Latitude: 48.8566, Longitude: 2.3522, Timezone: "Europe/Paris", DisplayName: "巴黎, 法国"},
	"sydney":    {Latitude: -33.8688, Longitude: 151.2093, Timezone: "Australia/Sydney", DisplayName: "Sydney, Australia"},
	"悉尼":        {Latitude: -33.8688, Longitude: 151.2093, Timezone: "Australia/Sydney", DisplayName: "悉尼, 澳大利亚"},
	"los angeles": {Latitude: 34.0522, Longitude: -118.2437, Timezone: "America/Los_Angeles", DisplayName: "Los Angeles, USA"},
	"洛杉矶":       {Latitude: 34.0522, Longitude: -118.2437, Timezone: "America/Los_Angeles", DisplayName: "洛杉矶, 美国"},
	"成都":        {Latitude: 30.5728, Longitude: 104.0668, Timezone: "Asia/Shanghai", DisplayName: "成都, 中国"},
	"武汉":        {Latitude: 30.5928, Longitude: 114.3055, Timezone: "Asia/Shanghai", DisplayName: "武汉, 中国"},
	"杭州":        {Latitude: 30.2741, Longitude: 120.1551, Timezone: "Asia/Shanghai", DisplayName: "杭州, 中国"},
	"南京":        {Latitude: 32.0603, Longitude: 118.7969, Timezone: "Asia/Shanghai", DisplayName: "南京, 中国"},
	"重庆":        {Latitude: 29.4316, Longitude: 106.9123, Timezone: "Asia/Shanghai", DisplayName: "重庆, 中国"},
	"天津":        {Latitude: 39.3434, Longitude: 117.3616, Timezone: "Asia/Shanghai", DisplayName: "天津, 中国"},
	"西安":        {Latitude: 34.3416, Longitude: 108.9398, Timezone: "Asia/Shanghai", DisplayName: "西安, 中国"},
	"长沙":        {Latitude: 28.2282, Longitude: 112.9388, Timezone: "Asia/Shanghai", DisplayName: "长沙, 中国"},
}

// nominatimResult represents the Nominatim API response
type nominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// Geocode converts a location name to geographic coordinates
func Geocode(locationName string) (*models.GeoLocation, error) {
	// Try builtin first (case-insensitive match via lowercase)
	lower := toLower(locationName)
	if loc, ok := builtinLocations[lower]; ok {
		return &loc, nil
	}
	// Also try exact match
	if loc, ok := builtinLocations[locationName]; ok {
		return &loc, nil
	}

	// Fall back to Nominatim API
	return geocodeNominatim(locationName)
}

func geocodeNominatim(locationName string) (*models.GeoLocation, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	u := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
		url.QueryEscape(locationName))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SwissephMCP/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("geocode decode failed: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("location not found: %s", locationName)
	}

	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)

	tz := guessTimezone(lon)

	return &models.GeoLocation{
		Latitude:    lat,
		Longitude:   lon,
		Timezone:    tz,
		DisplayName: results[0].DisplayName,
	}, nil
}

// guessTimezone provides a rough timezone from longitude
func guessTimezone(lon float64) string {
	offset := int(lon / 15.0)
	if offset == 0 {
		return "UTC"
	}
	if offset > 0 {
		return fmt.Sprintf("Etc/GMT-%d", offset)
	}
	return fmt.Sprintf("Etc/GMT+%d", -offset)
}

func toLower(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b = append(b, c)
	}
	return string(b)
}
