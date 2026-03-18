package dignity

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcDignity_Rulership(t *testing.T) {
	// Sun in Leo = rulership
	d := CalcDignity(models.PlanetSun, "Leo")
	if d.Score != 5 {
		t.Errorf("Sun in Leo score = %d, want 5", d.Score)
	}
	if len(d.Dignities) == 0 || d.Dignities[0] != Rulership {
		t.Error("Sun in Leo should have Rulership dignity")
	}
}

func TestCalcDignity_Exaltation(t *testing.T) {
	// Sun in Aries = exaltation
	d := CalcDignity(models.PlanetSun, "Aries")
	if !d.Exalted {
		t.Error("Sun in Aries should be exalted")
	}
	if d.Score != 4 {
		t.Errorf("Sun in Aries score = %d, want 4", d.Score)
	}
}

func TestCalcDignity_Detriment(t *testing.T) {
	// Sun in Aquarius = detriment
	d := CalcDignity(models.PlanetSun, "Aquarius")
	if !d.InDetriment {
		t.Error("Sun in Aquarius should be in detriment")
	}
	if d.Score != -5 {
		t.Errorf("Sun in Aquarius score = %d, want -5", d.Score)
	}
}

func TestCalcDignity_Fall(t *testing.T) {
	// Sun in Libra = fall
	d := CalcDignity(models.PlanetSun, "Libra")
	if !d.InFall {
		t.Error("Sun in Libra should be in fall")
	}
	if d.Score != -4 {
		t.Errorf("Sun in Libra score = %d, want -4", d.Score)
	}
}

func TestCalcDignity_Peregrine(t *testing.T) {
	// Sun in Gemini = no essential dignity
	d := CalcDignity(models.PlanetSun, "Gemini")
	if d.Score != 0 {
		t.Errorf("Sun in Gemini score = %d, want 0", d.Score)
	}
	if len(d.Dignities) != 0 {
		t.Error("Sun in Gemini should have no dignities")
	}
}

func TestCalcChartDignities(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Leo"},
		{PlanetID: models.PlanetMoon, Sign: "Cancer"},
		{PlanetID: models.PlanetMercury, Sign: "Gemini"},
	}
	dignities := CalcChartDignities(positions)
	if len(dignities) != 3 {
		t.Fatalf("Expected 3 dignities, got %d", len(dignities))
	}
	// Sun in Leo = rulership (5)
	if dignities[0].Score != 5 {
		t.Errorf("Sun in Leo: score = %d, want 5", dignities[0].Score)
	}
	// Moon in Cancer = rulership (5)
	if dignities[1].Score != 5 {
		t.Errorf("Moon in Cancer: score = %d, want 5", dignities[1].Score)
	}
	// Mercury in Gemini = rulership (5)
	if dignities[2].Score != 5 {
		t.Errorf("Mercury in Gemini: score = %d, want 5", dignities[2].Score)
	}
}

func TestSignRuler(t *testing.T) {
	tests := map[string]models.PlanetID{
		"Aries":  models.PlanetMars,
		"Taurus": models.PlanetVenus,
		"Leo":    models.PlanetSun,
		"Pisces": models.PlanetNeptune,
	}
	for sign, expected := range tests {
		got := SignRuler(sign)
		if got != expected {
			t.Errorf("SignRuler(%s) = %s, want %s", sign, got, expected)
		}
	}
}

func TestSignTraditionalRuler(t *testing.T) {
	// Scorpio traditional ruler = Mars (not Pluto)
	got := SignTraditionalRuler("Scorpio")
	if got != models.PlanetMars {
		t.Errorf("SignTraditionalRuler(Scorpio) = %s, want MARS", got)
	}
	// Aquarius traditional ruler = Saturn
	got = SignTraditionalRuler("Aquarius")
	if got != models.PlanetSaturn {
		t.Errorf("SignTraditionalRuler(Aquarius) = %s, want SATURN", got)
	}
}

func TestFindMutualReceptions(t *testing.T) {
	// Moon in Leo, Sun in Cancer = mutual reception by rulership
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Cancer"},
		{PlanetID: models.PlanetMoon, Sign: "Leo"},
	}
	receptions := FindMutualReceptions(positions)
	if len(receptions) != 1 {
		t.Fatalf("Expected 1 mutual reception, got %d", len(receptions))
	}
	if receptions[0].Type != "rulership" {
		t.Errorf("Type = %s, want rulership", receptions[0].Type)
	}
}

func TestFindMutualReceptions_None(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Aries"},
		{PlanetID: models.PlanetMoon, Sign: "Taurus"},
	}
	receptions := FindMutualReceptions(positions)
	if len(receptions) != 0 {
		t.Errorf("Expected 0 mutual receptions, got %d", len(receptions))
	}
}

func TestCalcSect(t *testing.T) {
	// Sun in day chart = in sect
	s := CalcSect(models.PlanetSun, true)
	if !s.InSect {
		t.Error("Sun should be in sect in day chart")
	}

	// Moon in day chart = out of sect
	s = CalcSect(models.PlanetMoon, true)
	if s.InSect {
		t.Error("Moon should NOT be in sect in day chart")
	}

	// Moon in night chart = in sect
	s = CalcSect(models.PlanetMoon, false)
	if !s.InSect {
		t.Error("Moon should be in sect in night chart")
	}

	// Mercury is always in sect
	s = CalcSect(models.PlanetMercury, true)
	if !s.InSect {
		t.Error("Mercury should always be in sect")
	}
	s = CalcSect(models.PlanetMercury, false)
	if !s.InSect {
		t.Error("Mercury should always be in sect")
	}
}

func TestCalcDignity_MercuryDualRulership(t *testing.T) {
	// Mercury rules both Gemini and Virgo
	d1 := CalcDignity(models.PlanetMercury, "Gemini")
	d2 := CalcDignity(models.PlanetMercury, "Virgo")
	if d1.Score != 5 {
		t.Errorf("Mercury in Gemini score = %d, want 5", d1.Score)
	}
	// Mercury in Virgo: rulership (5) + exaltation (4) = 9
	if d2.Score != 9 {
		t.Errorf("Mercury in Virgo score = %d, want 9 (rulership + exaltation)", d2.Score)
	}
}
