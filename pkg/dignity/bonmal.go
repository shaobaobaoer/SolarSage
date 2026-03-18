package dignity

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// BonMalCondition represents a bonification or maltreatment condition
type BonMalCondition string

const (
	CondBeneficTrine       BonMalCondition = "BENEFIC_TRINE"
	CondBeneficSextile     BonMalCondition = "BENEFIC_SEXTILE"
	CondBeneficConjunction BonMalCondition = "BENEFIC_CONJUNCTION"
	CondInBeneficSign      BonMalCondition = "IN_BENEFIC_SIGN"
	CondMaleficSquare      BonMalCondition = "MALEFIC_SQUARE"
	CondMaleficOpposition  BonMalCondition = "MALEFIC_OPPOSITION"
	CondMaleficConjunction BonMalCondition = "MALEFIC_CONJUNCTION"
	CondInMaleficSign      BonMalCondition = "IN_MALEFIC_SIGN"
	CondCombust            BonMalCondition = "COMBUST"
	CondUnderSunbeams      BonMalCondition = "UNDER_SUNBEAMS"
	CondBesieged           BonMalCondition = "BESIEGED"
)

// BonMalInfo holds the bonification and maltreatment analysis for a planet
type BonMalInfo struct {
	PlanetID      models.PlanetID `json:"planet_id"`
	Bonifications []BonMalDetail  `json:"bonifications"`
	Maltreatments []BonMalDetail  `json:"maltreatments"`
	NetScore      int             `json:"net_score"` // positive = bonified, negative = maltreated
}

// BonMalDetail holds one bonification or maltreatment condition
type BonMalDetail struct {
	Condition BonMalCondition `json:"condition"`
	Source    models.PlanetID `json:"source,omitempty"` // The planet causing it
	Score     int             `json:"score"`
}

// Benefic and malefic planet sets
var benefics = map[models.PlanetID]bool{
	models.PlanetJupiter: true,
	models.PlanetVenus:   true,
}

var malefics = map[models.PlanetID]bool{
	models.PlanetMars:   true,
	models.PlanetSaturn: true,
}

// Signs ruled by benefics (Jupiter: Sagittarius, Pisces; Venus: Taurus, Libra)
var beneficSigns = map[string]bool{
	"Sagittarius": true,
	"Pisces":      true,
	"Taurus":      true,
	"Libra":       true,
}

// Signs ruled by malefics (Mars: Aries, Scorpio; Saturn: Capricorn, Aquarius)
var maleficSigns = map[string]bool{
	"Aries":      true,
	"Scorpio":    true,
	"Capricorn":  true,
	"Aquarius":   true,
}

// Aspect orbs for bonification/maltreatment
const (
	orbConjunction = 8.0
	orbSextile     = 6.0
	orbSquare      = 7.0
	orbTrine       = 7.0
	orbOpposition  = 8.0
)

// Combustion thresholds
const (
	combustOrb       = 8.5
	combustOrbMoon   = 12.0
	underSunbeamsOrb = 17.0
)

// Scoring values
const (
	scoreBeneficTrine       = 3
	scoreBeneficSextile     = 2
	scoreBeneficConjunction = 2
	scoreInBeneficSign      = 1
	scoreMaleficSquare      = -3
	scoreMaleficOpposition  = -3
	scoreMaleficConjunction = -2
	scoreInMaleficSign      = -1
	scoreCombust            = -4
	scoreUnderSunbeams      = -2
	scoreBesieged           = -5
)

