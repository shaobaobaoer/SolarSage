package yoga

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

// makePos creates a synthetic sidereal position.
func makePos(id models.PlanetID, siderealLon float64) vedic.SiderealPosition {
	return vedic.SiderealPosition{
		PlanetPosition: models.PlanetPosition{
			PlanetID:  id,
			Longitude: siderealLon, // use same value for simplicity
			Sign:      models.SignFromLongitude(siderealLon),
		},
		SiderealLon:  siderealLon,
		SiderealSign: models.SignFromLongitude(siderealLon),
	}
}

func TestGajakesariYoga_Conjunction(t *testing.T) {
	// Jupiter conjunct Moon (same sign) -> kendra (house 1 from Moon)
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMoon, 100),    // Cancer area
		makePos(models.PlanetJupiter, 105), // Same sign
		makePos(models.PlanetSun, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Gajakesari Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Gajakesari Yoga when Jupiter conjoins Moon")
	}
}

func TestGajakesariYoga_Opposition(t *testing.T) {
	// Jupiter in 7th sign from Moon (180 degrees away) -> kendra
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMoon, 10),      // Aries (sign 0)
		makePos(models.PlanetJupiter, 190),  // Libra (sign 6) -> 7th from Moon
		makePos(models.PlanetSun, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Gajakesari Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Gajakesari Yoga when Jupiter is in 7th from Moon")
	}
}

func TestGajakesariYoga_NoYoga(t *testing.T) {
	// Jupiter in 3rd sign from Moon -> NOT a kendra
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMoon, 10),     // Aries (sign 0)
		makePos(models.PlanetJupiter, 70),  // Gemini (sign 2) -> 3rd from Moon
		makePos(models.PlanetSun, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	result := AnalyzeYogas(positions, nil, 0)
	for _, y := range result.Yogas {
		if y.Name == "Gajakesari Yoga" {
			t.Error("Should NOT detect Gajakesari Yoga when Jupiter is in 3rd from Moon")
		}
	}
}

func TestBudhadityaYoga(t *testing.T) {
	// Sun and Mercury conjunct
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMercury, 105), // Within 10 degrees
		makePos(models.PlanetMoon, 200),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetJupiter, 50),
		makePos(models.PlanetVenus, 150),
		makePos(models.PlanetSaturn, 250),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Budhaditya Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Budhaditya Yoga when Sun conjoins Mercury")
	}
}

func TestBudhadityaYoga_NotPresent(t *testing.T) {
	// Sun and Mercury far apart
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMercury, 200), // 100 degrees apart
		makePos(models.PlanetMoon, 50),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetJupiter, 250),
		makePos(models.PlanetVenus, 150),
		makePos(models.PlanetSaturn, 350),
	}

	result := AnalyzeYogas(positions, nil, 0)
	for _, y := range result.Yogas {
		if y.Name == "Budhaditya Yoga" {
			t.Error("Should NOT detect Budhaditya Yoga when Sun and Mercury are far apart")
		}
	}
}

func TestChandraMangalaYoga(t *testing.T) {
	// Moon and Mars conjunct
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMoon, 150),
		makePos(models.PlanetMars, 155), // Within 10 degrees
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMercury, 200),
		makePos(models.PlanetJupiter, 250),
		makePos(models.PlanetVenus, 300),
		makePos(models.PlanetSaturn, 350),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Chandra-Mangala Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Chandra-Mangala Yoga when Moon conjoins Mars")
	}
}

func TestMahapurushaYoga_Ruchaka(t *testing.T) {
	// Mars in Aries (own sign) in 1st house (kendra) -> Ruchaka Yoga
	// ASC at 0 Aries, Mars at 15 Aries
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMars, 15),     // Aries (own sign)
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMoon, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetJupiter, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	result := AnalyzeYogas(positions, nil, 0) // ASC at 0 Aries
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Ruchaka Yoga" {
			found = true
			if y.Category != YogaMahapurusha {
				t.Errorf("Ruchaka category = %s, want MAHAPURUSHA", y.Category)
			}
			break
		}
	}
	if !found {
		t.Error("Expected Ruchaka Yoga when Mars is in Aries in a kendra")
	}
}

func TestMahapurushaYoga_Hamsa(t *testing.T) {
	// Jupiter in Sagittarius (own sign) in 10th house (kendra)
	// ASC at 0 Pisces (index 11), Jupiter at 260 (Sagittarius, index 8)
	// House = (8 - 11 + 12) % 12 + 1 = 10
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetJupiter, 260), // Sagittarius (own sign)
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMoon, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	ascLon := 330.0 // Pisces
	result := AnalyzeYogas(positions, nil, ascLon)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Hamsa Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Hamsa Yoga when Jupiter is in Sagittarius in a kendra")
	}
}

