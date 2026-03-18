package sweph

import (
	"math"
	"testing"
)

func TestHousePos(t *testing.T) {
	// Get house data first
	jd := JulDay(2000, 1, 1, 12.0, true)
	hr, err := Houses(jd, 39.9042, 116.4074, HousePlacidus)
	if err != nil {
		t.Fatalf("Houses: %v", err)
	}

	eps, err := Obliquity(jd)
	if err != nil {
		t.Fatalf("Obliquity: %v", err)
	}

	// Test with Sun position
	r, err := CalcUT(jd, SE_SUN)
	if err != nil {
		t.Fatalf("CalcUT Sun: %v", err)
	}

	pos, err := HousePos(hr.ARMC, 39.9042, eps, HousePlacidus, r.Longitude, r.Latitude)
	if err != nil {
		t.Fatalf("HousePos: %v", err)
	}
	if pos < 1.0 || pos > 13.0 {
		t.Errorf("HousePos = %f, want 1.0-12.999", pos)
	}
}

func TestJulDay_Julian(t *testing.T) {
	// Test Julian calendar
	jd := JulDay(1582, 10, 4, 12.0, false)
	if jd < 2299150 || jd > 2299162 {
		t.Errorf("JulDay Julian calendar = %f, out of expected range", jd)
	}
}

func TestRevJul_Julian(t *testing.T) {
	jd := JulDay(1582, 10, 4, 12.0, false)
	y, m, d, h := RevJul(jd, false)
	if y != 1582 || m != 10 || d != 4 || math.Abs(h-12.0) > 0.001 {
		t.Errorf("RevJul Julian = %d-%d-%d %.3f, want 1582-10-4 12.000", y, m, d, h)
	}
}

func TestCalcUT_AllPlanets(t *testing.T) {
	jd := JulDay(2000, 1, 1, 12.0, true)
	planets := []int{
		SE_SUN, SE_MOON, SE_MERCURY, SE_VENUS, SE_MARS,
		SE_JUPITER, SE_SATURN, SE_URANUS, SE_NEPTUNE, SE_PLUTO,
		SE_CHIRON, SE_TRUE_NODE, SE_MEAN_NODE, SE_MEAN_APOG, SE_OSCU_APOG,
	}
	for _, p := range planets {
		r, err := CalcUT(jd, p)
		if err != nil {
			t.Errorf("CalcUT planet %d error: %v", p, err)
			continue
		}
		if r.Longitude < 0 || r.Longitude >= 360 {
			t.Errorf("Planet %d lon = %f, want 0-360", p, r.Longitude)
		}
	}
}

func TestInit_Reinit(t *testing.T) {
	// Re-initialize with same path should work
	Init("../../third_party/swisseph/ephe")
	r, err := CalcUT(JulDay(2000, 1, 1, 12.0, true), SE_SUN)
	if err != nil {
		t.Fatalf("CalcUT after re-init: %v", err)
	}
	if r.Longitude < 279 || r.Longitude > 281 {
		t.Errorf("Sun after re-init: lon=%f", r.Longitude)
	}
}
