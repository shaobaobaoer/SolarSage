package symbolic

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath := filepath.Join("..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	code := m.Run()
	sweph.Close()
	os.Exit(code)
}

// natalJD: 1990-Jun-15 00:30 UT
var natalJD = sweph.JulDay(1990, 6, 15, 0.5, true)

const (
	geoLat = 51.5074 // London
	geoLon = -0.1278
)

func baseInput(method DirectionMethod, age float64) SymbolicInput {
	return SymbolicInput{
		NatalJD:     natalJD,
		GeoLat:      geoLat,
		GeoLon:      geoLon,
		Age:         age,
		Method:      method,
		HouseSystem: models.HousePlacidus,
		OrbConfig:   models.DefaultOrbConfig(),
	}
}

func TestOneDegreePerYear(t *testing.T) {
	input := baseInput(MethodOneDegree, 10)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if result.Method != MethodOneDegree {
		t.Errorf("Method = %s, want ONE_DEGREE", result.Method)
	}
	if result.Rate != 1.0 {
		t.Errorf("Rate = %f, want 1.0", result.Rate)
	}
	if result.Age != 10 {
		t.Errorf("Age = %f, want 10", result.Age)
	}

	// Each planet should advance exactly 10 degrees
	for _, d := range result.Directions {
		expectedLon := sweph.NormalizeDegrees(d.NatalLon + 10.0)
		if math.Abs(d.DirectedLon-expectedLon) > 0.0001 {
			t.Errorf("%s: DirectedLon = %f, want %f (natal=%f + 10)",
				d.PlanetID, d.DirectedLon, expectedLon, d.NatalLon)
		}
		if math.Abs(d.ArcApplied-10.0) > 0.0001 {
			t.Errorf("%s: ArcApplied = %f, want 10.0", d.PlanetID, d.ArcApplied)
		}
	}
}

func TestNaibodArc(t *testing.T) {
	input := baseInput(MethodNaibod, 10)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if result.Rate != 0.98556 {
		t.Errorf("Rate = %f, want 0.98556", result.Rate)
	}

	expectedArc := 10 * 0.98556 // 9.8556
	for _, d := range result.Directions {
		if math.Abs(d.ArcApplied-expectedArc) > 0.0001 {
			t.Errorf("%s: ArcApplied = %f, want %f", d.PlanetID, d.ArcApplied, expectedArc)
		}
		expectedLon := sweph.NormalizeDegrees(d.NatalLon + expectedArc)
		if math.Abs(d.DirectedLon-expectedLon) > 0.0001 {
			t.Errorf("%s: DirectedLon = %f, want %f",
				d.PlanetID, d.DirectedLon, expectedLon)
		}
	}
}

func TestProfectionArc(t *testing.T) {
	input := baseInput(MethodProfection, 1)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if result.Rate != 30.0 {
		t.Errorf("Rate = %f, want 30.0", result.Rate)
	}

	// At age 1, each planet advances 30 degrees (one sign)
	for _, d := range result.Directions {
		if math.Abs(d.ArcApplied-30.0) > 0.0001 {
			t.Errorf("%s: ArcApplied = %f, want 30.0", d.PlanetID, d.ArcApplied)
		}
		expectedLon := sweph.NormalizeDegrees(d.NatalLon + 30.0)
		if math.Abs(d.DirectedLon-expectedLon) > 0.0001 {
			t.Errorf("%s: DirectedLon = %f, want %f",
				d.PlanetID, d.DirectedLon, expectedLon)
		}
	}
}

func TestCustomRate(t *testing.T) {
	input := baseInput(MethodCustom, 5)
	input.CustomRate = 2.5 // 2.5 deg/year
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if result.Rate != 2.5 {
		t.Errorf("Rate = %f, want 2.5", result.Rate)
	}

	expectedArc := 5 * 2.5 // 12.5
	for _, d := range result.Directions {
		if math.Abs(d.ArcApplied-expectedArc) > 0.0001 {
			t.Errorf("%s: ArcApplied = %f, want %f", d.PlanetID, d.ArcApplied, expectedArc)
		}
	}
}

func TestCustomRateInvalid(t *testing.T) {
	input := baseInput(MethodCustom, 5)
	input.CustomRate = 0 // invalid
	_, err := CalcSymbolicDirections(input)
	if err == nil {
		t.Fatal("expected error for zero custom rate")
	}

	input.CustomRate = -1 // invalid
	_, err = CalcSymbolicDirections(input)
	if err == nil {
		t.Fatal("expected error for negative custom rate")
	}
}

func TestNegativeAge(t *testing.T) {
	input := baseInput(MethodOneDegree, -5)
	_, err := CalcSymbolicDirections(input)
	if err == nil {
		t.Fatal("expected error for negative age")
	}
}

func TestDefaultPlanets(t *testing.T) {
	input := SymbolicInput{
		NatalJD: natalJD,
		GeoLat:  geoLat,
		GeoLon:  geoLon,
		Age:     10,
		Method:  MethodOneDegree,
		// Planets, OrbConfig, HouseSystem all empty/zero -> use defaults
	}
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections with defaults: %v", err)
	}

	if len(result.Directions) != 10 {
		t.Errorf("expected 10 default planets, got %d", len(result.Directions))
	}
}

