package transit

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthropic/swisseph-mcp/pkg/models"
	"github.com/anthropic/swisseph-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath := filepath.Join("..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	code := m.Run()
	sweph.Close()
	os.Exit(code)
}

var natalJD = sweph.JulDay(1990, 6, 15, 0.5, true)
var startJD = sweph.JulDay(2024, 1, 1, 4.0, true)
var orbs = models.DefaultOrbConfig()

var defaultPlanets = []models.PlanetID{
	models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
	models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
	models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
	models.PlanetPluto,
}

func TestCalcTransitEvents_TrNa(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMercury,
		},
		EventConfig: models.EventConfig{
			IncludeTrNa:         true,
			IncludeSignIngress:  true,
			IncludeHouseIngress: true,
			IncludeStation:      true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Tr-Na: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected some Tr-Na events, got 0")
	}

	// Verify events are sorted by JD
	for i := 1; i < len(events); i++ {
		if events[i].JD < events[i-1].JD {
			t.Errorf("Events not sorted: event[%d].JD=%f < event[%d].JD=%f", i, events[i].JD, i-1, events[i-1].JD)
			break
		}
	}

	// Count event types
	counts := make(map[models.EventType]int)
	for _, e := range events {
		counts[e.EventType]++
	}
	t.Logf("Event counts: %v (total %d)", counts, len(events))

	// Should have some aspect events
	if counts[models.EventAspectEnter]+counts[models.EventAspectExact]+counts[models.EventAspectLeave] == 0 {
		t.Error("Expected some aspect events")
	}

	// Verify all aspect events have chart_type = TRANSIT and target_chart_type = NATAL
	for _, e := range events {
		if e.EventType == models.EventAspectEnter || e.EventType == models.EventAspectExact || e.EventType == models.EventAspectLeave {
			if e.ChartType != models.ChartTransit {
				t.Errorf("Tr-Na event chart_type = %s, want TRANSIT", e.ChartType)
			}
			if e.TargetChartType != models.ChartNatal {
				t.Errorf("Tr-Na event target_chart_type = %s, want NATAL", e.TargetChartType)
			}
		}
	}

	// Verify age is reasonable (~33.5 years)
	for _, e := range events {
		if e.Age < 33 || e.Age > 34 {
			t.Errorf("Age = %f, expected ~33.5", e.Age)
			break
		}
	}
}

func TestCalcTransitEvents_TrTr(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 90, // 3 months for more Tr-Tr events
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMercury, models.PlanetVenus,
		},
		EventConfig: models.EventConfig{
			IncludeTrTr: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Tr-Tr: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected some Tr-Tr events in 3 months, got 0")
	}

	// All Tr-Tr events should have chart_type = TRANSIT and target_chart_type = TRANSIT
	for _, e := range events {
		if e.ChartType != models.ChartTransit {
			t.Errorf("Tr-Tr event chart_type = %s, want TRANSIT", e.ChartType)
		}
		if e.TargetChartType != models.ChartTransit {
			t.Errorf("Tr-Tr event target_chart_type = %s, want TRANSIT", e.TargetChartType)
		}
	}
	t.Logf("Tr-Tr events: %d", len(events))
}

func TestCalcTransitEvents_SpNa(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25, // 1 year
		TransitPlanets: []models.PlanetID{},
		ProgressionsConfig: &models.ProgressionsConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon,
			},
		},
		EventConfig: models.EventConfig{
			IncludeSpNa: true,
		},
		OrbConfigProgressions: models.OrbConfig{
			Conjunction: 1, Opposition: 1, Trine: 1, Square: 1,
			Sextile: 1, Quincunx: 0.5,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Sp-Na: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected some Sp-Na events in 1 year, got 0")
	}

	for _, e := range events {
		if e.ChartType != models.ChartProgressions {
			t.Errorf("Sp-Na chart_type = %s, want PROGRESSIONS", e.ChartType)
		}
		if e.TargetChartType != models.ChartNatal {
			t.Errorf("Sp-Na target_chart_type = %s, want NATAL", e.TargetChartType)
		}
	}
	t.Logf("Sp-Na events: %d", len(events))
}

func TestCalcTransitEvents_SaNa(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25,
		TransitPlanets: []models.PlanetID{},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon, models.PlanetMars,
				models.PlanetJupiter, models.PlanetSaturn,
			},
		},
		EventConfig: models.EventConfig{
			IncludeSaNa: true,
		},
		OrbConfigSolarArc: models.OrbConfig{
			Conjunction: 1, Opposition: 1, Trine: 1, Square: 1,
			Sextile: 1, Quincunx: 0.5,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Sa-Na: %v", err)
	}

	for _, e := range events {
		if e.ChartType != models.ChartSolarArc {
			t.Errorf("Sa-Na chart_type = %s, want SOLAR_ARC", e.ChartType)
		}
		if e.TargetChartType != models.ChartNatal {
			t.Errorf("Sa-Na target_chart_type = %s, want NATAL", e.TargetChartType)
		}
	}
	t.Logf("Sa-Na events: %d", len(events))
}

