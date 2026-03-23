package ashtakavarga

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// GocharaEntry holds the Ashtakavarga analysis for a transiting planet in one sign.
type GocharaEntry struct {
	TransitPlanet models.PlanetID `json:"transit_planet"`
	SignIndex     int             `json:"sign_index"`   // 0=Aries .. 11=Pisces (sidereal)
	SignName      string          `json:"sign_name"`
	Bindus        int             `json:"bindus"`        // bindu count from that planet's BAV
	IsAuspicious  bool            `json:"is_auspicious"` // true if Bindus >= 4
}

// gocharaThreshold is the classical BPHS threshold: a transit is auspicious
// when the planet's own BAV bindu count for that sign is 4 or more.
const gocharaThreshold = 4

var signNames = [12]string{
	"Aries", "Taurus", "Gemini", "Cancer",
	"Leo", "Virgo", "Libra", "Scorpio",
	"Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

// GocharaScore returns the BAV bindu count for a transiting planet in a given
// sidereal sign index (0=Aries … 11=Pisces). Returns -1 if the planet has no
// BAV table in the result (i.e. not a traditional planet).
func GocharaScore(result *AshtakavargaResult, planet models.PlanetID, signIdx int) int {
	if signIdx < 0 || signIdx > 11 {
		return -1
	}
	for _, t := range result.PlanetTables {
		if t.Planet == planet {
			return t.Bindus[signIdx]
		}
	}
	return -1
}

// IsGocharaAuspicious returns true if the transiting planet's BAV bindu count
// for the given sidereal sign index is >= 4 (the classical BPHS threshold).
func IsGocharaAuspicious(result *AshtakavargaResult, planet models.PlanetID, signIdx int) bool {
	return GocharaScore(result, planet, signIdx) >= gocharaThreshold
}

// GocharaForPlanet returns all 12 sign entries for a single transiting planet,
// showing the BAV bindus and auspiciousness for each sign in the zodiac.
func GocharaForPlanet(result *AshtakavargaResult, planet models.PlanetID) []GocharaEntry {
	entries := make([]GocharaEntry, 12)
	for i := 0; i < 12; i++ {
		bindus := GocharaScore(result, planet, i)
		auspicious := false
		if bindus >= gocharaThreshold {
			auspicious = true
		}
		entries[i] = GocharaEntry{
			TransitPlanet: planet,
			SignIndex:     i,
			SignName:      signNames[i],
			Bindus:        bindus,
			IsAuspicious:  auspicious,
		}
	}
	return entries
}

// GocharaAll returns the full Gochara table for all seven traditional planets:
// for each planet, all 12 signs with their BAV bindu counts and auspiciousness.
func GocharaAll(result *AshtakavargaResult) map[models.PlanetID][]GocharaEntry {
	out := make(map[models.PlanetID][]GocharaEntry, len(traditionalPlanets))
	for _, p := range traditionalPlanets {
		out[p] = GocharaForPlanet(result, p)
	}
	return out
}

// GocharaAtLongitude returns the Gochara entry for a transiting planet at a
// specific sidereal longitude (e.g. current transit position).
func GocharaAtLongitude(result *AshtakavargaResult, planet models.PlanetID, siderealLon float64) GocharaEntry {
	signIdx := signIndexFromLon(siderealLon)
	bindus := GocharaScore(result, planet, signIdx)
	return GocharaEntry{
		TransitPlanet: planet,
		SignIndex:     signIdx,
		SignName:      signNames[signIdx],
		Bindus:        bindus,
		IsAuspicious:  bindus >= gocharaThreshold,
	}
}

// AuspiciousSigns returns all sign indices where a transiting planet has
// bindus >= 4, i.e. the recommended signs for that planet's transit.
func AuspiciousSigns(result *AshtakavargaResult, planet models.PlanetID) []int {
	var signs []int
	for i := 0; i < 12; i++ {
		if GocharaScore(result, planet, i) >= gocharaThreshold {
			signs = append(signs, i)
		}
	}
	return signs
}
