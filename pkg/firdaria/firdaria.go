package firdaria

import "github.com/shaobaobaoer/solarsage-mcp/pkg/models"

// FirdariaResult holds the complete Firdaria analysis for a given age.
type FirdariaResult struct {
	IsDayBirth    bool             `json:"is_day_birth"`
	Periods       []FirdariaPeriod `json:"periods"`
	CurrentPeriod *FirdariaPeriod  `json:"current_period,omitempty"`
	CurrentSub    *SubPeriod       `json:"current_sub_period,omitempty"`
}

// FirdariaPeriod represents a major Firdaria period ruled by a single planet.
type FirdariaPeriod struct {
	Lord       models.PlanetID `json:"lord"`
	Years      float64         `json:"years"`
	StartAge   float64         `json:"start_age"`
	EndAge     float64         `json:"end_age"`
	SubPeriods []SubPeriod     `json:"sub_periods"`
}

// SubPeriod represents a sub-period within a major Firdaria period.
type SubPeriod struct {
	Lord     models.PlanetID `json:"lord"`
	StartAge float64         `json:"start_age"`
	EndAge   float64         `json:"end_age"`
}

// daySequence is the major period order for day births.
var daySequence = []struct {
	lord  models.PlanetID
	years float64
}{
	{models.PlanetSun, 10},
	{models.PlanetVenus, 8},
	{models.PlanetMercury, 13},
	{models.PlanetMoon, 9},
	{models.PlanetSaturn, 11},
	{models.PlanetJupiter, 12},
	{models.PlanetMars, 7},
	{models.PlanetNorthNodeTrue, 3},
	{models.PlanetSouthNode, 2},
}

// nightSequence is the major period order for night births.
var nightSequence = []struct {
	lord  models.PlanetID
	years float64
}{
	{models.PlanetMoon, 9},
	{models.PlanetSaturn, 11},
	{models.PlanetJupiter, 12},
	{models.PlanetMars, 7},
	{models.PlanetSun, 10},
	{models.PlanetVenus, 8},
	{models.PlanetMercury, 13},
	{models.PlanetNorthNodeTrue, 3},
	{models.PlanetSouthNode, 2},
}

// chaldeanOrder is the Chaldean order of planets (slowest to fastest, geocentric).
var chaldeanOrder = []models.PlanetID{
	models.PlanetSaturn,
	models.PlanetJupiter,
	models.PlanetMars,
	models.PlanetSun,
	models.PlanetVenus,
	models.PlanetMercury,
	models.PlanetMoon,
}

// chaldeanIndex returns the index of a planet in the Chaldean order, or -1.
func chaldeanIndex(id models.PlanetID) int {
	for i, p := range chaldeanOrder {
		if p == id {
			return i
		}
	}
	return -1
}

// subPeriodSequence returns the 7 sub-period lords starting from the major lord
// and proceeding in Chaldean order.
func subPeriodSequence(majorLord models.PlanetID) []models.PlanetID {
	idx := chaldeanIndex(majorLord)
	if idx < 0 {
		// Node periods: use Chaldean order starting from Saturn
		return append([]models.PlanetID{}, chaldeanOrder...)
	}
	seq := make([]models.PlanetID, 7)
	for i := 0; i < 7; i++ {
		seq[i] = chaldeanOrder[(idx+i)%7]
	}
	return seq
}

// buildPeriods constructs the full list of Firdaria periods for one 75-year cycle
// starting at the given base age offset.
func buildPeriods(seq []struct {
	lord  models.PlanetID
	years float64
}, baseAge float64) []FirdariaPeriod {
	periods := make([]FirdariaPeriod, len(seq))
	age := baseAge
	for i, s := range seq {
		endAge := age + s.years
		subLords := subPeriodSequence(s.lord)
		subDur := s.years / float64(len(subLords))
		subs := make([]SubPeriod, len(subLords))
		subAge := age
		for j, sl := range subLords {
			subs[j] = SubPeriod{
				Lord:     sl,
				StartAge: subAge,
				EndAge:   subAge + subDur,
			}
			subAge += subDur
		}
		periods[i] = FirdariaPeriod{
			Lord:       s.lord,
			Years:      s.years,
			StartAge:   age,
			EndAge:     endAge,
			SubPeriods: subs,
		}
		age = endAge
	}
	return periods
}

// CalcFirdaria calculates all Firdaria periods and identifies the current
// period and sub-period for the given age. The cycle repeats every 75 years.
func CalcFirdaria(isDayBirth bool, age float64) *FirdariaResult {
	seq := daySequence
	if !isDayBirth {
		seq = nightSequence
	}

	// Determine which 75-year cycle the age falls into
	cycleAge := age
	if cycleAge < 0 {
		cycleAge = 0
	}
	cycleNum := int(cycleAge / 75.0)
	baseAge := float64(cycleNum) * 75.0

	periods := buildPeriods(seq, baseAge)

	result := &FirdariaResult{
		IsDayBirth: isDayBirth,
		Periods:    periods,
	}

	// Identify current period and sub-period
	for i := range periods {
		p := &periods[i]
		if age >= p.StartAge && age < p.EndAge {
			result.CurrentPeriod = p
			for j := range p.SubPeriods {
				sp := &p.SubPeriods[j]
				if age >= sp.StartAge && age < sp.EndAge {
					result.CurrentSub = sp
					break
				}
			}
			break
		}
	}

	return result
}

// CalcFirdariaTimeline returns only those Firdaria periods that overlap the
// given age range [startAge, endAge].
func CalcFirdariaTimeline(isDayBirth bool, startAge, endAge float64) []FirdariaPeriod {
	seq := daySequence
	if !isDayBirth {
		seq = nightSequence
	}

	var result []FirdariaPeriod

	// Generate enough cycles to cover the range
	cycleStart := int(startAge / 75.0)
	if startAge < 0 {
		cycleStart = 0
	}
	cycleEnd := int(endAge/75.0) + 1

	for c := cycleStart; c < cycleEnd; c++ {
		baseAge := float64(c) * 75.0
		periods := buildPeriods(seq, baseAge)
		for _, p := range periods {
			if p.EndAge > startAge && p.StartAge < endAge {
				result = append(result, p)
			}
		}
	}

	return result
}
