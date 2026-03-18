package composite

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// DavisonChart holds the Davison relationship chart, which is a real chart
// cast for the time-midpoint and space-midpoint of two birth charts.
type DavisonChart struct {
	MidpointJD  float64                  `json:"midpoint_jd"`
	MidpointLat float64                  `json:"midpoint_lat"`
	MidpointLon float64                  `json:"midpoint_lon"`
	Planets     []models.PlanetPosition  `json:"planets"`
	Houses      []float64                `json:"houses"`
	Angles      models.AnglesInfo        `json:"angles"`
	Aspects     []models.AspectInfo      `json:"aspects"`
}

// CalcDavisonChart computes the Davison relationship chart.
// It calculates the midpoint in time (JD) and space (lat/lon) between two
// birth data sets, then casts a real natal chart for that midpoint moment
// and location.
func CalcDavisonChart(input CompositeInput) (*DavisonChart, error) {
	planets := input.Planets
	if len(planets) == 0 {
		planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}

	// Calculate midpoint in time and space
	midJD := (input.Person1JD + input.Person2JD) / 2
	midLat := (input.Person1Lat + input.Person2Lat) / 2
	midLon := (input.Person1Lon + input.Person2Lon) / 2

	// Cast a real natal chart for the midpoint time and location
	ci, err := chart.CalcSingleChart(midLat, midLon, midJD, planets, input.OrbConfig, input.HouseSystem)
	if err != nil {
		return nil, fmt.Errorf("davison chart calculation: %w", err)
	}

	return &DavisonChart{
		MidpointJD:  midJD,
		MidpointLat: midLat,
		MidpointLon: midLon,
		Planets:     ci.Planets,
		Houses:      ci.Houses,
		Angles:      ci.Angles,
		Aspects:     ci.Aspects,
	}, nil
}
