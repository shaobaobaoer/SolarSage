package lunar

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// bisectEps is the precision for bisection (~1 second)
const bisectEps = 1.0 / 86400.0

// Phase represents a lunar phase
type Phase string

const (
	PhaseNewMoon        Phase = "NEW_MOON"
	PhaseWaxingCrescent Phase = "WAXING_CRESCENT"
	PhaseFirstQuarter   Phase = "FIRST_QUARTER"
	PhaseWaxingGibbous  Phase = "WAXING_GIBBOUS"
	PhaseFullMoon       Phase = "FULL_MOON"
	PhaseWaningGibbous  Phase = "WANING_GIBBOUS"
	PhaseLastQuarter    Phase = "LAST_QUARTER"
	PhaseWaningCrescent Phase = "WANING_CRESCENT"
)

// PhaseInfo holds lunar phase data at a specific time
type PhaseInfo struct {
	Phase       Phase   `json:"phase"`
	PhaseName   string  `json:"phase_name"`
	PhaseAngle  float64 `json:"phase_angle"`  // Sun-Moon elongation (0-360)
	Illumination float64 `json:"illumination"` // 0.0 - 1.0
	MoonLon     float64 `json:"moon_longitude"`
	SunLon      float64 `json:"sun_longitude"`
	IsWaxing    bool    `json:"is_waxing"`
}

// LunarEvent represents a major lunar phase event (new/full/quarter)
type LunarEvent struct {
	Phase    Phase   `json:"phase"`
	JD       float64 `json:"jd"`
	MoonLon  float64 `json:"moon_longitude"`
	MoonSign string  `json:"moon_sign"`
	SunLon   float64 `json:"sun_longitude"`
	SunSign  string  `json:"sun_sign"`
}

// EclipseType represents the type of eclipse
type EclipseType string

const (
	EclipseSolarTotal   EclipseType = "SOLAR_TOTAL"
	EclipseSolarAnnular EclipseType = "SOLAR_ANNULAR"
	EclipseSolarPartial EclipseType = "SOLAR_PARTIAL"
	EclipseLunarTotal   EclipseType = "LUNAR_TOTAL"
	EclipseLunarPartial EclipseType = "LUNAR_PARTIAL"
	EclipseLunarPenumbral EclipseType = "LUNAR_PENUMBRAL"
)

// EclipseInfo holds eclipse data
type EclipseInfo struct {
	Type     EclipseType `json:"type"`
	JD       float64     `json:"jd"`
	MoonLon  float64     `json:"moon_longitude"`
	MoonSign string      `json:"moon_sign"`
	SunLon   float64     `json:"sun_longitude"`
	SunSign  string      `json:"sun_sign"`
	MoonLat  float64     `json:"moon_latitude"`
	Gamma    float64     `json:"gamma,omitempty"` // approximate eclipse magnitude indicator
}

// CalcLunarPhase returns the lunar phase at a given Julian Day
func CalcLunarPhase(jdUT float64) (*PhaseInfo, error) {
	moonLon, _, err := chart.CalcPlanetLongitude(models.PlanetMoon, jdUT)
	if err != nil {
		return nil, err
	}
	sunLon, _, err := chart.CalcPlanetLongitude(models.PlanetSun, jdUT)
	if err != nil {
		return nil, err
	}

	elongation := sweph.NormalizeDegrees(moonLon - sunLon)
	illumination := (1 - math.Cos(elongation*math.Pi/180)) / 2

	phase, name := phaseFromElongation(elongation)

	return &PhaseInfo{
		Phase:        phase,
		PhaseName:    name,
		PhaseAngle:   elongation,
		Illumination: illumination,
		MoonLon:      moonLon,
		SunLon:       sunLon,
		IsWaxing:     elongation < 180,
	}, nil
}

// FindLunarPhases finds all major lunar phase events in a date range
func FindLunarPhases(startJD, endJD float64) ([]LunarEvent, error) {
	var events []LunarEvent

	// Scan for new moons (elongation crosses 0°) and full moons (crosses 180°)
	// Also first quarter (~90°) and last quarter (~270°)
	targets := []struct {
		angle float64
		phase Phase
	}{
		{0, PhaseNewMoon},
		{90, PhaseFirstQuarter},
		{180, PhaseFullMoon},
		{270, PhaseLastQuarter},
	}

	step := 1.0 // 1 day steps
	jd := startJD

	prevElong, err := elongation(jd)
	if err != nil {
		return nil, err
	}

	for jd+step <= endJD+step {
		jd += step
		currElong, err := elongation(jd)
		if err != nil {
			continue
		}

		for _, t := range targets {
			if crossesAngle(prevElong, currElong, t.angle) {
				exactJD, err := bisectPhase(jd-step, jd, t.angle)
				if err != nil {
					continue
				}
				if exactJD < startJD || exactJD > endJD {
					continue
				}

				moonLon, _, _ := chart.CalcPlanetLongitude(models.PlanetMoon, exactJD)
				sunLon, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, exactJD)

				events = append(events, LunarEvent{
					Phase:    t.phase,
					JD:       exactJD,
					MoonLon:  moonLon,
					MoonSign: models.SignFromLongitude(moonLon),
					SunLon:   sunLon,
					SunSign:  models.SignFromLongitude(sunLon),
				})
			}
		}

		prevElong = currElong
	}

	return events, nil
}

