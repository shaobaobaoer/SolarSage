package avastha

import "github.com/shaobaobaoer/solarsage-mcp/pkg/models"

// Avastha represents a planetary state/condition
type Avastha string

const (
	Lajjita   Avastha = "LAJJITA"   // Ashamed — planet in debilitation or enemy sign
	Garvita   Avastha = "GARVITA"   // Proud — planet in own sign or exaltation
	Kshudhita Avastha = "KSHUDHITA" // Hungry — planet in enemy sign aspected by malefic
	Trishita  Avastha = "TRISHITA"  // Thirsty — planet in water sign, malefic aspect, no benefic
	Mudita    Avastha = "MUDITA"    // Happy — planet in friend's sign with benefic aspect
	Kshobhita Avastha = "KSHOBHITA" // Agitated — conjunct Sun and aspected by malefic
)

// AvasthaPlanetInfo holds the input data needed to evaluate avasthas for one planet.
type AvasthaPlanetInfo struct {
	PlanetID  models.PlanetID
	Longitude float64 // sidereal longitude
	SignIndex int     // 0-11
}

// AvasthResult holds the avastha evaluation for a single planet.
type AvasthResult struct {
	PlanetID models.PlanetID `json:"planet_id"`
	Avasthas []Avastha       `json:"avasthas"` // a planet can have multiple states
}

// sign rulership in Vedic astrology (using sidereal signs 0-11)
var signRuler = [12]models.PlanetID{
	models.PlanetMars,    // 0 Aries
	models.PlanetVenus,   // 1 Taurus
	models.PlanetMercury, // 2 Gemini
	models.PlanetMoon,    // 3 Cancer
	models.PlanetSun,     // 4 Leo
	models.PlanetMercury, // 5 Virgo
	models.PlanetVenus,   // 6 Libra
	models.PlanetMars,    // 7 Scorpio
	models.PlanetJupiter, // 8 Sagittarius
	models.PlanetSaturn,  // 9 Capricorn
	models.PlanetSaturn,  // 10 Aquarius
	models.PlanetJupiter, // 11 Pisces
}

// exaltation signs (sidereal sign index)
var exaltationSign = map[models.PlanetID]int{
	models.PlanetSun:     0,  // Aries
	models.PlanetMoon:    1,  // Taurus
	models.PlanetMars:    9,  // Capricorn
	models.PlanetMercury: 5,  // Virgo
	models.PlanetJupiter: 3,  // Cancer
	models.PlanetVenus:   11, // Pisces
	models.PlanetSaturn:  6,  // Libra
}

// debilitation signs (opposite of exaltation)
var debilitationSign = map[models.PlanetID]int{
	models.PlanetSun:     6,  // Libra
	models.PlanetMoon:    7,  // Scorpio
	models.PlanetMars:    3,  // Cancer
	models.PlanetMercury: 11, // Pisces
	models.PlanetJupiter: 9,  // Capricorn
	models.PlanetVenus:   5,  // Virgo
	models.PlanetSaturn:  0,  // Aries
}

// natural friendship map (Parashari scheme)
// 1 = friend, 0 = neutral, -1 = enemy
var naturalFriendship = map[models.PlanetID]map[models.PlanetID]int{
	models.PlanetSun: {
		models.PlanetMoon: 1, models.PlanetMars: 1, models.PlanetJupiter: 1,
		models.PlanetVenus: -1, models.PlanetSaturn: -1, models.PlanetMercury: 0,
	},
	models.PlanetMoon: {
		models.PlanetSun: 1, models.PlanetMercury: 1,
		models.PlanetMars: 0, models.PlanetJupiter: 0, models.PlanetVenus: 0, models.PlanetSaturn: 0,
	},
	models.PlanetMars: {
		models.PlanetSun: 1, models.PlanetMoon: 1, models.PlanetJupiter: 1,
		models.PlanetVenus: 0, models.PlanetSaturn: 0, models.PlanetMercury: -1,
	},
	models.PlanetMercury: {
		models.PlanetSun: 1, models.PlanetVenus: 1,
		models.PlanetMars: 0, models.PlanetJupiter: 0, models.PlanetSaturn: 0, models.PlanetMoon: -1,
	},
	models.PlanetJupiter: {
		models.PlanetSun: 1, models.PlanetMoon: 1, models.PlanetMars: 1,
		models.PlanetSaturn: 0, models.PlanetMercury: -1, models.PlanetVenus: -1,
	},
	models.PlanetVenus: {
		models.PlanetMercury: 1, models.PlanetSaturn: 1,
		models.PlanetMars: 0, models.PlanetJupiter: 0, models.PlanetSun: -1, models.PlanetMoon: -1,
	},
	models.PlanetSaturn: {
		models.PlanetMercury: 1, models.PlanetVenus: 1,
		models.PlanetJupiter: 0, models.PlanetSun: -1, models.PlanetMoon: -1, models.PlanetMars: -1,
	},
}

// isMalefic returns true for natural malefics (Mars, Saturn, Sun)
func isMalefic(p models.PlanetID) bool {
	return p == models.PlanetMars || p == models.PlanetSaturn || p == models.PlanetSun
}

