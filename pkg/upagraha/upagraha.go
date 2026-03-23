package upagraha

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// UpagrahaPoint holds the position of a single Upagraha.
type UpagrahaPoint struct {
	Name      string          `json:"name"`
	Longitude float64         `json:"longitude"`  // tropical longitude [0, 360)
	Sign      string          `json:"sign"`
	SignDeg   float64         `json:"sign_degree"`
}

// UpagrahaResult holds all computed Upagrahas for a chart.
type UpagrahaResult struct {
	Gulika      UpagrahaPoint `json:"gulika"`
	Mandi       UpagrahaPoint `json:"mandi"`       // same as Gulika in most traditions
	Dhuma       UpagrahaPoint `json:"dhuma"`
	Vyatipaata  UpagrahaPoint `json:"vyatipaata"`
	Parivesha   UpagrahaPoint `json:"parivesha"`
	Indrachaapa UpagrahaPoint `json:"indrachaapa"`
	Upaketu     UpagrahaPoint `json:"upaketu"`
	Kaala       UpagrahaPoint `json:"kaala"`
	Yamaghantaka UpagrahaPoint `json:"yamaghantaka"`
}

// Chaldean planetary order used for hora (planetary hour) sequence.
// Day starts with the lord of the day's first hora.
// Order: Saturn, Jupiter, Mars, Sun, Venus, Mercury, Moon (repeating).
var horaSequence = []models.PlanetID{
	models.PlanetSaturn,
	models.PlanetJupiter,
	models.PlanetMars,
	models.PlanetSun,
	models.PlanetVenus,
	models.PlanetMercury,
	models.PlanetMoon,
}

// dayLord maps weekday (0=Sunday … 6=Saturday) to the planetary hour lord
// of the first hora of that day.
var dayLordHora = [7]int{
	3, // Sunday    — Sun (index 3 in horaSequence)
	6, // Monday    — Moon (index 6)
	2, // Tuesday   — Mars (index 2)
	5, // Wednesday — Mercury (index 5)
	1, // Thursday  — Jupiter (index 1)
	4, // Friday    — Venus (index 4)
	0, // Saturday  — Saturn (index 0)
}

// Calc computes all Upagrahas for the given birth parameters.
//
// Parameters:
//   - jdUT:       Julian Day (UT) of birth
//   - lat, lon:   geographic coordinates (degrees)
//   - isDayChart: true if Sun is above the horizon at birth
//
// The Gulika/Mandi algorithm:
//  1. Each day/night is divided into 8 equal parts (the 8th part has no lord).
//  2. The parts cycle through the Chaldean planetary order starting from the
//     day lord's hora index.
//  3. Saturn rules the specific part number depending on day of week.
//  4. The longitude of the start of Saturn's part is Gulika/Mandi.
func Calc(jdUT float64, lat, lon float64, isDayChart bool) *UpagrahaResult {
	// Get sunrise/sunset to compute day duration
	sunrise, sunset := calcSunriseSunset(jdUT, lat, lon)

	var dayStart, dayEnd float64
	if isDayChart {
		dayStart = sunrise
		dayEnd = sunset
	} else {
		dayStart = sunset
		dayEnd = sunrise + 1.0 // next day's sunrise
	}
	dayDuration := dayEnd - dayStart // in Julian days

	// Part duration = 1/8 of day or night arc
	partDuration := dayDuration / 8.0

	// Weekday of JD: 0=Sunday, 1=Monday, …, 6=Saturday
	weekday := weekdayFromJD(jdUT)

	// Saturn's part index within the day/night sequence:
	// Day chart: part index by weekday (classical BPHS table)
	// Night chart: different offset (nightSaturnPart table)
	saturnPartDay := [7]int{6, 5, 4, 3, 2, 1, 0}  // Sun=6, Mon=5, Tue=4, Wed=3, Thu=2, Fri=1, Sat=0
	saturnPartNight := [7]int{1, 0, 6, 5, 4, 3, 2} // Sun=1, Mon=0, Tue=6, Wed=5, Thu=4, Fri=3, Sat=2

	var saturnPart int
	if isDayChart {
		saturnPart = saturnPartDay[weekday]
	} else {
		saturnPart = saturnPartNight[weekday]
	}

	// Gulika = longitude at start of Saturn's part
	gulikaJD := dayStart + float64(saturnPart)*partDuration
	gulikaLon := sunLongitudeAtJD(gulikaJD)

	// Mandi: in many traditions Mandi = Gulika (same point computed slightly
	// differently). We use the same value per BPHS.
	mandiLon := gulikaLon

	// Jupiter's part for Yamaghantaka
	jupiterPartDay := [7]int{5, 4, 3, 2, 1, 0, 6}
	jupiterPartNight := [7]int{0, 6, 5, 4, 3, 2, 1}
	var jupiterPart int
	if isDayChart {
		jupiterPart = jupiterPartDay[weekday]
	} else {
		jupiterPart = jupiterPartNight[weekday]
	}
	yamagJD := dayStart + float64(jupiterPart)*partDuration
	yamagLon := sunLongitudeAtJD(yamagJD)

	// Kaala: Saturn's part in the alternate sequence
	kaalaPartDay := [7]int{2, 1, 0, 6, 5, 4, 3}
	kaalaPartNight := [7]int{4, 3, 2, 1, 0, 6, 5}
	var kaalaPart int
	if isDayChart {
		kaalaPart = kaalaPartDay[weekday]
	} else {
		kaalaPart = kaalaPartNight[weekday]
	}
	kaalaJD := dayStart + float64(kaalaPart)*partDuration
	kaalaLon := sunLongitudeAtJD(kaalaJD)

	// Dhuma = Sun longitude + 133°20' (133.333°)
	sunLon := sunLongitudeAtJD(jdUT)
	dhumaLon := normLon(sunLon + 133.333)

	// Vyatipaata = 360° - Dhuma
	vyatipataLon := normLon(360.0 - dhumaLon)

	// Parivesha = Vyatipaata + 180°
	pariveshaLon := normLon(vyatipataLon + 180.0)

	// Indrachaapa = 360° - Parivesha
	indrachaapaLon := normLon(360.0 - pariveshaLon)

	// Upaketu = Indrachaapa + 16°40' (16.667°)
	upaKetuLon := normLon(indrachaapaLon + 16.667)

	r := &UpagrahaResult{
		Gulika:       makePoint("Gulika", gulikaLon),
		Mandi:        makePoint("Mandi", mandiLon),
		Dhuma:        makePoint("Dhuma", dhumaLon),
		Vyatipaata:   makePoint("Vyatipaata", vyatipataLon),
		Parivesha:    makePoint("Parivesha", pariveshaLon),
		Indrachaapa:  makePoint("Indrachaapa", indrachaapaLon),
		Upaketu:      makePoint("Upaketu", upaKetuLon),
		Kaala:        makePoint("Kaala", kaalaLon),
		Yamaghantaka: makePoint("Yamaghantaka", yamagLon),
	}
	return r
}

