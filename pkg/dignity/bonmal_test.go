package dignity

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcBonMal_Combustion(t *testing.T) {
	// Mercury at 100 degrees, Sun at 105 degrees = 5 degree separation (combust)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 105, Sign: "Cancer"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondCombust {
			found = true
			if m.Score != scoreCombust {
				t.Errorf("Combust score = %d, want %d", m.Score, scoreCombust)
			}
		}
	}
	if !found {
		t.Error("Expected combustion condition for Mercury 5 degrees from Sun")
	}
}

func TestCalcBonMal_CombustionMoonWiderOrb(t *testing.T) {
	// Moon at 100 degrees, Sun at 110 degrees = 10 degree separation
	// Regular combust orb (8.5) would miss, but Moon uses 12 degree orb
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMoon, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 110, Sign: "Cancer"},
	}
	info := CalcBonMal(models.PlanetMoon, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondCombust {
			found = true
		}
	}
	if !found {
		t.Error("Expected combustion for Moon 10 degrees from Sun (within 12 degree orb)")
	}
}

func TestCalcBonMal_UnderSunbeams(t *testing.T) {
	// Venus at 100, Sun at 115 = 15 degree separation (under sunbeams, not combust)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetVenus, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 115, Sign: "Leo"},
	}
	info := CalcBonMal(models.PlanetVenus, positions)

	foundSunbeams := false
	foundCombust := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondUnderSunbeams {
			foundSunbeams = true
		}
		if m.Condition == CondCombust {
			foundCombust = true
		}
	}
	if !foundSunbeams {
		t.Error("Expected under sunbeams for Venus 15 degrees from Sun")
	}
	if foundCombust {
		t.Error("Should NOT be combust at 15 degree separation")
	}
}

func TestCalcBonMal_SunNotCombustItself(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetMercury, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetSun, positions)
	for _, m := range info.Maltreatments {
		if m.Condition == CondCombust || m.Condition == CondUnderSunbeams {
			t.Error("Sun should not be combust or under sunbeams itself")
		}
	}
}

func TestCalcBonMal_BeneficTrine(t *testing.T) {
	// Target at 0 degrees, Venus at 120 degrees = exact trine
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 0, Sign: "Aries"},
		{PlanetID: models.PlanetVenus, Longitude: 120, Sign: "Leo"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"}, // far away Sun
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, b := range info.Bonifications {
		if b.Condition == CondBeneficTrine && b.Source == models.PlanetVenus {
			found = true
			if b.Score != scoreBeneficTrine {
				t.Errorf("Benefic trine score = %d, want %d", b.Score, scoreBeneficTrine)
			}
		}
	}
	if !found {
		t.Error("Expected benefic trine from Venus at 120 degree separation")
	}
}

func TestCalcBonMal_BeneficSextile(t *testing.T) {
	// Target at 10 degrees, Jupiter at 70 degrees = exact sextile
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 10, Sign: "Aries"},
		{PlanetID: models.PlanetJupiter, Longitude: 70, Sign: "Gemini"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, b := range info.Bonifications {
		if b.Condition == CondBeneficSextile && b.Source == models.PlanetJupiter {
			found = true
		}
	}
	if !found {
		t.Error("Expected benefic sextile from Jupiter at 60 degree separation")
	}
}

func TestCalcBonMal_BeneficConjunction(t *testing.T) {
	// Target at 50 degrees, Jupiter at 53 degrees = conjunction within orb
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 50, Sign: "Taurus"},
		{PlanetID: models.PlanetJupiter, Longitude: 53, Sign: "Taurus"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, b := range info.Bonifications {
		if b.Condition == CondBeneficConjunction && b.Source == models.PlanetJupiter {
			found = true
		}
	}
	if !found {
		t.Error("Expected benefic conjunction from Jupiter at 3 degree separation")
	}
}

func TestCalcBonMal_MaleficSquare(t *testing.T) {
	// Target at 0 degrees, Mars at 90 degrees = exact square
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 0, Sign: "Aries"},
		{PlanetID: models.PlanetMars, Longitude: 90, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondMaleficSquare && m.Source == models.PlanetMars {
			found = true
			if m.Score != scoreMaleficSquare {
				t.Errorf("Malefic square score = %d, want %d", m.Score, scoreMaleficSquare)
			}
		}
	}
	if !found {
		t.Error("Expected malefic square from Mars at 90 degree separation")
	}
}

func TestCalcBonMal_MaleficOpposition(t *testing.T) {
	// Target at 10, Saturn at 190 = exact opposition
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 10, Sign: "Aries"},
		{PlanetID: models.PlanetSaturn, Longitude: 190, Sign: "Libra"},
		{PlanetID: models.PlanetSun, Longitude: 300, Sign: "Aquarius"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondMaleficOpposition && m.Source == models.PlanetSaturn {
			found = true
		}
	}
	if !found {
		t.Error("Expected malefic opposition from Saturn at 180 degree separation")
	}
}

