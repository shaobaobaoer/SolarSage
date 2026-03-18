package render

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func makeTestChart() *models.ChartInfo {
	return &models.ChartInfo{
		Planets: []models.PlanetPosition{
			{PlanetID: models.PlanetSun, Longitude: 280},
			{PlanetID: models.PlanetMoon, Longitude: 120},
			{PlanetID: models.PlanetMars, Longitude: 45},
		},
		Houses: []float64{15, 45, 75, 105, 135, 165, 195, 225, 255, 285, 315, 345},
		Angles: models.AnglesInfo{ASC: 15, MC: 285, DSC: 195, IC: 105},
		Aspects: []models.AspectInfo{
			{PlanetA: "SUN", PlanetB: "MOON", AspectType: models.AspectTrine},
		},
	}
}

func TestCalcChartWheel(t *testing.T) {
	chart := makeTestChart()
	wheel := CalcChartWheel(chart, 0.4)

	if len(wheel.Planets) != 3 {
		t.Errorf("Expected 3 planets, got %d", len(wheel.Planets))
	}
	if len(wheel.Houses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(wheel.Houses))
	}
	if len(wheel.Aspects) != 1 {
		t.Errorf("Expected 1 aspect, got %d", len(wheel.Aspects))
	}

	// All planet positions should be within the unit square
	for _, p := range wheel.Planets {
		if p.Position.X < 0 || p.Position.X > 1 || p.Position.Y < 0 || p.Position.Y > 1 {
			t.Errorf("Planet %s position out of bounds: (%.2f, %.2f)",
				p.PlanetID, p.Position.X, p.Position.Y)
		}
	}

	// Center should be (0.5, 0.5)
	if wheel.Center.X != 0.5 || wheel.Center.Y != 0.5 {
		t.Errorf("Center = (%.2f, %.2f), want (0.5, 0.5)", wheel.Center.X, wheel.Center.Y)
	}
}

func TestCalcChartWheel_DefaultRadius(t *testing.T) {
	chart := makeTestChart()
	wheel := CalcChartWheel(chart, 0)
	if wheel.Radius != 0.4 {
		t.Errorf("Default radius = %.2f, want 0.4", wheel.Radius)
	}
}

func TestLonToAngle(t *testing.T) {
	// ASC at 0°: the ASC should be at 180° on the wheel
	angle := lonToAngle(0, 0)
	if math.Abs(angle-180) > 0.01 {
		t.Errorf("lonToAngle(0, 0) = %.2f, want 180", angle)
	}

	// Planet at same longitude as ASC should be at 180°
	angle = lonToAngle(15, 15)
	if math.Abs(angle-180) > 0.01 {
		t.Errorf("lonToAngle(15, 15) = %.2f, want 180", angle)
	}
}

func TestPolarToCartesian(t *testing.T) {
	center := Point{0.5, 0.5}

	// 0° (right) should give (0.5+r, 0.5)
	pt := polarToCartesian(center, 0.4, 0)
	if math.Abs(pt.X-0.9) > 0.01 || math.Abs(pt.Y-0.5) > 0.01 {
		t.Errorf("0° = (%.2f, %.2f), want (0.9, 0.5)", pt.X, pt.Y)
	}

	// 90° (up in screen coords) should give (0.5, 0.1)
	pt = polarToCartesian(center, 0.4, 90)
	if math.Abs(pt.X-0.5) > 0.01 || math.Abs(pt.Y-0.1) > 0.01 {
		t.Errorf("90° = (%.2f, %.2f), want (0.5, 0.1)", pt.X, pt.Y)
	}
}

func TestCalcSignSegments(t *testing.T) {
	segments := CalcSignSegments(0, 0.4)
	if len(segments) != 12 {
		t.Fatalf("Expected 12 segments, got %d", len(segments))
	}

	// First sign should be Aries
	if segments[0].Sign != "Aries" {
		t.Errorf("First sign = %s, want Aries", segments[0].Sign)
	}

	// All midpoints should be in unit square
	for _, s := range segments {
		if s.MidPoint.X < 0 || s.MidPoint.X > 1 || s.MidPoint.Y < 0 || s.MidPoint.Y > 1 {
			t.Errorf("Sign %s midpoint out of bounds: (%.2f, %.2f)",
				s.Sign, s.MidPoint.X, s.MidPoint.Y)
		}
	}
}

func TestCalcChartWheel_HouseAngles(t *testing.T) {
	chart := makeTestChart()
	wheel := CalcChartWheel(chart, 0.4)

	// First house cusp angle should correspond to ASC
	if len(wheel.Houses) > 0 {
		ascAngle := lonToAngle(15, 15) // ASC at 15°
		if math.Abs(wheel.Houses[0].Angle-ascAngle) > 0.01 {
			t.Errorf("First house angle = %.2f, ASC angle = %.2f", wheel.Houses[0].Angle, ascAngle)
		}
	}
}
