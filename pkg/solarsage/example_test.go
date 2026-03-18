package solarsage_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/solarsage"
)

func init() {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	solarsage.Init(ephePath)
}

func ExampleNatalChart() {
	chart, err := solarsage.NatalChart(51.5074, -0.1278, "1990-06-15T14:30:00Z")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, p := range chart.Planets {
		fmt.Println(p)
	}
	// Output is dynamic based on ephemeris data
}

func ExampleMoonPhase() {
	phase, err := solarsage.MoonPhase("2025-03-18T12:00:00Z")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("%s (%.0f%% illuminated)\n", phase.PhaseName, phase.Illumination*100)
	// Output is dynamic
}

func ExamplePlanetPosition() {
	pos, err := solarsage.PlanetPosition("Venus", "2025-01-01T00:00:00Z")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(pos)
	// Output is dynamic
}

func ExampleCompatibility() {
	score, err := solarsage.Compatibility(
		51.5074, -0.1278, "1990-06-15T14:30:00Z",
		40.7128, -74.006, "1992-03-22T08:00:00Z",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("Compatibility: %.0f%%\n", score.Compatibility)
	// Output is dynamic
}

func ExampleParseDatetime() {
	jd, err := solarsage.ParseDatetime("2000-01-01T12:00:00Z")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("JD: %.1f\n", jd)
	// Output:
	// JD: 2451545.0
}

func ExampleParsePlanet() {
	pid, err := solarsage.ParsePlanet("Venus")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(pid)
	// Output:
	// VENUS
}