func TestMahapurushaYoga_NotInKendra(t *testing.T) {
	// Mars in Aries (own sign) but in 2nd house (not kendra)
	// ASC at 0 Pisces (sign 11), Mars at 15 Aries (sign 0) -> house 2
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMars, 15), // Aries, house 2 from Pisces
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMoon, 200),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetJupiter, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	ascLon := 330.0 // Pisces
	result := AnalyzeYogas(positions, nil, ascLon)
	for _, y := range result.Yogas {
		if y.Name == "Ruchaka Yoga" {
			t.Error("Should NOT detect Ruchaka when Mars is not in a kendra")
		}
	}
}

func TestRajaYoga_KendraTrikona(t *testing.T) {
	// ASC at 0 Aries. Lord of 1st = Mars, lord of 5th = Sun (Leo).
	// Mars (kendra lord) conjoins Sun (trikona lord).
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetMars, 200),
		makePos(models.PlanetSun, 205), // Conjunct Mars
		makePos(models.PlanetMoon, 100),
		makePos(models.PlanetMercury, 150),
		makePos(models.PlanetJupiter, 300),
		makePos(models.PlanetVenus, 350),
		makePos(models.PlanetSaturn, 50),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Raja Yoga" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Raja Yoga when kendra and trikona lords conjoin")
	}
}

func TestDhanaYoga_Lords2And11(t *testing.T) {
	// ASC at 0 Aries.
	// 2nd house = Taurus -> lord = Venus
	// 11th house = Aquarius -> lord = Saturn
	// Venus and Saturn conjunct
	positions := []vedic.SiderealPosition{
		makePos(models.PlanetVenus, 180),
		makePos(models.PlanetSaturn, 185), // Conjunct Venus
		makePos(models.PlanetSun, 100),
		makePos(models.PlanetMoon, 50),
		makePos(models.PlanetMars, 300),
		makePos(models.PlanetMercury, 250),
		makePos(models.PlanetJupiter, 350),
	}

	result := AnalyzeYogas(positions, nil, 0)
	found := false
	for _, y := range result.Yogas {
		if y.Name == "Dhana Yoga" && y.Category == YogaDhana {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected Dhana Yoga when lords of 2nd and 11th conjoin")
	}
}

func TestAnalyzeYogas_WithRealChart(t *testing.T) {
	sc, err := vedic.CalcSiderealChart(51.5074, -0.1278, 2451545.0, vedic.AyanamsaLahiri)
	if err != nil {
		t.Fatalf("CalcSiderealChart: %v", err)
	}

	ascLon := vedic.TropicalToSidereal(sc.Angles.ASC, sc.AyanamsaValue)
	result := AnalyzeYogas(sc.Planets, sc.Houses, ascLon)

	// Should return a non-nil analysis
	if result == nil {
		t.Fatal("AnalyzeYogas returned nil")
	}

	// Log found yogas for inspection
	for _, y := range result.Yogas {
		t.Logf("Found: %s (%s) - %s [%s]", y.Name, y.Category, y.Description, y.Strength)
	}
}

func TestHouseOf(t *testing.T) {
	tests := []struct {
		lon    float64
		ascLon float64
		want   int
	}{
		{0, 0, 1},      // Same sign as ASC = house 1
		{30, 0, 2},     // Next sign = house 2
		{90, 0, 4},     // 4th sign = house 4
		{270, 0, 10},   // 10th sign = house 10
		{0, 30, 12},    // One sign before ASC = house 12
	}
	for _, tt := range tests {
		got := houseOf(tt.lon, tt.ascLon)
		if got != tt.want {
			t.Errorf("houseOf(%.0f, %.0f) = %d, want %d", tt.lon, tt.ascLon, got, tt.want)
		}
	}
}

func TestIsKendra(t *testing.T) {
	kendras := []int{1, 4, 7, 10}
	for _, h := range kendras {
		if !isKendra(h) {
			t.Errorf("isKendra(%d) = false, want true", h)
		}
	}
	nonKendras := []int{2, 3, 5, 6, 8, 9, 11, 12}
	for _, h := range nonKendras {
		if isKendra(h) {
			t.Errorf("isKendra(%d) = true, want false", h)
		}
	}
}

func TestIsTrikona(t *testing.T) {
	trikonas := []int{1, 5, 9}
	for _, h := range trikonas {
		if !isTrikona(h) {
			t.Errorf("isTrikona(%d) = false, want true", h)
		}
	}
}

func TestConjunction(t *testing.T) {
	if !conjunction(100, 105) {
		t.Error("conjunction(100, 105) should be true")
	}
	if conjunction(100, 120) {
		t.Error("conjunction(100, 120) should be false")
	}
	// Wrap-around: 355 and 3 are 8 degrees apart
	if !conjunction(355, 3) {
		t.Error("conjunction(355, 3) should be true (wrap-around)")
	}
}
