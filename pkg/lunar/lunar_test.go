package lunar

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

const j2000 = 2451545.0

func TestCalcLunarPhase(t *testing.T) {
	info, err := CalcLunarPhase(j2000)
	if err != nil {
		t.Fatalf("CalcLunarPhase: %v", err)
	}
	if info.Phase == "" {
		t.Error("Phase is empty")
	}
	if info.PhaseName == "" {
		t.Error("PhaseName is empty")
	}
	if info.Illumination < 0 || info.Illumination > 1 {
		t.Errorf("Illumination out of range: %.4f", info.Illumination)
	}
	if info.PhaseAngle < 0 || info.PhaseAngle >= 360 {
		t.Errorf("PhaseAngle out of range: %.4f", info.PhaseAngle)
	}
}

func TestFindLunarPhases_OneMonth(t *testing.T) {
	// One month should have ~4 major phases (new, first quarter, full, last quarter)
	phases, err := FindLunarPhases(j2000, j2000+30)
	if err != nil {
		t.Fatalf("FindLunarPhases: %v", err)
	}

	if len(phases) < 3 || len(phases) > 6 {
		t.Errorf("Expected 3-6 phases in 30 days, got %d", len(phases))
	}

	// Check all phases have valid data
	for _, p := range phases {
		if p.MoonSign == "" {
			t.Error("Phase has empty moon sign")
		}
		if p.JD < j2000 || p.JD > j2000+30 {
			t.Errorf("Phase JD out of range: %.2f", p.JD)
		}
	}
}

func TestFindLunarPhases_OneYear(t *testing.T) {
	// One year should have ~49 major phases (~12-13 each of new/full/quarters)
	phases, err := FindLunarPhases(j2000, j2000+365.25)
	if err != nil {
		t.Fatalf("FindLunarPhases: %v", err)
	}

	newCount, fullCount := 0, 0
	for _, p := range phases {
		switch p.Phase {
		case PhaseNewMoon:
			newCount++
		case PhaseFullMoon:
			fullCount++
		}
	}

	if newCount < 12 || newCount > 14 {
		t.Errorf("Expected 12-13 new moons, got %d", newCount)
	}
	if fullCount < 12 || fullCount > 14 {
		t.Errorf("Expected 12-13 full moons, got %d", fullCount)
	}
}

func TestFindEclipses_OneYear(t *testing.T) {
	// 2000: should find some eclipses (there are typically 4-7 per year)
	eclipses, err := FindEclipses(j2000, j2000+365.25)
	if err != nil {
		t.Fatalf("FindEclipses: %v", err)
	}

	if len(eclipses) < 2 || len(eclipses) > 8 {
		t.Errorf("Expected 2-8 eclipses in 2000, got %d", len(eclipses))
	}

	for _, e := range eclipses {
		if e.Type == "" {
			t.Error("Eclipse has empty type")
		}
		if e.MoonSign == "" {
			t.Error("Eclipse has empty moon sign")
		}
		t.Logf("Eclipse: %s at JD %.2f (%s Moon, lat=%.4f)", e.Type, e.JD, e.MoonSign, e.MoonLat)
	}
}

func TestNextNewMoon(t *testing.T) {
	jd, err := NextNewMoon(j2000)
	if err != nil {
		t.Fatalf("NextNewMoon: %v", err)
	}
	if jd <= j2000 {
		t.Error("Next new moon should be after search JD")
	}
	if jd > j2000+32 {
		t.Error("Next new moon should be within ~30 days")
	}
}

func TestNextFullMoon(t *testing.T) {
	jd, err := NextFullMoon(j2000)
	if err != nil {
		t.Fatalf("NextFullMoon: %v", err)
	}
	if jd <= j2000 {
		t.Error("Next full moon should be after search JD")
	}
	if jd > j2000+32 {
		t.Error("Next full moon should be within ~30 days")
	}
}

func TestClassifySolarEclipse(t *testing.T) {
	if classifySolarEclipse(0.1) != EclipseSolarTotal {
		t.Error("0.1 should be total")
	}
	if classifySolarEclipse(0.5) != EclipseSolarAnnular {
		t.Error("0.5 should be annular")
	}
	if classifySolarEclipse(1.0) != EclipseSolarPartial {
		t.Error("1.0 should be partial")
	}
}

func TestClassifyLunarEclipse(t *testing.T) {
	if classifyLunarEclipse(0.1) != EclipseLunarTotal {
		t.Error("0.1 should be total")
	}
	if classifyLunarEclipse(0.5) != EclipseLunarPartial {
		t.Error("0.5 should be partial")
	}
	if classifyLunarEclipse(0.9) != EclipseLunarPenumbral {
		t.Error("0.9 should be penumbral")
	}
}

func TestPhaseFromElongation(t *testing.T) {
	tests := []struct {
		angle float64
		phase Phase
	}{
		{0, PhaseNewMoon},
		{45, PhaseWaxingCrescent},
		{90, PhaseFirstQuarter},
		{135, PhaseWaxingGibbous},
		{180, PhaseFullMoon},
		{225, PhaseWaningGibbous},
		{270, PhaseLastQuarter},
		{315, PhaseWaningCrescent},
	}
	for _, tt := range tests {
		p, _ := phaseFromElongation(tt.angle)
		if p != tt.phase {
			t.Errorf("phaseFromElongation(%.0f) = %s, want %s", tt.angle, p, tt.phase)
		}
	}
}
