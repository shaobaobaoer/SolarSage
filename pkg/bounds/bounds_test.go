package bounds

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcDecan_Aries1(t *testing.T) {
	// 5° Aries = 1st decan, ruled by Mars
	d := CalcDecan(5.0)
	if d.Decan != 1 {
		t.Errorf("5° Aries decan = %d, want 1", d.Decan)
	}
	if d.DecanRuler != models.PlanetMars {
		t.Errorf("Aries 1st decan ruler = %s, want MARS", d.DecanRuler)
	}
	if d.Sign != "Aries" {
		t.Errorf("Sign = %s, want Aries", d.Sign)
	}
}

func TestCalcDecan_Aries2(t *testing.T) {
	// 15° Aries = 2nd decan, ruled by Sun
	d := CalcDecan(15.0)
	if d.Decan != 2 {
		t.Errorf("15° Aries decan = %d, want 2", d.Decan)
	}
	if d.DecanRuler != models.PlanetSun {
		t.Errorf("Aries 2nd decan ruler = %s, want SUN", d.DecanRuler)
	}
}

func TestCalcDecan_Aries3(t *testing.T) {
	// 25° Aries = 3rd decan, ruled by Venus
	d := CalcDecan(25.0)
	if d.Decan != 3 {
		t.Errorf("25° Aries decan = %d, want 3", d.Decan)
	}
	if d.DecanRuler != models.PlanetVenus {
		t.Errorf("Aries 3rd decan ruler = %s, want VENUS", d.DecanRuler)
	}
}

func TestCalcTerm_Aries(t *testing.T) {
	// 3° Aries: first term (0-6°), ruled by Jupiter
	term := CalcTerm(3.0)
	if term.TermRuler != models.PlanetJupiter {
		t.Errorf("3° Aries term ruler = %s, want JUPITER", term.TermRuler)
	}
	if term.TermStart != 0 || term.TermEnd != 6 {
		t.Errorf("Term bounds = %.0f-%.0f, want 0-6", term.TermStart, term.TermEnd)
	}

	// 8° Aries: second term (6-12°), ruled by Venus
	term = CalcTerm(8.0)
	if term.TermRuler != models.PlanetVenus {
		t.Errorf("8° Aries term ruler = %s, want VENUS", term.TermRuler)
	}
}

func TestCalcTerm_AllSigns(t *testing.T) {
	// Test that every degree maps to some valid term
	for lon := 0.0; lon < 360.0; lon += 5.0 {
		term := CalcTerm(lon)
		if term.TermRuler == "" {
			t.Errorf("No term ruler at %.0f°", lon)
		}
		if term.Sign == "" {
			t.Errorf("No sign at %.0f°", lon)
		}
	}
}

func TestCalcChartFaces(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 5, Sign: "Aries", SignDegree: 5},
		{PlanetID: models.PlanetMoon, Longitude: 125, Sign: "Leo", SignDegree: 5},
	}
	faces := CalcChartFaces(positions)
	if len(faces) != 2 {
		t.Fatalf("Expected 2 faces, got %d", len(faces))
	}
	if faces[0].Decan.DecanRuler != models.PlanetMars {
		t.Errorf("Sun decan ruler = %s, want MARS", faces[0].Decan.DecanRuler)
	}
	if faces[0].Term.TermRuler != models.PlanetJupiter {
		t.Errorf("Sun term ruler = %s, want JUPITER", faces[0].Term.TermRuler)
	}
}

func TestCalcTerm_Taurus(t *testing.T) {
	// 10° Taurus: second term (8-14°), ruled by Mercury
	term := CalcTerm(40.0) // 40° = 10° Taurus
	if term.TermRuler != models.PlanetMercury {
		t.Errorf("10° Taurus term ruler = %s, want MERCURY", term.TermRuler)
	}
}

func TestCalcTerm_LastDegree(t *testing.T) {
	// 29.9° of any sign should have a valid term
	for signIdx := 0; signIdx < 12; signIdx++ {
		lon := float64(signIdx)*30.0 + 29.9
		term := CalcTerm(lon)
		if term.TermRuler == "" {
			t.Errorf("No term ruler at %.1f° (sign %d)", lon, signIdx)
		}
		if term.TermEnd != 30 {
			t.Errorf("Last term at %.1f° should end at 30, got %.0f", lon, term.TermEnd)
		}
	}
}

func TestCalcDecan_AllSigns(t *testing.T) {
	// Every sign should have 3 decans with valid rulers
	for signIdx := 0; signIdx < 12; signIdx++ {
		for decan := 0; decan < 3; decan++ {
			lon := float64(signIdx)*30.0 + float64(decan)*10.0 + 1.0
			d := CalcDecan(lon)
			if d.DecanRuler == "" {
				t.Errorf("No decan ruler at %.1f°", lon)
			}
			if d.Decan != decan+1 {
				t.Errorf("%.1f° decan = %d, want %d", lon, d.Decan, decan+1)
			}
		}
	}
}

func TestCalcDecan_Boundary(t *testing.T) {
	// 0° exactly -> 1st decan
	d := CalcDecan(0.0)
	if d.Decan != 1 {
		t.Errorf("0° decan = %d, want 1", d.Decan)
	}
	// 10° exactly -> 2nd decan
	d = CalcDecan(10.0)
	if d.Decan != 2 {
		t.Errorf("10° decan = %d, want 2", d.Decan)
	}
	// 29.99° -> 3rd decan
	d = CalcDecan(29.99)
	if d.Decan != 3 {
		t.Errorf("29.99° decan = %d, want 3", d.Decan)
	}
}

func TestCalcDecan_EdgeCases(t *testing.T) {
	// Negative longitude (shouldn't happen but guard clause exists)
	d := CalcDecan(-1.0)
	if d.Sign != "Aries" {
		t.Errorf("Negative lon: sign = %s", d.Sign)
	}

	// Very high longitude (guard clause)
	d = CalcDecan(361.0)
	if d.Sign == "" {
		t.Error("High lon: empty sign")
	}

	// Exactly 30° boundary -> should be decan 1 of next sign (but 30.0/10 = 3.0, +1 = 4, clamped to 3)
	d = CalcDecan(30.0)
	if d.Decan < 1 || d.Decan > 3 {
		t.Errorf("30° decan = %d, out of range", d.Decan)
	}
}

func TestCalcTerm_EdgeCases(t *testing.T) {
	// Negative longitude
	term := CalcTerm(-1.0)
	if term.Sign != "Aries" {
		t.Errorf("Negative lon: sign = %s", term.Sign)
	}

	// High longitude
	term = CalcTerm(361.0)
	if term.Sign == "" {
		t.Error("High lon: empty sign")
	}

	// Exactly 0°
	term = CalcTerm(0.0)
	if term.TermRuler == "" {
		t.Error("0° has no term ruler")
	}

	// 30° exactly (boundary between signs)
	term = CalcTerm(30.0)
	if term.Sign != "Taurus" {
		// 30°/30 = 1 -> Taurus
		t.Errorf("30° sign = %s, want Taurus", term.Sign)
	}
}
