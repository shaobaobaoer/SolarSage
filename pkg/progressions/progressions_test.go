package progressions

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
var transitJD = sweph.JulDay(2024, 1, 1, 4.0, true) // age ~33.5

func TestSecondaryProgressionJD(t *testing.T) {
	pJD := SecondaryProgressionJD(natalJD, transitJD)

	// age ~33.5 years → ~33.5 days after natal
	expectedDays := (transitJD - natalJD) / JulianYear
	expectedPJD := natalJD + expectedDays
	if math.Abs(pJD-expectedPJD) > 0.001 {
		t.Errorf("SecondaryProgressionJD = %f, expected %f", pJD, expectedPJD)
	}

	// Progressed JD should be natal + ~33.5 days
	diff := pJD - natalJD
	if diff < 33 || diff > 34 {
		t.Errorf("Progressed JD is %f days after natal, expected ~33.5", diff)
	}
}

func TestAge(t *testing.T) {
	age := Age(natalJD, transitJD)
	if age < 33 || age > 34 {
		t.Errorf("Age = %f, expected ~33.5", age)
	}
}

func TestCalcProgressedLongitude_Sun(t *testing.T) {
	lon, speed, err := CalcProgressedLongitude(models.PlanetSun, natalJD, transitJD)
	if err != nil {
		t.Fatalf("CalcProgressedLongitude: %v", err)
	}

	// Natal Sun is at ~83.67° (Gemini). After ~33.5 years (=33.5 days progressed),
	// Sun moves ~1°/day, so progressed Sun ≈ 83.67 + 33.5 ≈ ~117° (Cancer)
	if lon < 100 || lon > 130 {
		t.Errorf("Progressed Sun lon = %f, expected ~117 (Cancer)", lon)
	}

	// Speed should be ~1°/day actual / 365.25 ≈ 0.00274°/real day
	if speed < 0.002 || speed > 0.003 {
		t.Errorf("Progressed Sun speed = %f, expected ~0.00274", speed)
	}
}

func TestCalcProgressedLongitude_Moon(t *testing.T) {
	lon, speed, err := CalcProgressedLongitude(models.PlanetMoon, natalJD, transitJD)
	if err != nil {
		t.Fatalf("CalcProgressedLongitude Moon: %v", err)
	}

	// Moon moves ~13°/day, so after 33.5 days progressed ≈ 33.5 * 13 = ~435.5°
	// Natal Moon at ~339°, so progressed ≈ 339 + 435.5 = ~774.5 → normalized ~54.5° (Taurus/Gemini)
	if lon < 0 || lon >= 360 {
		t.Errorf("Progressed Moon lon = %f, out of range", lon)
	}

	// Progressed Moon speed ≈ 12-15/365.25 ≈ 0.033-0.041°/day
	if speed < 0.03 || speed > 0.045 {
		t.Errorf("Progressed Moon speed = %f, expected 0.03-0.045", speed)
	}
}

func TestSolarArcOffset(t *testing.T) {
	offset, err := SolarArcOffset(natalJD, transitJD)
	if err != nil {
		t.Fatalf("SolarArcOffset: %v", err)
	}

	// After ~33.5 years, solar arc ≈ 33.5° (sun moves ~1°/year)
	if offset < 31 || offset > 36 {
		t.Errorf("Solar arc offset = %f, expected ~33.5", offset)
	}
}

func TestCalcSolarArcLongitude(t *testing.T) {
	// Mars natal = ~10.70°, solar arc ≈ 33.5°, directed ≈ 44.2°
	lon, speed, err := CalcSolarArcLongitude(models.PlanetMars, natalJD, transitJD)
	if err != nil {
		t.Fatalf("CalcSolarArcLongitude Mars: %v", err)
	}

	if lon < 40 || lon > 50 {
		t.Errorf("Solar arc Mars lon = %f, expected ~44 (Taurus)", lon)
	}

	// Speed ≈ 0.00274°/day
	if speed < 0.002 || speed > 0.003 {
		t.Errorf("Solar arc speed = %f, expected ~0.00274", speed)
	}
}

func TestCalcSolarArcLongitude_Sun(t *testing.T) {
	// Solar arc Sun = natal Sun + offset = same as progressed Sun position (approximately)
	saLon, _, err := CalcSolarArcLongitude(models.PlanetSun, natalJD, transitJD)
	if err != nil {
		t.Fatalf("CalcSolarArcLongitude Sun: %v", err)
	}

	pLon, _, err := CalcProgressedLongitude(models.PlanetSun, natalJD, transitJD)
	if err != nil {
		t.Fatalf("CalcProgressedLongitude Sun: %v", err)
	}

	// Solar arc Sun and progressed Sun should be very close
	diff := math.Abs(saLon - pLon)
	if diff > 180 {
		diff = 360 - diff
	}
	if diff > 0.5 {
		t.Errorf("Solar arc Sun (%f) vs progressed Sun (%f): diff=%f, expected < 0.5", saLon, pLon, diff)
	}
}
