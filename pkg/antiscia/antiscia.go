package antiscia

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// AntisciaPoint represents a planet's antiscion or contra-antiscion
type AntisciaPoint struct {
	PlanetID        models.PlanetID `json:"planet_id"`
	OriginalLon     float64         `json:"original_longitude"`
	AntisciaLon     float64         `json:"antiscia_longitude"`
	AntisciaSign    string          `json:"antiscia_sign"`
	AntiSciaDeg     float64         `json:"antiscia_sign_degree"`
	ContraAntisciaLon  float64      `json:"contra_antiscia_longitude"`
	ContraAntisciaSign string       `json:"contra_antiscia_sign"`
	ContraAntiSciaDeg  float64      `json:"contra_antiscia_sign_degree"`
}

// AntisciaPair represents two planets whose antiscia are conjunct
type AntisciaPair struct {
	PlanetA     models.PlanetID `json:"planet_a"`
	PlanetB     models.PlanetID `json:"planet_b"`
	Type        string          `json:"type"` // "antiscia" or "contra_antiscia"
	Orb         float64         `json:"orb"`
}

// CalcAntiscia computes the antiscion of an ecliptic longitude.
// The antiscion mirrors a point across the Cancer-Capricorn axis (solstice axis).
// Formula: antiscion = 360° - longitude + 180° = (180° - longitude) mod 360 ...
// Actually: antiscion mirrors across the 0°Can/0°Cap axis.
// If longitude is X, the antiscion is at (360° - X) reflected over the Cancer/Capricorn axis.
// Standard formula: antiscion(lon) = (360° - lon) mod 360° ... no.
// Correct: The solstice axis runs from 0° Cancer (90°) to 0° Capricorn (270°).
// Mirror across this axis: antiscion = 180° - lon (normalized to [0, 360))
// This means: Aries ↔ Virgo, Taurus ↔ Leo, Gemini ↔ Cancer,
//             Libra ↔ Pisces, Scorpio ↔ Aquarius, Sagittarius ↔ Capricorn
func CalcAntiscia(lon float64) float64 {
	return sweph.NormalizeDegrees(180 - lon)
}

// CalcContraAntiscia computes the contra-antiscion (mirror across equinox axis).
// The equinox axis runs from 0° Aries (0°) to 0° Libra (180°).
// contra_antiscion = 360° - lon
func CalcContraAntiscia(lon float64) float64 {
	return sweph.NormalizeDegrees(360 - lon)
}

// CalcChartAntiscia computes antiscia and contra-antiscia for all positions
func CalcChartAntiscia(positions []models.PlanetPosition) []AntisciaPoint {
	points := make([]AntisciaPoint, len(positions))
	for i, p := range positions {
		aLon := CalcAntiscia(p.Longitude)
		caLon := CalcContraAntiscia(p.Longitude)
		points[i] = AntisciaPoint{
			PlanetID:           p.PlanetID,
			OriginalLon:        p.Longitude,
			AntisciaLon:        aLon,
			AntisciaSign:       models.SignFromLongitude(aLon),
			AntiSciaDeg:        models.SignDegreeFromLongitude(aLon),
			ContraAntisciaLon:  caLon,
			ContraAntisciaSign: models.SignFromLongitude(caLon),
			ContraAntiSciaDeg:  models.SignDegreeFromLongitude(caLon),
		}
	}
	return points
}

// FindAntisciaPairs finds planets whose antiscia or contra-antiscia are conjunct
// to another planet (within the given orb)
func FindAntisciaPairs(positions []models.PlanetPosition, orb float64) []AntisciaPair {
	if orb <= 0 {
		orb = 2.0
	}

	var pairs []AntisciaPair

	for i := 0; i < len(positions); i++ {
		aLon := CalcAntiscia(positions[i].Longitude)
		caLon := CalcContraAntiscia(positions[i].Longitude)

		for j := i + 1; j < len(positions); j++ {
			targetLon := positions[j].Longitude

			// Check antiscia conjunction
			diff := angleDiff(aLon, targetLon)
			if diff <= orb {
				pairs = append(pairs, AntisciaPair{
					PlanetA: positions[i].PlanetID,
					PlanetB: positions[j].PlanetID,
					Type:    "antiscia",
					Orb:     diff,
				})
			}

			// Check contra-antiscia conjunction
			diff = angleDiff(caLon, targetLon)
			if diff <= orb {
				pairs = append(pairs, AntisciaPair{
					PlanetA: positions[i].PlanetID,
					PlanetB: positions[j].PlanetID,
					Type:    "contra_antiscia",
					Orb:     diff,
				})
			}
		}
	}
	return pairs
}

func angleDiff(a, b float64) float64 {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}