// angleDiff returns the shortest angular distance between two longitudes
func angleDiff(lon1, lon2 float64) float64 {
	diff := math.Abs(lon1 - lon2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// CalcBonMal analyzes bonification and maltreatment for a single planet
func CalcBonMal(target models.PlanetID, positions []models.PlanetPosition) BonMalInfo {
	info := BonMalInfo{
		PlanetID:      target,
		Bonifications: []BonMalDetail{},
		Maltreatments: []BonMalDetail{},
	}

	// Find target position
	var targetPos *models.PlanetPosition
	for i := range positions {
		if positions[i].PlanetID == target {
			targetPos = &positions[i]
			break
		}
	}
	if targetPos == nil {
		return info
	}

	// Find Sun position (for combustion checks)
	var sunPos *models.PlanetPosition
	for i := range positions {
		if positions[i].PlanetID == models.PlanetSun {
			sunPos = &positions[i]
			break
		}
	}

	// Find Mars and Saturn positions (for besiegement)
	var marsPos, saturnPos *models.PlanetPosition
	for i := range positions {
		if positions[i].PlanetID == models.PlanetMars {
			marsPos = &positions[i]
		}
		if positions[i].PlanetID == models.PlanetSaturn {
			saturnPos = &positions[i]
		}
	}

	// Check aspects from each other planet
	for _, p := range positions {
		if p.PlanetID == target {
			continue
		}

		diff := angleDiff(targetPos.Longitude, p.Longitude)

		if benefics[p.PlanetID] {
			// Benefic trine
			if math.Abs(diff-120) <= orbTrine {
				d := BonMalDetail{Condition: CondBeneficTrine, Source: p.PlanetID, Score: scoreBeneficTrine}
				info.Bonifications = append(info.Bonifications, d)
				info.NetScore += d.Score
			}
			// Benefic sextile
			if math.Abs(diff-60) <= orbSextile {
				d := BonMalDetail{Condition: CondBeneficSextile, Source: p.PlanetID, Score: scoreBeneficSextile}
				info.Bonifications = append(info.Bonifications, d)
				info.NetScore += d.Score
			}
			// Benefic conjunction
			if diff <= orbConjunction {
				d := BonMalDetail{Condition: CondBeneficConjunction, Source: p.PlanetID, Score: scoreBeneficConjunction}
				info.Bonifications = append(info.Bonifications, d)
				info.NetScore += d.Score
			}
		}

		if malefics[p.PlanetID] {
			// Malefic square
			if math.Abs(diff-90) <= orbSquare {
				d := BonMalDetail{Condition: CondMaleficSquare, Source: p.PlanetID, Score: scoreMaleficSquare}
				info.Maltreatments = append(info.Maltreatments, d)
				info.NetScore += d.Score
			}
			// Malefic opposition
			if math.Abs(diff-180) <= orbOpposition {
				d := BonMalDetail{Condition: CondMaleficOpposition, Source: p.PlanetID, Score: scoreMaleficOpposition}
				info.Maltreatments = append(info.Maltreatments, d)
				info.NetScore += d.Score
			}
			// Malefic conjunction
			if diff <= orbConjunction {
				d := BonMalDetail{Condition: CondMaleficConjunction, Source: p.PlanetID, Score: scoreMaleficConjunction}
				info.Maltreatments = append(info.Maltreatments, d)
				info.NetScore += d.Score
			}
		}
	}

	// Check sign-based conditions
	if beneficSigns[targetPos.Sign] {
		d := BonMalDetail{Condition: CondInBeneficSign, Score: scoreInBeneficSign}
		info.Bonifications = append(info.Bonifications, d)
		info.NetScore += d.Score
	}
	if maleficSigns[targetPos.Sign] {
		d := BonMalDetail{Condition: CondInMaleficSign, Score: scoreInMaleficSign}
		info.Maltreatments = append(info.Maltreatments, d)
		info.NetScore += d.Score
	}

	// Combustion and Under Sunbeams (target must not be the Sun itself)
	if target != models.PlanetSun && sunPos != nil {
		sunDiff := angleDiff(targetPos.Longitude, sunPos.Longitude)

		// Determine combustion orb (Moon uses 12 degrees)
		combOrb := combustOrb
		if target == models.PlanetMoon {
			combOrb = combustOrbMoon
		}

		if sunDiff <= combOrb {
			d := BonMalDetail{Condition: CondCombust, Source: models.PlanetSun, Score: scoreCombust}
			info.Maltreatments = append(info.Maltreatments, d)
			info.NetScore += d.Score
		} else if sunDiff <= underSunbeamsOrb {
			d := BonMalDetail{Condition: CondUnderSunbeams, Source: models.PlanetSun, Score: scoreUnderSunbeams}
			info.Maltreatments = append(info.Maltreatments, d)
			info.NetScore += d.Score
		}
	}

	// Besiegement: target between Mars and Saturn by longitude (both within orb)
	if target != models.PlanetMars && target != models.PlanetSaturn && marsPos != nil && saturnPos != nil {
		if isBesieged(targetPos.Longitude, marsPos.Longitude, saturnPos.Longitude) {
			d := BonMalDetail{Condition: CondBesieged, Score: scoreBesieged}
			info.Maltreatments = append(info.Maltreatments, d)
			info.NetScore += d.Score
		}
	}

	return info
}

// isBesieged checks if targetLon is between marsLon and saturnLon on the
// zodiac circle, with both malefics within 12 degrees of the target.
func isBesieged(targetLon, marsLon, saturnLon float64) bool {
	// Both must be within reasonable orb (use conjunction orb)
	marsDiff := angleDiff(targetLon, marsLon)
	saturnDiff := angleDiff(targetLon, saturnLon)

	if marsDiff > orbConjunction || saturnDiff > orbConjunction {
		return false
	}

	// Check that target is between Mars and Saturn on the zodiac circle.
	// Normalize positions relative to Mars.
	normTarget := math.Mod(targetLon-marsLon+360, 360)
	normSaturn := math.Mod(saturnLon-marsLon+360, 360)

	// Target should be between 0 (Mars) and normSaturn on the shorter arc
	if normSaturn <= 180 {
		// Saturn is ahead of Mars (short arc going forward)
		return normTarget > 0 && normTarget < normSaturn
	}
	// Saturn is behind Mars (short arc going backward)
	return normTarget > normSaturn && normTarget < 360
}

// CalcChartBonMal analyzes bonification and maltreatment for all planets in a chart
func CalcChartBonMal(positions []models.PlanetPosition) []BonMalInfo {
	results := make([]BonMalInfo, 0, len(positions))
	for _, p := range positions {
		results = append(results, CalcBonMal(p.PlanetID, positions))
	}
	return results
}
