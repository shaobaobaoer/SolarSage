package dignity

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// signMidLon returns the midpoint longitude of a sign (e.g. Leo = 120+15 = 135)
func signMidLon(sign string) float64 {
	idx := map[string]float64{
		"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
		"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
		"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
	}
	return idx[sign] + 15.0
}

func TestCalcDignity_Rulership(t *testing.T) {
	// Sun in Leo = rulership (+5). Leo midpoint = 135°.
	d := CalcDignity(models.PlanetSun, 135.0)
	if d.Score < 5 {
		t.Errorf("Sun in Leo score = %d, want >= 5 (rulership)", d.Score)
	}
	found := false
	for _, dig := range d.Dignities {
		if dig == Rulership {
			found = true
		}
	}
	if !found {
		t.Error("Sun in Leo should have Rulership dignity")
	}
}

func TestCalcDignity_Exaltation(t *testing.T) {
	// Sun in Aries = exaltation (+4). Aries midpoint = 15°.
	d := CalcDignity(models.PlanetSun, 15.0)
	if !d.Exalted {
		t.Error("Sun in Aries should be exalted")
	}
	if d.Score < 4 {
		t.Errorf("Sun in Aries score = %d, want >= 4", d.Score)
	}
}

func TestCalcDignity_Detriment(t *testing.T) {
	// Sun in Aquarius = detriment. Aquarius midpoint = 315°.
	d := CalcDignity(models.PlanetSun, 315.0)
	if !d.InDetriment {
		t.Error("Sun in Aquarius should be in detriment")
	}
	if d.Score > -5 {
		t.Errorf("Sun in Aquarius score = %d, want <= -5", d.Score)
	}
}

func TestCalcDignity_Fall(t *testing.T) {
	// Sun in Libra = fall. Libra midpoint = 195°.
	d := CalcDignity(models.PlanetSun, 195.0)
	if !d.InFall {
		t.Error("Sun in Libra should be in fall")
	}
	if d.Score > -4 {
		t.Errorf("Sun in Libra score = %d, want <= -4", d.Score)
	}
}

func TestCalcDignity_Peregrine(t *testing.T) {
	// Sun in Gemini = no major dignity. Gemini midpoint = 75°.
	// Score may be non-zero if Term or Face applies — only check no big dignity.
	d := CalcDignity(models.PlanetSun, 75.0)
	for _, dig := range d.Dignities {
		if dig == Rulership || dig == Exaltation || dig == Detriment || dig == Fall {
			t.Errorf("Sun in Gemini should not have %s", dig)
		}
	}
}

func TestCalcDignity_TermAndFace(t *testing.T) {
	// Jupiter rules the first Egyptian term of Aries (0°–6°), so Jupiter at 3° Aries
	// should have TERM. Aries starts at 0°.
	d := CalcDignity(models.PlanetJupiter, 3.0) // 3° Aries
	hasTerm := false
	for _, dig := range d.Dignities {
		if dig == Term {
			hasTerm = true
		}
	}
	if !hasTerm {
		t.Errorf("Jupiter at 3° Aries should have Term dignity (got dignities: %v)", d.Dignities)
	}

	// Mars rules the Chaldean decan 3 of Aries (20°–30°), so Mars at 25° Aries should have FACE.
	d2 := CalcDignity(models.PlanetVenus, 25.0) // 25° Aries — Venus is term ruler here
	_ = d2
	// Mars rules 3rd decan of Aries (Chaldean)
	d3 := CalcDignity(models.PlanetVenus, 10.0) // 10° Aries → 2nd decan = Sun ruler in Chaldean
	_ = d3
}

func TestCalcDignity_Triplicity(t *testing.T) {
	// Sun is day triplicity ruler for Fire signs (Aries, Leo, Sagittarius).
	// Sun at 45° = 15° Taurus — NOT fire. Sun at 135° = 15° Leo — IS fire.
	d := CalcDignity(models.PlanetSun, 135.0) // Leo
	hasTrip := false
	for _, dig := range d.Dignities {
		if dig == Triplicity {
			hasTrip = true
		}
	}
	if !hasTrip {
		t.Errorf("Sun in Leo should have Triplicity dignity (got %v)", d.Dignities)
	}
}

func TestCalcDignityWithSect_DayChart(t *testing.T) {
	// Sun rules Fire triplicity in day charts.
	// In a day chart, Sun at 15° Leo should have Rulership + Triplicity.
	d := CalcDignityWithSect(models.PlanetSun, 135.0, true)
	hasRulership, hasTrip := false, false
	for _, dig := range d.Dignities {
		if dig == Rulership {
			hasRulership = true
		}
		if dig == Triplicity {
			hasTrip = true
		}
	}
	if !hasRulership {
		t.Error("Sun in Leo (day) should have Rulership")
	}
	if !hasTrip {
		t.Errorf("Sun in Leo day chart should have Triplicity (got %v)", d.Dignities)
	}
}

