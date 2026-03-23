package vedic

import (
	"math"
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

func TestGetAyanamsa_Lahiri(t *testing.T) {
	aya, err := GetAyanamsa(j2000, AyanamsaLahiri)
	if err != nil {
		t.Fatalf("GetAyanamsa: %v", err)
	}
	// Lahiri at J2000.0 should be ~23.85°
	if aya < 23.5 || aya > 24.2 {
		t.Errorf("Lahiri ayanamsa = %.4f, expected ~23.85", aya)
	}
}

func TestGetAyanamsa_Precession(t *testing.T) {
	aya1, _ := GetAyanamsa(j2000, AyanamsaLahiri)
	aya2, _ := GetAyanamsa(j2000+365.25*100, AyanamsaLahiri)
	// 100 years later, ayanamsa should be ~1.4° more
	diff := aya2 - aya1
	if diff < 1.3 || diff > 1.5 {
		t.Errorf("100-year precession = %.4f, expected ~1.4", diff)
	}
}

func TestGetAyanamsa_InvalidSystem(t *testing.T) {
	_, err := GetAyanamsa(j2000, "INVALID")
	if err == nil {
		t.Error("Expected error for invalid ayanamsa system")
	}
}

func TestTropicalToSidereal(t *testing.T) {
	// Tropical 280° with ayanamsa 24° = sidereal 256°
	sid := TropicalToSidereal(280, 24)
	if math.Abs(sid-256) > 0.01 {
		t.Errorf("TropicalToSidereal(280, 24) = %.4f, want 256", sid)
	}

	// Wrapping: tropical 10° with ayanamsa 24° = 346°
	sid = TropicalToSidereal(10, 24)
	if math.Abs(sid-346) > 0.01 {
		t.Errorf("TropicalToSidereal(10, 24) = %.4f, want 346", sid)
	}
}

func TestCalcNakshatra(t *testing.T) {
	// 0° sidereal = Ashwini
	name, pada, lord := CalcNakshatra(0)
	if name != "Ashwini" {
		t.Errorf("0° nakshatra = %s, want Ashwini", name)
	}
	if pada != 1 {
		t.Errorf("0° pada = %d, want 1", pada)
	}
	_ = lord

	// 120° = Magha (10th nakshatra, starts at 120°)
	name, _, _ = CalcNakshatra(120)
	if name != "Magha" {
		t.Errorf("120° nakshatra = %s, want Magha", name)
	}

	// 350° = Revati (27th, starts at 346.667°)
	name, _, _ = CalcNakshatra(350)
	if name != "Revati" {
		t.Errorf("350° nakshatra = %s, want Revati", name)
	}
}

func TestCalcNakshatra_AllPadas(t *testing.T) {
	span := 360.0 / 27.0
	padaSpan := span / 4.0

	// Test each pada of Ashwini
	for p := 1; p <= 4; p++ {
		lon := float64(p-1) * padaSpan
		_, pada, _ := CalcNakshatra(lon + 0.1)
		if pada != p {
			t.Errorf("Ashwini pada at %.1f° = %d, want %d", lon+0.1, pada, p)
		}
	}
}

func TestCalcSiderealChart(t *testing.T) {
	sc, err := CalcSiderealChart(51.5074, -0.1278, j2000, AyanamsaLahiri)
	if err != nil {
		t.Fatalf("CalcSiderealChart: %v", err)
	}

	if len(sc.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(sc.Planets))
	}

	// All sidereal longitudes should differ from tropical by ~ayanamsa
	for _, p := range sc.Planets {
		diff := p.Longitude - p.SiderealLon
		if diff < 0 {
			diff += 360
		}
		if math.Abs(diff-sc.AyanamsaValue) > 0.1 {
			t.Errorf("%s: tropical-sidereal diff = %.4f, ayanamsa = %.4f",
				p.PlanetID, diff, sc.AyanamsaValue)
		}

		// Nakshatra should be non-empty
		if p.Nakshatra == "" {
			t.Errorf("%s has empty nakshatra", p.PlanetID)
		}
		if p.NakshatraPada < 1 || p.NakshatraPada > 4 {
			t.Errorf("%s pada = %d, out of range", p.PlanetID, p.NakshatraPada)
		}
	}
}

