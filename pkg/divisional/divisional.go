package divisional

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

// VargaType identifies a Vedic divisional chart.
type VargaType string

const (
	VargaRasi           VargaType = "D1"  // Natal chart (1/1)
	VargaHora           VargaType = "D2"  // Wealth (1/2)
	VargaDrekkana       VargaType = "D3"  // Siblings (1/3)
	VargaChaturthamsa   VargaType = "D4"  // Fortune, property (1/4)
	VargaSaptamsa       VargaType = "D7"  // Children (1/7)
	VargaNavamsa        VargaType = "D9"  // Spouse, dharma (1/9)
	VargaDasamsa        VargaType = "D10" // Career (1/10)
	VargaDwadasamsa     VargaType = "D12" // Parents (1/12)
	VargaShodasamsa     VargaType = "D16" // Vehicles, comforts (1/16)
	VargaVimsamsa       VargaType = "D20" // Spiritual progress (1/20)
	VargaSiddhamsa      VargaType = "D24" // Education (1/24)
	VargaSaptavimsamsa  VargaType = "D27" // Strength (1/27)
	VargaTrimsamsa      VargaType = "D30" // Misfortune (1/30)
	VargaKhavedamsa     VargaType = "D40" // Auspicious effects (1/40)
	VargaAkshavedamsa   VargaType = "D45" // General indications (1/45)
	VargaShashtiamsa    VargaType = "D60" // Past life karma (1/60)
)

// vargaDivision maps each VargaType to its numeric divisor.
var vargaDivision = map[VargaType]int{
	VargaRasi:          1,
	VargaHora:          2,
	VargaDrekkana:      3,
	VargaChaturthamsa:  4,
	VargaSaptamsa:      7,
	VargaNavamsa:       9,
	VargaDasamsa:       10,
	VargaDwadasamsa:    12,
	VargaShodasamsa:    16,
	VargaVimsamsa:      20,
	VargaSiddhamsa:     24,
	VargaSaptavimsamsa: 27,
	VargaTrimsamsa:     30,
	VargaKhavedamsa:    40,
	VargaAkshavedamsa:  45,
	VargaShashtiamsa:   60,
}

// VargaDescription returns a human-readable name for a varga type.
var VargaDescription = map[VargaType]string{
	VargaRasi:          "Rasi (Natal)",
	VargaHora:          "Hora (Wealth)",
	VargaDrekkana:      "Drekkana (Siblings)",
	VargaChaturthamsa:  "Chaturthamsa (Fortune)",
	VargaSaptamsa:      "Saptamsa (Children)",
	VargaNavamsa:       "Navamsa (Spouse/Dharma)",
	VargaDasamsa:       "Dasamsa (Career)",
	VargaDwadasamsa:    "Dwadasamsa (Parents)",
	VargaShodasamsa:    "Shodasamsa (Vehicles)",
	VargaVimsamsa:      "Vimsamsa (Spiritual)",
	VargaSiddhamsa:     "Siddhamsa (Education)",
	VargaSaptavimsamsa: "Saptavimsamsa (Strength)",
	VargaTrimsamsa:     "Trimsamsa (Misfortune)",
	VargaKhavedamsa:    "Khavedamsa (Auspicious)",
	VargaAkshavedamsa:  "Akshavedamsa (General)",
	VargaShashtiamsa:   "Shashtiamsa (Past Life)",
}

// DivisionalPosition holds a planet's position in a divisional chart.
type DivisionalPosition struct {
	PlanetID    models.PlanetID `json:"planet_id"`
	TropicalLon float64        `json:"tropical_longitude"`
	SiderealLon float64        `json:"sidereal_longitude"`
	VargaLon    float64        `json:"varga_longitude"`
	VargaSign   string         `json:"varga_sign"`
	VargaDegree float64        `json:"varga_degree"`
}

// DivisionalChart holds a complete divisional chart result.
type DivisionalChart struct {
	Varga     VargaType             `json:"varga"`
	Ayanamsa  string                `json:"ayanamsa"`
	Positions []DivisionalPosition  `json:"positions"`
}

