package primary

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// DirectionKey selects the rate at which the arc of direction converts to years.
type DirectionKey string

const (
	KeyNaibod   DirectionKey = "NAIBOD"    // 0.98556 deg/year (Sun's mean daily motion)
	KeyPtolemy  DirectionKey = "PTOLEMY"   // Same as Naibod (classical)
	KeySolarArc DirectionKey = "SOLAR_ARC" // Actual Sun's yearly motion at birth
)

// DirectionType distinguishes direct (zodiacal) from converse directions.
type DirectionType string

const (
	DirectDirect   DirectionType = "DIRECT"   // In order of signs (zodiacal)
	DirectConverse DirectionType = "CONVERSE" // Against order of signs
)

// NaibodRate is the Naibod key in degrees per year: 360 / 365.25 ≈ 0.98556.
const NaibodRate = 0.98556

// PrimaryDirection represents a single primary direction hit.
type PrimaryDirection struct {
	Promissor    string            `json:"promissor"`
	Significator string            `json:"significator"`
	AspectType   models.AspectType `json:"aspect_type"`
	Arc          float64           `json:"arc"`
	AgeExact     float64           `json:"age_exact"`
	Direction    DirectionType     `json:"direction"`
	Key          DirectionKey      `json:"direction_key"`
}

// PrimaryDirectionInput holds the parameters for a primary direction calculation.
type PrimaryDirectionInput struct {
	NatalJD     float64
	GeoLat      float64
	GeoLon      float64
	Planets     []models.PlanetID
	Aspects     []models.AspectType
	Key         DirectionKey
	MaxAge      float64
	HouseSystem models.HouseSystem
}

// PrimaryDirectionResult holds the computed primary directions.
type PrimaryDirectionResult struct {
	Directions []PrimaryDirection `json:"directions"`
}

// defaultAspects used when none are specified.
var defaultAspects = []models.AspectType{
	models.AspectConjunction,
	models.AspectOpposition,
	models.AspectSquare,
	models.AspectTrine,
	models.AspectSextile,
}

// aspectAngle returns the angle for a standard aspect type.
func aspectAngle(at models.AspectType) float64 {
	for _, a := range models.StandardAspects {
		if a.Type == at {
			return a.Angle
		}
	}
	return 0
}

// deg2rad converts degrees to radians.
func deg2rad(d float64) float64 { return d * math.Pi / 180 }

// rad2deg converts radians to degrees.
func rad2deg(r float64) float64 { return r * 180 / math.Pi }

// eclipticToRA converts ecliptic longitude (assuming lat=0) to right ascension.
func eclipticToRA(lon, obliquity float64) float64 {
	lonRad := deg2rad(lon)
	epsRad := deg2rad(obliquity)
	ra := math.Atan2(math.Sin(lonRad)*math.Cos(epsRad), math.Cos(lonRad))
	return sweph.NormalizeDegrees(rad2deg(ra))
}

// eclipticToDec converts ecliptic longitude (assuming lat=0) to declination.
func eclipticToDec(lon, obliquity float64) float64 {
	lonRad := deg2rad(lon)
	epsRad := deg2rad(obliquity)
	dec := math.Asin(math.Sin(lonRad) * math.Sin(epsRad))
	return rad2deg(dec)
}

// ascensionalDifference returns the ascensional difference: AD = arcsin(tan(dec) * tan(geoLat)).
// Returns 0 if the computation would be out of domain (circumpolar).
func ascensionalDifference(dec, geoLat float64) float64 {
	val := math.Tan(deg2rad(dec)) * math.Tan(deg2rad(geoLat))
	if val > 1 {
		val = 1
	} else if val < -1 {
		val = -1
	}
	return rad2deg(math.Asin(val))
}

// semiArc returns the semi-diurnal arc for a point with given declination at
// the given geographic latitude. If the point is below the horizon (nocturnal),
// the semi-nocturnal arc is returned instead.
func semiArc(dec, geoLat float64) float64 {
	ad := ascensionalDifference(dec, geoLat)
	sa := 90 + ad // diurnal semi-arc
	if sa < 0 {
		sa = -sa
	}
	return sa
}

// semiArcNocturnal returns the semi-nocturnal arc.
func semiArcNocturnal(dec, geoLat float64) float64 {
	ad := ascensionalDifference(dec, geoLat)
	sa := 90 - ad
	if sa < 0 {
		sa = -sa
	}
	return sa
}

// arcToYears converts an arc in degrees to years using the specified key.
// For SOLAR_ARC, solarArcRate must be provided (Sun's actual daily motion in degrees).
func arcToYears(arc float64, key DirectionKey, solarArcRate float64) float64 {
	switch key {
	case KeySolarArc:
		if solarArcRate <= 0 {
			return arc / NaibodRate
		}
		return arc / solarArcRate
	default: // Naibod and Ptolemy
		return arc / NaibodRate
	}
}

// pointData holds pre-computed equatorial data for a chart point.
type pointData struct {
	id   string  // planet ID or special point ID
	lon  float64 // ecliptic longitude
	ra   float64 // right ascension
	dec  float64 // declination
	sa   float64 // diurnal semi-arc
	saNt float64 // nocturnal semi-arc
}

