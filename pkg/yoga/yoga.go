package yoga

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

// YogaCategory classifies a yoga by its life-theme.
type YogaCategory string

const (
	YogaRaja       YogaCategory = "RAJA"       // Power, authority
	YogaDhana      YogaCategory = "DHANA"      // Wealth
	YogaArishta    YogaCategory = "ARISHTA"    // Difficulty
	YogaMahapurusha YogaCategory = "MAHAPURUSHA" // Great person
	YogaNabhasa    YogaCategory = "NABHASA"    // Pattern-based
	YogaOther      YogaCategory = "OTHER"
)

// YogaResult describes a single detected yoga.
type YogaResult struct {
	Name        string          `json:"name"`
	Category    YogaCategory    `json:"category"`
	Description string          `json:"description"`
	Planets     []models.PlanetID `json:"planets_involved"`
	Strength    string          `json:"strength"` // "strong", "moderate", "weak"
}

// YogaAnalysis holds all detected yogas.
type YogaAnalysis struct {
	Yogas []YogaResult `json:"yogas"`
}

// ---------- sign / house helpers ----------

// signIdx returns the sidereal sign index (0=Aries .. 11=Pisces).
func signIdx(lon float64) int {
	idx := int(lon / 30.0)
	if idx < 0 {
		idx += 12
	}
	if idx > 11 {
		idx = 11
	}
	return idx
}

// houseOf returns the whole-sign house number (1-12) a sidereal longitude
// occupies, given the sidereal ASC longitude.
func houseOf(lon, ascLon float64) int {
	ascSign := signIdx(ascLon)
	pSign := signIdx(lon)
	h := ((pSign - ascSign + 12) % 12) + 1
	return h
}

// isKendra returns true if the house is 1, 4, 7, or 10.
func isKendra(house int) bool {
	return house == 1 || house == 4 || house == 7 || house == 10
}

// isTrikona returns true if the house is 1, 5, or 9.
func isTrikona(house int) bool {
	return house == 1 || house == 5 || house == 9
}

// conjunction returns true if two longitudes are within 10 degrees.
func conjunction(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff <= 10
}

// traditionalRuler returns the traditional (Vedic) ruler of a sign index (0-11).
var traditionalRulers = [12]models.PlanetID{
	models.PlanetMars,    // Aries
	models.PlanetVenus,   // Taurus
	models.PlanetMercury, // Gemini
	models.PlanetMoon,    // Cancer
	models.PlanetSun,     // Leo
	models.PlanetMercury, // Virgo
	models.PlanetVenus,   // Libra
	models.PlanetMars,    // Scorpio
	models.PlanetJupiter, // Sagittarius
	models.PlanetSaturn,  // Capricorn
	models.PlanetSaturn,  // Aquarius
	models.PlanetJupiter, // Pisces
}

// ownOrExalted returns true if the planet is in its own sign or exaltation sign.
func ownOrExalted(planet models.PlanetID, sIdx int) bool {
	// Own sign check
	if traditionalRulers[sIdx] == planet {
		return true
	}
	// Exaltation check
	exaltSign, ok := exaltationSign[planet]
	if ok && exaltSign == sIdx {
		return true
	}
	return false
}

// exaltationSign maps each planet to the sign index where it is exalted.
var exaltationSign = map[models.PlanetID]int{
	models.PlanetSun:     0,  // Aries
	models.PlanetMoon:    1,  // Taurus
	models.PlanetMars:    9,  // Capricorn
	models.PlanetMercury: 5,  // Virgo
	models.PlanetJupiter: 3,  // Cancer
	models.PlanetVenus:   11, // Pisces
	models.PlanetSaturn:  6,  // Libra
}

// ---------- position map helper ----------

type posMap struct {
	positions []vedic.SiderealPosition
	byID      map[models.PlanetID]*vedic.SiderealPosition
	ascLon    float64
}

func newPosMap(positions []vedic.SiderealPosition, siderealASC float64) *posMap {
	pm := &posMap{
		positions: positions,
		byID:      make(map[models.PlanetID]*vedic.SiderealPosition),
		ascLon:    siderealASC,
	}
	for i := range positions {
		pm.byID[positions[i].PlanetID] = &positions[i]
	}
	return pm
}

func (pm *posMap) get(id models.PlanetID) (float64, int, bool) {
	p, ok := pm.byID[id]
	if !ok {
		return 0, 0, false
	}
	return p.SiderealLon, houseOf(p.SiderealLon, pm.ascLon), true
}

// ---------- yoga detection functions ----------