func TestCalcBonMal_Besieged(t *testing.T) {
	// Mercury at 100, Mars at 95, Saturn at 105 = besieged
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetMars, Longitude: 95, Sign: "Cancer"},
		{PlanetID: models.PlanetSaturn, Longitude: 105, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 250, Sign: "Sagittarius"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondBesieged {
			found = true
			if m.Score != scoreBesieged {
				t.Errorf("Besieged score = %d, want %d", m.Score, scoreBesieged)
			}
		}
	}
	if !found {
		t.Error("Expected besiegement with Mercury between Mars and Saturn")
	}
}

func TestCalcBonMal_NotBesiegedTooFar(t *testing.T) {
	// Mercury at 100, Mars at 80, Saturn at 120 = both > 8 degree orb
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetMars, Longitude: 80, Sign: "Gemini"},
		{PlanetID: models.PlanetSaturn, Longitude: 120, Sign: "Leo"},
		{PlanetID: models.PlanetSun, Longitude: 250, Sign: "Sagittarius"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	for _, m := range info.Maltreatments {
		if m.Condition == CondBesieged {
			t.Error("Should NOT be besieged when malefics are > 8 degrees away")
		}
	}
}

func TestCalcBonMal_InBeneficSign(t *testing.T) {
	// Mercury in Taurus (ruled by Venus = benefic sign)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 40, Sign: "Taurus"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, b := range info.Bonifications {
		if b.Condition == CondInBeneficSign {
			found = true
			if b.Score != scoreInBeneficSign {
				t.Errorf("In benefic sign score = %d, want %d", b.Score, scoreInBeneficSign)
			}
		}
	}
	if !found {
		t.Error("Expected in-benefic-sign for Mercury in Taurus")
	}
}

func TestCalcBonMal_InMaleficSign(t *testing.T) {
	// Mercury in Aries (ruled by Mars = malefic sign)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 5, Sign: "Aries"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondInMaleficSign {
			found = true
		}
	}
	if !found {
		t.Error("Expected in-malefic-sign for Mercury in Aries")
	}
}

func TestCalcBonMal_NetScore(t *testing.T) {
	// Mercury at 0 Aries with Venus trine (120) and Mars square (90)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 0, Sign: "Aries"},
		{PlanetID: models.PlanetVenus, Longitude: 120, Sign: "Leo"},
		{PlanetID: models.PlanetMars, Longitude: 90, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 200, Sign: "Libra"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	// Expected: benefic trine +3, malefic square -3, in malefic sign (Aries) -1 = -1
	expectedNet := scoreBeneficTrine + scoreMaleficSquare + scoreInMaleficSign
	if info.NetScore != expectedNet {
		t.Errorf("NetScore = %d, want %d", info.NetScore, expectedNet)
	}
}

func TestCalcChartBonMal(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 120, Sign: "Leo"},
		{PlanetID: models.PlanetMoon, Longitude: 60, Sign: "Gemini"},
		{PlanetID: models.PlanetMercury, Longitude: 115, Sign: "Leo"},
	}
	results := CalcChartBonMal(positions)
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}
	if results[0].PlanetID != models.PlanetSun {
		t.Errorf("First result planet = %s, want SUN", results[0].PlanetID)
	}
	if results[2].PlanetID != models.PlanetMercury {
		t.Errorf("Third result planet = %s, want MERCURY", results[2].PlanetID)
	}
}

func TestCalcBonMal_TargetNotFound(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 100, Sign: "Cancer"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)
	if len(info.Bonifications) != 0 || len(info.Maltreatments) != 0 {
		t.Error("Expected empty results when target not in positions")
	}
}

func TestAngleDiff(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{0, 120, 120},
		{350, 10, 20},
		{10, 350, 20},
		{180, 0, 180},
		{0, 0, 0},
		{90, 270, 180},
	}
	for _, tc := range tests {
		got := angleDiff(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("angleDiff(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestCalcBonMal_MaleficConjunction(t *testing.T) {
	// Mercury at 100, Saturn at 103 = conjunction
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetMercury, Longitude: 100, Sign: "Cancer"},
		{PlanetID: models.PlanetSaturn, Longitude: 103, Sign: "Cancer"},
		{PlanetID: models.PlanetSun, Longitude: 250, Sign: "Sagittarius"},
	}
	info := CalcBonMal(models.PlanetMercury, positions)

	found := false
	for _, m := range info.Maltreatments {
		if m.Condition == CondMaleficConjunction && m.Source == models.PlanetSaturn {
			found = true
			if m.Score != scoreMaleficConjunction {
				t.Errorf("Malefic conjunction score = %d, want %d", m.Score, scoreMaleficConjunction)
			}
		}
	}
	if !found {
		t.Error("Expected malefic conjunction from Saturn at 3 degree separation")
	}
}
