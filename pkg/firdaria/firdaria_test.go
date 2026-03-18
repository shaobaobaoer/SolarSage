package firdaria

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestDayBirth_Age0_SunPeriod(t *testing.T) {
	r := CalcFirdaria(true, 0)
	if r.CurrentPeriod == nil {
		t.Fatal("CurrentPeriod is nil at age 0")
	}
	if r.CurrentPeriod.Lord != models.PlanetSun {
		t.Errorf("Day birth age 0: lord = %s, want SUN", r.CurrentPeriod.Lord)
	}
	if r.CurrentPeriod.Years != 10 {
		t.Errorf("Day birth Sun period years = %v, want 10", r.CurrentPeriod.Years)
	}
}

func TestNightBirth_Age0_MoonPeriod(t *testing.T) {
	r := CalcFirdaria(false, 0)
	if r.CurrentPeriod == nil {
		t.Fatal("CurrentPeriod is nil at age 0")
	}
	if r.CurrentPeriod.Lord != models.PlanetMoon {
		t.Errorf("Night birth age 0: lord = %s, want MOON", r.CurrentPeriod.Lord)
	}
	if r.CurrentPeriod.Years != 9 {
		t.Errorf("Night birth Moon period years = %v, want 9", r.CurrentPeriod.Years)
	}
}

func TestTotalYears_DayBirth(t *testing.T) {
	r := CalcFirdaria(true, 0)
	var total float64
	for _, p := range r.Periods {
		total += p.Years
	}
	if total != 75 {
		t.Errorf("Day birth total years = %v, want 75", total)
	}
}

func TestTotalYears_NightBirth(t *testing.T) {
	r := CalcFirdaria(false, 0)
	var total float64
	for _, p := range r.Periods {
		total += p.Years
	}
	if total != 75 {
		t.Errorf("Night birth total years = %v, want 75", total)
	}
}

func TestSubPeriods_SumToMajor(t *testing.T) {
	r := CalcFirdaria(true, 0)
	for _, p := range r.Periods {
		var subTotal float64
		for _, sp := range p.SubPeriods {
			subTotal += sp.EndAge - sp.StartAge
		}
		if math.Abs(subTotal-p.Years) > 1e-9 {
			t.Errorf("Period %s: sub-period total = %v, want %v", p.Lord, subTotal, p.Years)
		}
	}
}

func TestSubPeriods_CountIs7(t *testing.T) {
	r := CalcFirdaria(true, 0)
	for _, p := range r.Periods {
		if len(p.SubPeriods) != 7 {
			t.Errorf("Period %s: sub-period count = %d, want 7", p.Lord, len(p.SubPeriods))
		}
	}
}

func TestSubPeriod_StartsFromMajorLord(t *testing.T) {
	r := CalcFirdaria(true, 0)
	// Sun period: sub-periods start from Sun in Chaldean order
	sunPeriod := r.Periods[0]
	if sunPeriod.SubPeriods[0].Lord != models.PlanetSun {
		t.Errorf("Sun period first sub-lord = %s, want SUN", sunPeriod.SubPeriods[0].Lord)
	}
	// Venus period: sub-periods start from Venus
	venusPeriod := r.Periods[1]
	if venusPeriod.SubPeriods[0].Lord != models.PlanetVenus {
		t.Errorf("Venus period first sub-lord = %s, want VENUS", venusPeriod.SubPeriods[0].Lord)
	}
}

func TestSubPeriod_ChaldeanOrder(t *testing.T) {
	r := CalcFirdaria(true, 0)
	// Saturn period: starts from Saturn, wraps around
	// Expected: Saturn, Jupiter, Mars, Sun, Venus, Mercury, Moon
	saturnPeriod := r.Periods[4] // Saturn is 5th in day sequence
	expected := []models.PlanetID{
		models.PlanetSaturn, models.PlanetJupiter, models.PlanetMars,
		models.PlanetSun, models.PlanetVenus, models.PlanetMercury, models.PlanetMoon,
	}
	for i, exp := range expected {
		if saturnPeriod.SubPeriods[i].Lord != exp {
			t.Errorf("Saturn sub-period[%d] = %s, want %s", i, saturnPeriod.SubPeriods[i].Lord, exp)
		}
	}
}

func TestCurrentPeriod_MidLife(t *testing.T) {
	// Day birth, age 35: Sun(0-10) Venus(10-18) Mercury(18-31) Moon(31-40)
	r := CalcFirdaria(true, 35)
	if r.CurrentPeriod == nil {
		t.Fatal("CurrentPeriod is nil at age 35")
	}
	if r.CurrentPeriod.Lord != models.PlanetMoon {
		t.Errorf("Day birth age 35: lord = %s, want MOON", r.CurrentPeriod.Lord)
	}
}

func TestCurrentSubPeriod(t *testing.T) {
	// Day birth, age 0: Sun period, first sub-period is Sun
	r := CalcFirdaria(true, 0)
	if r.CurrentSub == nil {
		t.Fatal("CurrentSub is nil at age 0")
	}
	if r.CurrentSub.Lord != models.PlanetSun {
		t.Errorf("Day birth age 0: sub-lord = %s, want SUN", r.CurrentSub.Lord)
	}
}

