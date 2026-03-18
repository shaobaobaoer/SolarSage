package returns

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// bisectEps is the precision for bisection (~1 second)
const bisectEps = 1.0 / 86400.0

// ReturnChart holds the return chart data along with return-specific metadata
type ReturnChart struct {
	ReturnJD    float64           `json:"return_jd"`
	ReturnType  string            `json:"return_type"` // "solar" or "lunar"
	NatalJD     float64           `json:"natal_jd"`
	Age         float64           `json:"age"`
	Chart       *models.ChartInfo `json:"chart"`
	PlanetLon   float64           `json:"planet_longitude"`
	NatalLon    float64           `json:"natal_longitude"`
	IsRetrograde bool             `json:"is_retrograde,omitempty"`
}

// ReturnInput configures a return chart calculation
type ReturnInput struct {
	NatalJD     float64
	NatalLat    float64
	NatalLon    float64
	SearchJD    float64 // Start searching from this JD
	Planets     []models.PlanetID
	OrbConfig   models.OrbConfig
	HouseSystem models.HouseSystem
}

// CalcSolarReturn finds the exact moment the Sun returns to its natal longitude
// after SearchJD, and calculates a full chart for that moment.
func CalcSolarReturn(input ReturnInput) (*ReturnChart, error) {
	return calcReturn(input, models.PlanetSun, "solar")
}

// CalcLunarReturn finds the exact moment the Moon returns to its natal longitude
// after SearchJD, and calculates a full chart for that moment.
func CalcLunarReturn(input ReturnInput) (*ReturnChart, error) {
	return calcReturn(input, models.PlanetMoon, "lunar")
}

// CalcPlanetReturn finds the exact moment a planet returns to its natal longitude.
// Works for any non-retrograding period. For outer planets, handles retrograde passes.
func CalcPlanetReturn(input ReturnInput, planet models.PlanetID) (*ReturnChart, error) {
	return calcReturn(input, planet, "planetary")
}

// CalcSolarReturnSeries calculates multiple consecutive solar returns
func CalcSolarReturnSeries(input ReturnInput, count int) ([]*ReturnChart, error) {
	return calcReturnSeries(input, models.PlanetSun, "solar", count)
}

// CalcLunarReturnSeries calculates multiple consecutive lunar returns
func CalcLunarReturnSeries(input ReturnInput, count int) ([]*ReturnChart, error) {
	return calcReturnSeries(input, models.PlanetMoon, "lunar", count)
}

func calcReturnSeries(input ReturnInput, planet models.PlanetID, returnType string, count int) ([]*ReturnChart, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	results := make([]*ReturnChart, 0, count)
	searchJD := input.SearchJD

	for i := 0; i < count; i++ {
		inp := input
		inp.SearchJD = searchJD
		rc, err := calcReturn(inp, planet, returnType)
		if err != nil {
			return results, fmt.Errorf("return #%d: %w", i+1, err)
		}
		results = append(results, rc)
		// Advance past this return to find the next one
		if planet == models.PlanetMoon {
			searchJD = rc.ReturnJD + 25 // Moon cycle ~27.3 days
		} else if planet == models.PlanetSun {
			searchJD = rc.ReturnJD + 360 // Sun cycle ~365.25 days
		} else {
			searchJD = rc.ReturnJD + 30 // generic advance
		}
	}
	return results, nil
}

