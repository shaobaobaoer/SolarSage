package transit

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// TestGetStations_CacheHit tests that GetStations returns cached stations without calling calcFn
func TestGetStations_CacheHit(t *testing.T) {
	ctx := &CalcContext{
		StartJD:      2461072.166667,
		EndJD:        2461438.166655,
		StationCache: make(map[string][]StationInfo),
	}

	// Pre-populate cache
	key := "SATURN"
	prepopulated := []StationInfo{
		{JD: 2461100.0, IsDirecting: true},
		{JD: 2461200.0, IsDirecting: false},
	}
	ctx.StationCache[key] = prepopulated

	// Call with calcFn that should NOT be called
	calcFn := func(jd float64) (float64, float64, error) {
		t.Error("calcFn should not be called on cache hit")
		return 0, 0, nil
	}

	stations := ctx.GetStations(key, calcFn, models.PlanetSaturn)

	// Should return prepopulated stations
	if len(stations) != 2 {
		t.Errorf("Expected 2 cached stations, got %d", len(stations))
	}

	if stations[0].JD != 2461100.0 || stations[1].JD != 2461200.0 {
		t.Error("Returned stations don't match cached values")
	}
}

// TestGetStations_CacheMissAndPopulate tests cache miss triggers computation
func TestGetStations_CacheMissAndPopulate(t *testing.T) {
	ctx := &CalcContext{
		StartJD:      2461072.166667,
		EndJD:        2461438.166655,
		StationCache: make(map[string][]StationInfo),
	}

	key := "JUPITER"

	// Mock calcFn that simulates stations via speed sign changes
	callCount := 0
	calcFn := func(jd float64) (float64, float64, error) {
		callCount++
		days := jd - 2461072.166667
		// Create speed sign changes to simulate retrograde cycle
		var speed float64
		if days < 180 {
			speed = 0.08 // Direct
		} else if days < 220 {
			speed = -0.05 // Retrograde
		} else {
			speed = 0.08 // Direct again
		}
		return 319.5 + days*0.08, speed, nil
	}

	stations := ctx.GetStations(key, calcFn, models.PlanetJupiter)

	// Should have stations due to speed sign changes
	if len(stations) < 2 {
		t.Errorf("Expected at least 2 stations from speed changes, got %d", len(stations))
	}

	// Verify cache was populated
	if cached, ok := ctx.StationCache[key]; !ok {
		t.Error("Expected stations to be cached")
	} else if len(cached) != len(stations) {
		t.Errorf("Cache mismatch: expected %d stations, got %d", len(stations), len(cached))
	}

	// Second call should use cache (calcFn not called again)
	callCountBefore := callCount
	stations2 := ctx.GetStations(key, calcFn, models.PlanetJupiter)

	if callCount != callCountBefore {
		t.Error("Second call should use cache, calcFn should not be called again")
	}

	if len(stations2) != len(stations) {
		t.Errorf("Cached stations mismatch: expected %d, got %d", len(stations), len(stations2))
	}
}

// TestGetStations_MultiplePlanets tests cache isolation between different planets
func TestGetStations_MultiplePlanets(t *testing.T) {
	ctx := &CalcContext{
		StartJD:      2461072.166667,
		EndJD:        2461438.166655,
		StationCache: make(map[string][]StationInfo),
	}

	// Get stations for Jupiter
	jupiterStations := ctx.GetStations("JUPITER", func(jd float64) (float64, float64, error) {
		return 319.5, 0.08, nil // No stations (constant speed)
	}, models.PlanetJupiter)

	// Get stations for Saturn
	saturnStations := ctx.GetStations("SATURN", func(jd float64) (float64, float64, error) {
		return 13.5, 0.03, nil // No stations (constant speed)
	}, models.PlanetSaturn)

	// Both should be cached (even if empty)
	if len(ctx.StationCache) != 2 {
		t.Errorf("Expected 2 cached entries, got %d", len(ctx.StationCache))
	}

	// Both should have no stations (constant speed = no sign changes)
	if len(jupiterStations) != 0 {
		t.Errorf("Jupiter should have 0 stations with constant speed, got %d", len(jupiterStations))
	}
	if len(saturnStations) != 0 {
		t.Errorf("Saturn should have 0 stations with constant speed, got %d", len(saturnStations))
	}
}

// TestBuildCalcContext tests the context builder with valid input
func TestBuildCalcContext(t *testing.T) {
	input := TransitCalcInput{
		NatalChart: NatalChartConfig{
			Lat:     39.9042,
			Lon:     116.4074,
			JD:      2451545.0,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		},
		TimeRange: TimeRangeConfig{
			StartJD: 2461072.166667,
			EndJD:   2461438.166655,
		},
		HouseSystem: models.HousePlacidus,
	}

	ctx, err := buildCalcContext(input)
	if err != nil {
		t.Fatalf("buildCalcContext failed: %v", err)
	}

	// Verify context fields
	if ctx.NatalJD != 2451545.0 {
		t.Errorf("Expected NatalJD 2451545.0, got %f", ctx.NatalJD)
	}

	if ctx.StartJD != 2461072.166667 {
		t.Errorf("Expected StartJD 2461072.166667, got %f", ctx.StartJD)
	}

	if ctx.EndJD != 2461438.166655 {
		t.Errorf("Expected EndJD 2461438.166655, got %f", ctx.EndJD)
	}

	// Should have 12 houses
	if len(ctx.NatalHouses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(ctx.NatalHouses))
	}

	// Should have 2 natal refs (Sun + Moon)
	if len(ctx.NatalRefs) != 2 {
		t.Errorf("Expected 2 natal refs, got %d", len(ctx.NatalRefs))
	}

	// Station cache should be initialized
	if ctx.StationCache == nil {
		t.Error("StationCache should be initialized")
	}
}

// TestBuildNatalRefs tests natal reference point collection
func TestBuildNatalRefs(t *testing.T) {
	input := TransitCalcInput{
		NatalChart: NatalChartConfig{
			Lat:     39.9042,
			Lon:     116.4074,
			JD:      2451545.0,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon, models.PlanetMercury},
			Points:  []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		HouseSystem: models.HousePlacidus,
	}

	// Calculate houses first (required for special points)
	ctx, err := buildCalcContext(input)
	if err != nil {
		t.Fatalf("Failed to build context: %v", err)
	}

	refs := buildNatalRefs(input, ctx.NatalHouses)

	// Should have 5 refs: 3 planets + 2 special points
	if len(refs) != 5 {
		t.Errorf("Expected 5 natal refs, got %d", len(refs))
	}

	// Verify ref structure
	for i, ref := range refs {
		if ref.ID == "" {
			t.Errorf("Ref %d has empty ID", i)
		}
		if ref.Longitude < 0 || ref.Longitude >= 360 {
			t.Errorf("Ref %d longitude out of range: %f", i, ref.Longitude)
		}
		if ref.ChartType != models.ChartNatal {
			t.Errorf("Ref %d should be ChartNatal, got %v", i, ref.ChartType)
		}
	}
}
