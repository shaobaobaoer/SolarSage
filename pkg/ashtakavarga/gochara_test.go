package ashtakavarga

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

// buildTestResult creates a minimal AshtakavargaResult for testing Gochara.
// Uses a known chart: all planets in Aries (0°) for simplicity.
func buildTestResult() *AshtakavargaResult {
	// All planets at 0° sidereal (Aries 0°) for predictable BAV.
	positions := []vedic.SiderealPosition{
		{SiderealLon: 0},   // Sun
		{SiderealLon: 0},   // Moon
		{SiderealLon: 0},   // Mars
		{SiderealLon: 0},   // Mercury
		{SiderealLon: 0},   // Jupiter
		{SiderealLon: 0},   // Venus
		{SiderealLon: 0},   // Saturn
	}
	for i, pid := range traditionalPlanets {
		positions[i].PlanetID = pid
		positions[i].PlanetPosition.PlanetID = pid
	}
	return CalcAshtakavarga(positions, 0)
}

func TestGocharaScore_ValidPlanet(t *testing.T) {
	result := buildTestResult()
	score := GocharaScore(result, models.PlanetSun, 0)
	if score < 0 || score > 8 {
		t.Errorf("GocharaScore(Sun, Aries) = %d, want 0-8", score)
	}
}

func TestGocharaScore_InvalidSign(t *testing.T) {
	result := buildTestResult()
	score := GocharaScore(result, models.PlanetSun, 13)
	if score != -1 {
		t.Errorf("GocharaScore with out-of-range sign = %d, want -1", score)
	}
}

func TestGocharaScore_UnknownPlanet(t *testing.T) {
	result := buildTestResult()
	score := GocharaScore(result, models.PlanetUranus, 0)
	if score != -1 {
		t.Errorf("GocharaScore(Uranus) = %d, want -1 (not a traditional planet)", score)
	}
}

func TestIsGocharaAuspicious(t *testing.T) {
	result := buildTestResult()
	// Check that the function doesn't panic and returns a bool
	for signIdx := 0; signIdx < 12; signIdx++ {
		auspicious := IsGocharaAuspicious(result, models.PlanetJupiter, signIdx)
		score := GocharaScore(result, models.PlanetJupiter, signIdx)
		expected := score >= gocharaThreshold
		if auspicious != expected {
			t.Errorf("IsGocharaAuspicious(Jupiter, sign %d) = %v, want %v (bindus=%d)",
				signIdx, auspicious, expected, score)
		}
	}
}

func TestGocharaForPlanet(t *testing.T) {
	result := buildTestResult()
	entries := GocharaForPlanet(result, models.PlanetSaturn)
	if len(entries) != 12 {
		t.Fatalf("GocharaForPlanet: got %d entries, want 12", len(entries))
	}
	for i, e := range entries {
		if e.SignIndex != i {
			t.Errorf("entries[%d].SignIndex = %d, want %d", i, e.SignIndex, i)
		}
		if e.SignName != signNames[i] {
			t.Errorf("entries[%d].SignName = %s, want %s", i, e.SignName, signNames[i])
		}
		if e.TransitPlanet != models.PlanetSaturn {
			t.Errorf("entries[%d].TransitPlanet = %s, want SATURN", i, e.TransitPlanet)
		}
		if e.Bindus < 0 || e.Bindus > 8 {
			t.Errorf("entries[%d].Bindus = %d out of valid range", i, e.Bindus)
		}
	}
}

func TestGocharaAll(t *testing.T) {
	result := buildTestResult()
	all := GocharaAll(result)
	if len(all) != len(traditionalPlanets) {
		t.Errorf("GocharaAll: got %d planets, want %d", len(all), len(traditionalPlanets))
	}
	for _, p := range traditionalPlanets {
		entries, ok := all[p]
		if !ok {
			t.Errorf("GocharaAll missing planet %s", p)
			continue
		}
		if len(entries) != 12 {
			t.Errorf("GocharaAll[%s]: got %d entries, want 12", p, len(entries))
		}
	}
}

func TestGocharaAtLongitude(t *testing.T) {
	result := buildTestResult()
	// Saturn at 45° sidereal = Taurus (sign index 1)
	entry := GocharaAtLongitude(result, models.PlanetSaturn, 45.0)
	if entry.SignIndex != 1 {
		t.Errorf("GocharaAtLongitude(Saturn, 45°) sign index = %d, want 1 (Taurus)", entry.SignIndex)
	}
	if entry.SignName != "Taurus" {
		t.Errorf("GocharaAtLongitude sign name = %s, want Taurus", entry.SignName)
	}
	expectedBindus := GocharaScore(result, models.PlanetSaturn, 1)
	if entry.Bindus != expectedBindus {
		t.Errorf("GocharaAtLongitude bindus = %d, want %d", entry.Bindus, expectedBindus)
	}
}

func TestAuspiciousSigns(t *testing.T) {
	result := buildTestResult()
	signs := AuspiciousSigns(result, models.PlanetJupiter)
	// Verify all returned signs actually have bindus >= 4
	for _, s := range signs {
		score := GocharaScore(result, models.PlanetJupiter, s)
		if score < gocharaThreshold {
			t.Errorf("AuspiciousSigns returned sign %d with only %d bindus", s, score)
		}
	}
	// Verify no auspicious sign was missed
	for i := 0; i < 12; i++ {
		score := GocharaScore(result, models.PlanetJupiter, i)
		if score >= gocharaThreshold {
			found := false
			for _, s := range signs {
				if s == i {
					found = true
				}
			}
			if !found {
				t.Errorf("AuspiciousSigns missed sign %d (score=%d)", i, score)
			}
		}
	}
}

func TestGocharaSAVConsistency(t *testing.T) {
	result := buildTestResult()
	// SAV for any sign must equal the sum of all 7 planet BAV bindus for that sign
	for signIdx := 0; signIdx < 12; signIdx++ {
		sum := 0
		for _, p := range traditionalPlanets {
			sum += GocharaScore(result, p, signIdx)
		}
		if result.SAV[signIdx] != sum {
			t.Errorf("SAV[%d] = %d, but sum of BAV = %d", signIdx, result.SAV[signIdx], sum)
		}
	}
}