func TestCalcSiderealChart_DifferentAyanamsas(t *testing.T) {
	systems := []Ayanamsa{AyanamsaLahiri, AyanamsaRaman, AyanamsaFaganBradley}
	for _, sys := range systems {
		sc, err := CalcSiderealChart(51.5074, -0.1278, j2000, sys)
		if err != nil {
			t.Errorf("%s: %v", sys, err)
			continue
		}
		if sc.AyanamsaValue < 20 || sc.AyanamsaValue > 26 {
			t.Errorf("%s ayanamsa = %.4f, out of expected range", sys, sc.AyanamsaValue)
		}
	}
}

func TestCalcVimshottariDasha(t *testing.T) {
	periods := CalcVimshottariDasha(100) // Moon at 100° sidereal
	if len(periods) != len(dashaSequence) {
		t.Fatalf("Expected %d periods, got %d", len(dashaSequence), len(periods))
	}

	// First period should be partial (startAge = 0)
	if periods[0].StartAge != 0 {
		t.Errorf("First period startAge = %.2f, want 0", periods[0].StartAge)
	}

	// Total should be close to 120 years (Vimshottari full cycle)
	total := 0.0
	for _, p := range periods {
		total += p.Years
		if p.Lord == "" {
			t.Error("Period has empty lord")
		}
	}
	// Total won't be exactly 120 because first period is partial
	if total < 100 || total > 121 {
		t.Errorf("Total dasha years = %.2f, expected ~120", total)
	}
}

func TestCalcVimshottariDasha_MoonPosition(t *testing.T) {
	// Moon in Ashwini (0-13.333°) => lord is Sun (Ketu proxy)
	periods := CalcVimshottariDasha(5)
	if periods[0].Lord != models.PlanetSun {
		t.Errorf("First dasha lord for Moon at 5° = %s, expected SUN (Ketu)", periods[0].Lord)
	}
}

func TestCalcAshtottariDasha(t *testing.T) {
	periods := CalcAshtottariDasha(100) // Moon at 100° sidereal
	if len(periods) != len(ashtottariSequence) {
		t.Fatalf("Expected %d periods, got %d", len(ashtottariSequence), len(periods))
	}

	// First period starts at age 0
	if periods[0].StartAge != 0 {
		t.Errorf("First period startAge = %.2f, want 0", periods[0].StartAge)
	}

	// Each period has a valid lord
	for _, p := range periods {
		if p.Lord == "" {
			t.Error("Ashtottari period has empty lord")
		}
		if p.Years <= 0 {
			t.Errorf("Ashtottari period years = %.2f, must be > 0", p.Years)
		}
	}

	// Total should be close to 108 years (first period is partial so < 108)
	total := 0.0
	for _, p := range periods {
		total += p.Years
	}
	if total < 85 || total > 109 {
		t.Errorf("Ashtottari total years = %.2f, expected close to 108", total)
	}
}

func TestCalcAshtottariDasha_FullCycleAt0(t *testing.T) {
	// Moon exactly at 0° sidereal → start of Ashwini → Sun slot
	periods := CalcAshtottariDasha(0)
	if periods[0].Lord != models.PlanetSun {
		t.Errorf("Ashtottari at 0°: first lord = %s, want SUN", periods[0].Lord)
	}
	// First period should be full (fraction used = 0 → remaining = full 6 years)
	if math.Abs(periods[0].Years-6) > 0.01 {
		t.Errorf("Ashtottari at 0°: first period years = %.2f, want 6", periods[0].Years)
	}
}

func TestCalcYoginiDasha(t *testing.T) {
	periods := CalcYoginiDasha(100)
	if len(periods) != len(yoginiData) {
		t.Fatalf("Expected %d yogini periods, got %d", len(yoginiData), len(periods))
	}

	// First period starts at age 0
	if periods[0].StartAge != 0 {
		t.Errorf("First yogini period startAge = %.2f, want 0", periods[0].StartAge)
	}

	// Each period has a valid Yogini name and lord
	for _, p := range periods {
		if p.Yogini == "" {
			t.Error("Yogini period has empty name")
		}
		if p.Lord == "" {
			t.Error("Yogini period has empty lord")
		}
		if p.Years <= 0 {
			t.Errorf("Yogini period years = %.2f, must be > 0", p.Years)
		}
	}

	// Total should be close to 36 years (first period may be partial)
	total := 0.0
	for _, p := range periods {
		total += p.Years
	}
	if total < 28 || total > 37 {
		t.Errorf("Yogini total years = %.2f, expected close to 36", total)
	}
}

