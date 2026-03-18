package profection

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/dignity"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// ProfectionResult holds the annual profection analysis
type ProfectionResult struct {
	Age             int             `json:"age"`
	ProfectedSign   string          `json:"profected_sign"`
	ProfectedHouse  int             `json:"profected_house"`
	TimeLord        models.PlanetID `json:"time_lord"`
	TraditionalLord models.PlanetID `json:"traditional_time_lord"`
	ProfectedASC    float64         `json:"profected_asc"`
}

// CalcAnnualProfection computes the annual profection for a given age.
// The profected sign advances one sign per year from the natal ASC sign.
func CalcAnnualProfection(natalASC float64, natalHouses []float64, age int) ProfectionResult {
	// Profection advances 30° per year
	profectedASC := sweph.NormalizeDegrees(natalASC + float64(age)*30.0)
	profectedSign := models.SignFromLongitude(profectedASC)
	profectedHouse := (age % 12) + 1

	return ProfectionResult{
		Age:             age,
		ProfectedSign:   profectedSign,
		ProfectedHouse:  profectedHouse,
		TimeLord:        dignity.SignRuler(profectedSign),
		TraditionalLord: dignity.SignTraditionalRuler(profectedSign),
		ProfectedASC:    profectedASC,
	}
}

// MonthlyProfection returns the monthly sub-profection within an annual profection year.
// Each month advances one sign from the annual profected sign.
type MonthlyProfection struct {
	Month           int             `json:"month"`
	ProfectedSign   string          `json:"profected_sign"`
	TimeLord        models.PlanetID `json:"time_lord"`
	TraditionalLord models.PlanetID `json:"traditional_time_lord"`
}

// CalcMonthlyProfections returns 12 monthly profections for the given year
func CalcMonthlyProfections(natalASC float64, age int) []MonthlyProfection {
	annualASC := sweph.NormalizeDegrees(natalASC + float64(age)*30.0)
	months := make([]MonthlyProfection, 12)

	for i := 0; i < 12; i++ {
		monthASC := sweph.NormalizeDegrees(annualASC + float64(i)*30.0)
		sign := models.SignFromLongitude(monthASC)
		months[i] = MonthlyProfection{
			Month:           i + 1,
			ProfectedSign:   sign,
			TimeLord:        dignity.SignRuler(sign),
			TraditionalLord: dignity.SignTraditionalRuler(sign),
		}
	}
	return months
}

// ProfectionTimeline returns annual profections for a range of ages
func ProfectionTimeline(natalASC float64, natalHouses []float64, startAge, endAge int) []ProfectionResult {
	results := make([]ProfectionResult, 0, endAge-startAge+1)
	for age := startAge; age <= endAge; age++ {
		results = append(results, CalcAnnualProfection(natalASC, natalHouses, age))
	}
	return results
}
