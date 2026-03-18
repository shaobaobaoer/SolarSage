package sweph

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Find ephemeris path
	ephePath := filepath.Join("..", "..", "third_party", "swisseph", "ephe")
	Init(ephePath)
	code := m.Run()
	Close()
	os.Exit(code)
}

func TestJulDay(t *testing.T) {
	// J2000.0 = 2000-01-01 12:00:00 TT = JD 2451545.0
	jd := JulDay(2000, 1, 1, 12.0, true)
	if math.Abs(jd-2451545.0) > 0.001 {
		t.Errorf("JulDay(2000,1,1,12) = %f, want ~2451545.0", jd)
	}

	// 1990-06-15 00:30 UTC
	jd2 := JulDay(1990, 6, 15, 0.5, true)
	if jd2 < 2448057 || jd2 > 2448058 {
		t.Errorf("JulDay(1990,6,15,0.5) = %f, out of expected range", jd2)
	}
}

func TestRevJul(t *testing.T) {
	jd := JulDay(2000, 1, 1, 12.0, true)
	y, m, d, h := RevJul(jd, true)
	if y != 2000 || m != 1 || d != 1 || math.Abs(h-12.0) > 0.001 {
		t.Errorf("RevJul(%f) = %d-%d-%d %.3f, want 2000-1-1 12.000", jd, y, m, d, h)
	}
}

func TestCalcUT_Sun(t *testing.T) {
	// J2000.0: Sun should be near 280° (Capricorn)
	jd := JulDay(2000, 1, 1, 12.0, true)
	r, err := CalcUT(jd, SE_SUN)
	if err != nil {
		t.Fatalf("CalcUT Sun: %v", err)
	}
	if r.Longitude < 279 || r.Longitude > 281 {
		t.Errorf("Sun at J2000: lon=%f, expected ~280", r.Longitude)
	}
	if r.IsRetrograde {
		t.Error("Sun should never be retrograde")
	}
	if r.SpeedLong < 0.9 || r.SpeedLong > 1.1 {
		t.Errorf("Sun speed=%f, expected ~1.0 deg/day", r.SpeedLong)
	}
}

func TestCalcUT_Moon(t *testing.T) {
	jd := JulDay(2000, 1, 1, 12.0, true)
	r, err := CalcUT(jd, SE_MOON)
	if err != nil {
		t.Fatalf("CalcUT Moon: %v", err)
	}
	if r.Longitude < 0 || r.Longitude >= 360 {
		t.Errorf("Moon longitude out of range: %f", r.Longitude)
	}
	// Moon moves ~13 deg/day
	if r.SpeedLong < 10 || r.SpeedLong > 16 {
		t.Errorf("Moon speed=%f, expected 10-16 deg/day", r.SpeedLong)
	}
}

func TestCalcUT_Mars_Retrograde(t *testing.T) {
	// Mars retrograde: roughly around 2024-12-06 to 2025-02-24
	// At 2025-01-15, Mars should be retrograde
	jd := JulDay(2025, 1, 15, 0.0, true)
	r, err := CalcUT(jd, SE_MARS)
	if err != nil {
		t.Fatalf("CalcUT Mars: %v", err)
	}
	if !r.IsRetrograde {
		t.Errorf("Mars should be retrograde around 2025-01-15, speed=%f", r.SpeedLong)
	}
}

func TestHouses(t *testing.T) {
	jd := JulDay(1990, 6, 15, 0.5, true)
	hr, err := Houses(jd, 39.9042, 116.4074, HousePlacidus)
	if err != nil {
		t.Fatalf("Houses: %v", err)
	}
	if hr.ASC < 0 || hr.ASC >= 360 {
		t.Errorf("ASC out of range: %f", hr.ASC)
	}
	if hr.MC < 0 || hr.MC >= 360 {
		t.Errorf("MC out of range: %f", hr.MC)
	}
	// Cusps should be in order (mostly), cusp[1] is 1st house
	for i := 1; i <= 12; i++ {
		if hr.Cusps[i] < 0 || hr.Cusps[i] >= 360 {
			t.Errorf("Cusp[%d] out of range: %f", i, hr.Cusps[i])
		}
	}
}

func TestDeltaT(t *testing.T) {
	jd := JulDay(2000, 1, 1, 12.0, true)
	dt := DeltaT(jd)
	// DeltaT around J2000 is about 63.8 seconds = ~0.000738 days
	if dt < 0.0005 || dt > 0.001 {
		t.Errorf("DeltaT at J2000 = %f days, expected ~0.000738", dt)
	}
}

func TestObliquity(t *testing.T) {
	jd := JulDay(2000, 1, 1, 12.0, true)
	eps, err := Obliquity(jd)
	if err != nil {
		t.Fatalf("Obliquity: %v", err)
	}
	// Obliquity at J2000 is ~23.44°
	if math.Abs(eps-23.44) > 0.1 {
		t.Errorf("Obliquity = %f, expected ~23.44", eps)
	}
}

func TestNormalizeDegrees(t *testing.T) {
	tests := []struct {
		in, want float64
	}{
		{0, 0},
		{360, 0},
		{-90, 270},
		{450, 90},
		{-360, 0},
		{720, 0},
		{180, 180},
		{359.999, 359.999},
	}
	for _, tt := range tests {
		got := NormalizeDegrees(tt.in)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("NormalizeDegrees(%f) = %f, want %f", tt.in, got, tt.want)
		}
	}
}
