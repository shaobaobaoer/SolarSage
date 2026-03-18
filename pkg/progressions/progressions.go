package progressions

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// JulianYear is the length of one Julian year in days
const JulianYear = 365.25

// SecondaryProgressionJD converts a real-time JD to the progressed JD.
// Secondary progressions: 1 day of real time after birth = 1 year of progressed life.
// progressedJD = natalJD + (transitJD - natalJD) / JulianYear
func SecondaryProgressionJD(natalJD, transitJD float64) float64 {
	return natalJD + (transitJD-natalJD)/JulianYear
}

// CalcProgressedLongitude returns the ecliptic longitude and speed of a planet
// at the secondary progressed time corresponding to the given transit JD.
// The returned speed is in progressed degrees per real day (i.e. actual speed / JulianYear).
func CalcProgressedLongitude(planet models.PlanetID, natalJD, transitJD float64) (lon, speed float64, err error) {
	pJD := SecondaryProgressionJD(natalJD, transitJD)
	lon, rawSpeed, err := chart.CalcPlanetLongitude(planet, pJD)
	if err != nil {
		return 0, 0, fmt.Errorf("progressions calc %s: %w", planet, err)
	}
	// Speed in real-time units: actual speed (°/day at progressed date) / JulianYear
	speed = rawSpeed / JulianYear
	return lon, speed, nil
}

// SolarArcOffset returns the solar arc offset in degrees for a given transit JD.
// Solar Arc = Sun's progressed longitude - Sun's natal longitude
func SolarArcOffset(natalJD, transitJD float64) (float64, error) {
	natalSun, err := sweph.CalcUT(natalJD, sweph.SE_SUN)
	if err != nil {
		return 0, fmt.Errorf("solar arc natal sun: %w", err)
	}

	progressedJD := SecondaryProgressionJD(natalJD, transitJD)
	progressedSun, err := sweph.CalcUT(progressedJD, sweph.SE_SUN)
	if err != nil {
		return 0, fmt.Errorf("solar arc progressed sun: %w", err)
	}

	return progressedSun.Longitude - natalSun.Longitude, nil
}

// CalcSolarArcLongitude returns the solar arc directed position of a natal planet.
// Solar Arc Direction = natal longitude + solar arc offset
// Speed is approximately the sun's progressed speed / JulianYear
func CalcSolarArcLongitude(planet models.PlanetID, natalJD, transitJD float64) (lon, speed float64, err error) {
	// Get natal planet position
	natalLon, _, err := chart.CalcPlanetLongitude(planet, natalJD)
	if err != nil {
		return 0, 0, err
	}

	offset, err := SolarArcOffset(natalJD, transitJD)
	if err != nil {
		return 0, 0, err
	}

	lon = sweph.NormalizeDegrees(natalLon + offset)

	// Solar arc speed ≈ Sun's progressed speed / JulianYear ≈ ~0.00274°/day
	progressedJD := SecondaryProgressionJD(natalJD, transitJD)
	_, sunSpeed, err := chart.CalcPlanetLongitude(models.PlanetSun, progressedJD)
	if err != nil {
		return lon, 0, nil // non-fatal
	}
	speed = sunSpeed / JulianYear

	return lon, speed, nil
}

// Age returns the age in Julian years at a given transit JD relative to natal JD
func Age(natalJD, transitJD float64) float64 {
	return (transitJD - natalJD) / JulianYear
}

// eclipticToRA converts ecliptic longitude to right ascension (for bodies on the ecliptic, lat=0)
func eclipticToRA(lon, obliquity float64) float64 {
	lonRad := lon * math.Pi / 180
	epsRad := obliquity * math.Pi / 180
	ra := math.Atan2(math.Sin(lonRad)*math.Cos(epsRad), math.Cos(lonRad))
	return sweph.NormalizeDegrees(ra * 180 / math.Pi)
}

