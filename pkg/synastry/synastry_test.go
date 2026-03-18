package synastry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

const j2000 = 2451545.0

func TestCalcSynastryFromCharts(t *testing.T) {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}
	orbs := models.DefaultOrbConfig()

	chart1, err := chart.CalcSingleChart(51.5074, -0.1278, j2000, planets, orbs, models.HousePlacidus)
	if err != nil {
		t.Fatalf("chart1: %v", err)
	}
	chart2, err := chart.CalcSingleChart(40.7128, -74.006, j2000+365, planets, orbs, models.HousePlacidus)
	if err != nil {
		t.Fatalf("chart2: %v", err)
	}

	score := CalcSynastryFromCharts(chart1.Planets, chart2.Planets, orbs)
	if score == nil {
		t.Fatal("Score is nil")
	}

	if score.Compatibility < 0 || score.Compatibility > 100 {
		t.Errorf("Compatibility out of range: %.2f", score.Compatibility)
	}
	if score.Harmony < 0 {
		t.Errorf("Harmony should be non-negative: %.2f", score.Harmony)
	}
	if score.Tension < 0 {
		t.Errorf("Tension should be non-negative: %.2f", score.Tension)
	}
	if len(score.TopAspects) == 0 {
		t.Error("Expected at least some top aspects")
	}
	if len(score.CategoryScores) == 0 {
		t.Error("Expected at least some category scores")
	}

	t.Logf("Synastry: compatibility=%.1f%%, harmony=%.2f, tension=%.2f, total=%.2f",
		score.Compatibility, score.Harmony, score.Tension, score.TotalScore)
	for _, c := range score.CategoryScores {
		t.Logf("  %s: %.2f (%d aspects)", c.Category, c.Score, c.Count)
	}
}

func TestCalcSynastryScore_SamePerson(t *testing.T) {
	// Same chart with itself should have high compatibility
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetVenus, models.PlanetMars,
	}
	orbs := models.DefaultOrbConfig()

	chart1, _ := chart.CalcSingleChart(51.5074, -0.1278, j2000, planets, orbs, models.HousePlacidus)
	score := CalcSynastryFromCharts(chart1.Planets, chart1.Planets, orbs)

	// Same chart should have very high harmony (mostly conjunctions)
	if score.Harmony <= 0 {
		t.Errorf("Same chart harmony = %.2f, expected positive", score.Harmony)
	}
}

func TestCalcSynastryScore_EmptyAspects(t *testing.T) {
	// Empty input should return neutral score
	score := CalcSynastryScore(nil)
	if score.Compatibility != 50 {
		t.Errorf("Empty synastry compatibility = %.2f, want 50", score.Compatibility)
	}
}
