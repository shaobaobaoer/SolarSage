package dignity

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Dignity represents an essential dignity type
type Dignity string

const (
	Rulership  Dignity = "RULERSHIP"
	Exaltation Dignity = "EXALTATION"
	Detriment  Dignity = "DETRIMENT"
	Fall       Dignity = "FALL"
)

// DignityInfo holds the essential dignity analysis for a planet
type DignityInfo struct {
	PlanetID   models.PlanetID `json:"planet_id"`
	Sign       string          `json:"sign"`
	Dignities  []Dignity       `json:"dignities,omitempty"`
	Score      int             `json:"score"`
	Ruler      models.PlanetID `json:"ruler"`
	Exalted    bool            `json:"exalted"`
	InDetriment bool           `json:"in_detriment"`
	InFall     bool            `json:"in_fall"`
}

// MutualReceptionInfo holds mutual reception between two planets
type MutualReceptionInfo struct {
	PlanetA models.PlanetID `json:"planet_a"`
	SignA   string          `json:"sign_a"`
	PlanetB models.PlanetID `json:"planet_b"`
	SignB   string          `json:"sign_b"`
	Type    string          `json:"type"` // "rulership" or "exaltation"
}

// signIndex maps sign name to 0-11 index
var signIndex = map[string]int{
	"Aries": 0, "Taurus": 1, "Gemini": 2, "Cancer": 3,
	"Leo": 4, "Virgo": 5, "Libra": 6, "Scorpio": 7,
	"Sagittarius": 8, "Capricorn": 9, "Aquarius": 10, "Pisces": 11,
}

// rulershipMap: sign -> ruling planet (traditional + modern)
var rulershipMap = map[string]models.PlanetID{
	"Aries":       models.PlanetMars,
	"Taurus":      models.PlanetVenus,
	"Gemini":      models.PlanetMercury,
	"Cancer":      models.PlanetMoon,
	"Leo":         models.PlanetSun,
	"Virgo":       models.PlanetMercury,
	"Libra":       models.PlanetVenus,
	"Scorpio":     models.PlanetPluto,
	"Sagittarius": models.PlanetJupiter,
	"Capricorn":   models.PlanetSaturn,
	"Aquarius":    models.PlanetUranus,
	"Pisces":      models.PlanetNeptune,
}

// traditionalRulerMap: sign -> traditional ruler (pre-modern planets)
var traditionalRulerMap = map[string]models.PlanetID{
	"Aries":       models.PlanetMars,
	"Taurus":      models.PlanetVenus,
	"Gemini":      models.PlanetMercury,
	"Cancer":      models.PlanetMoon,
	"Leo":         models.PlanetSun,
	"Virgo":       models.PlanetMercury,
	"Libra":       models.PlanetVenus,
	"Scorpio":     models.PlanetMars,
	"Sagittarius": models.PlanetJupiter,
	"Capricorn":   models.PlanetSaturn,
	"Aquarius":    models.PlanetSaturn,
	"Pisces":      models.PlanetJupiter,
}

// exaltationMap: planet -> sign of exaltation
var exaltationMap = map[models.PlanetID]string{
	models.PlanetSun:     "Aries",
	models.PlanetMoon:    "Taurus",
	models.PlanetMercury: "Virgo",
	models.PlanetVenus:   "Pisces",
	models.PlanetMars:    "Capricorn",
	models.PlanetJupiter: "Cancer",
	models.PlanetSaturn:  "Libra",
	models.PlanetUranus:  "Scorpio",
	models.PlanetNeptune: "Leo",    // Some traditions use Aquarius/Cancer
	models.PlanetPluto:   "Aries",  // Some traditions use Leo
}

// detrimentMap: planet -> sign(s) of detriment (opposite of rulership)
var detrimentMap = map[models.PlanetID][]string{
	models.PlanetSun:     {"Aquarius"},
	models.PlanetMoon:    {"Capricorn"},
	models.PlanetMercury: {"Sagittarius", "Pisces"},
	models.PlanetVenus:   {"Aries", "Scorpio"},
	models.PlanetMars:    {"Taurus", "Libra"},
	models.PlanetJupiter: {"Gemini", "Virgo"},
	models.PlanetSaturn:  {"Cancer", "Leo"},
	models.PlanetUranus:  {"Leo"},
	models.PlanetNeptune: {"Virgo"},
	models.PlanetPluto:   {"Taurus"},
}

// fallMap: planet -> sign of fall (opposite of exaltation)
var fallMap = map[models.PlanetID]string{
	models.PlanetSun:     "Libra",
	models.PlanetMoon:    "Scorpio",
	models.PlanetMercury: "Pisces",
	models.PlanetVenus:   "Virgo",
	models.PlanetMars:    "Cancer",
	models.PlanetJupiter: "Capricorn",
	models.PlanetSaturn:  "Aries",
	models.PlanetUranus:  "Taurus",
	models.PlanetNeptune: "Aquarius",
	models.PlanetPluto:   "Libra",
}

