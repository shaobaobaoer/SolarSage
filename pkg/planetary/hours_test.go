package planetary

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

const j2000 = 2451545.0 // 2000-01-01 12:00 UT

func TestCalcPlanetaryHours(t *testing.T) {
	// London, J2000.0 (a Saturday, ruled by Saturn)
	day, err := CalcPlanetaryHours(j2000, 51.5074, -0.1278)
	if err != nil {
		t.Fatalf("CalcPlanetaryHours: %v", err)
	}

	if len(day.Hours) != 24 {
		t.Errorf("Expected 24 hours, got %d", len(day.Hours))
	}

	// First 12 hours should be day hours
	for i := 0; i < 12; i++ {
		if !day.Hours[i].IsDay {
			t.Errorf("Hour %d should be a day hour", i+1)
		}
	}
	// Last 12 hours should be night hours
	for i := 12; i < 24; i++ {
		if day.Hours[i].IsDay {
			t.Errorf("Hour %d should be a night hour", i+1)
		}
	}

	// Sunrise should be before sunset
	if day.Sunrise >= day.Sunset {
		t.Errorf("Sunrise (%.6f) >= Sunset (%.6f)", day.Sunrise, day.Sunset)
	}

	// Next sunrise should be after sunset
	if day.NextSunrise <= day.Sunset {
		t.Errorf("NextSunrise (%.6f) <= Sunset (%.6f)", day.NextSunrise, day.Sunset)
	}

	// First hour of the day should be ruled by the day ruler
	if day.Hours[0].Planet != day.DayRuler {
		t.Errorf("First hour planet = %s, day ruler = %s", day.Hours[0].Planet, day.DayRuler)
	}

	// Hours should be contiguous
	for i := 1; i < 24; i++ {
		if day.Hours[i].StartJD < day.Hours[i-1].EndJD-0.0001 {
			t.Errorf("Hour %d starts before hour %d ends", i+1, i)
		}
	}
}

func TestCurrentPlanetaryHour(t *testing.T) {
	// Noon at London should fall in a daytime hour
	h, err := CurrentPlanetaryHour(j2000, 51.5074, -0.1278)
	if err != nil {
		t.Fatalf("CurrentPlanetaryHour: %v", err)
	}
	if !h.IsDay {
		t.Error("Noon should be a daytime hour")
	}
	if h.Planet == "" {
		t.Error("Hour planet should not be empty")
	}
}

func TestCalcPlanetaryHours_DifferentLocations(t *testing.T) {
	// Test with equatorial location (Singapore)
	day, err := CalcPlanetaryHours(j2000, 1.3521, 103.8198)
	if err != nil {
		t.Fatalf("Singapore: %v", err)
	}
	if len(day.Hours) != 24 {
		t.Errorf("Expected 24 hours, got %d", len(day.Hours))
	}

	// Day and night durations should be roughly equal near equator
	dayDur := day.Sunset - day.Sunrise
	nightDur := day.NextSunrise - day.Sunset
	ratio := dayDur / nightDur
	if ratio < 0.8 || ratio > 1.2 {
		t.Errorf("Equatorial day/night ratio = %.2f, expected ~1.0", ratio)
	}
}

func TestJulianWeekday(t *testing.T) {
	// 2000-01-01 is a Saturday (weekday 6)
	wd := julianWeekday(2000, 1, 1)
	if wd != 6 {
		t.Errorf("2000-01-01 weekday = %d, want 6 (Saturday)", wd)
	}
}