// FindEclipses finds solar and lunar eclipses in a date range.
// Solar eclipses occur at new moons when Moon is near the nodes.
// Lunar eclipses occur at full moons when Moon is near the nodes.
func FindEclipses(startJD, endJD float64) ([]EclipseInfo, error) {
	phases, err := FindLunarPhases(startJD, endJD)
	if err != nil {
		return nil, err
	}

	var eclipses []EclipseInfo

	for _, p := range phases {
		if p.Phase != PhaseNewMoon && p.Phase != PhaseFullMoon {
			continue
		}

		// Get Moon's ecliptic latitude at this moment
		moonResult, err := sweph.CalcUT(p.JD, sweph.SE_MOON)
		if err != nil {
			continue
		}
		moonLat := moonResult.Latitude

		// Get node position for proximity check
		nodeLon, _, _ := chart.CalcPlanetLongitude(models.PlanetNorthNodeTrue, p.JD)
		moonNodeDist := angleDiff(p.MoonLon, nodeLon)
		southNodeDist := angleDiff(p.MoonLon, sweph.NormalizeDegrees(nodeLon+180))
		nearNode := math.Min(moonNodeDist, southNodeDist)

		absLat := math.Abs(moonLat)

		if p.Phase == PhaseNewMoon {
			// Solar eclipse conditions: Moon near node (within ~18°) and small latitude
			if nearNode <= 18.5 && absLat < 1.6 {
				eclType := classifySolarEclipse(absLat)
				eclipses = append(eclipses, EclipseInfo{
					Type:     eclType,
					JD:       p.JD,
					MoonLon:  p.MoonLon,
					MoonSign: p.MoonSign,
					SunLon:   p.SunLon,
					SunSign:  p.SunSign,
					MoonLat:  moonLat,
					Gamma:    absLat,
				})
			}
		} else { // Full Moon
			// Lunar eclipse conditions: Moon near node (within ~12°) and small latitude
			if nearNode <= 12.5 && absLat < 1.1 {
				eclType := classifyLunarEclipse(absLat)
				eclipses = append(eclipses, EclipseInfo{
					Type:     eclType,
					JD:       p.JD,
					MoonLon:  p.MoonLon,
					MoonSign: p.MoonSign,
					SunLon:   p.SunLon,
					SunSign:  p.SunSign,
					MoonLat:  moonLat,
					Gamma:    absLat,
				})
			}
		}
	}

	return eclipses, nil
}

func elongation(jd float64) (float64, error) {
	moonLon, _, err := chart.CalcPlanetLongitude(models.PlanetMoon, jd)
	if err != nil {
		return 0, err
	}
	sunLon, _, err := chart.CalcPlanetLongitude(models.PlanetSun, jd)
	if err != nil {
		return 0, err
	}
	return sweph.NormalizeDegrees(moonLon - sunLon), nil
}

func crossesAngle(prev, curr, target float64) bool {
	// Normalize relative to target
	p := sweph.NormalizeDegrees(prev - target)
	c := sweph.NormalizeDegrees(curr - target)
	// Crossing happens when prev is just before 0° (in 350-360 range)
	// and curr is just after 0° (in 0-10 range), meaning elongation crossed the target
	// We only detect forward crossing (increasing elongation passes target)
	return p > 180 && c <= 180 && (360-p+c) < 20
}

func bisectPhase(lo, hi, targetAngle float64) (float64, error) {
	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		midElong, err := elongation(mid)
		if err != nil {
			return 0, err
		}
		diff := sweph.NormalizeDegrees(midElong - targetAngle)
		if diff > 180 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2, nil
}

func phaseFromElongation(e float64) (Phase, string) {
	switch {
	case e < 22.5 || e >= 337.5:
		return PhaseNewMoon, "New Moon"
	case e < 67.5:
		return PhaseWaxingCrescent, "Waxing Crescent"
	case e < 112.5:
		return PhaseFirstQuarter, "First Quarter"
	case e < 157.5:
		return PhaseWaxingGibbous, "Waxing Gibbous"
	case e < 202.5:
		return PhaseFullMoon, "Full Moon"
	case e < 247.5:
		return PhaseWaningGibbous, "Waning Gibbous"
	case e < 292.5:
		return PhaseLastQuarter, "Last Quarter"
	default:
		return PhaseWaningCrescent, "Waning Crescent"
	}
}

func classifySolarEclipse(absLat float64) EclipseType {
	if absLat < 0.4 {
		return EclipseSolarTotal
	}
	if absLat < 0.8 {
		return EclipseSolarAnnular
	}
	return EclipseSolarPartial
}

func classifyLunarEclipse(absLat float64) EclipseType {
	if absLat < 0.3 {
		return EclipseLunarTotal
	}
	if absLat < 0.7 {
		return EclipseLunarPartial
	}
	return EclipseLunarPenumbral
}

func angleDiff(a, b float64) float64 {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// NextNewMoon finds the next new moon after the given JD
func NextNewMoon(jdUT float64) (float64, error) {
	phases, err := FindLunarPhases(jdUT, jdUT+35)
	if err != nil {
		return 0, err
	}
	for _, p := range phases {
		if p.Phase == PhaseNewMoon && p.JD > jdUT {
			return p.JD, nil
		}
	}
	return 0, fmt.Errorf("new moon not found within 35 days")
}

// NextFullMoon finds the next full moon after the given JD
func NextFullMoon(jdUT float64) (float64, error) {
	phases, err := FindLunarPhases(jdUT, jdUT+35)
	if err != nil {
		return 0, err
	}
	for _, p := range phases {
		if p.Phase == PhaseFullMoon && p.JD > jdUT {
			return p.JD, nil
		}
	}
	return 0, fmt.Errorf("full moon not found within 35 days")
}