func TestDirectedAngles(t *testing.T) {
	input := baseInput(MethodOneDegree, 15)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if len(result.Angles) != 4 {
		t.Fatalf("expected 4 directed angles (ASC,MC,DSC,IC), got %d", len(result.Angles))
	}

	for _, a := range result.Angles {
		expectedLon := sweph.NormalizeDegrees(a.NatalLon + 15.0)
		if math.Abs(a.DirectedLon-expectedLon) > 0.0001 {
			t.Errorf("%s: DirectedLon = %f, want %f (natal=%f + 15)",
				a.PointID, a.DirectedLon, expectedLon, a.NatalLon)
		}
		if a.DirectedSign == "" {
			t.Errorf("%s: DirectedSign is empty", a.PointID)
		}
		if math.Abs(a.ArcApplied-15.0) > 0.0001 {
			t.Errorf("%s: ArcApplied = %f, want 15.0", a.PointID, a.ArcApplied)
		}
	}

	// Verify all four angle IDs are present
	ids := map[string]bool{}
	for _, a := range result.Angles {
		ids[a.PointID] = true
	}
	for _, expected := range []string{"ASC", "MC", "DSC", "IC"} {
		if !ids[expected] {
			t.Errorf("missing directed angle %s", expected)
		}
	}
}

func TestAspectDetection(t *testing.T) {
	// At age 0, directed = natal, so all directed planets are conjunct their natal counterpart
	input := baseInput(MethodOneDegree, 0)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	// We should find at least one conjunction (Dir_SUN conjunct SUN, etc.)
	conjCount := 0
	for _, a := range result.Aspects {
		if a.AspectType == models.AspectConjunction && a.Orb < 0.001 {
			conjCount++
		}
	}
	if conjCount < 10 {
		t.Errorf("at age 0 expected at least 10 exact conjunctions (directed=natal), got %d", conjCount)
	}
}

func TestAspectDetectionNonZeroAge(t *testing.T) {
	// At age 90 with one-degree, arc=90 which is a square.
	// Directed planet should be square to its natal position.
	input := baseInput(MethodOneDegree, 90)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	// Every directed planet is exactly 90 degrees from its natal position,
	// so we should find squares between Dir_X and X.
	squareCount := 0
	for _, a := range result.Aspects {
		if a.AspectType == models.AspectSquare && a.Orb < 0.001 {
			squareCount++
		}
	}
	if squareCount < 10 {
		t.Errorf("at age 90 expected at least 10 exact squares, got %d", squareCount)
	}
}

func TestSignCalculation(t *testing.T) {
	input := baseInput(MethodOneDegree, 10)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	for _, d := range result.Directions {
		expectedSign := models.SignFromLongitude(d.DirectedLon)
		if d.DirectedSign != expectedSign {
			t.Errorf("%s: DirectedSign = %s, want %s (lon=%f)",
				d.PlanetID, d.DirectedSign, expectedSign, d.DirectedLon)
		}
		expectedDeg := models.SignDegreeFromLongitude(d.DirectedLon)
		if math.Abs(d.DirectedDeg-expectedDeg) > 0.0001 {
			t.Errorf("%s: DirectedDeg = %f, want %f",
				d.PlanetID, d.DirectedDeg, expectedDeg)
		}
	}
}

func TestZeroAge(t *testing.T) {
	input := baseInput(MethodOneDegree, 0)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	// At age 0, directed = natal
	for _, d := range result.Directions {
		if math.Abs(d.DirectedLon-d.NatalLon) > 0.0001 {
			t.Errorf("%s: at age 0, DirectedLon=%f should equal NatalLon=%f",
				d.PlanetID, d.DirectedLon, d.NatalLon)
		}
	}
}

func TestMethodRate(t *testing.T) {
	tests := []struct {
		method DirectionMethod
		want   float64
	}{
		{MethodOneDegree, 1.0},
		{MethodNaibod, 0.98556},
		{MethodProfection, 30.0},
		{MethodCustom, 1.0}, // default fallback for custom
		{"UNKNOWN", 1.0},    // unknown falls through to default
	}

	for _, tt := range tests {
		got := tt.method.Rate()
		if math.Abs(got-tt.want) > 0.00001 {
			t.Errorf("%s.Rate() = %f, want %f", tt.method, got, tt.want)
		}
	}
}

func TestWrapAround360(t *testing.T) {
	// Profection at age 12 => arc = 360, should wrap to same position
	input := baseInput(MethodProfection, 12)
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	for _, d := range result.Directions {
		if math.Abs(d.DirectedLon-d.NatalLon) > 0.001 {
			t.Errorf("%s: at 360 arc, DirectedLon=%f should equal NatalLon=%f",
				d.PlanetID, d.DirectedLon, d.NatalLon)
		}
	}
}

func TestSpecificPlanets(t *testing.T) {
	input := baseInput(MethodOneDegree, 10)
	input.Planets = []models.PlanetID{models.PlanetSun, models.PlanetMoon}
	result, err := CalcSymbolicDirections(input)
	if err != nil {
		t.Fatalf("CalcSymbolicDirections: %v", err)
	}

	if len(result.Directions) != 2 {
		t.Errorf("expected 2 directions, got %d", len(result.Directions))
	}

	ids := map[models.PlanetID]bool{}
	for _, d := range result.Directions {
		ids[d.PlanetID] = true
	}
	if !ids[models.PlanetSun] || !ids[models.PlanetMoon] {
		t.Errorf("expected SUN and MOON in directions, got %v", ids)
	}
}
