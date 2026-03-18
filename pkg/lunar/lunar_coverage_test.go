package lunar

import "testing"

func TestCalcLunarPhase_MultipleDates(t *testing.T) {
	// Test at various points in the lunar cycle to cover all phase branches
	dates := []float64{
		2451545.0,  // J2000
		2451550.0,  // ~5 days later
		2451552.0,  // ~7 days later
		2451556.0,  // ~11 days later
		2451560.0,  // ~15 days later (near full moon)
		2451563.0,  // ~18 days later
		2451567.0,  // ~22 days later
		2451572.0,  // ~27 days later
	}
	for _, jd := range dates {
		phase, err := CalcLunarPhase(jd)
		if err != nil {
			t.Errorf("CalcLunarPhase(%f): %v", jd, err)
			continue
		}
		if phase.PhaseName == "" {
			t.Errorf("CalcLunarPhase(%f): empty phase name", jd)
		}
		if phase.Illumination < 0 || phase.Illumination > 1 {
			t.Errorf("CalcLunarPhase(%f): illumination %f out of range", jd, phase.Illumination)
		}
		if phase.PhaseAngle < 0 || phase.PhaseAngle > 360 {
			t.Errorf("CalcLunarPhase(%f): phase angle %f out of range", jd, phase.PhaseAngle)
		}
	}
}

func TestFindLunarPhases_ShortRange(t *testing.T) {
	// Very short range - may find zero phases
	phases, err := FindLunarPhases(2451545.0, 2451546.0)
	if err != nil {
		t.Fatalf("FindLunarPhases short range: %v", err)
	}
	_ = phases // may be empty, that's OK
}

func TestFindEclipses_ShortRange(t *testing.T) {
	// Short range with no eclipses expected
	eclipses, err := FindEclipses(2451545.0, 2451560.0)
	if err != nil {
		t.Fatalf("FindEclipses short range: %v", err)
	}
	_ = eclipses
}
