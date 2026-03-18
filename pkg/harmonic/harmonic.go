package harmonic

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/internal/aspect"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// HarmonicChart holds a harmonic (divisional) chart
type HarmonicChart struct {
	Harmonic int                     `json:"harmonic"`
	Planets  []models.PlanetPosition `json:"planets"`
	Aspects  []models.AspectInfo     `json:"aspects"`
}

// CalcHarmonicChart computes the Nth harmonic chart.
// Each planet's longitude is multiplied by N, then taken mod 360.
// The harmonic chart reveals N-fold symmetry in the natal chart.
func CalcHarmonicChart(lat, lon, jdUT float64, harmonic int, planets []models.PlanetID, orbs models.OrbConfig, hsys models.HouseSystem) (*HarmonicChart, error) {
	if harmonic < 1 || harmonic > 180 {
		return nil, fmt.Errorf("harmonic must be between 1 and 180, got %d", harmonic)
	}

	if len(planets) == 0 {
		planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}

	// Get natal positions first
	natalChart, err := chart.CalcSingleChart(lat, lon, jdUT, planets, orbs, hsys)
	if err != nil {
		return nil, fmt.Errorf("natal chart: %w", err)
	}

	// Multiply each longitude by the harmonic number
	positions := make([]models.PlanetPosition, len(natalChart.Planets))
	bodies := make([]aspect.Body, len(natalChart.Planets))

	for i, p := range natalChart.Planets {
		hLon := sweph.NormalizeDegrees(p.Longitude * float64(harmonic))
		positions[i] = models.PlanetPosition{
			PlanetID:     p.PlanetID,
			Longitude:    hLon,
			Latitude:     p.Latitude,
			Speed:        p.Speed * float64(harmonic),
			IsRetrograde: p.IsRetrograde,
			Sign:         models.SignFromLongitude(hLon),
			SignDegree:   models.SignDegreeFromLongitude(hLon),
			House:        p.House, // house from natal chart
		}
		bodies[i] = aspect.Body{
			ID:        string(p.PlanetID),
			Longitude: hLon,
			Speed:     p.Speed * float64(harmonic),
		}
	}

	// Calculate aspects in the harmonic chart
	aspects := aspect.FindAspects(bodies, bodies, orbs, true)

	return &HarmonicChart{
		Harmonic: harmonic,
		Planets:  positions,
		Aspects:  aspects,
	}, nil
}

// CommonHarmonics returns standard harmonic numbers used in astrology
func CommonHarmonics() map[int]string {
	return map[int]string{
		1:  "Natal (base chart)",
		2:  "Opposition (polarity)",
		3:  "Trine (creativity)",
		4:  "Square (tension, challenge)",
		5:  "Quintile (talent, creative power)",
		7:  "Septile (inspiration, fate)",
		8:  "Semi-square (friction, drive)",
		9:  "Novile (completion, initiation)",
		12: "Semi-sextile/quincunx (adjustment)",
		16: "16th harmonic (refinement)",
	}
}
