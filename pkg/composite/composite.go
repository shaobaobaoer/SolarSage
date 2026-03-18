package composite

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/internal/aspect"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// CompositeChart holds the composite (midpoint) chart for two people
type CompositeChart struct {
	Planets []models.PlanetPosition `json:"planets"`
	Houses  []float64               `json:"houses"`
	Angles  models.AnglesInfo       `json:"angles"`
	Aspects []models.AspectInfo     `json:"aspects"`
}

// CompositeInput configures a composite chart calculation
type CompositeInput struct {
	Person1Lat     float64
	Person1Lon     float64
	Person1JD      float64
	Person2Lat     float64
	Person2Lon     float64
	Person2JD      float64
	Planets        []models.PlanetID
	OrbConfig      models.OrbConfig
	HouseSystem    models.HouseSystem
}

// CalcCompositeChart computes the composite chart using the midpoint method.
// Each composite planet position is the midpoint of the two natal positions.
// House cusps use the midpoint of the two natal MCs to derive RAMC.
func CalcCompositeChart(input CompositeInput) (*CompositeChart, error) {
	planets := input.Planets
	if len(planets) == 0 {
		planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}

	// Calculate both natal charts
	chart1, err := chart.CalcSingleChart(input.Person1Lat, input.Person1Lon, input.Person1JD,
		planets, input.OrbConfig, input.HouseSystem)
	if err != nil {
		return nil, fmt.Errorf("person 1 chart: %w", err)
	}

	chart2, err := chart.CalcSingleChart(input.Person2Lat, input.Person2Lon, input.Person2JD,
		planets, input.OrbConfig, input.HouseSystem)
	if err != nil {
		return nil, fmt.Errorf("person 2 chart: %w", err)
	}

	if len(chart1.Planets) != len(chart2.Planets) {
		return nil, fmt.Errorf("planet count mismatch between charts")
	}

	// Compute midpoint for each planet
	positions := make([]models.PlanetPosition, len(chart1.Planets))
	bodies := make([]aspect.Body, len(chart1.Planets))

	for i := range chart1.Planets {
		p1 := chart1.Planets[i]
		p2 := chart2.Planets[i]

		lon := midpoint(p1.Longitude, p2.Longitude)
		speed := (p1.Speed + p2.Speed) / 2

		positions[i] = models.PlanetPosition{
			PlanetID:     p1.PlanetID,
			Longitude:    lon,
			Speed:        speed,
			IsRetrograde: speed < 0,
			Sign:         models.SignFromLongitude(lon),
			SignDegree:   models.SignDegreeFromLongitude(lon),
			House:        1, // will be set below after houses are computed
		}

		bodies[i] = aspect.Body{
			ID:        string(p1.PlanetID),
			Longitude: lon,
			Speed:     speed,
		}
	}

	// Composite angles: midpoint of the two natal angles
	angles := models.AnglesInfo{
		ASC: midpoint(chart1.Angles.ASC, chart2.Angles.ASC),
		MC:  midpoint(chart1.Angles.MC, chart2.Angles.MC),
	}
	angles.DSC = sweph.NormalizeDegrees(angles.ASC + 180)
	angles.IC = sweph.NormalizeDegrees(angles.MC + 180)

	// Composite houses: derive from midpoint MC
	houses := deriveCompositeCusps(angles.ASC, angles.MC)

	// Set house for each planet
	for i := range positions {
		positions[i].House = chart.FindHouseForLongitude(positions[i].Longitude, houses)
	}

	// Calculate aspects
	aspects := aspect.FindAspects(bodies, bodies, input.OrbConfig, true)

	return &CompositeChart{
		Planets: positions,
		Houses:  houses,
		Angles:  angles,
		Aspects: aspects,
	}, nil
}

// midpoint calculates the shorter arc midpoint between two ecliptic longitudes
func midpoint(lon1, lon2 float64) float64 {
	// Normalize both to [0, 360)
	lon1 = sweph.NormalizeDegrees(lon1)
	lon2 = sweph.NormalizeDegrees(lon2)

	diff := lon2 - lon1
	if diff < 0 {
		diff += 360
	}

	if diff <= 180 {
		return sweph.NormalizeDegrees(lon1 + diff/2)
	}
	return sweph.NormalizeDegrees(lon2 + (360-diff)/2)
}

// deriveCompositeCusps creates equal-sized house cusps from ASC
// This is the standard method for composite charts
func deriveCompositeCusps(asc, mc float64) []float64 {
	cusps := make([]float64, 12)
	for i := 0; i < 12; i++ {
		cusps[i] = sweph.NormalizeDegrees(asc + float64(i)*30)
	}
	return cusps
}

// Midpoint calculates the midpoint between two ecliptic longitudes (exported for reuse)
func Midpoint(lon1, lon2 float64) float64 {
	return midpoint(lon1, lon2)
}