// CalcPrimaryDirections computes primary direction hits for all promissor-significator-aspect
// combinations within the configured age range.
func CalcPrimaryDirections(input PrimaryDirectionInput) (*PrimaryDirectionResult, error) {
	// Apply defaults
	if input.Key == "" {
		input.Key = KeyNaibod
	}
	if input.MaxAge <= 0 {
		input.MaxAge = 100
	}
	if len(input.Aspects) == 0 {
		input.Aspects = defaultAspects
	}
	if len(input.Planets) == 0 {
		input.Planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}
	if input.HouseSystem == "" {
		input.HouseSystem = models.HousePlacidus
	}

	// Get obliquity
	obliquity, err := sweph.Obliquity(input.NatalJD)
	if err != nil {
		return nil, fmt.Errorf("primary: obliquity: %w", err)
	}

	// Get house cusps and angles
	hsysChar := models.HouseSystemToChar(input.HouseSystem)
	hr, err := sweph.Houses(input.NatalJD, input.GeoLat, input.GeoLon, hsysChar)
	if err != nil {
		return nil, fmt.Errorf("primary: houses: %w", err)
	}

	// Build point data for all planets
	var points []pointData
	for _, pid := range input.Planets {
		sweID, ok := models.PlanetToSweID(pid)
		if !ok {
			continue
		}
		res, err := sweph.CalcUT(input.NatalJD, sweID)
		if err != nil {
			return nil, fmt.Errorf("primary: calc %s: %w", pid, err)
		}
		lon := res.Longitude
		ra := eclipticToRA(lon, obliquity)
		dec := eclipticToDec(lon, obliquity)
		points = append(points, pointData{
			id:   string(pid),
			lon:  lon,
			ra:   ra,
			dec:  dec,
			sa:   semiArc(dec, input.GeoLat),
			saNt: semiArcNocturnal(dec, input.GeoLat),
		})
	}

	// Add angles as significators: MC and ASC
	mcRA := eclipticToRA(hr.MC, obliquity)
	mcDec := eclipticToDec(hr.MC, obliquity)
	ascRA := eclipticToRA(hr.ASC, obliquity)
	ascDec := eclipticToDec(hr.ASC, obliquity)

	anglePoints := []pointData{
		{
			id:   string(models.PointMC),
			lon:  hr.MC,
			ra:   mcRA,
			dec:  mcDec,
			sa:   semiArc(mcDec, input.GeoLat),
			saNt: semiArcNocturnal(mcDec, input.GeoLat),
		},
		{
			id:   string(models.PointASC),
			lon:  hr.ASC,
			ra:   ascRA,
			dec:  ascDec,
			sa:   semiArc(ascDec, input.GeoLat),
			saNt: semiArcNocturnal(ascDec, input.GeoLat),
		},
	}

	// All significators = planets + angles
	significators := append(points, anglePoints...)
	// All promissors = planets + angles
	promissors := append(points, anglePoints...)

	// Solar arc rate for SOLAR_ARC key
	var solarArcRate float64
	if input.Key == KeySolarArc {
		sunRes, err := sweph.CalcUT(input.NatalJD, sweph.SE_SUN)
		if err == nil {
			solarArcRate = math.Abs(sunRes.SpeedLong) // deg/day ≈ yearly rate in primary directions
		}
	}

	maxArc := input.MaxAge * NaibodRate
	if input.Key == KeySolarArc && solarArcRate > 0 {
		maxArc = input.MaxAge * solarArcRate
	}

	var directions []PrimaryDirection

	for _, prom := range promissors {
		for _, sig := range significators {
			// Skip same point
			if prom.id == sig.id {
				continue
			}

			for _, aspType := range input.Aspects {
				angle := aspectAngle(aspType)

				// Apply aspect to promissor's RA
				for _, dir := range []struct {
					dt    DirectionType
					raAdj float64
				}{
					{DirectDirect, angle},
					{DirectConverse, -angle},
				} {
					promRA := sweph.NormalizeDegrees(prom.ra + dir.raAdj)

					// Ptolemy semi-arc method:
					// arc = (SA_sig / SA_prom) * (promRA_aspected - sig.ra)
					// Use diurnal semi-arcs
					if prom.sa == 0 || sig.sa == 0 {
						continue
					}

					rawArc := promRA - sig.ra
					// Normalize to [-180, 180)
					if rawArc > 180 {
						rawArc -= 360
					} else if rawArc < -180 {
						rawArc += 360
					}

					arc := math.Abs((sig.sa / prom.sa) * rawArc)

					if arc <= 0 || arc > maxArc {
						continue
					}

					age := arcToYears(arc, input.Key, solarArcRate)
					if age <= 0 || age > input.MaxAge {
						continue
					}

					directions = append(directions, PrimaryDirection{
						Promissor:    prom.id,
						Significator: sig.id,
						AspectType:   aspType,
						Arc:          math.Round(arc*10000) / 10000,
						AgeExact:     math.Round(age*10000) / 10000,
						Direction:    dir.dt,
						Key:          input.Key,
					})
				}
			}
		}
	}

	return &PrimaryDirectionResult{Directions: directions}, nil
}