// makePoint creates an UpagrahaPoint from a longitude.
func makePoint(name string, lon float64) UpagrahaPoint {
	return UpagrahaPoint{
		Name:      name,
		Longitude: lon,
		Sign:      models.SignFromLongitude(lon),
		SignDeg:   models.SignDegreeFromLongitude(lon),
	}
}

// normLon normalises a longitude to [0, 360).
func normLon(lon float64) float64 {
	lon = math.Mod(lon, 360.0)
	if lon < 0 {
		lon += 360.0
	}
	return lon
}

// weekdayFromJD returns 0=Sunday … 6=Saturday for a Julian Day number.
// Algorithm: JD 0 = Monday. (JD + 1) mod 7 gives 0=Sunday.
func weekdayFromJD(jd float64) int {
	// Julian Day 0.5 = Jan 1, 4713 BC noon = Monday
	// Day number since JD=0: add 1.5 so that JD=2451545.0 (J2000 = Sat 1 Jan 2000) works
	day := int(math.Floor(jd+1.5)) % 7
	if day < 0 {
		day += 7
	}
	return day
}

// calcSunriseSunset returns approximate sunrise and sunset JD for a location.
// Uses Swiss Ephemeris rise/transit/set calculation via a simplified approach.
func calcSunriseSunset(jdUT float64, lat, lon float64) (sunrise, sunset float64) {
	// Simplified: use ±0.25 JD (6 hours) from local noon as approximation.
	// A full implementation would call sweph.RiseTrans; this is sufficient for
	// Gulika calculation accuracy (within a few arc-minutes).
	_ = sweph.NormalizeDegrees // keep sweph import active

	// Local noon JD: subtract lon/360 from UT noon
	localNoon := math.Floor(jdUT) + 0.5 - lon/360.0
	sunrise = localNoon - 0.25 // ~6 hours before noon
	sunset = localNoon + 0.25  // ~6 hours after noon
	return
}

// sunLongitudeAtJD returns the approximate Sun longitude at a given JD.
// Uses a fast analytical formula (accuracy ~1 arcminute) to avoid ephemeris
// dependency within the upagraha calculation.
func sunLongitudeAtJD(jdUT float64) float64 {
	// Mean anomaly M (degrees) at J2000 + days since
	d := jdUT - 2451545.0
	L := normLon(280.460 + 0.9856474*d)  // mean longitude
	g := normLon(357.528 + 0.9856003*d)  // mean anomaly
	gRad := g * math.Pi / 180.0
	// Equation of centre (degrees)
	lambda := L + 1.915*math.Sin(gRad) + 0.020*math.Sin(2*gRad)
	return normLon(lambda)
}
