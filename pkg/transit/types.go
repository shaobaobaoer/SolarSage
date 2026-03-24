package transit

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// TransitChartConfig holds configuration for transit chart calculations.
type TransitChartConfig struct {
	Lat, Lon    float64               // Transit location (affects ASC/MC special points)
	Planets     []models.PlanetID     // Transit planets to calculate
	Points      []models.SpecialPointID // Transit special points to calculate
	Orbs        models.OrbConfig      // Orb configuration for this chart
	HouseSystem models.HouseSystem    // House system (shared across all charts)
}

// ProgressionsChartConfig holds configuration for secondary progressions chart calculations.
type ProgressionsChartConfig struct {
	Planets     []models.PlanetID     // Progressed planets to calculate
	Points      []models.SpecialPointID // Progressed special points to calculate
	Orbs        models.OrbConfig      // Orb configuration for this chart
	Lat, Lon    float64               // Location for progressed special points
	HouseSystem models.HouseSystem    // House system (shared across all charts)
}

// SolarArcChartConfig holds configuration for solar arc chart calculations.
type SolarArcChartConfig struct {
	Planets     []models.PlanetID     // Solar arc planets to calculate
	Points      []models.SpecialPointID // Solar arc special points to calculate
	Orbs        models.OrbConfig      // Orb configuration for this chart
	Lat, Lon    float64               // Location for solar arc special points
	HouseSystem models.HouseSystem    // House system (shared across all charts)
}

// ChartSetConfig holds all chart configurations for transit calculation.
type ChartSetConfig struct {
	Transit      *TransitChartConfig      // nil means disabled
	Progressions *ProgressionsChartConfig // nil means disabled
	SolarArc     *SolarArcChartConfig     // nil means disabled
}

// NatalChartConfig holds configuration for the natal chart (fixed reference).
type NatalChartConfig struct {
	Lat, Lon float64               // Birth location
	JD       float64               // Birth moment
	Planets  []models.PlanetID     // Natal planets to include
	Points   []models.SpecialPointID // Natal special points to include
}

// TimeRangeConfig holds the time range for transit calculation.
type TimeRangeConfig struct {
	StartJD float64
	EndJD   float64
}

// EventFilterConfig holds flags for which event types to include.
type EventFilterConfig struct {
	// Event types
	Station      bool
	SignIngress  bool
	HouseIngress bool
	VoidOfCourse bool

	// Aspect combinations
	TrNa bool // Transit → Natal
	TrTr bool // Transit → Transit
	TrSp bool // Transit → Progressions
	TrSa bool // Transit → SolarArc
	SpNa bool // Progressions → Natal
	SpSp bool // Progressions → Progressions
	SaNa bool // SolarArc → Natal
}

// DefaultEventFilterConfig returns a config with all events enabled.
func DefaultEventFilterConfig() EventFilterConfig {
	return EventFilterConfig{
		Station:      true,
		SignIngress:  true,
		HouseIngress: true,
		VoidOfCourse: true,
		TrNa:         true,
		TrTr:         true,
		TrSp:         true,
		TrSa:         true,
		SpNa:         true,
		SpSp:         true,
		SaNa:         true,
	}
}