func TestCalcTransitEvents_Station(t *testing.T) {
	// Mercury stations ~3-4 times per year
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25,
		TransitPlanets: []models.PlanetID{models.PlanetMercury},
		EventConfig: models.EventConfig{
			IncludeStation: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Station: %v", err)
	}

	stationCount := 0
	for _, e := range events {
		if e.EventType == models.EventStation {
			stationCount++
			if e.StationType != models.StationRetrograde && e.StationType != models.StationDirect {
				t.Errorf("Invalid station type: %s", e.StationType)
			}
		}
	}
	if stationCount < 4 {
		t.Errorf("Mercury stations in 1 year = %d, expected >= 4", stationCount)
	}
	t.Logf("Mercury station events: %d", stationCount)
}

func TestCalcTransitEvents_SignIngress(t *testing.T) {
	// Sun enters a new sign approximately every 30 days
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		EventConfig: models.EventConfig{
			IncludeSignIngress: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents SignIngress: %v", err)
	}

	ingressCount := 0
	for _, e := range events {
		if e.EventType == models.EventSignIngress {
			ingressCount++
			if e.FromSign == e.ToSign {
				t.Errorf("Sign ingress from %s to %s (same sign)", e.FromSign, e.ToSign)
			}
		}
	}
	// Sun should cross ~12 sign boundaries in a year
	if ingressCount < 11 || ingressCount > 13 {
		t.Errorf("Sun sign ingresses in 1 year = %d, expected ~12", ingressCount)
	}
	t.Logf("Sun sign ingress events: %d", ingressCount)
}

func TestCalcTransitEvents_VoidOfCourse(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 14, // 2 weeks
		TransitPlanets: []models.PlanetID{models.PlanetMoon},
		EventConfig: models.EventConfig{
			IncludeTrNa:         true,
			IncludeSignIngress:  true,
			IncludeVoidOfCourse: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents VOC: %v", err)
	}

	vocCount := 0
	for _, e := range events {
		if e.EventType == models.EventVoidOfCourse {
			vocCount++
			if e.VoidStartJD >= e.VoidEndJD {
				t.Errorf("VOC start (%f) >= end (%f)", e.VoidStartJD, e.VoidEndJD)
			}
			if e.LastAspectType == "" {
				t.Error("VOC LastAspectType is empty")
			}
			if e.NextSign == "" {
				t.Error("VOC NextSign is empty")
			}
			duration := (e.VoidEndJD - e.VoidStartJD) * 24
			if duration > 48 {
				t.Errorf("VOC duration = %.1f hours, seems too long", duration)
			}
		}
	}
	// Moon changes sign every ~2.5 days, but genuine VOC (no new aspect after last leave)
	// is less common when many natal planets are present. Expect at least 1 in 2 weeks.
	if vocCount < 1 {
		t.Errorf("VOC events in 2 weeks = %d, expected >= 1", vocCount)
	}
	t.Logf("VOC events: %d", vocCount)
}

func TestCalcTransitEvents_EmptyConfig(t *testing.T) {
	// All event types disabled → should return empty
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		EventConfig: models.EventConfig{}, // all false
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("Empty config: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events with all disabled, got %d", len(events))
	}
}

func TestAngleDiffToAspect(t *testing.T) {
	tests := []struct {
		lon1, lon2, aspectAngle float64
		wantAbs                 float64 // we check absolute value since sign depends on direction
	}{
		{100, 100, 0, 0},       // exact conjunction
		{100, 280, 180, 0},     // exact opposition
		{100, 220, 120, 0},     // exact trine
		{105, 100, 0, 5},       // 5° off conjunction
		{95, 100, 0, 5},        // 5° off conjunction (other side)
	}
	for _, tt := range tests {
		got := angleDiffToAspect(tt.lon1, tt.lon2, tt.aspectAngle)
		if math.Abs(math.Abs(got)-tt.wantAbs) > 0.1 {
			t.Errorf("angleDiffToAspect(%f, %f, %f) = %f, want |%f|", tt.lon1, tt.lon2, tt.aspectAngle, got, tt.wantAbs)
		}
	}
}

func TestWrapAngle(t *testing.T) {
	tests := []struct {
		in, want float64
	}{
		{0, 0},
		{180, -180},
		{-180, -180},
		{90, 90},
		{-90, -90},
		{270, -90},
		{-270, 90},
	}
	for _, tt := range tests {
		got := wrapAngle(tt.in)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("wrapAngle(%f) = %f, want %f", tt.in, got, tt.want)
		}
	}
}

