package report

import (
	"os"
	"path/filepath"
	"testing"

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

func TestGenerateNatalReport(t *testing.T) {
	report, err := GenerateNatalReport(51.5074, -0.1278, j2000)
	if err != nil {
		t.Fatalf("GenerateNatalReport: %v", err)
	}

	// Chart
	if report.Chart == nil {
		t.Fatal("Chart is nil")
	}
	if len(report.Chart.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(report.Chart.Planets))
	}

	// Lunar phase
	if report.LunarPhase == nil {
		t.Error("LunarPhase is nil")
	}

	// Dignities
	if len(report.Dignities) != 10 {
		t.Errorf("Expected 10 dignities, got %d", len(report.Dignities))
	}

	// Sect
	if len(report.Sect) != 10 {
		t.Errorf("Expected 10 sect entries, got %d", len(report.Sect))
	}

	// Dispositors
	if report.Dispositors == nil {
		t.Error("Dispositors is nil")
	}
	if len(report.Dispositors.Dispositors) != 10 {
		t.Errorf("Expected 10 dispositor entries, got %d", len(report.Dispositors.Dispositors))
	}

	// Faces
	if len(report.Faces) != 10 {
		t.Errorf("Expected 10 faces, got %d", len(report.Faces))
	}

	// Lots
	if len(report.Lots) == 0 {
		t.Error("Expected at least some lots")
	}

	// Element balance
	total := 0
	for _, v := range report.ElementBalance {
		total += v
	}
	if total != 10 {
		t.Errorf("Element balance total = %d, want 10", total)
	}

	// Modality balance
	total = 0
	for _, v := range report.ModalityBalance {
		total += v
	}
	if total != 10 {
		t.Errorf("Modality balance total = %d, want 10", total)
	}

	// Hemisphere balance
	if len(report.HemisphereBalance) != 4 {
		t.Errorf("Expected 4 hemisphere entries, got %d", len(report.HemisphereBalance))
	}
}

func TestCalcElementBalance(t *testing.T) {
	// No need for Swiss Ephemeris
	planets := []models.PlanetPosition{
		{Sign: "Aries"},   // Fire
		{Sign: "Taurus"},  // Earth
		{Sign: "Gemini"},  // Air
		{Sign: "Cancer"},  // Water
		{Sign: "Leo"},     // Fire
	}
	elements := calcElementBalance(planets)
	if elements["Fire"] != 2 {
		t.Errorf("Fire = %d, want 2", elements["Fire"])
	}
	if elements["Earth"] != 1 {
		t.Errorf("Earth = %d, want 1", elements["Earth"])
	}
}

func TestCalcModalityBalance(t *testing.T) {
	planets := []models.PlanetPosition{
		{Sign: "Aries"},       // Cardinal
		{Sign: "Taurus"},      // Fixed
		{Sign: "Gemini"},      // Mutable
		{Sign: "Cancer"},      // Cardinal
		{Sign: "Sagittarius"}, // Mutable
	}
	mods := calcModalityBalance(planets)
	if mods["Cardinal"] != 2 {
		t.Errorf("Cardinal = %d, want 2", mods["Cardinal"])
	}
	if mods["Mutable"] != 2 {
		t.Errorf("Mutable = %d, want 2", mods["Mutable"])
	}
}
