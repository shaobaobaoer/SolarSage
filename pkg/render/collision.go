package render

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// CollisionRadius is the minimum angular separation (degrees) between two
// planet glyph centers before they are considered overlapping.
// Corresponds to roughly half a glyph diameter at normal chart sizes.
const CollisionRadius = 3.0

// maxNudgeIterations is the safety cap on the recursive collision resolution
// loop to prevent infinite looping when the chart is extremely crowded.
const maxNudgeIterations = 48

// ResolveCollisions takes a slice of PlanetGlyphs (each with an Angle in
// degrees on the chart wheel) and returns a new slice where no two glyph
// centers are closer than CollisionRadius degrees.
//
// Algorithm (ported from AstroChart utils.ts assemble/placePointsInCollision):
//  1. For each pair of glyphs, compute angular distance.
//  2. If distance < CollisionRadius, nudge the later glyph by +CollisionRadius.
//  3. Repeat until no collisions remain or the iteration cap is reached.
//
// The original (true ecliptic) angle is preserved in Angle; the display
// position after nudging is stored in DisplayAngle. Callers should draw
// a pointer line from DisplayAngle back to Angle when they differ.
func ResolveCollisions(glyphs []PlanetGlyph) []PlanetGlyph {
	if len(glyphs) == 0 {
		return glyphs
	}

	// Copy and initialise DisplayAngle = Angle
	result := make([]PlanetGlyph, len(glyphs))
	for i, g := range glyphs {
		result[i] = g
		result[i].DisplayAngle = g.Angle
		result[i].Displaced = false
	}

	for iter := 0; iter < maxNudgeIterations; iter++ {
		collisionFound := false
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				dist := angularDistance(result[i].DisplayAngle, result[j].DisplayAngle)
				if dist < CollisionRadius {
					collisionFound = true
					// Nudge j forward by CollisionRadius degrees
					result[j].DisplayAngle = normAngle(result[j].DisplayAngle + CollisionRadius)
					result[j].Displaced = result[j].DisplayAngle != result[j].Angle
				}
			}
		}
		if !collisionFound {
			break
		}
	}

	// Recalculate x/y positions from DisplayAngle
	center := Point{0.5, 0.5}
	for i := range result {
		radius := result[i].Ring
		result[i].Position = polarToCartesian(center, radius, result[i].DisplayAngle)
	}

	return result
}

// angularDistance returns the shortest arc distance between two angles (0–180).
func angularDistance(a, b float64) float64 {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// normAngle normalises an angle to [0, 360).
func normAngle(a float64) float64 {
	a = math.Mod(a, 360.0)
	if a < 0 {
		a += 360
	}
	return a
}

// PointerLine describes a line from a displaced glyph back to its true position.
type PointerLine struct {
	PlanetID    models.PlanetID `json:"planet_id"`
	TrueAngle   float64         `json:"true_angle"`    // original ecliptic angle
	DisplayAngle float64        `json:"display_angle"` // nudged display angle
	TruePoint   Point           `json:"true_point"`
	DisplayPoint Point          `json:"display_point"`
}

// CalcPointerLines returns pointer lines for all displaced planets after
// collision resolution. Pass the ring radius used for planet placement.
func CalcPointerLines(glyphs []PlanetGlyph, trueRing, displayRing float64) []PointerLine {
	center := Point{0.5, 0.5}
	var lines []PointerLine
	for _, g := range glyphs {
		if !g.Displaced {
			continue
		}
		lines = append(lines, PointerLine{
			PlanetID:     g.PlanetID,
			TrueAngle:    g.Angle,
			DisplayAngle: g.DisplayAngle,
			TruePoint:    polarToCartesian(center, trueRing, g.Angle),
			DisplayPoint: polarToCartesian(center, displayRing, g.DisplayAngle),
		})
	}
	return lines
}
