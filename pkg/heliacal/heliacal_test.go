package heliacal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

// J2000.0 = 2000-01-01 12:00 UT
const j2000 = 2451545.0

func TestVenusHeliacalRising(t *testing.T) {
	// Search for Venus heliacal rising starting from J2000
	evt, err := NextHeliacalRising(models.PlanetVenus, 48.85, 2.35, j2000)
	if err != nil {
		t.Fatalf("NextHeliacalRising(Venus): %v", err)
	}
	if evt.JDStart <= j2000 {
		t.Errorf("JDStart (%.2f) should be after startJD (%.2f)", evt.JDStart, j2000)
	}
	if evt.JDOptimum < evt.JDStart {
		t.Errorf("JDOptimum (%.2f) should be >= JDStart (%.2f)", evt.JDOptimum, evt.JDStart)
	}
	if evt.Planet != models.PlanetVenus {
		t.Errorf("expected Venus, got %s", evt.Planet)
	}
	if evt.EventType != HeliacalRising {
		t.Errorf("expected HELIACAL_RISING, got %s", evt.EventType)
	}
	t.Logf("Venus heliacal rising: JD %.4f (optimum %.4f, end %.4f)", evt.JDStart, evt.JDOptimum, evt.JDEnd)
}

func TestCalcHeliacalEvents_WithinRange(t *testing.T) {
	startJD := j2000
	endJD := j2000 + 365.25 // one year

	result, err := CalcHeliacalEvents(48.85, 2.35, 0, startJD, endJD, nil)
	if err != nil {
		t.Fatalf("CalcHeliacalEvents: %v", err)
	}

	if len(result.Events) == 0 {
		t.Fatal("expected at least one heliacal event in a year")
	}

	for _, evt := range result.Events {
		if evt.JDStart < startJD || evt.JDStart > endJD {
			t.Errorf("event JD %.2f outside range [%.2f, %.2f]", evt.JDStart, startJD, endJD)
		}
		if evt.Planet == "" {
			t.Error("event has empty planet")
		}
		if evt.EventType == "" {
			t.Error("event has empty event type")
		}
	}

	t.Logf("Found %d heliacal events in one year", len(result.Events))
}

func TestNextHeliacalRising_AfterStartJD(t *testing.T) {
	planets := []models.PlanetID{
		models.PlanetMercury,
		models.PlanetMars,
		models.PlanetJupiter,
		models.PlanetSaturn,
	}

	for _, p := range planets {
		evt, err := NextHeliacalRising(p, 40.71, -74.01, j2000)
		if err != nil {
			t.Errorf("NextHeliacalRising(%s): %v", p, err)
			continue
		}
		if evt.JDStart <= j2000 {
			t.Errorf("%s: JDStart (%.2f) should be after startJD (%.2f)", p, evt.JDStart, j2000)
		}
		t.Logf("%s heliacal rising: JD %.4f", p, evt.JDStart)
	}
}

func TestInvalidPlanet_Sun(t *testing.T) {
	_, err := NextHeliacalRising(models.PlanetSun, 48.85, 2.35, j2000)
	if err == nil {
		t.Error("expected error for Sun, got nil")
	}
}

func TestInvalidPlanet_Moon(t *testing.T) {
	_, err := NextHeliacalRising(models.PlanetMoon, 48.85, 2.35, j2000)
	if err == nil {
		t.Error("expected error for Moon, got nil")
	}
}

func TestCalcHeliacalEvents_InvalidRange(t *testing.T) {
	_, err := CalcHeliacalEvents(48.85, 2.35, 0, j2000+100, j2000, nil)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestCalcHeliacalEvents_SpecificPlanets(t *testing.T) {
	startJD := j2000
	endJD := j2000 + 365.25

	result, err := CalcHeliacalEvents(48.85, 2.35, 0, startJD, endJD, []models.PlanetID{models.PlanetVenus})
	if err != nil {
		t.Fatalf("CalcHeliacalEvents(Venus only): %v", err)
	}

	for _, evt := range result.Events {
		if evt.Planet != models.PlanetVenus {
			t.Errorf("expected only Venus events, got %s", evt.Planet)
		}
	}

	if len(result.Events) == 0 {
		t.Error("expected at least one Venus event in a year")
	}
	t.Logf("Found %d Venus heliacal events in one year", len(result.Events))
}

func TestEventTypeName(t *testing.T) {
	tests := []struct {
		et   EventType
		want string
	}{
		{HeliacalRising, "Heliacal Rising"},
		{HeliacalSetting, "Heliacal Setting"},
		{EveningFirst, "Evening First"},
		{MorningLast, "Morning Last"},
	}
	for _, tt := range tests {
		got := EventTypeName(tt.et)
		if got != tt.want {
			t.Errorf("EventTypeName(%s) = %q, want %q", tt.et, got, tt.want)
		}
	}
}
