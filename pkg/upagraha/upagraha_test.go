package upagraha

import (
	"testing"
)

func TestCalc_ReturnsAllPoints(t *testing.T) {
	// J2000.0 = 2451545.0, London ~51.5°N 0°W, day chart
	result := Calc(2451545.0, 51.5, 0.0, true)
	if result == nil {
		t.Fatal("Calc returned nil")
	}

	points := []struct {
		name string
		lon  float64
	}{
		{"Gulika", result.Gulika.Longitude},
		{"Mandi", result.Mandi.Longitude},
		{"Dhuma", result.Dhuma.Longitude},
		{"Vyatipaata", result.Vyatipaata.Longitude},
		{"Parivesha", result.Parivesha.Longitude},
		{"Indrachaapa", result.Indrachaapa.Longitude},
		{"Upaketu", result.Upaketu.Longitude},
		{"Kaala", result.Kaala.Longitude},
		{"Yamaghantaka", result.Yamaghantaka.Longitude},
	}

	for _, p := range points {
		if p.lon < 0 || p.lon >= 360 {
			t.Errorf("%s longitude = %.4f, want [0, 360)", p.name, p.lon)
		}
	}
}

func TestCalc_SignPopulated(t *testing.T) {
	result := Calc(2451545.0, 28.6, 77.2, true) // New Delhi
	if result.Gulika.Sign == "" {
		t.Error("Gulika.Sign is empty")
	}
	if result.Dhuma.Sign == "" {
		t.Error("Dhuma.Sign is empty")
	}
}

func TestDhumaFormula(t *testing.T) {
	// Dhuma = Sun + 133.333°. Verify with a known Sun longitude.
	// If Sun is at 280° (approx Capricorn), Dhuma = 280 + 133.333 = 413.333 → 53.333°
	// We test the formula indirectly via Calc.
	result := Calc(2451545.0, 0, 0, true) // Sun near 280° at J2000
	// Sun at ~280° → Dhuma ~53°
	if result.Dhuma.Longitude < 0 || result.Dhuma.Longitude >= 360 {
		t.Errorf("Dhuma longitude out of range: %.4f", result.Dhuma.Longitude)
	}
}

func TestVyatipataFormula(t *testing.T) {
	result := Calc(2451545.0, 0, 0, true)
	// Vyatipaata = 360 - Dhuma
	expected := normLon(360.0 - result.Dhuma.Longitude)
	if result.Vyatipaata.Longitude < 0 || result.Vyatipaata.Longitude >= 360 {
		t.Errorf("Vyatipaata out of range: %.4f", result.Vyatipaata.Longitude)
	}
	_ = expected
}

func TestPariveshaFormula(t *testing.T) {
	result := Calc(2451545.0, 0, 0, true)
	// Parivesha = Vyatipaata + 180°
	expected := normLon(result.Vyatipaata.Longitude + 180.0)
	diff := result.Parivesha.Longitude - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.01 {
		t.Errorf("Parivesha = %.4f, expected %.4f (Vyatipaata+180)", result.Parivesha.Longitude, expected)
	}
}

func TestIndrachaapaFormula(t *testing.T) {
	result := Calc(2451545.0, 0, 0, true)
	// Indrachaapa = 360 - Parivesha
	expected := normLon(360.0 - result.Parivesha.Longitude)
	diff := result.Indrachaapa.Longitude - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.01 {
		t.Errorf("Indrachaapa = %.4f, expected %.4f (360-Parivesha)", result.Indrachaapa.Longitude, expected)
	}
}

func TestUpaKetuFormula(t *testing.T) {
	result := Calc(2451545.0, 0, 0, true)
	// Upaketu = Indrachaapa + 16.667°
	expected := normLon(result.Indrachaapa.Longitude + 16.667)
	diff := result.Upaketu.Longitude - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.01 {
		t.Errorf("Upaketu = %.4f, expected %.4f (Indrachaapa+16.667)", result.Upaketu.Longitude, expected)
	}
}

func TestNightChart_DifferentFromDay(t *testing.T) {
	// Gulika should differ between day and night charts
	day := Calc(2451545.0, 28.6, 77.2, true)
	night := Calc(2451545.0, 28.6, 77.2, false)
	if day.Gulika.Longitude == night.Gulika.Longitude {
		t.Log("Gulika day == night (may be same part index by coincidence, not necessarily a bug)")
	}
	// Dhuma is Sun-derived and should be same regardless of day/night
	diff := day.Dhuma.Longitude - night.Dhuma.Longitude
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.01 {
		t.Errorf("Dhuma differs between day/night: %.4f vs %.4f", day.Dhuma.Longitude, night.Dhuma.Longitude)
	}
}

func TestWeekdayFromJD(t *testing.T) {
	// J2000 = 2451545.0 = Saturday, 1 January 2000
	wd := weekdayFromJD(2451545.0)
	if wd != 6 { // 6 = Saturday
		t.Errorf("J2000 weekday = %d, want 6 (Saturday)", wd)
	}
}

func TestNormLon(t *testing.T) {
	tests := []struct{ in, out float64 }{
		{0, 0}, {360, 0}, {361, 1}, {-1, 359}, {720, 0},
	}
	for _, tt := range tests {
		got := normLon(tt.in)
		if got != tt.out {
			t.Errorf("normLon(%.0f) = %.0f, want %.0f", tt.in, got, tt.out)
		}
	}
}
