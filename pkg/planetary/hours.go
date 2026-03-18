package planetary

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// PlanetaryHour represents a single planetary hour
type PlanetaryHour struct {
	Number    int             `json:"number"`     // 1-24
	Planet    models.PlanetID `json:"planet"`
	StartJD   float64         `json:"start_jd"`
	EndJD     float64         `json:"end_jd"`
	IsDay     bool            `json:"is_day"`
}

// PlanetaryDay holds the planetary day ruler and all 24 hours
type PlanetaryDay struct {
	DayRuler    models.PlanetID `json:"day_ruler"`
	Sunrise     float64         `json:"sunrise_jd"`
	Sunset      float64         `json:"sunset_jd"`
	NextSunrise float64         `json:"next_sunrise_jd"`
	Hours       []PlanetaryHour `json:"hours"`
}

// Chaldean order of planets (used for planetary hours)
var chaldeanOrder = []models.PlanetID{
	models.PlanetSaturn,
	models.PlanetJupiter,
	models.PlanetMars,
	models.PlanetSun,
	models.PlanetVenus,
	models.PlanetMercury,
	models.PlanetMoon,
}

// dayRulerOrder maps weekday (0=Sunday) to Chaldean index
var dayRulerOrder = []int{3, 6, 2, 5, 1, 4, 0} // Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn

// CalcPlanetaryHours computes the 24 planetary hours for a given date and location.
// Uses actual sunrise/sunset times calculated via Sun position.
func CalcPlanetaryHours(jdUT, lat, lon float64) (*PlanetaryDay, error) {
	// Find sunrise and sunset for this day
	sunrise, err := findSunrise(jdUT, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("sunrise: %w", err)
	}
	sunset, err := findSunset(jdUT, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("sunset: %w", err)
	}
	nextSunrise, err := findSunrise(jdUT+1, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("next sunrise: %w", err)
	}

	// Day ruler from weekday
	y, m, d, _ := sweph.RevJul(sunrise, true)
	weekday := julianWeekday(y, m, d)
	dayRulerIdx := dayRulerOrder[weekday]

	dayDuration := sunset - sunrise
	nightDuration := nextSunrise - sunset
	dayHourLen := dayDuration / 12
	nightHourLen := nightDuration / 12

	hours := make([]PlanetaryHour, 24)
	planetIdx := dayRulerIdx

	for i := 0; i < 24; i++ {
		var start, end float64
		var isDay bool

		if i < 12 {
			start = sunrise + float64(i)*dayHourLen
			end = start + dayHourLen
			isDay = true
		} else {
			nightIdx := i - 12
			start = sunset + float64(nightIdx)*nightHourLen
			end = start + nightHourLen
			isDay = false
		}

		hours[i] = PlanetaryHour{
			Number:  i + 1,
			Planet:  chaldeanOrder[planetIdx%7],
			StartJD: start,
			EndJD:   end,
			IsDay:   isDay,
		}
		planetIdx++
	}

	return &PlanetaryDay{
		DayRuler:    chaldeanOrder[dayRulerIdx],
		Sunrise:     sunrise,
		Sunset:      sunset,
		NextSunrise: nextSunrise,
		Hours:       hours,
	}, nil
}

// CurrentPlanetaryHour returns the planetary hour active at the given JD
func CurrentPlanetaryHour(jdUT, lat, lon float64) (*PlanetaryHour, error) {
	day, err := CalcPlanetaryHours(jdUT, lat, lon)
	if err != nil {
		return nil, err
	}
	for _, h := range day.Hours {
		if jdUT >= h.StartJD && jdUT < h.EndJD {
			return &h, nil
		}
	}
	return nil, fmt.Errorf("JD not within calculated planetary hours range")
}

// findSunrise finds the approximate JD of sunrise for a given day and location.
// Uses iterative refinement on the Sun's altitude.
func findSunrise(jdUT, lat, lon float64) (float64, error) {
	return findSunEvent(jdUT, lat, lon, true)
}

// findSunset finds the approximate JD of sunset for a given day and location.
func findSunset(jdUT, lat, lon float64) (float64, error) {
	return findSunEvent(jdUT, lat, lon, false)
}

// findSunEvent finds sunrise (rising=true) or sunset (rising=false)
// using bisection on the Sun's altitude above the horizon.
// Standard refraction correction: -0.8333 degrees below geometric horizon.
func findSunEvent(jdUT, lat, lon float64, rising bool) (float64, error) {
	// Start from noon
	y, m, d, _ := sweph.RevJul(jdUT, true)
	noon := sweph.JulDay(y, m, d, 12.0, true)

	var lo, hi float64
	if rising {
		lo = noon - 0.5 // midnight before
		hi = noon        // noon
	} else {
		lo = noon        // noon
		hi = noon + 0.5  // midnight after
	}

	// Bisect to find when altitude crosses -0.8333°
	const targetAlt = -0.8333
	const eps = 1.0 / 86400.0 // 1 second

	for hi-lo > eps {
		mid := (lo + hi) / 2
		alt := sunAltitude(mid, lat, lon)
		if rising {
			if alt < targetAlt {
				lo = mid
			} else {
				hi = mid
			}
		} else {
			if alt > targetAlt {
				lo = mid
			} else {
				hi = mid
			}
		}
	}
	return (lo + hi) / 2, nil
}

// sunAltitude computes the topocentric altitude of the Sun
func sunAltitude(jdUT, lat, lon float64) float64 {
	result, err := sweph.CalcUT(jdUT, sweph.SE_SUN)
	if err != nil {
		return 0
	}
	sunLon := result.Longitude

	// Get sidereal time approximation
	// GMST = 280.46061837 + 360.98564736629 * (JD - 2451545.0) + lon
	gmst := 280.46061837 + 360.98564736629*(jdUT-2451545.0)
	lst := math.Mod(gmst+lon, 360)
	if lst < 0 {
		lst += 360
	}

	// Hour angle
	ra := sunLon // simplified: RA ≈ ecliptic longitude for Sun
	ha := lst - ra

	// Convert to radians
	latRad := lat * math.Pi / 180
	decRad := 0.0 // approximate: use ecliptic latitude

	// Sun's declination from ecliptic longitude
	eps := 23.4393 // approximate obliquity
	epsRad := eps * math.Pi / 180
	lonRad := sunLon * math.Pi / 180
	decRad = math.Asin(math.Sin(epsRad) * math.Sin(lonRad))

	haRad := ha * math.Pi / 180

	// Altitude = asin(sin(lat)*sin(dec) + cos(lat)*cos(dec)*cos(ha))
	alt := math.Asin(math.Sin(latRad)*math.Sin(decRad) + math.Cos(latRad)*math.Cos(decRad)*math.Cos(haRad))
	return alt * 180 / math.Pi
}

// julianWeekday returns the day of week (0=Sunday, 6=Saturday) for a date
func julianWeekday(year, month, day int) int {
	jd := sweph.JulDay(year, month, day, 12.0, true)
	return int(math.Mod(jd+1.5, 7))
}