func calcReturn(input ReturnInput, planet models.PlanetID, returnType string) (*ReturnChart, error) {
	// Get natal planet longitude
	natalLon, _, err := chart.CalcPlanetLongitude(planet, input.NatalJD)
	if err != nil {
		return nil, fmt.Errorf("natal %s position: %w", planet, err)
	}

	// Find the return JD via scanning + bisection
	returnJD, err := findReturnJD(planet, natalLon, input.SearchJD)
	if err != nil {
		return nil, fmt.Errorf("find %s return: %w", returnType, err)
	}

	// Get planet info at return moment
	returnLon, speed, err := chart.CalcPlanetLongitude(planet, returnJD)
	if err != nil {
		return nil, err
	}

	// Calculate full chart at return moment
	planets := input.Planets
	if len(planets) == 0 {
		planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}

	chartInfo, err := chart.CalcSingleChart(
		input.NatalLat, input.NatalLon, returnJD,
		planets, input.OrbConfig, input.HouseSystem,
	)
	if err != nil {
		return nil, fmt.Errorf("return chart: %w", err)
	}

	age := (returnJD - input.NatalJD) / 365.25

	return &ReturnChart{
		ReturnJD:     returnJD,
		ReturnType:   returnType,
		NatalJD:      input.NatalJD,
		Age:          age,
		Chart:        chartInfo,
		PlanetLon:    returnLon,
		NatalLon:     natalLon,
		IsRetrograde: speed < 0,
	}, nil
}

// findReturnJD locates the exact JD when a planet returns to targetLon after startJD.
// Uses adaptive scanning followed by bisection for 1-second precision.
func findReturnJD(planet models.PlanetID, targetLon, startJD float64) (float64, error) {
	// Determine scan step based on planet speed
	step := scanStep(planet)
	maxDays := maxScanDays(planet)

	prevLon, _, err := chart.CalcPlanetLongitude(planet, startJD)
	if err != nil {
		return 0, err
	}
	prevDiff := normDiff(prevLon, targetLon)

	for d := step; d <= maxDays; d += step {
		jd := startJD + d
		lon, _, err := chart.CalcPlanetLongitude(planet, jd)
		if err != nil {
			continue
		}
		currDiff := normDiff(lon, targetLon)

		// Detect sign change in the difference (crossing the target longitude)
		if prevDiff*currDiff < 0 && math.Abs(prevDiff-currDiff) < 180 {
			// Bisect to find exact crossing
			lo, hi := jd-step, jd
			for hi-lo > bisectEps {
				mid := (lo + hi) / 2
				midLon, _, err := chart.CalcPlanetLongitude(planet, mid)
				if err != nil {
					break
				}
				midDiff := normDiff(midLon, targetLon)
				if prevDiff*midDiff < 0 {
					hi = mid
				} else {
					lo = mid
					prevDiff = midDiff
				}
			}
			return (lo + hi) / 2, nil
		}
		prevDiff = currDiff
	}

	return 0, fmt.Errorf("return not found within %.0f days", maxDays)
}

// normDiff returns a signed angular difference normalized to [-180, 180)
func normDiff(lon, target float64) float64 {
	d := lon - target
	d = math.Mod(d+180, 360)
	if d < 0 {
		d += 360
	}
	return d - 180
}

// scanStep returns the daily step size for scanning, based on planet
func scanStep(planet models.PlanetID) float64 {
	switch planet {
	case models.PlanetMoon:
		return 0.5 // Moon moves ~13°/day
	case models.PlanetSun, models.PlanetMercury, models.PlanetVenus:
		return 1.0
	case models.PlanetMars:
		return 2.0
	default:
		return 5.0 // slow outer planets
	}
}

// maxScanDays returns the maximum number of days to scan for a return
func maxScanDays(planet models.PlanetID) float64 {
	switch planet {
	case models.PlanetMoon:
		return 32 // ~27.3 day cycle
	case models.PlanetSun:
		return 370 // ~365.25 day cycle
	case models.PlanetMercury:
		return 400 // ~88 day orbit, but synodic ~116
	case models.PlanetVenus:
		return 600 // ~225 day orbit, synodic ~584
	case models.PlanetMars:
		return 800 // ~687 day orbit
	case models.PlanetJupiter:
		return 4400 // ~11.86 year orbit
	case models.PlanetSaturn:
		return 10800 // ~29.46 year orbit
	default:
		return 31000 // ~84 years for Uranus, etc.
	}
}
