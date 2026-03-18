package symbolic

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/internal/aspect"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// DirectionMethod identifies the symbolic direction formula.
type DirectionMethod string

const (
	// MethodOneDegree advances each position by 1.0 per year.
	MethodOneDegree DirectionMethod = "ONE_DEGREE"
	// MethodNaibod uses the Sun's mean daily motion (0.98556/year).
	MethodNaibod DirectionMethod = "NAIBOD"
	// MethodProfection advances each position by 30 per year.
	MethodProfection DirectionMethod = "PROFECTION"
	// MethodCustom uses the caller-supplied CustomRate.
	MethodCustom DirectionMethod = "CUSTOM"
)

// Rate returns the degrees-per-year for the method.
// For MethodCustom it returns the default 1.0; callers should use
// SymbolicInput.CustomRate instead.
func (m DirectionMethod) Rate() float64 {
	switch m {
	case MethodOneDegree:
		return 1.0
	case MethodNaibod:
		return 0.98556
	case MethodProfection:
		return 30.0
	default:
		return 1.0
	}
}

// SymbolicDirection holds the directed position for a single planet.
type SymbolicDirection struct {
	PlanetID    models.PlanetID `json:"planet_id"`
	NatalLon    float64         `json:"natal_longitude"`
	DirectedLon float64         `json:"directed_longitude"`
	DirectedSign string         `json:"directed_sign"`
	DirectedDeg float64         `json:"directed_sign_degree"`
	ArcApplied  float64         `json:"arc_applied"`
}

// DirectedAngle holds the directed position for an angle (ASC/MC/DSC/IC).
type DirectedAngle struct {
	PointID     string  `json:"point_id"`
	NatalLon    float64 `json:"natal_longitude"`
	DirectedLon float64 `json:"directed_longitude"`
	DirectedSign string `json:"directed_sign"`
	DirectedDeg float64 `json:"directed_sign_degree"`
	ArcApplied  float64 `json:"arc_applied"`
}

// SymbolicDirectionResult is the complete output of CalcSymbolicDirections.
type SymbolicDirectionResult struct {
	Method     DirectionMethod     `json:"method"`
	Rate       float64             `json:"rate_per_year"`
	Age        float64             `json:"age"`
	Directions []SymbolicDirection `json:"directions"`
	Angles     []DirectedAngle     `json:"angles"`
	Aspects    []models.AspectInfo `json:"aspects"`
}

// SymbolicInput configures a symbolic directions calculation.
type SymbolicInput struct {
	NatalJD     float64
	GeoLat      float64
	GeoLon      float64
	Age         float64 // age in years
	Method      DirectionMethod
	CustomRate  float64 // only used when Method == MethodCustom
	Planets     []models.PlanetID
	OrbConfig   models.OrbConfig
	HouseSystem models.HouseSystem
}

// defaultPlanets is the standard 10-planet set.
var defaultPlanets = []models.PlanetID{
	models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
	models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
	models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
	models.PlanetPluto,
}

// CalcSymbolicDirections computes symbolic directions for the given input.
//
// It calculates the natal chart, applies the arc (age x rate) to every planet
// and angle, and finds aspects between directed and natal positions.
func CalcSymbolicDirections(input SymbolicInput) (*SymbolicDirectionResult, error) {
	// Determine rate
	rate := input.Method.Rate()
	if input.Method == MethodCustom {
		if input.CustomRate <= 0 {
			return nil, fmt.Errorf("symbolic: custom rate must be positive, got %f", input.CustomRate)
		}
		rate = input.CustomRate
	}

	if input.Age < 0 {
		return nil, fmt.Errorf("symbolic: age must be non-negative, got %f", input.Age)
	}

	// Use default planets when none specified
	planets := input.Planets
	if len(planets) == 0 {
		planets = defaultPlanets
	}

	// Use default orbs if all zeros
	orbs := input.OrbConfig
	if orbs.Conjunction == 0 && orbs.Opposition == 0 && orbs.Trine == 0 {
		orbs = models.DefaultOrbConfig()
	}

	hsys := input.HouseSystem
	if hsys == "" {
		hsys = models.HousePlacidus
	}

	// Calculate natal chart
	natalChart, err := chart.CalcSingleChart(input.GeoLat, input.GeoLon, input.NatalJD, planets, orbs, hsys)
	if err != nil {
		return nil, fmt.Errorf("symbolic: natal chart: %w", err)
	}

	arc := input.Age * rate

	// Build directed planet positions and body lists for aspect calculation
	var directions []SymbolicDirection
	var natalBodies []aspect.Body
	var directedBodies []aspect.Body

	for _, p := range natalChart.Planets {
		dirLon := sweph.NormalizeDegrees(p.Longitude + arc)
		directions = append(directions, SymbolicDirection{
			PlanetID:     p.PlanetID,
			NatalLon:     p.Longitude,
			DirectedLon:  dirLon,
			DirectedSign: models.SignFromLongitude(dirLon),
			DirectedDeg:  models.SignDegreeFromLongitude(dirLon),
			ArcApplied:   arc,
		})
		natalBodies = append(natalBodies, aspect.Body{
			ID:        string(p.PlanetID),
			Longitude: p.Longitude,
			Speed:     p.Speed,
		})
		directedBodies = append(directedBodies, aspect.Body{
			ID:        "Dir_" + string(p.PlanetID),
			Longitude: dirLon,
			Speed:     0, // symbolic positions have no intrinsic speed
		})
	}

	// Directed angles
	var angles []DirectedAngle
	for _, pt := range []struct {
		id  string
		lon float64
	}{
		{"ASC", natalChart.Angles.ASC},
		{"MC", natalChart.Angles.MC},
		{"DSC", natalChart.Angles.DSC},
		{"IC", natalChart.Angles.IC},
	} {
		dirLon := sweph.NormalizeDegrees(pt.lon + arc)
		angles = append(angles, DirectedAngle{
			PointID:      pt.id,
			NatalLon:     pt.lon,
			DirectedLon:  dirLon,
			DirectedSign: models.SignFromLongitude(dirLon),
			DirectedDeg:  models.SignDegreeFromLongitude(dirLon),
			ArcApplied:   arc,
		})
		// Add natal angle to natal bodies for aspect detection
		natalBodies = append(natalBodies, aspect.Body{
			ID:        pt.id,
			Longitude: pt.lon,
			Speed:     0,
		})
	}

	// Find aspects between directed planets and natal planets+angles
	// sameSet=false because directed and natal are distinct sets
	aspects := aspect.FindAspects(directedBodies, natalBodies, orbs, false)

	return &SymbolicDirectionResult{
		Method:     input.Method,
		Rate:       rate,
		Age:        input.Age,
		Directions: directions,
		Angles:     angles,
		Aspects:    aspects,
	}, nil
}