func TestShortestAngle(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{0, 0, 0},
		{0, 180, 180},
		{10, 350, 20},
		{350, 10, 20},
		{0, 90, 90},
	}
	for _, tt := range tests {
		got := shortestAngle(tt.a, tt.b)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("shortestAngle(%f, %f) = %f, want %f", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestBuildMonoIntervals(t *testing.T) {
	stations := []StationInfo{
		{JD: 10, IsDirecting: false},
		{JD: 20, IsDirecting: true},
	}
	intervals := buildMonoIntervals(5, 25, stations)
	if len(intervals) != 3 {
		t.Fatalf("Expected 3 intervals, got %d", len(intervals))
	}
	if intervals[0].Start != 5 || intervals[0].End != 10 {
		t.Errorf("Interval 0: %v", intervals[0])
	}
	if intervals[1].Start != 10 || intervals[1].End != 20 {
		t.Errorf("Interval 1: %v", intervals[1])
	}
	if intervals[2].Start != 20 || intervals[2].End != 25 {
		t.Errorf("Interval 2: %v", intervals[2])
	}
}

func TestCalcTransitEvents_TrSp(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 90, // 3 months
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMars,
		},
		ProgressionsConfig: &models.ProgressionsConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon,
			},
		},
		EventConfig: models.EventConfig{
			IncludeTrSp: true,
		},
		OrbConfigTransit: models.OrbConfig{
			Conjunction: 2, Opposition: 2, Trine: 2, Square: 2,
			Sextile: 2, Quincunx: 1,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Tr-Sp: %v", err)
	}

	for _, e := range events {
		if e.ChartType != models.ChartTransit {
			t.Errorf("Tr-Sp chart_type = %s, want TRANSIT", e.ChartType)
		}
		if e.TargetChartType != models.ChartProgressions {
			t.Errorf("Tr-Sp target_chart_type = %s, want PROGRESSIONS", e.TargetChartType)
		}
	}
	t.Logf("Tr-Sp events: %d", len(events))
}

func TestCalcTransitEvents_TrSa(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 90,
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMars,
		},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon, models.PlanetMars,
			},
		},
		EventConfig: models.EventConfig{
			IncludeTrSa: true,
		},
		OrbConfigTransit: models.OrbConfig{
			Conjunction: 2, Opposition: 2, Trine: 2, Square: 2,
			Sextile: 2, Quincunx: 1,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Tr-Sa: %v", err)
	}

	for _, e := range events {
		if e.ChartType != models.ChartTransit {
			t.Errorf("Tr-Sa chart_type = %s, want TRANSIT", e.ChartType)
		}
		if e.TargetChartType != models.ChartSolarArc {
			t.Errorf("Tr-Sa target_chart_type = %s, want SOLAR_ARC", e.TargetChartType)
		}
	}
	t.Logf("Tr-Sa events: %d", len(events))
}

func TestCalcTransitEvents_SpSp(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25 * 3, // 3 years for slow-moving progressed planets
		TransitPlanets: []models.PlanetID{},
		ProgressionsConfig: &models.ProgressionsConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			},
		},
		EventConfig: models.EventConfig{
			IncludeSpSp: true,
		},
		OrbConfigProgressions: models.OrbConfig{
			Conjunction: 1, Opposition: 1, Trine: 1, Square: 1,
			Sextile: 1,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("CalcTransitEvents Sp-Sp: %v", err)
	}

	for _, e := range events {
		if e.ChartType != models.ChartProgressions {
			t.Errorf("Sp-Sp chart_type = %s, want PROGRESSIONS", e.ChartType)
		}
		if e.TargetChartType != models.ChartProgressions {
			t.Errorf("Sp-Sp target_chart_type = %s, want PROGRESSIONS", e.TargetChartType)
		}
	}
	t.Logf("Sp-Sp events: %d", len(events))
}

func TestCalcTransitEvents_WithChiron(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{models.PlanetChiron},
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{models.PlanetChiron},
		EventConfig: models.EventConfig{
			IncludeTrNa:    true,
			IncludeStation: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("Chiron transit: %v", err)
	}
	t.Logf("Chiron events: %d", len(events))
}

func TestCalcTransitEvents_WithNodes(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{models.PlanetNorthNodeTrue, models.PlanetSouthNode},
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 90,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		EventConfig: models.EventConfig{
			IncludeTrNa: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("Node transit: %v", err)
	}
	// Sun should aspect the natal nodes in 3 months
	if len(events) == 0 {
		t.Error("Expected Sun to aspect natal nodes in 3 months")
	}
	t.Logf("Node transit events: %d", len(events))
}

