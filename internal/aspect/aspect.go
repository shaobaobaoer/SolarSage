package aspect

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// AngleDiff returns the shortest angular distance between two ecliptic longitudes (0-180)
func AngleDiff(lon1, lon2 float64) float64 {
	diff := math.Abs(lon1 - lon2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// SignedAngleDiff returns signed angular difference (lon1 - lon2), normalized to [-180, 180)
func SignedAngleDiff(lon1, lon2 float64) float64 {
	d := lon1 - lon2
	d = math.Mod(d+180, 360)
	if d < 0 {
		d += 360
	}
	return d - 180
}

// Body represents a celestial body or point for aspect calculation
type Body struct {
	ID        string
	Longitude float64
	Speed     float64 // degrees/day, used for applying/separating
}

// FindAspects finds all aspects between two sets of bodies
func FindAspects(bodiesA, bodiesB []Body, orbs models.OrbConfig, sameSet bool) []models.AspectInfo {
	var aspects []models.AspectInfo

	for i, a := range bodiesA {
		startJ := 0
		if sameSet {
			startJ = i + 1
		}
		for j := startJ; j < len(bodiesB); j++ {
			b := bodiesB[j]
			if sameSet && a.ID == b.ID {
				continue
			}

			angle := AngleDiff(a.Longitude, b.Longitude)

			// Get aspect definitions from config
		aspectDefs := orbs.GetAspectDefs()

		for _, def := range aspectDefs {
				if !def.Enabled {
					continue
				}
				diff := math.Abs(angle - def.Angle)
				isApplying := computeApplying(a, b, def.Angle)

				// Use entering (applying) or exiting (separating) orb based on aspect direction
				var orb float64
				if isApplying {
					orb = def.EnteringOrb
				} else {
					orb = def.ExitingOrb
				}

				if diff <= orb {
					// Find the aspect type from StandardAspects for the name
					var aspectType models.AspectType
					for _, sa := range models.StandardAspects {
						if math.Abs(sa.Angle-def.Angle) < 1.0 {
							aspectType = sa.Type
							break
						}
					}
					if aspectType == "" {
						aspectType = models.AspectType(def.Name)
					}
					aspects = append(aspects, models.AspectInfo{
						PlanetA:     a.ID,
						PlanetB:     b.ID,
						AspectType:  aspectType,
						AspectAngle: def.Angle,
						ActualAngle: angle,
						Orb:         diff,
						IsApplying:  isApplying,
					})
					break // one aspect per pair
				}
			}
		}
	}

	return aspects
}

// FindCrossAspects finds aspects between inner and outer chart bodies
func FindCrossAspects(inner, outer []Body, orbs models.OrbConfig) []models.CrossAspectInfo {
	var aspects []models.CrossAspectInfo
	aspectDefs := orbs.GetAspectDefs()

	for _, a := range inner {
		for _, b := range outer {
			angle := AngleDiff(a.Longitude, b.Longitude)

			for _, def := range aspectDefs {
				if !def.Enabled {
					continue
				}
				diff := math.Abs(angle - def.Angle)
				isApplying := computeApplying(
					Body{ID: a.ID, Longitude: a.Longitude, Speed: a.Speed},
					Body{ID: b.ID, Longitude: b.Longitude, Speed: b.Speed},
					def.Angle,
				)

				// Use entering (applying) or exiting (separating) orb based on aspect direction
				var orb float64
				if isApplying {
					orb = def.EnteringOrb
				} else {
					orb = def.ExitingOrb
				}

				if diff <= orb {
					// Find the aspect type from StandardAspects for the name
					var aspectType models.AspectType
					for _, sa := range models.StandardAspects {
						if math.Abs(sa.Angle-def.Angle) < 1.0 {
							aspectType = sa.Type
							break
						}
					}
					if aspectType == "" {
						aspectType = models.AspectType(def.Name)
					}
					aspects = append(aspects, models.CrossAspectInfo{
						InnerBody:   a.ID,
						OuterBody:   b.ID,
						AspectType:  aspectType,
						AspectAngle: def.Angle,
						ActualAngle: angle,
						Orb:         diff,
						IsApplying:  isApplying,
					})
					break
				}
			}
		}
	}

	return aspects
}

// computeApplying determines if the aspect is applying (getting tighter) or separating
func computeApplying(a, b Body, targetAngle float64) bool {
	// Current angular separation
	currentSep := AngleDiff(a.Longitude, b.Longitude)
	currentOrb := math.Abs(currentSep - targetAngle)

	// Project forward a small amount (0.01 day)
	dt := 0.01
	futureA := a.Longitude + a.Speed*dt
	futureB := b.Longitude + b.Speed*dt
	futureSep := AngleDiff(futureA, futureB)
	futureOrb := math.Abs(futureSep - targetAngle)

	return futureOrb < currentOrb
}