// isBenefic returns true for natural benefics (Jupiter, Venus, well-associated Mercury, waxing Moon)
// For simplicity we treat Jupiter, Venus as always benefic.
func isBenefic(p models.PlanetID) bool {
	return p == models.PlanetJupiter || p == models.PlanetVenus
}

// isWaterSign returns true for Cancer(3), Scorpio(7), Pisces(11)
func isWaterSign(signIdx int) bool {
	return signIdx == 3 || signIdx == 7 || signIdx == 11
}

// CalcAvasthas evaluates the six Shayana Avasthas for each planet.
// Parameters:
//   - planets: planet positions (sidereal)
//   - aspects: list of aspects between planets (used for malefic/benefic aspect checks)
func CalcAvasthas(planets []AvasthaPlanetInfo, aspects []AspectRef) []AvasthResult {
	var results []AvasthResult

	aspectedByMalefic := buildAspectMap(aspects, true)
	aspectedByBenefic := buildAspectMap(aspects, false)

	for _, p := range planets {
		// Skip nodes and outer planets
		if !isTraditionalPlanet(p.PlanetID) {
			continue
		}

		var avasthas []Avastha
		signIdx := p.SignIndex
		ruler := signRuler[signIdx]

		isOwn := ruler == p.PlanetID
		isExalted := exaltationSign[p.PlanetID] == signIdx
		isDebilitated := debilitationSign[p.PlanetID] == signIdx
		isEnemy := isEnemySign(p.PlanetID, signIdx)
		isFriend := isFriendSign(p.PlanetID, signIdx)
		hasMaleficAspect := aspectedByMalefic[p.PlanetID]
		hasBeneficAspect := aspectedByBenefic[p.PlanetID]

		// Garvita: in own sign or exaltation
		if isOwn || isExalted {
			avasthas = append(avasthas, Garvita)
		}

		// Lajjita: in debilitation or enemy sign
		if isDebilitated || isEnemy {
			avasthas = append(avasthas, Lajjita)
		}

		// Kshudhita: in enemy sign AND aspected by malefic
		if isEnemy && hasMaleficAspect {
			avasthas = append(avasthas, Kshudhita)
		}

		// Trishita: in water sign, aspected by malefic, NOT aspected by benefic
		if isWaterSign(signIdx) && hasMaleficAspect && !hasBeneficAspect {
			avasthas = append(avasthas, Trishita)
		}

		// Mudita: in friend's sign with benefic aspect
		if isFriend && hasBeneficAspect {
			avasthas = append(avasthas, Mudita)
		}

		// Kshobhita: conjunct Sun AND aspected by malefic (non-Sun)
		if p.PlanetID != models.PlanetSun && isConjunctSun(p.PlanetID, aspects) && hasMaleficAspect {
			avasthas = append(avasthas, Kshobhita)
		}

		results = append(results, AvasthResult{
			PlanetID: p.PlanetID,
			Avasthas: avasthas,
		})
	}

	return results
}

// AspectRef is a minimal aspect reference for avastha evaluation.
type AspectRef struct {
	PlanetA models.PlanetID
	PlanetB models.PlanetID
	Type    string // "conjunction", "opposition", "trine", "square", "sextile"
}

func isTraditionalPlanet(p models.PlanetID) bool {
	switch p {
	case models.PlanetSun, models.PlanetMoon, models.PlanetMars,
		models.PlanetMercury, models.PlanetJupiter, models.PlanetVenus, models.PlanetSaturn:
		return true
	}
	return false
}

func isEnemySign(planet models.PlanetID, signIdx int) bool {
	ruler := signRuler[signIdx]
	if ruler == planet {
		return false
	}
	f, ok := naturalFriendship[planet]
	if !ok {
		return false
	}
	return f[ruler] == -1
}

func isFriendSign(planet models.PlanetID, signIdx int) bool {
	ruler := signRuler[signIdx]
	if ruler == planet {
		return false // own sign, not "friend's sign"
	}
	f, ok := naturalFriendship[planet]
	if !ok {
		return false
	}
	return f[ruler] == 1
}

// buildAspectMap returns a set of planets that are aspected by malefic (if wantMalefic=true)
// or benefic (if wantMalefic=false) planets.
func buildAspectMap(aspects []AspectRef, wantMalefic bool) map[models.PlanetID]bool {
	result := make(map[models.PlanetID]bool)
	for _, a := range aspects {
		if wantMalefic {
			if isMalefic(a.PlanetA) {
				result[a.PlanetB] = true
			}
			if isMalefic(a.PlanetB) {
				result[a.PlanetA] = true
			}
		} else {
			if isBenefic(a.PlanetA) {
				result[a.PlanetB] = true
			}
			if isBenefic(a.PlanetB) {
				result[a.PlanetA] = true
			}
		}
	}
	return result
}

func isConjunctSun(planet models.PlanetID, aspects []AspectRef) bool {
	for _, a := range aspects {
		if a.Type != "conjunction" {
			continue
		}
		if (a.PlanetA == planet && a.PlanetB == models.PlanetSun) ||
			(a.PlanetB == planet && a.PlanetA == models.PlanetSun) {
			return true
		}
	}
	return false
}