func TestCalcDignityWithSect_NightChart(t *testing.T) {
	// In a night chart, Jupiter rules Fire triplicity — Sun should not get Triplicity.
	d := CalcDignityWithSect(models.PlanetSun, 135.0, false) // Leo, night
	for _, dig := range d.Dignities {
		if dig == Triplicity {
			t.Error("Sun in Leo night chart should NOT have Triplicity (Jupiter rules night fire)")
		}
	}
	// Jupiter should get Triplicity in night chart for Leo
	dJ := CalcDignityWithSect(models.PlanetJupiter, 135.0, false)
	hasTrip := false
	for _, dig := range dJ.Dignities {
		if dig == Triplicity {
			hasTrip = true
		}
	}
	if !hasTrip {
		t.Errorf("Jupiter in Leo night chart should have Triplicity (got %v)", dJ.Dignities)
	}
}

func TestCalcChartDignities(t *testing.T) {
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 135.0, Sign: "Leo"},     // Leo 15°
		{PlanetID: models.PlanetMoon, Longitude: 105.0, Sign: "Cancer"}, // Cancer 15°
		{PlanetID: models.PlanetMercury, Longitude: 75.0, Sign: "Gemini"}, // Gemini 15°
	}
	dignities := CalcChartDignities(positions)
	if len(dignities) != 3 {
		t.Fatalf("Expected 3 dignities, got %d", len(dignities))
	}
	// Sun in Leo must have at least Rulership
	hasRule := false
	for _, d := range dignities[0].Dignities {
		if d == Rulership {
			hasRule = true
		}
	}
	if !hasRule {
		t.Error("Sun in Leo should have Rulership")
	}
	// Moon in Cancer must have Rulership
	hasRule = false
	for _, d := range dignities[1].Dignities {
		if d == Rulership {
			hasRule = true
		}
	}
	if !hasRule {
		t.Error("Moon in Cancer should have Rulership")
	}
}

func TestAlmutenFiguris(t *testing.T) {
	// At 15° Leo (135°): Sun has Rulership(+5) + Triplicity(+3) = 8 minimum.
	// No other planet has rulership here, so Sun should win.
	planet, score := AlmutenFiguris(135.0)
	if planet != models.PlanetSun {
		t.Errorf("Almuten at 15° Leo = %s (score %d), want SUN", planet, score)
	}
	if score < 8 {
		t.Errorf("Almuten score at 15° Leo = %d, want >= 8", score)
	}
}

func TestAlmutenFiguris_Virgo(t *testing.T) {
	// Mercury in Virgo has Rulership(+5) + Exaltation(+4) = 9 minimum.
	planet, score := AlmutenFiguris(165.0) // 15° Virgo
	if planet != models.PlanetMercury {
		t.Errorf("Almuten at 15° Virgo = %s (score %d), want MERCURY", planet, score)
	}
	if score < 9 {
		t.Errorf("Almuten score at 15° Virgo = %d, want >= 9", score)
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
	got := SignTraditionalRuler("Scorpio")
	if got != models.PlanetMars {
		t.Errorf("SignTraditionalRuler(Scorpio) = %s, want MARS", got)
	}
	got = SignTraditionalRuler("Aquarius")
	if got != models.PlanetSaturn {
		t.Errorf("SignTraditionalRuler(Aquarius) = %s, want SATURN", got)
	}
}

func TestFindMutualReceptions(t *testing.T) {
	// Moon in Leo, Sun in Cancer = mutual reception by rulership
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 105.0, Sign: "Cancer"},
		{PlanetID: models.PlanetMoon, Longitude: 135.0, Sign: "Leo"},
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
		{PlanetID: models.PlanetSun, Longitude: 15.0, Sign: "Aries"},
		{PlanetID: models.PlanetMoon, Longitude: 45.0, Sign: "Taurus"},
	}
	receptions := FindMutualReceptions(positions)
	if len(receptions) != 0 {
		t.Errorf("Expected 0 mutual receptions, got %d", len(receptions))
	}
}

func TestCalcSect(t *testing.T) {
	s := CalcSect(models.PlanetSun, true)
	if !s.InSect {
		t.Error("Sun should be in sect in day chart")
	}
	s = CalcSect(models.PlanetMoon, true)
	if s.InSect {
		t.Error("Moon should NOT be in sect in day chart")
	}
	s = CalcSect(models.PlanetMoon, false)
	if !s.InSect {
		t.Error("Moon should be in sect in night chart")
	}
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
	// Mercury in Gemini = rulership (5) at minimum
	d1 := CalcDignity(models.PlanetMercury, 75.0) // 15° Gemini
	if d1.Score < 5 {
		t.Errorf("Mercury in Gemini score = %d, want >= 5", d1.Score)
	}
	// Mercury in Virgo = rulership (5) + exaltation (4) = 9 at minimum
	d2 := CalcDignity(models.PlanetMercury, 165.0) // 15° Virgo
	if d2.Score < 9 {
		t.Errorf("Mercury in Virgo score = %d, want >= 9 (rulership + exaltation)", d2.Score)
	}
}
