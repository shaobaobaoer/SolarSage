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
// progressedJD = natalJDE + (transitJD - natalJDE) / JulianYear
// NOTE: Uses JDE (Ephemeris Time) to match Solar Fire's convention
func SecondaryProgressionJD(natalJD, transitJD float64) float64 {
	// Convert natal JD to JDE by adding ΔT
	natalDeltaT := sweph.DeltaT(natalJD)
	natalJDE := natalJD + natalDeltaT
	
	return natalJDE + (transitJD-natalJDE)/JulianYear
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

	// For secondary progression, the "progressed" Sun position is calculated
	// at the progressed JD, which is natalJDE + (transitJD - natalJDE) / JulianYear
	// The SecondaryProgressionJD function returns JDE, so we need to convert to UT
	progressedJDE := SecondaryProgressionJD(natalJD, transitJD)
	progressedDeltaT := sweph.DeltaT(progressedJDE)
	progressedJD := progressedJDE - progressedDeltaT
	
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

	// Solar arc speed ≈ Sun's progressed speed / SolarFireProgressionFactor ≈ ~0.00274°/day
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

// CalcProgressedAngles computes progressed ASC and MC using Solar Arc in Longitude.
// Parameters:
//   - natalMCOverrideForMC: if non-zero, use this for progressed MC calculation
//   - natalMCOverrideForASC: controls ASC derivation method:
//     - > 0: use this MC value for MC→RAMC→ASC chain
//     - == -1: use sweph-computed MC for MC→RAMC→ASC chain
//     - == 0: fall back to natalMCOverrideForMC
//   - natalASCOverride: controls ASC progression method:
//     - > 0: use DIRECT solar arc to ASC (progASC = natalASC + solarArc)
//     - == -1: use Solar Arc in Right Ascension method (SF style)
//     - == 0: use MC→RAMC→ASC chain (traditional)
//
// For Solar Fire compatibility:
//   - MC progression uses SF meta's precise MC value
//   - ASC progression uses Solar Arc in Right Ascension method
func CalcProgressedAngles(natalJD, transitJD, geoLat, geoLon float64, hsys models.HouseSystem, natalMCOverrideForMC, natalMCOverrideForASC, natalASCOverride float64) (asc, mc float64, err error) {
	// Get obliquity at natal epoch
	eps, err := sweph.Obliquity(natalJD)
	if err != nil {
		return 0, 0, fmt.Errorf("obliquity: %w", err)
	}

	// Get natal MC from sweph (always computed for fallback)
	hsysChar := models.HouseSystemToChar(hsys)
	natalHR, err := sweph.Houses(natalJD, geoLat, geoLon, hsysChar)
	if err != nil {
		return 0, 0, fmt.Errorf("natal houses: %w", err)
	}
	natalMCSweph := natalHR.MC

	// Get solar arc offset (progressed Sun - natal Sun in ecliptic longitude)
	offset, err := SolarArcOffset(natalJD, transitJD)
	if err != nil {
		return 0, 0, fmt.Errorf("solar arc: %w", err)
	}

	// Progressed MC: use override if provided, otherwise sweph value
	var natalMCForMC float64
	if natalMCOverrideForMC != 0 {
		natalMCForMC = natalMCOverrideForMC
	} else {
		natalMCForMC = natalMCSweph
	}
	mc = sweph.NormalizeDegrees(natalMCForMC + offset)

	// Progressed ASC: three methods available
	//
	// Method 1 (Direct): progASC = natalASC + solarArc
	//   Simple ecliptic longitude addition.
	//
	// Method 2 (Solar Arc in Right Ascension - SF style):
	//   Convert natal ASC to RA, add solar arc offset, convert back to longitude.
	//   This accounts for the obliquity effect on the ASC.
	//
	// Method 3 (Traditional): MC→RAMC→ASC chain
	//   progMC = natalMC + solarArc
	//   progRAMC = MC_to_RAMC(progMC)
	//   progASC = RAMC_to_ASC(progRAMC)

	if natalASCOverride > 0 {
		// Method 1: Direct solar arc to ASC
		asc = sweph.NormalizeDegrees(natalASCOverride + offset)
	} else if natalASCOverride == -1 {
		// Method 2: Solar Arc in Right Ascension (SF style)
		// Get natal ASC from sweph
		natalASC := natalHR.ASC
		// Convert natal ASC to RA
		natalASCRad := natalASC * math.Pi / 180
		epsRad := eps * math.Pi / 180
		natalRA := math.Atan2(math.Sin(natalASCRad)*math.Cos(epsRad)-math.Tan(geoLat*math.Pi/180)*math.Sin(epsRad), math.Cos(natalASCRad))
		// Add solar arc offset to RA
		progRA := sweph.NormalizeDegrees(natalRA*180/math.Pi + offset)
		// Convert back to ecliptic longitude
		progRARad := progRA * math.Pi / 180
		latRad := geoLat * math.Pi / 180
		asc = sweph.NormalizeDegrees(math.Atan2(math.Sin(progRARad)*math.Cos(epsRad)+math.Tan(latRad)*math.Sin(epsRad), math.Cos(progRARad)) * 180 / math.Pi)
	} else {
		// Method 3: MC→RAMC→ASC chain (traditional)
		var natalMCForASC float64
		if natalMCOverrideForASC > 0 {
			natalMCForASC = natalMCOverrideForASC
		} else if natalMCOverrideForASC == -1 {
			natalMCForASC = natalMCSweph
		} else if natalMCOverrideForMC != 0 {
			natalMCForASC = natalMCOverrideForMC
		} else {
			natalMCForASC = natalMCSweph
		}

		progMCForASC := sweph.NormalizeDegrees(natalMCForASC + offset)
		mcRad := progMCForASC * math.Pi / 180
		epsRad := eps * math.Pi / 180
		progRAMC := math.Atan2(math.Sin(mcRad)*math.Cos(epsRad), math.Cos(mcRad))
		progRAMCDeg := sweph.NormalizeDegrees(progRAMC * 180 / math.Pi)

		asc = ascFromRAMC(progRAMCDeg, eps, geoLat)
	}

	return asc, mc, nil
}

// CalcProgressedSpecialPoint returns the progressed longitude of a special point (ASC/MC)
// using the Solar Arc in Longitude method.
// natalMCOverrideForMC: override for MC progression
// natalMCOverrideForASC: override for ASC derivation (MC→RAMC→ASC chain, only if natalASCOverride == 0)
// natalASCOverride: controls ASC progression method:
//   - > 0: direct solar arc to ASC
//   - == -1: Solar Arc in Right Ascension method
//   - == 0: MC→RAMC→ASC chain (traditional)
func CalcProgressedSpecialPoint(sp models.SpecialPointID, natalJD, transitJD, geoLat, geoLon float64, hsys models.HouseSystem, natalMCOverrideForMC, natalMCOverrideForASC, natalASCOverride float64) (float64, error) {
	asc, mc, err := CalcProgressedAngles(natalJD, transitJD, geoLat, geoLon, hsys, natalMCOverrideForMC, natalMCOverrideForASC, natalASCOverride)
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