// mcFromRAMC computes MC ecliptic longitude from RAMC and obliquity
func mcFromRAMC(ramc, obliquity float64) float64 {
	ramcRad := ramc * math.Pi / 180
	epsRad := obliquity * math.Pi / 180
	mc := math.Atan2(math.Sin(ramcRad), math.Cos(ramcRad)*math.Cos(epsRad))
	return sweph.NormalizeDegrees(mc * 180 / math.Pi)
}

// ascFromRAMC computes ASC ecliptic longitude from RAMC, obliquity, and geographic latitude
func ascFromRAMC(ramc, obliquity, geoLat float64) float64 {
	ramcRad := ramc * math.Pi / 180
	epsRad := obliquity * math.Pi / 180
	latRad := geoLat * math.Pi / 180
	// Standard formula: atan2(-cos(RAMC), sin(eps)*tan(lat) + cos(eps)*sin(RAMC))
	// gives DSC. Negate both args to get ASC (rotate 180 deg).
	y := math.Cos(ramcRad)
	x := -(math.Sin(epsRad)*math.Tan(latRad) + math.Cos(epsRad)*math.Sin(ramcRad))
	return sweph.NormalizeDegrees(math.Atan2(y, x) * 180 / math.Pi)
}

// CalcProgressedAngles computes progressed ASC and MC.
// Uses solar arc in longitude for MC, then derives ASC from the RAMC of the progressed MC.
// Uses the standard Solar Arc in RA method for progressed angles.
func CalcProgressedAngles(natalJD, transitJD, geoLat, geoLon float64, hsys models.HouseSystem) (asc, mc float64, err error) {
	// Get obliquity at natal epoch
	eps, err := sweph.Obliquity(natalJD)
	if err != nil {
		return 0, 0, fmt.Errorf("obliquity: %w", err)
	}

	// Get natal MC
	hsysChar := models.HouseSystemToChar(hsys)
	natalHR, err := sweph.Houses(natalJD, geoLat, geoLon, hsysChar)
	if err != nil {
		return 0, 0, fmt.Errorf("natal houses: %w", err)
	}
	natalMC := natalHR.MC

	// Get solar arc offset (progressed Sun - natal Sun in ecliptic longitude)
	offset, err := SolarArcOffset(natalJD, transitJD)
	if err != nil {
		return 0, 0, fmt.Errorf("solar arc: %w", err)
	}

	// Progressed MC = natal MC + solar arc offset
	mc = sweph.NormalizeDegrees(natalMC + offset)

	// Derive RAMC from progressed MC: tan(RAMC) = tan(MC) * cos(eps)
	mcRad := mc * math.Pi / 180
	epsRad := eps * math.Pi / 180
	progRAMC := math.Atan2(math.Sin(mcRad)*math.Cos(epsRad), math.Cos(mcRad))
	progRAMCDeg := sweph.NormalizeDegrees(progRAMC * 180 / math.Pi)

	// ASC from progressed RAMC
	asc = ascFromRAMC(progRAMCDeg, eps, geoLat)

	return asc, mc, nil
}

// CalcProgressedSpecialPoint returns the progressed longitude of a special point (ASC/MC)
// using the Q1/Solar Arc in RA method.
func CalcProgressedSpecialPoint(sp models.SpecialPointID, natalJD, transitJD, geoLat, geoLon float64, hsys models.HouseSystem) (float64, error) {
	asc, mc, err := CalcProgressedAngles(natalJD, transitJD, geoLat, geoLon, hsys)
	if err != nil {
		return 0, err
	}
	switch sp {
	case models.PointASC:
		return asc, nil
	case models.PointMC:
		return mc, nil
	case models.PointDSC:
		return sweph.NormalizeDegrees(asc + 180), nil
	case models.PointIC:
		return sweph.NormalizeDegrees(mc + 180), nil
	default:
		return 0, fmt.Errorf("progressed special point %s not supported", sp)
	}
}
