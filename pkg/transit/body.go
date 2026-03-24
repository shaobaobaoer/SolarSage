package transit

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// MovingBody represents a celestial body that moves through time and is used in transit calculations.
// It unifies transit planets, progressed planets, solar arc planets, and special points.
type MovingBody struct {
	ID            string         // Unique identifier (planet ID or special point ID)
	ChartType     models.ChartType // Chart type (TRANSIT, PROGRESSIONS, SOLAR_ARC)
	CalcFn        bodyCalcFunc   // Calculation function returning (longitude, speed, error)
	Orbs          models.OrbConfig // Orb configuration for aspect detection
	CanRetrograde bool           // Whether this body can go retrograde (controls station detection)
}

// buildTransitBodies creates moving bodies for transit chart.
// Returns nil if transit config is nil.
func buildTransitBodies(cfg *TransitChartConfig) []MovingBody {
	if cfg == nil {
		return nil
	}

	var bodies []MovingBody

	// Add transit planets
	for _, planet := range cfg.Planets {
		calcFn := makeTransitCalcFn(planet)
		bodies = append(bodies, MovingBody{
			ID:            string(planet),
			ChartType:     models.ChartTransit,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: canRetrograde(planet),
		})
	}

	// Add transit special points
	for _, sp := range cfg.Points {
		calcFn := makeTransitSpecialPointCalcFn(sp, cfg.Lat, cfg.Lon, cfg.HouseSystem)
		bodies = append(bodies, MovingBody{
			ID:            string(sp),
			ChartType:     models.ChartTransit,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: false, // Special points don't retrograde
		})
	}

	return bodies
}

// buildProgressionBodies creates moving bodies for secondary progressions chart.
// Returns nil if progressions config is nil.
func buildProgressionBodies(cfg *ProgressionsChartConfig, natalJD float64) []MovingBody {
	if cfg == nil {
		return nil
	}

	var bodies []MovingBody

	// Add progressed planets
	for _, planet := range cfg.Planets {
		calcFn := makeProgressionsCalcFn(planet, natalJD)
		bodies = append(bodies, MovingBody{
			ID:            string(planet),
			ChartType:     models.ChartProgressions,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: true, // Progressed planets can retrograde
		})
	}

	// Add progressed special points
	for _, sp := range cfg.Points {
		calcFn := makeProgressionsSpecialPointCalcFn(sp, cfg.Lat, cfg.Lon, natalJD, cfg.HouseSystem)
		bodies = append(bodies, MovingBody{
			ID:            string(sp),
			ChartType:     models.ChartProgressions,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: false,
		})
	}

	return bodies
}

// buildSolarArcBodies creates moving bodies for solar arc chart.
// Returns nil if solar arc config is nil.
func buildSolarArcBodies(cfg *SolarArcChartConfig, natalJD float64) []MovingBody {
	if cfg == nil {
		return nil
	}

	var bodies []MovingBody

	// Add solar arc planets
	for _, planet := range cfg.Planets {
		calcFn := makeSolarArcCalcFn(planet, natalJD)
		bodies = append(bodies, MovingBody{
			ID:            string(planet),
			ChartType:     models.ChartSolarArc,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: false, // Solar arc planets don't retrograde
		})
	}

	// Add solar arc special points
	for _, sp := range cfg.Points {
		calcFn := makeSolarArcSpecialPointCalcFn(sp, cfg.Lat, cfg.Lon, natalJD, cfg.HouseSystem)
		bodies = append(bodies, MovingBody{
			ID:            string(sp),
			ChartType:     models.ChartSolarArc,
			CalcFn:        calcFn,
			Orbs:          cfg.Orbs,
			CanRetrograde: false,
		})
	}

	return bodies
}

// canRetrograde returns whether a planet can go retrograde.
// Sun and Moon never retrograde.
func canRetrograde(planet models.PlanetID) bool {
	return planet != models.PlanetSun && planet != models.PlanetMoon
}

// Helper functions for creating calc functions (moved from transit.go)

func makeTransitCalcFn(planet models.PlanetID) bodyCalcFunc {
	return func(jd float64) (float64, float64, error) {
		return chart.CalcPlanetLongitude(planet, jd)
	}
}

func makeProgressionsCalcFn(planet models.PlanetID, natalJD float64) bodyCalcFunc {
	return func(jd float64) (float64, float64, error) {
		return progressions.CalcProgressedLongitude(planet, natalJD, jd)
	}
}

func makeSolarArcCalcFn(planet models.PlanetID, natalJD float64) bodyCalcFunc {
	return func(jd float64) (float64, float64, error) {
		return progressions.CalcSolarArcLongitude(planet, natalJD, jd)
	}
}

func makeTransitSpecialPointCalcFn(sp models.SpecialPointID, lat, lon float64, hsys models.HouseSystem) bodyCalcFunc {
	lonFn := func(jd float64) (float64, error) {
		return chart.CalcSpecialPointLongitude(sp, lat, lon, jd, hsys)
	}
	return func(jd float64) (float64, float64, error) {
		spLon, err := lonFn(jd)
		if err != nil {
			return 0, 0, err
		}
		return spLon, numericalSpeed(lonFn, jd, 0.01), nil
	}
}

func makeProgressionsSpecialPointCalcFn(sp models.SpecialPointID, lat, lon float64, natalJD float64, hsys models.HouseSystem) bodyCalcFunc {
	lonFn := func(jd float64) (float64, error) {
		return progressions.CalcProgressedSpecialPoint(sp, natalJD, jd, lat, lon, hsys)
	}
	return func(jd float64) (float64, float64, error) {
		spLon, err := lonFn(jd)
		if err != nil {
			return 0, 0, err
		}
		return spLon, numericalSpeed(lonFn, jd, 1.0), nil
	}
}

func makeSolarArcSpecialPointCalcFn(sp models.SpecialPointID, lat, lon float64, natalJD float64, hsys models.HouseSystem) bodyCalcFunc {
	// Pre-compute natal special point longitude (fixed)
	natalSpLon, _ := chart.CalcSpecialPointLongitude(sp, lat, lon, natalJD, hsys)
	return func(jd float64) (float64, float64, error) {
		offset, err := progressions.SolarArcOffset(natalJD, jd)
		if err != nil {
			return 0, 0, err
		}
		directed := sweph.NormalizeDegrees(natalSpLon + offset)
		// Speed ~ sun's progressed speed / JulianYear
		pJD := progressions.SecondaryProgressionJD(natalJD, jd)
		_, sunSpeed, _ := chart.CalcPlanetLongitude(models.PlanetSun, pJD)
		speed := sunSpeed / progressions.JulianYear
		return directed, speed, nil
	}
}
