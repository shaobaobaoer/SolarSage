// Package render provides coordinate calculations for chart wheel visualization.
// It converts chart data into x/y positions suitable for SVG, Canvas, or any
// 2D rendering system. All coordinates are normalized to a unit circle (0-1 range).
package render

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Point represents a 2D coordinate
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// PlanetGlyph holds rendering data for a planet on the chart wheel
type PlanetGlyph struct {
	PlanetID models.PlanetID `json:"planet_id"`
	Position Point           `json:"position"`
	Angle    float64         `json:"angle"` // degrees from ASC (counter-clockwise)
	Ring     float64         `json:"ring"`  // 0=center, 1=outer
}

// HouseLine holds rendering data for a house cusp line
type HouseLine struct {
	House int   `json:"house"`
	Start Point `json:"start"`
	End   Point `json:"end"`
	Angle float64 `json:"angle"`
}

// AspectLine holds rendering data for an aspect line between two planets
type AspectLine struct {
	PlanetA    string `json:"planet_a"`
	PlanetB    string `json:"planet_b"`
	Start      Point  `json:"start"`
	End        Point  `json:"end"`
	AspectType models.AspectType `json:"aspect_type"`
}

// ChartWheel holds all rendering coordinates for a chart wheel
type ChartWheel struct {
	Planets  []PlanetGlyph `json:"planets"`
	Houses   []HouseLine   `json:"houses"`
	Aspects  []AspectLine  `json:"aspects"`
	ASCAngle float64       `json:"asc_angle"`
	Center   Point         `json:"center"`
	Radius   float64       `json:"radius"`
}

// CalcChartWheel generates rendering coordinates for a chart wheel.
// The ASC is placed at the left (9 o'clock position) and signs proceed counter-clockwise.
// radius controls the overall size (0.5 for unit square centering).
func CalcChartWheel(chartInfo *models.ChartInfo, radius float64) *ChartWheel {
	if radius <= 0 {
		radius = 0.4
	}

	center := Point{0.5, 0.5}
	ascLon := chartInfo.Angles.ASC

	wheel := &ChartWheel{
		Center:   center,
		Radius:   radius,
		ASCAngle: ascLon,
	}

	// Planet positions on the wheel
	planetRing := radius * 0.75 // planets at 75% of outer radius
	for _, p := range chartInfo.Planets {
		angle := lonToAngle(p.Longitude, ascLon)
		pt := polarToCartesian(center, planetRing, angle)
		wheel.Planets = append(wheel.Planets, PlanetGlyph{
			PlanetID: p.PlanetID,
			Position: pt,
			Angle:    angle,
			Ring:     0.75,
		})
	}

	// House cusp lines
	for i, cusp := range chartInfo.Houses {
		angle := lonToAngle(cusp, ascLon)
		inner := polarToCartesian(center, radius*0.55, angle)
		outer := polarToCartesian(center, radius, angle)
		wheel.Houses = append(wheel.Houses, HouseLine{
			House: i + 1,
			Start: inner,
			End:   outer,
			Angle: angle,
		})
	}

	// Aspect lines between planets
	aspectRing := radius * 0.5
	for _, a := range chartInfo.Aspects {
		aAngle := findPlanetAngle(a.PlanetA, chartInfo.Planets, ascLon)
		bAngle := findPlanetAngle(a.PlanetB, chartInfo.Planets, ascLon)
		start := polarToCartesian(center, aspectRing, aAngle)
		end := polarToCartesian(center, aspectRing, bAngle)
		wheel.Aspects = append(wheel.Aspects, AspectLine{
			PlanetA:    a.PlanetA,
			PlanetB:    a.PlanetB,
			Start:      start,
			End:        end,
			AspectType: a.AspectType,
		})
	}

	return wheel
}

// lonToAngle converts an ecliptic longitude to a chart wheel angle in degrees.
// ASC is at 180° (left/9 o'clock), signs proceed counter-clockwise.
func lonToAngle(lon, ascLon float64) float64 {
	// Offset so ASC is at 180° (left side)
	return math.Mod(180-(lon-ascLon)+360, 360)
}

// polarToCartesian converts polar coordinates to x/y
func polarToCartesian(center Point, radius, angleDeg float64) Point {
	rad := angleDeg * math.Pi / 180
	return Point{
		X: center.X + radius*math.Cos(rad),
		Y: center.Y - radius*math.Sin(rad), // Y inverted for screen coords
	}
}

func findPlanetAngle(id string, planets []models.PlanetPosition, ascLon float64) float64 {
	for _, p := range planets {
		if string(p.PlanetID) == id {
			return lonToAngle(p.Longitude, ascLon)
		}
	}
	return 0
}

// SignSegment holds the angular extent of a zodiac sign on the wheel
type SignSegment struct {
	Sign       string  `json:"sign"`
	StartAngle float64 `json:"start_angle"`
	EndAngle   float64 `json:"end_angle"`
	MidAngle   float64 `json:"mid_angle"`
	MidPoint   Point   `json:"mid_point"`
}

// CalcSignSegments returns the 12 zodiac sign segments for the chart wheel.
func CalcSignSegments(ascLon, radius float64) []SignSegment {
	center := Point{0.5, 0.5}
	labelRing := radius * 0.9
	segments := make([]SignSegment, 12)

	for i := 0; i < 12; i++ {
		signStart := float64(i) * 30.0
		signEnd := signStart + 30.0
		startAngle := lonToAngle(signStart, ascLon)
		endAngle := lonToAngle(signEnd, ascLon)
		midAngle := lonToAngle(signStart+15, ascLon)
		midPt := polarToCartesian(center, labelRing, midAngle)

		segments[i] = SignSegment{
			Sign:       models.ZodiacSigns[i],
			StartAngle: startAngle,
			EndAngle:   endAngle,
			MidAngle:   midAngle,
			MidPoint:   midPt,
		}
	}
	return segments
}