func TestCalcYoginiDasha_At0(t *testing.T) {
	// Moon at 0° → Ashwini → Mangala (Moon, 1 year)
	periods := CalcYoginiDasha(0)
	if periods[0].Yogini != YoginiMangala {
		t.Errorf("Yogini at 0°: first yogini = %s, want Mangala", periods[0].Yogini)
	}
	if periods[0].Lord != models.PlanetMoon {
		t.Errorf("Yogini at 0°: first lord = %s, want MOON", periods[0].Lord)
	}
	// Full period at start of nakshatra = 1 year
	if math.Abs(periods[0].Years-1) > 0.01 {
		t.Errorf("Yogini at 0°: first years = %.2f, want 1", periods[0].Years)
	}
}

func TestCalcYoginiDasha_SequenceOrder(t *testing.T) {
	// Verify each period's Yogini name matches the expected sequence
	periods := CalcYoginiDasha(0)
	expected := []YoginiName{
		YoginiMangala, YoginiPingala, YoginiBhramari, YoginiBhadrika,
		YoginiUlka, YoginiSiddha, YoginiSankata, YoginaDhanya,
	}
	for i, p := range periods {
		if p.Yogini != expected[i] {
			t.Errorf("Yogini[%d] = %s, want %s", i, p.Yogini, expected[i])
		}
	}
}

func TestCalcCharaDasha(t *testing.T) {
	// Aries lagna (odd sign), Mars in Gemini (sign 2)
	planetSigns := map[models.PlanetID]int{
		models.PlanetSun:     4,
		models.PlanetMoon:    1,
		models.PlanetMars:    2,
		models.PlanetMercury: 5,
		models.PlanetJupiter: 8,
		models.PlanetVenus:   6,
		models.PlanetSaturn:  9,
	}

	periods := CalcCharaDasha(0, planetSigns) // Aries lagna
	if len(periods) != 12 {
		t.Fatalf("expected 12 periods, got %d", len(periods))
	}

	// First period should be Aries (odd lagna → zodiacal order)
	if periods[0].Sign != "Aries" {
		t.Errorf("first period sign = %s, want Aries", periods[0].Sign)
	}
	if periods[0].StartAge != 0 {
		t.Errorf("first period startAge = %.2f, want 0", periods[0].StartAge)
	}

	// Aries ruler is Mars, Mars is in Gemini (sign 2).
	// Direct count from Aries(0) to Gemini(2) = 2 signs → 2 years
	// (Actually it's inclusive counting from sign 0 forward to sign 2: distance = 2)
	if periods[0].Years < 1 || periods[0].Years > 12 {
		t.Errorf("Aries dasha years = %d, out of range 1-12", periods[0].Years)
	}

	// Second period should be Taurus (zodiacal order for odd lagna)
	if periods[1].Sign != "Taurus" {
		t.Errorf("second period sign = %s, want Taurus", periods[1].Sign)
	}

	// All periods should have valid years
	totalYears := 0
	for _, p := range periods {
		if p.Years < 1 || p.Years > 12 {
			t.Errorf("sign %s has %d years, out of range 1-12", p.Sign, p.Years)
		}
		totalYears += p.Years
	}
	// Total should be between 12 and 144
	if totalYears < 12 || totalYears > 144 {
		t.Errorf("total years = %d, out of reasonable range", totalYears)
	}
}

func TestCalcCharaDasha_EvenLagna(t *testing.T) {
	// Taurus lagna (even sign) → reverse order
	planetSigns := map[models.PlanetID]int{
		models.PlanetSun:     4,
		models.PlanetMoon:    1,
		models.PlanetMars:    0,
		models.PlanetMercury: 2,
		models.PlanetJupiter: 8,
		models.PlanetVenus:   6,
		models.PlanetSaturn:  9,
	}

	periods := CalcCharaDasha(1, planetSigns) // Taurus lagna
	if len(periods) != 12 {
		t.Fatalf("expected 12 periods, got %d", len(periods))
	}

	// First should be Taurus
	if periods[0].Sign != "Taurus" {
		t.Errorf("first period = %s, want Taurus", periods[0].Sign)
	}
	// Second should be Aries (reverse order from Taurus)
	if periods[1].Sign != "Aries" {
		t.Errorf("second period = %s, want Aries (reverse order)", periods[1].Sign)
	}
}

func TestCalcCharaDasha_RulerInOwnSign(t *testing.T) {
	// Mars in Aries → ruler in own sign → 12 years
	planetSigns := map[models.PlanetID]int{
		models.PlanetMars: 0, // Mars in Aries
	}

	periods := CalcCharaDasha(0, planetSigns)
	if periods[0].Years != 12 {
		t.Errorf("ruler in own sign should give 12 years, got %d", periods[0].Years)
	}
}