// checkMahapurushaYogas checks for the five Mahapurusha yogas.
func checkMahapurushaYogas(pm *posMap) []YogaResult {
	type mpYoga struct {
		planet models.PlanetID
		name   string
		desc   string
	}
	yogas := []mpYoga{
		{models.PlanetMars, "Ruchaka", "Mars in own/exalted sign in a kendra: courage, leadership, physical prowess"},
		{models.PlanetMercury, "Bhadra", "Mercury in own/exalted sign in a kendra: intelligence, eloquence, learning"},
		{models.PlanetJupiter, "Hamsa", "Jupiter in own/exalted sign in a kendra: wisdom, virtue, spiritual inclination"},
		{models.PlanetVenus, "Malavya", "Venus in own/exalted sign in a kendra: beauty, luxury, artistic talent"},
		{models.PlanetSaturn, "Shasha", "Saturn in own/exalted sign in a kendra: authority, discipline, organizational power"},
	}

	var results []YogaResult
	for _, y := range yogas {
		lon, house, ok := pm.get(y.planet)
		if !ok {
			continue
		}
		sIdx := signIdx(lon)
		if ownOrExalted(y.planet, sIdx) && isKendra(house) {
			strength := "strong"
			results = append(results, YogaResult{
				Name:        fmt.Sprintf("%s Yoga", y.name),
				Category:    YogaMahapurusha,
				Description: y.desc,
				Planets:     []models.PlanetID{y.planet},
				Strength:    strength,
			})
		}
	}
	return results
}

// checkRajaYogas checks whether lords of kendras and trikonas conjoin.
func checkRajaYogas(pm *posMap) []YogaResult {
	ascSign := signIdx(pm.ascLon)

	// Compute kendra and trikona house signs and their lords
	kendraHouses := []int{1, 4, 7, 10}
	trikonaHouses := []int{1, 5, 9}

	kendraLords := make(map[models.PlanetID]int)
	trikonaLords := make(map[models.PlanetID]int)

	for _, h := range kendraHouses {
		sIdx := (ascSign + h - 1) % 12
		lord := traditionalRulers[sIdx]
		kendraLords[lord] = h
	}
	for _, h := range trikonaHouses {
		sIdx := (ascSign + h - 1) % 12
		lord := traditionalRulers[sIdx]
		trikonaLords[lord] = h
	}

	var results []YogaResult
	// Find pairs where a kendra lord and a trikona lord conjoin
	for kLord, kHouse := range kendraLords {
		for tLord, tHouse := range trikonaLords {
			if kLord == tLord {
				// Same planet is lord of both kendra and trikona = yoga by itself
				lonK, houseK, ok := pm.get(kLord)
				if !ok {
					continue
				}
				_ = lonK
				results = append(results, YogaResult{
					Name:        "Raja Yoga",
					Category:    YogaRaja,
					Description: fmt.Sprintf("Lord of houses %d and %d (%s) forms Raja Yoga in house %d", kHouse, tHouse, kLord, houseK),
					Planets:     []models.PlanetID{kLord},
					Strength:    "moderate",
				})
				continue
			}
			lonK, _, okK := pm.get(kLord)
			lonT, _, okT := pm.get(tLord)
			if !okK || !okT {
				continue
			}
			if conjunction(lonK, lonT) {
				strength := "strong"
				if !isKendra(houseOf(lonK, pm.ascLon)) {
					strength = "moderate"
				}
				results = append(results, YogaResult{
					Name:     "Raja Yoga",
					Category: YogaRaja,
					Description: fmt.Sprintf("Lord of kendra %d (%s) conjoins lord of trikona %d (%s)",
						kHouse, kLord, tHouse, tLord),
					Planets:  []models.PlanetID{kLord, tLord},
					Strength: strength,
				})
			}
		}
	}
	return results
}

// checkDhanaYogas checks for wealth combinations.
func checkDhanaYogas(pm *posMap) []YogaResult {
	ascSign := signIdx(pm.ascLon)
	var results []YogaResult

	// Lord of 2nd and 11th in conjunction
	lord2 := traditionalRulers[(ascSign+1)%12]
	lord11 := traditionalRulers[(ascSign+10)%12]
	lon2, _, ok2 := pm.get(lord2)
	lon11, _, ok11 := pm.get(lord11)
	if ok2 && ok11 && lord2 != lord11 && conjunction(lon2, lon11) {
		results = append(results, YogaResult{
			Name:        "Dhana Yoga",
			Category:    YogaDhana,
			Description: fmt.Sprintf("Lord of 2nd (%s) conjoins lord of 11th (%s): wealth accumulation", lord2, lord11),
			Planets:     []models.PlanetID{lord2, lord11},
			Strength:    "strong",
		})
	}

	// Lord of 9th in the 2nd house
	lord9 := traditionalRulers[(ascSign+8)%12]
	lon9, house9, ok9 := pm.get(lord9)
	_ = lon9
	if ok9 && house9 == 2 {
		results = append(results, YogaResult{
			Name:        "Dhana Yoga",
			Category:    YogaDhana,
			Description: fmt.Sprintf("Lord of 9th (%s) in the 2nd house: fortune through luck and dharma", lord9),
			Planets:     []models.PlanetID{lord9},
			Strength:    "moderate",
		})
	}

	return results
}