func TestCycleRestart_Age75(t *testing.T) {
	// Age 75 should restart the cycle
	r := CalcFirdaria(true, 75)
	if r.CurrentPeriod == nil {
		t.Fatal("CurrentPeriod is nil at age 75")
	}
	if r.CurrentPeriod.Lord != models.PlanetSun {
		t.Errorf("Day birth age 75: lord = %s, want SUN (cycle restart)", r.CurrentPeriod.Lord)
	}
	if r.CurrentPeriod.StartAge != 75 {
		t.Errorf("Day birth age 75: start_age = %v, want 75", r.CurrentPeriod.StartAge)
	}
}

func TestTimeline_Range(t *testing.T) {
	periods := CalcFirdariaTimeline(true, 10, 20)
	if len(periods) == 0 {
		t.Fatal("Timeline returned no periods for age 10-20")
	}
	for _, p := range periods {
		if p.EndAge <= 10 || p.StartAge >= 20 {
			t.Errorf("Period %s (%.0f-%.0f) outside range 10-20", p.Lord, p.StartAge, p.EndAge)
		}
	}
}

func TestTimeline_CoversEntireRange(t *testing.T) {
	periods := CalcFirdariaTimeline(true, 0, 75)
	if len(periods) != 9 {
		t.Errorf("Full cycle timeline: got %d periods, want 9", len(periods))
	}
}

func TestDayBirth_PeriodSequence(t *testing.T) {
	r := CalcFirdaria(true, 0)
	expected := []models.PlanetID{
		models.PlanetSun, models.PlanetVenus, models.PlanetMercury,
		models.PlanetMoon, models.PlanetSaturn, models.PlanetJupiter,
		models.PlanetMars, models.PlanetNorthNodeTrue, models.PlanetSouthNode,
	}
	if len(r.Periods) != len(expected) {
		t.Fatalf("Day birth: got %d periods, want %d", len(r.Periods), len(expected))
	}
	for i, exp := range expected {
		if r.Periods[i].Lord != exp {
			t.Errorf("Day birth period[%d] = %s, want %s", i, r.Periods[i].Lord, exp)
		}
	}
}

func TestNightBirth_PeriodSequence(t *testing.T) {
	r := CalcFirdaria(false, 0)
	expected := []models.PlanetID{
		models.PlanetMoon, models.PlanetSaturn, models.PlanetJupiter,
		models.PlanetMars, models.PlanetSun, models.PlanetVenus,
		models.PlanetMercury, models.PlanetNorthNodeTrue, models.PlanetSouthNode,
	}
	if len(r.Periods) != len(expected) {
		t.Fatalf("Night birth: got %d periods, want %d", len(r.Periods), len(expected))
	}
	for i, exp := range expected {
		if r.Periods[i].Lord != exp {
			t.Errorf("Night birth period[%d] = %s, want %s", i, r.Periods[i].Lord, exp)
		}
	}
}

func TestNodePeriods_UseChaldeanOrder(t *testing.T) {
	// North Node period: sub-periods should use Chaldean order from Saturn
	r := CalcFirdaria(true, 0)
	nnPeriod := r.Periods[7] // North Node is 8th
	expected := []models.PlanetID{
		models.PlanetSaturn, models.PlanetJupiter, models.PlanetMars,
		models.PlanetSun, models.PlanetVenus, models.PlanetMercury, models.PlanetMoon,
	}
	for i, exp := range expected {
		if nnPeriod.SubPeriods[i].Lord != exp {
			t.Errorf("NorthNode sub-period[%d] = %s, want %s", i, nnPeriod.SubPeriods[i].Lord, exp)
		}
	}
}

func TestIsDayBirth_Flag(t *testing.T) {
	day := CalcFirdaria(true, 0)
	if !day.IsDayBirth {
		t.Error("IsDayBirth should be true for day birth")
	}
	night := CalcFirdaria(false, 0)
	if night.IsDayBirth {
		t.Error("IsDayBirth should be false for night birth")
	}
}

func TestPeriods_Contiguous(t *testing.T) {
	r := CalcFirdaria(true, 0)
	for i := 1; i < len(r.Periods); i++ {
		if r.Periods[i].StartAge != r.Periods[i-1].EndAge {
			t.Errorf("Gap between period %d and %d: %.2f != %.2f",
				i-1, i, r.Periods[i-1].EndAge, r.Periods[i].StartAge)
		}
	}
}

func TestSubPeriods_Contiguous(t *testing.T) {
	r := CalcFirdaria(true, 0)
	for _, p := range r.Periods {
		for j := 1; j < len(p.SubPeriods); j++ {
			if math.Abs(p.SubPeriods[j].StartAge-p.SubPeriods[j-1].EndAge) > 1e-9 {
				t.Errorf("Period %s: gap between sub %d and %d", p.Lord, j-1, j)
			}
		}
	}
}