// SignRuler returns the modern ruler of a zodiac sign
func SignRuler(sign string) models.PlanetID {
	return rulershipMap[sign]
}

// SignTraditionalRuler returns the traditional ruler of a zodiac sign
func SignTraditionalRuler(sign string) models.PlanetID {
	return traditionalRulerMap[sign]
}

// CalcDignity computes the essential dignity for a planet in a given sign
func CalcDignity(planet models.PlanetID, sign string) DignityInfo {
	info := DignityInfo{
		PlanetID: planet,
		Sign:     sign,
		Ruler:    rulershipMap[sign],
	}

	// Check rulership: planet rules the sign it's in
	if rulershipMap[sign] == planet || traditionalRulerMap[sign] == planet {
		info.Dignities = append(info.Dignities, Rulership)
		info.Score += 5
	}

	// Check exaltation
	if exSign, ok := exaltationMap[planet]; ok && exSign == sign {
		info.Dignities = append(info.Dignities, Exaltation)
		info.Exalted = true
		info.Score += 4
	}

	// Check detriment
	if signs, ok := detrimentMap[planet]; ok {
		for _, s := range signs {
			if s == sign {
				info.Dignities = append(info.Dignities, Detriment)
				info.InDetriment = true
				info.Score -= 5
				break
			}
		}
	}

	// Check fall
	if fSign, ok := fallMap[planet]; ok && fSign == sign {
		info.Dignities = append(info.Dignities, Fall)
		info.InFall = true
		info.Score -= 4
	}

	return info
}

// CalcChartDignities computes essential dignities for all planets in a chart
func CalcChartDignities(positions []models.PlanetPosition) []DignityInfo {
	dignities := make([]DignityInfo, 0, len(positions))
	for _, p := range positions {
		d := CalcDignity(p.PlanetID, p.Sign)
		dignities = append(dignities, d)
	}
	return dignities
}

// FindMutualReceptions finds mutual receptions between planets
// A mutual reception occurs when two planets are in each other's ruling signs.
func FindMutualReceptions(positions []models.PlanetPosition) []MutualReceptionInfo {
	var receptions []MutualReceptionInfo

	// Build planet -> sign map
	planetSign := make(map[models.PlanetID]string)
	for _, p := range positions {
		planetSign[p.PlanetID] = p.Sign
	}

	// Check all pairs for rulership mutual reception
	planets := make([]models.PlanetID, 0, len(positions))
	for _, p := range positions {
		planets = append(planets, p.PlanetID)
	}

	for i := 0; i < len(planets); i++ {
		for j := i + 1; j < len(planets); j++ {
			a, b := planets[i], planets[j]
			signA, signB := planetSign[a], planetSign[b]

			// Rulership mutual reception: A rules B's sign AND B rules A's sign
			if (rulershipMap[signB] == a || traditionalRulerMap[signB] == a) &&
				(rulershipMap[signA] == b || traditionalRulerMap[signA] == b) {
				receptions = append(receptions, MutualReceptionInfo{
					PlanetA: a, SignA: signA,
					PlanetB: b, SignB: signB,
					Type: "rulership",
				})
			}

			// Exaltation mutual reception: A is exalted in B's sign AND B is exalted in A's sign
			exA, okA := exaltationMap[a]
			exB, okB := exaltationMap[b]
			if okA && okB && exA == signB && exB == signA {
				receptions = append(receptions, MutualReceptionInfo{
					PlanetA: a, SignA: signA,
					PlanetB: b, SignB: signB,
					Type: "exaltation",
				})
			}
		}
	}

	return receptions
}

// Sect determines if a planet is in sect (diurnal/nocturnal alignment)
type SectInfo struct {
	PlanetID models.PlanetID `json:"planet_id"`
	IsDayChart bool          `json:"is_day_chart"`
	InSect   bool            `json:"in_sect"`
}

// diurnalPlanets are planets that prefer day charts
var diurnalPlanets = map[models.PlanetID]bool{
	models.PlanetSun:     true,
	models.PlanetJupiter: true,
	models.PlanetSaturn:  true,
}

// nocturnalPlanets are planets that prefer night charts
var nocturnalPlanets = map[models.PlanetID]bool{
	models.PlanetMoon:  true,
	models.PlanetVenus: true,
	models.PlanetMars:  true,
}

// CalcSect determines if a planet is in sect
func CalcSect(planet models.PlanetID, isDayChart bool) SectInfo {
	info := SectInfo{
		PlanetID:   planet,
		IsDayChart: isDayChart,
	}

	if isDayChart {
		info.InSect = diurnalPlanets[planet]
	} else {
		info.InSect = nocturnalPlanets[planet]
	}

	// Mercury is a neutral planet - considered in sect in either
	if planet == models.PlanetMercury {
		info.InSect = true
	}

	return info
}
