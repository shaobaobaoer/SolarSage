package lots

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// LotResult holds a calculated Arabic lot/part
type LotResult struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Sign      string  `json:"sign"`
	SignDeg   float64 `json:"sign_degree"`
	Formula   string  `json:"formula"`
}

// lotDef defines a standard Arabic lot
type lotDef struct {
	Name       string
	DayBodyA   string // added in day
	DayBodyB   string // subtracted in day
	NightBodyA string // added at night (if different)
	NightBodyB string // subtracted at night
	Reverses   bool   // whether day/night formula reverses
}

// standardLots defines the most commonly used Arabic lots
var standardLots = []lotDef{
	{Name: "Lot of Fortune", DayBodyA: "MOON", DayBodyB: "SUN", Reverses: true},
	{Name: "Lot of Spirit", DayBodyA: "SUN", DayBodyB: "MOON", Reverses: true},
	{Name: "Lot of Eros", DayBodyA: "VENUS", DayBodyB: "SPIRIT", Reverses: true},
	{Name: "Lot of Necessity", DayBodyA: "FORTUNE", DayBodyB: "MERCURY", Reverses: true},
	{Name: "Lot of Courage", DayBodyA: "FORTUNE", DayBodyB: "MARS", Reverses: true},
	{Name: "Lot of Victory", DayBodyA: "SPIRIT", DayBodyB: "JUPITER", Reverses: true},
	{Name: "Lot of Nemesis", DayBodyA: "FORTUNE", DayBodyB: "SATURN", Reverses: true},
	{Name: "Lot of Marriage (M)", DayBodyA: "VENUS", DayBodyB: "SATURN", Reverses: true},
	{Name: "Lot of Marriage (F)", DayBodyA: "SATURN", DayBodyB: "VENUS", Reverses: true},
	{Name: "Lot of Children", DayBodyA: "SATURN", DayBodyB: "JUPITER", Reverses: true},
	{Name: "Lot of Father", DayBodyA: "SUN", DayBodyB: "SATURN", Reverses: true},
	{Name: "Lot of Mother", DayBodyA: "MOON", DayBodyB: "VENUS", Reverses: true},
	{Name: "Lot of Siblings", DayBodyA: "SATURN", DayBodyB: "MERCURY", Reverses: true},
	{Name: "Lot of Disease", DayBodyA: "MARS", DayBodyB: "SATURN", Reverses: true},
	{Name: "Lot of Death", DayBodyA: "SATURN", DayBodyB: "MOON", NightBodyA: "MOON", NightBodyB: "SATURN", Reverses: false},
}

// CalcStandardLots computes all standard Arabic lots for a chart
func CalcStandardLots(positions []models.PlanetPosition, asc float64, isDayChart bool) []LotResult {
	// Build lookup map
	lonMap := make(map[string]float64)
	for _, p := range positions {
		lonMap[string(p.PlanetID)] = p.Longitude
	}

	var results []LotResult

	// First compute Fortune and Spirit since other lots may reference them
	fortuneLon := calcLotValue(asc, lonMap, "MOON", "SUN", isDayChart, true)
	spiritLon := calcLotValue(asc, lonMap, "SUN", "MOON", isDayChart, true)
	lonMap["FORTUNE"] = fortuneLon
	lonMap["SPIRIT"] = spiritLon

	for _, lot := range standardLots {
		bodyA, bodyB := lot.DayBodyA, lot.DayBodyB
		reverses := lot.Reverses

		if !isDayChart {
			if lot.NightBodyA != "" {
				bodyA, bodyB = lot.NightBodyA, lot.NightBodyB
			} else if reverses {
				bodyA, bodyB = bodyB, bodyA
			}
		}

		lonA, okA := lonMap[bodyA]
		lonB, okB := lonMap[bodyB]
		if !okA || !okB {
			continue
		}

		lon := sweph.NormalizeDegrees(asc + lonA - lonB)
		formula := fmt.Sprintf("ASC + %s - %s", models.BodyDisplayName(bodyA), models.BodyDisplayName(bodyB))
		if !isDayChart && reverses && lot.NightBodyA == "" {
			formula += " (night reversal)"
		}

		results = append(results, LotResult{
			Name:      lot.Name,
			Longitude: lon,
			Sign:      models.SignFromLongitude(lon),
			SignDeg:   models.SignDegreeFromLongitude(lon),
			Formula:   formula,
		})
	}

	return results
}

// CalcCustomLot computes a custom Arabic lot: ASC + bodyA - bodyB
func CalcCustomLot(asc, lonA, lonB float64, isDayChart, reverseAtNight bool) float64 {
	if !isDayChart && reverseAtNight {
		lonA, lonB = lonB, lonA
	}
	return sweph.NormalizeDegrees(asc + lonA - lonB)
}

func calcLotValue(asc float64, lonMap map[string]float64, bodyA, bodyB string, isDayChart, reverses bool) float64 {
	a, b := bodyA, bodyB
	if !isDayChart && reverses {
		a, b = b, a
	}
	lonA, okA := lonMap[a]
	lonB, okB := lonMap[b]
	if !okA || !okB {
		return 0
	}
	return sweph.NormalizeDegrees(asc + lonA - lonB)
}
