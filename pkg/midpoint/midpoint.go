package midpoint

import (
	"math"
	"sort"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// MidpointEntry represents a single midpoint between two bodies
type MidpointEntry struct {
	BodyA     string  `json:"body_a"`
	BodyB     string  `json:"body_b"`
	Longitude float64 `json:"longitude"`
	Sign      string  `json:"sign"`
	SignDeg   float64 `json:"sign_degree"`
}

// MidpointAspect represents a planet on a midpoint axis
type MidpointAspect struct {
	Planet      string  `json:"planet"`
	PlanetLon   float64 `json:"planet_longitude"`
	MidpointOf  [2]string `json:"midpoint_of"`
	MidpointLon float64 `json:"midpoint_longitude"`
	Orb         float64 `json:"orb"`
	HarmonicDiv int     `json:"harmonic_div"` // 1=conjunction/opp, 2=square, 4=semi-square/sesqui
}

// MidpointTree holds all midpoints and planet-on-midpoint activations
type MidpointTree struct {
	Midpoints     []MidpointEntry  `json:"midpoints"`
	Activations   []MidpointAspect `json:"activations"`
	SortedDial90  []DialEntry      `json:"sorted_dial_90"`
}

// DialEntry is a point on the 90-degree Cosmobiology sort
type DialEntry struct {
	ID       string  `json:"id"`
	Dial90   float64 `json:"dial_90"`
	Original float64 `json:"original_longitude"`
}

// CalcMidpoints computes all midpoints and finds activations
func CalcMidpoints(positions []models.PlanetPosition, orb float64) *MidpointTree {
	if orb <= 0 {
		orb = 1.5 // default midpoint orb
	}

	tree := &MidpointTree{}

	// Compute all midpoints
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			mp := midpoint(positions[i].Longitude, positions[j].Longitude)
			tree.Midpoints = append(tree.Midpoints, MidpointEntry{
				BodyA:     string(positions[i].PlanetID),
				BodyB:     string(positions[j].PlanetID),
				Longitude: mp,
				Sign:      models.SignFromLongitude(mp),
				SignDeg:   models.SignDegreeFromLongitude(mp),
			})
		}
	}

	// Check each planet against each midpoint for activations
	for _, p := range positions {
		for _, mp := range tree.Midpoints {
			// Skip midpoints involving this planet
			if mp.BodyA == string(p.PlanetID) || mp.BodyB == string(p.PlanetID) {
				continue
			}
			// Check hard aspects to midpoint: conjunction/opposition (div=1), square (div=2), semi-square (div=4)
			for _, div := range []int{1, 2, 4} {
				diff := harmonicDiff(p.Longitude, mp.Longitude, div)
				if diff <= orb {
					tree.Activations = append(tree.Activations, MidpointAspect{
						Planet:      string(p.PlanetID),
						PlanetLon:   p.Longitude,
						MidpointOf:  [2]string{mp.BodyA, mp.BodyB},
						MidpointLon: mp.Longitude,
						Orb:         diff,
						HarmonicDiv: div,
					})
				}
			}
		}
	}

	// Build 90° dial sort (Cosmobiology)
	for _, p := range positions {
		tree.SortedDial90 = append(tree.SortedDial90, DialEntry{
			ID:       string(p.PlanetID),
			Dial90:   math.Mod(p.Longitude, 90),
			Original: p.Longitude,
		})
	}
	sort.Slice(tree.SortedDial90, func(i, j int) bool {
		return tree.SortedDial90[i].Dial90 < tree.SortedDial90[j].Dial90
	})

	return tree
}

// midpoint returns the shorter arc midpoint
func midpoint(lon1, lon2 float64) float64 {
	lon1 = sweph.NormalizeDegrees(lon1)
	lon2 = sweph.NormalizeDegrees(lon2)
	diff := lon2 - lon1
	if diff < 0 {
		diff += 360
	}
	if diff <= 180 {
		return sweph.NormalizeDegrees(lon1 + diff/2)
	}
	return sweph.NormalizeDegrees(lon2 + (360-diff)/2)
}

// harmonicDiff checks the angular difference on a specific harmonic
// div=1: 0°/180° axis, div=2: 90° axis, div=4: 45° axis
func harmonicDiff(lon, mpLon float64, div int) float64 {
	angle := math.Abs(lon - mpLon)
	if angle > 180 {
		angle = 360 - angle
	}
	modAngle := math.Mod(angle, 180.0/float64(div))
	if modAngle > 90.0/float64(div) {
		modAngle = 180.0/float64(div) - modAngle
	}
	return modAngle
}