func TestCalcTransitEvents_HouseIngress(t *testing.T) {
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 365.25,
		TransitPlanets: []models.PlanetID{models.PlanetSun},
		EventConfig: models.EventConfig{
			IncludeHouseIngress: true,
		},
		OrbConfigTransit: orbs,
		HouseSystem:      models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("House ingress: %v", err)
	}

	houseCount := 0
	for _, e := range events {
		if e.EventType == models.EventHouseIngress {
			houseCount++
			if e.FromHouse == e.ToHouse {
				t.Errorf("House ingress from %d to %d (same house)", e.FromHouse, e.ToHouse)
			}
			if e.ToHouse < 1 || e.ToHouse > 12 {
				t.Errorf("Invalid house: %d", e.ToHouse)
			}
		}
	}
	// Sun should cross ~12 house cusps in a year
	if houseCount < 10 || houseCount > 14 {
		t.Errorf("Sun house ingresses in 1 year = %d, expected ~12", houseCount)
	}
	t.Logf("Sun house ingress events: %d", houseCount)
}

func TestCalcTransitEvents_FullPipeline(t *testing.T) {
	// Test all event types simultaneously
	events, err := CalcTransitEvents(TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      natalJD,
		NatalPlanets: defaultPlanets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        startJD + 30,
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars,
		},
		ProgressionsConfig: &models.ProgressionsConfig{
			Enabled: true,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMoon},
		},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{models.PlanetSun, models.PlanetMars},
		},
		EventConfig:           models.DefaultEventConfig(),
		OrbConfigTransit:      orbs,
		OrbConfigProgressions: models.OrbConfig{Conjunction: 1, Opposition: 1, Trine: 1, Square: 1, Sextile: 1},
		OrbConfigSolarArc:     models.OrbConfig{Conjunction: 1, Opposition: 1, Trine: 1, Square: 1, Sextile: 1},
		HouseSystem:           models.HousePlacidus,
	})
	if err != nil {
		t.Fatalf("Full pipeline: %v", err)
	}

	counts := make(map[models.EventType]int)
	chartTypes := make(map[models.ChartType]int)
	for _, e := range events {
		counts[e.EventType]++
		if e.EventType == models.EventAspectEnter || e.EventType == models.EventAspectExact || e.EventType == models.EventAspectLeave {
			chartTypes[e.ChartType]++
		}
	}

	t.Logf("Full pipeline events: %d", len(events))
	t.Logf("Event types: %v", counts)
	t.Logf("Aspect chart types: %v", chartTypes)

	// Should have events from multiple chart types
	if chartTypes[models.ChartTransit] == 0 {
		t.Error("Expected TRANSIT aspect events")
	}

	// Verify sorted
	for i := 1; i < len(events); i++ {
		if events[i].JD < events[i-1].JD {
			t.Errorf("Events not sorted at index %d", i)
			break
		}
	}
}

func TestIntersectIntervals(t *testing.T) {
	a := []MonoInterval{{Start: 0, End: 10}, {Start: 15, End: 25}}
	b := []MonoInterval{{Start: 5, End: 20}}
	result := intersectIntervals(a, b)
	if len(result) != 2 {
		t.Fatalf("Expected 2 intervals, got %d: %v", len(result), result)
	}
	if result[0].Start != 5 || result[0].End != 10 {
		t.Errorf("Interval 0: %v, expected [5,10]", result[0])
	}
	if result[1].Start != 15 || result[1].End != 20 {
		t.Errorf("Interval 1: %v, expected [15,20]", result[1])
	}
}

func TestIntersectIntervals_NoOverlap(t *testing.T) {
	a := []MonoInterval{{Start: 0, End: 5}}
	b := []MonoInterval{{Start: 10, End: 15}}
	result := intersectIntervals(a, b)
	if len(result) != 0 {
		t.Errorf("Expected no intervals, got %d", len(result))
	}
}

func TestAdaptiveStep(t *testing.T) {
	// Fast movers
	step := adaptiveStep(1.0, 0.5)
	if step != fineStep {
		t.Errorf("Fast movers step = %f, want %f", step, fineStep)
	}
	// Slow movers (progressions)
	step = adaptiveStep(0.05, 0.0)
	if step != 2.0 {
		t.Errorf("Slow movers step = %f, want 2.0", step)
	}
	// Very slow movers (solar arc)
	step = adaptiveStep(0.003, 0.0)
	if step != 7.0 {
		t.Errorf("Very slow movers step = %f, want 7.0", step)
	}
}

func TestBuildMonoIntervals_NoStations(t *testing.T) {
	intervals := buildMonoIntervals(5, 25, nil)
	if len(intervals) != 1 {
		t.Fatalf("Expected 1 interval, got %d", len(intervals))
	}
	if intervals[0].Start != 5 || intervals[0].End != 25 {
		t.Errorf("Single interval: %v", intervals[0])
	}
}