// checkGajakesariYoga checks if Jupiter is in a kendra from the Moon.
func checkGajakesariYoga(pm *posMap) []YogaResult {
	moonPos, ok := pm.byID[models.PlanetMoon]
	if !ok {
		return nil
	}
	jupPos, ok := pm.byID[models.PlanetJupiter]
	if !ok {
		return nil
	}

	moonSign := signIdx(moonPos.SiderealLon)
	jupSign := signIdx(jupPos.SiderealLon)
	diff := ((jupSign - moonSign) + 12) % 12

	// Kendra from Moon: 0 (1st), 3 (4th), 6 (7th), 9 (10th)
	if diff == 0 || diff == 3 || diff == 6 || diff == 9 {
		houseFromMoon := diff + 1
		if diff == 3 {
			houseFromMoon = 4
		} else if diff == 6 {
			houseFromMoon = 7
		} else if diff == 9 {
			houseFromMoon = 10
		} else {
			houseFromMoon = 1
		}
		return []YogaResult{{
			Name:        "Gajakesari Yoga",
			Category:    YogaOther,
			Description: fmt.Sprintf("Jupiter in house %d from Moon: wisdom, fame, lasting reputation", houseFromMoon),
			Planets:     []models.PlanetID{models.PlanetJupiter, models.PlanetMoon},
			Strength:    "strong",
		}}
	}
	return nil
}

// checkBudhadityaYoga checks for Sun-Mercury conjunction.
func checkBudhadityaYoga(pm *posMap) []YogaResult {
	sunLon, _, okS := pm.get(models.PlanetSun)
	merLon, _, okM := pm.get(models.PlanetMercury)
	if !okS || !okM {
		return nil
	}
	if conjunction(sunLon, merLon) {
		return []YogaResult{{
			Name:        "Budhaditya Yoga",
			Category:    YogaOther,
			Description: "Sun-Mercury conjunction: intelligence, communication skill, administrative ability",
			Planets:     []models.PlanetID{models.PlanetSun, models.PlanetMercury},
			Strength:    "moderate",
		}}
	}
	return nil
}

// checkChandraMangalaYoga checks for Moon-Mars conjunction.
func checkChandraMangalaYoga(pm *posMap) []YogaResult {
	moonLon, _, okMo := pm.get(models.PlanetMoon)
	marsLon, _, okMa := pm.get(models.PlanetMars)
	if !okMo || !okMa {
		return nil
	}
	if conjunction(moonLon, marsLon) {
		return []YogaResult{{
			Name:        "Chandra-Mangala Yoga",
			Category:    YogaOther,
			Description: "Moon-Mars conjunction: earning ability, courage, enterprising nature",
			Planets:     []models.PlanetID{models.PlanetMoon, models.PlanetMars},
			Strength:    "moderate",
		}}
	}
	return nil
}

// AnalyzeYogas detects Vedic yogas from sidereal positions and house cusps.
// The houses parameter is not currently used (whole-sign houses are derived
// from siderealASC), but is accepted for future compatibility with
// quadrant-based house systems.
func AnalyzeYogas(positions []vedic.SiderealPosition, houses []float64, siderealASC float64) *YogaAnalysis {
	pm := newPosMap(positions, siderealASC)
	analysis := &YogaAnalysis{}

	analysis.Yogas = append(analysis.Yogas, checkMahapurushaYogas(pm)...)
	analysis.Yogas = append(analysis.Yogas, checkRajaYogas(pm)...)
	analysis.Yogas = append(analysis.Yogas, checkDhanaYogas(pm)...)
	analysis.Yogas = append(analysis.Yogas, checkGajakesariYoga(pm)...)
	analysis.Yogas = append(analysis.Yogas, checkBudhadityaYoga(pm)...)
	analysis.Yogas = append(analysis.Yogas, checkChandraMangalaYoga(pm)...)

	return analysis
}
