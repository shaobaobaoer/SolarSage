package dignity

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/bounds"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Dignity represents an essential dignity type
type Dignity string

// The five essential dignities of traditional Hellenistic/Medieval astrology,
// scored per Lilly's Christian Astrology weighting system.
const (
	Rulership  Dignity = "RULERSHIP"   // +5: planet rules the sign
	Exaltation Dignity = "EXALTATION"  // +4: planet is exalted in the sign
	Triplicity Dignity = "TRIPLICITY"  // +3: planet rules the element triplicity
	Term       Dignity = "TERM"        // +2: planet rules the Egyptian term (bound)
	Face       Dignity = "FACE"        // +1: planet rules the Chaldean decan (face)
	Detriment  Dignity = "DETRIMENT"   // -5: planet is in detriment
	Fall       Dignity = "FALL"        // -4: planet is in fall
)

// DignityInfo holds the complete essential dignity analysis for a planet
type DignityInfo struct {
	PlanetID      models.PlanetID `json:"planet_id"`
	Longitude     float64         `json:"longitude"`
	Sign          string          `json:"sign"`
	Dignities     []Dignity       `json:"dignities,omitempty"`
	Score         int             `json:"score"`
	Ruler         models.PlanetID `json:"ruler"`         // modern ruler of the sign
	TraditionalRuler models.PlanetID `json:"traditional_ruler"` // pre-modern ruler
	Exalted       bool            `json:"exalted"`
	InDetriment   bool            `json:"in_detriment"`
	InFall        bool            `json:"in_fall"`
	TriplicityRuler models.PlanetID `json:"triplicity_ruler,omitempty"`
	TermRuler     models.PlanetID `json:"term_ruler,omitempty"`
	FaceRuler     models.PlanetID `json:"face_ruler,omitempty"`
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

// triplicityDay maps sign -> day triplicity ruler (traditional, Dorotheus system)
// Fire: Sun (day) / Jupiter (night)
// Earth: Venus (day) / Moon (night)
// Air: Saturn (day) / Mercury (night)
// Water: Venus (day) / Mars (night)
var triplicityDay = map[string]models.PlanetID{
	"Aries":       models.PlanetSun,
	"Leo":         models.PlanetSun,
	"Sagittarius": models.PlanetSun,
	"Taurus":      models.PlanetVenus,
	"Virgo":       models.PlanetVenus,
	"Capricorn":   models.PlanetVenus,
	"Gemini":      models.PlanetSaturn,
	"Libra":       models.PlanetSaturn,
	"Aquarius":    models.PlanetSaturn,
	"Cancer":      models.PlanetVenus,
	"Scorpio":     models.PlanetVenus,
	"Pisces":      models.PlanetVenus,
}

// triplicityNight maps sign -> night triplicity ruler (Dorotheus system)
var triplicityNight = map[string]models.PlanetID{
	"Aries":       models.PlanetJupiter,
	"Leo":         models.PlanetJupiter,
	"Sagittarius": models.PlanetJupiter,
	"Taurus":      models.PlanetMoon,
	"Virgo":       models.PlanetMoon,
	"Capricorn":   models.PlanetMoon,
	"Gemini":      models.PlanetMercury,
	"Libra":       models.PlanetMercury,
	"Aquarius":    models.PlanetMercury,
	"Cancer":      models.PlanetMars,
	"Scorpio":     models.PlanetMars,
	"Pisces":      models.PlanetMars,
}

// SignRuler returns the modern ruler of a zodiac sign
func SignRuler(sign string) models.PlanetID {
	return rulershipMap[sign]
}

// SignTraditionalRuler returns the traditional ruler of a zodiac sign
func SignTraditionalRuler(sign string) models.PlanetID {
	return traditionalRulerMap[sign]
}

// CalcDignity computes the full five-tier essential dignity for a planet at a
// given ecliptic longitude. The longitude is required for Term and Face lookups.
func CalcDignity(planet models.PlanetID, longitude float64) DignityInfo {
	sign := models.SignFromLongitude(longitude)
	info := DignityInfo{
		PlanetID:         planet,
		Longitude:        longitude,
		Sign:             sign,
		Ruler:            rulershipMap[sign],
		TraditionalRuler: traditionalRulerMap[sign],
	}

	// --- Tier 1: Rulership (+5) ---
	if rulershipMap[sign] == planet || traditionalRulerMap[sign] == planet {
		info.Dignities = append(info.Dignities, Rulership)
		info.Score += 5
	}

	// --- Tier 2: Exaltation (+4) ---
	if exSign, ok := exaltationMap[planet]; ok && exSign == sign {
		info.Dignities = append(info.Dignities, Exaltation)
		info.Exalted = true
		info.Score += 4
	}

	// --- Tier 3: Triplicity (+3) ---
	// Day/night triplicity ruler: day ruler is used unless caller specifies
	// explicitly. CalcDignityWithSect provides sect-aware triplicity.
	if dayRuler, ok := triplicityDay[sign]; ok && dayRuler == planet {
		info.Dignities = append(info.Dignities, Triplicity)
		info.TriplicityRuler = dayRuler
		info.Score += 3
	} else if nightRuler, ok := triplicityNight[sign]; ok && nightRuler == planet {
		info.Dignities = append(info.Dignities, Triplicity)
		info.TriplicityRuler = nightRuler
		info.Score += 3
	}

	// --- Tier 4: Term / Bound (+2) — Egyptian terms from pkg/bounds ---
	termInfo := bounds.CalcTerm(longitude)
	info.TermRuler = termInfo.TermRuler
	if termInfo.TermRuler == planet {
		info.Dignities = append(info.Dignities, Term)
		info.Score += 2
	}

	// --- Tier 5: Face / Decan (+1) — Chaldean decans from pkg/bounds ---
	decanInfo := bounds.CalcDecan(longitude)
	info.FaceRuler = decanInfo.DecanRuler
	if decanInfo.DecanRuler == planet {
		info.Dignities = append(info.Dignities, Face)
		info.Score += 1
	}

	// --- Detriment (-5) ---
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

	// --- Fall (-4) ---
	if fSign, ok := fallMap[planet]; ok && fSign == sign {
		info.Dignities = append(info.Dignities, Fall)
		info.InFall = true
		info.Score -= 4
	}

	return info
}

// CalcDignityWithSect computes essential dignity with sect-aware triplicity.
// isDayChart should be true if the Sun is above the horizon at birth time.
func CalcDignityWithSect(planet models.PlanetID, longitude float64, isDayChart bool) DignityInfo {
	info := CalcDignity(planet, longitude)
	sign := info.Sign

	// Re-evaluate triplicity with correct sect ruler
	// Remove any triplicity already added by CalcDignity (sect-neutral)
	filtered := info.Dignities[:0]
	for _, d := range info.Dignities {
		if d != Triplicity {
			filtered = append(filtered, d)
		} else {
			info.Score -= 3 // undo the neutral assignment
		}
	}
	info.Dignities = filtered
	info.TriplicityRuler = ""

	var sectRuler models.PlanetID
	if isDayChart {
		sectRuler = triplicityDay[sign]
	} else {
		sectRuler = triplicityNight[sign]
	}
	info.TriplicityRuler = sectRuler
	if sectRuler == planet {
		info.Dignities = append(info.Dignities, Triplicity)
		info.Score += 3
	}

	return info
}

// CalcChartDignities computes full five-tier essential dignities for all planets.
// Positions must include the Longitude field.
func CalcChartDignities(positions []models.PlanetPosition) []DignityInfo {
	dignities := make([]DignityInfo, 0, len(positions))
	for _, p := range positions {
		d := CalcDignity(p.PlanetID, p.Longitude)
		dignities = append(dignities, d)
	}
	return dignities
}

// AlmutenFiguris returns the planet with the highest five-tier dignity score
// at the given ecliptic longitude — the "Lord of the Geniture" in Hellenistic tradition.
// In case of a tie, the planet with the higher natural strength order is preferred.
func AlmutenFiguris(longitude float64) (models.PlanetID, int) {
	candidates := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter, models.PlanetSaturn,
	}

	var best models.PlanetID
	bestScore := -999
	for _, p := range candidates {
		d := CalcDignity(p, longitude)
		if d.Score > bestScore {
			bestScore = d.Score
			best = p
		}
	}
	return best, bestScore
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