// navamsaStartOffset returns the starting sign index for Navamsa (D9)
// based on the element of the natal sign.
//
//	Fire signs  (Aries=0, Leo=4, Sagittarius=8):  start from Aries (0)
//	Earth signs (Taurus=1, Virgo=5, Capricorn=9): start from Capricorn (9)
//	Air signs   (Gemini=2, Libra=6, Aquarius=10): start from Libra (6)
//	Water signs (Cancer=3, Scorpio=7, Pisces=11): start from Cancer (3)
func navamsaStartOffset(signIdx int) int {
	switch signIdx % 4 {
	case 0: // fire
		return 0
	case 1: // earth
		return 9
	case 2: // air
		return 6
	case 3: // water
		return 3
	}
	return 0
}

// CalcVargaPosition computes the divisional-chart longitude for a given
// sidereal longitude and division number.
//
// For Navamsa (division=9) the element-specific starting offsets are applied.
// For all other vargas the general Parashari formula is used:
//
//	vargaSignIdx = (signIdx * division + partIdx) % 12
func CalcVargaPosition(siderealLon float64, division int) float64 {
	siderealLon = sweph.NormalizeDegrees(siderealLon)
	signIdx := int(siderealLon / 30.0)
	signDeg := math.Mod(siderealLon, 30.0)

	if division <= 0 {
		return siderealLon
	}

	partSize := 30.0 / float64(division)
	partIdx := int(signDeg / partSize)
	if partIdx >= division {
		partIdx = division - 1
	}

	var vargaSignIdx int
	if division == 9 {
		// Navamsa uses element-based starting offset
		startSign := navamsaStartOffset(signIdx)
		vargaSignIdx = (startSign + partIdx) % 12
	} else {
		// General Parashari formula
		vargaSignIdx = (signIdx*division + partIdx) % 12
	}

	// The degree within the varga sign is the fractional position
	// within the subdivision, scaled back to 0-30.
	fracInPart := signDeg - float64(partIdx)*partSize
	vargaDeg := fracInPart * (30.0 / partSize)

	return float64(vargaSignIdx)*30.0 + vargaDeg
}

// CalcNavamsaPosition is a convenience wrapper for the D9 chart.
func CalcNavamsaPosition(siderealLon float64) float64 {
	return CalcVargaPosition(siderealLon, 9)
}

// CalcDivisionalChart computes a full divisional chart for the given birth data.
func CalcDivisionalChart(lat, lon, jdUT float64, varga VargaType, ayanamsa vedic.Ayanamsa) (*DivisionalChart, error) {
	division, ok := vargaDivision[varga]
	if !ok {
		return nil, fmt.Errorf("unsupported varga type: %s", varga)
	}

	// Compute the tropical natal chart
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}
	tropicalChart, err := chart.CalcSingleChart(lat, lon, jdUT, planets,
		models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		return nil, fmt.Errorf("chart calculation failed: %w", err)
	}

	// Get ayanamsa value
	ayanamsaVal, err := vedic.GetAyanamsa(jdUT, ayanamsa)
	if err != nil {
		return nil, err
	}

	positions := make([]DivisionalPosition, len(tropicalChart.Planets))
	for i, p := range tropicalChart.Planets {
		sidLon := vedic.TropicalToSidereal(p.Longitude, ayanamsaVal)
		vargaLon := CalcVargaPosition(sidLon, division)

		positions[i] = DivisionalPosition{
			PlanetID:    p.PlanetID,
			TropicalLon: p.Longitude,
			SiderealLon: sidLon,
			VargaLon:    vargaLon,
			VargaSign:   models.SignFromLongitude(vargaLon),
			VargaDegree: models.SignDegreeFromLongitude(vargaLon),
		}
	}

	return &DivisionalChart{
		Varga:     varga,
		Ayanamsa:  string(ayanamsa),
		Positions: positions,
	}, nil
}
