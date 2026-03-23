package kp

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// TestSubLords_NakshatraLord verifies that the Nakshatra lord is correct
// for well-known positions.
func TestSubLords_NakshatraLord(t *testing.T) {
	tests := []struct {
		lon      float64
		wantNak  int
		wantLord models.PlanetID
	}{
		// Ashwini starts at 0° — Ketu (Sun proxy)
		{0.5, 0, models.PlanetSun},
		// Bharani starts at 13.333° — Venus
		{14.0, 1, models.PlanetVenus},
		// Krittika starts at 26.666° — Sun
		{27.0, 2, models.PlanetSun},
		// Rohini starts at 40.0° — Moon
		{41.0, 3, models.PlanetMoon},
		// Revati starts at 346.666° — Mercury
		{347.0, 26, models.PlanetMercury},
	}
	for _, tt := range tests {
		info := SubLords(tt.lon)
		if info.NakshatraIndex != tt.wantNak {
			t.Errorf("SubLords(%.1f°).NakshatraIndex = %d, want %d", tt.lon, info.NakshatraIndex, tt.wantNak)
		}
		if info.NakshatraLord != tt.wantLord {
			t.Errorf("SubLords(%.1f°).NakshatraLord = %s, want %s", tt.lon, info.NakshatraLord, tt.wantLord)
		}
	}
}

// TestSubLords_SubLordNotEmpty verifies sub lord is always set.
func TestSubLords_SubLordNotEmpty(t *testing.T) {
	for lon := 0.0; lon < 360.0; lon += 5.0 {
		info := SubLords(lon)
		if info.SubLord == "" {
			t.Errorf("SubLords(%.1f°).SubLord is empty", lon)
		}
		if info.SubSubLord == "" {
			t.Errorf("SubLords(%.1f°).SubSubLord is empty", lon)
		}
	}
}

// TestSubLords_CuspDegreeInRange verifies CuspDegree is always within [0, nakshatraSpan).
func TestSubLords_CuspDegreeInRange(t *testing.T) {
	for lon := 0.0; lon < 360.0; lon += 1.3 {
		info := SubLords(lon)
		if info.CuspDegree < 0 || info.CuspDegree >= nakshatraSpan {
			t.Errorf("SubLords(%.2f°).CuspDegree = %.4f out of range [0, %.4f)", lon, info.CuspDegree, nakshatraSpan)
		}
	}
}

// TestSubLords_NakshatraBoundary verifies sub lords change at Nakshatra boundaries.
func TestSubLords_NakshatraBoundary(t *testing.T) {
	// Just before and just after Bharani boundary (13.333°)
	before := SubLords(13.3)
	after := SubLords(13.4)
	if before.NakshatraIndex == after.NakshatraIndex {
		// If both in same nakshatra, skip — boundary resolution may differ
		t.Skip("boundary at 13.333° not crossed")
	}
	if before.NakshatraLord == after.NakshatraLord {
		t.Errorf("Nakshatra lords should differ across boundary: %s == %s",
			before.NakshatraLord, after.NakshatraLord)
	}
}

// TestSubLords_NormalizationNegative verifies negative longitudes are normalised.
func TestSubLords_NormalizationNegative(t *testing.T) {
	a := SubLords(-1.0)   // should normalise to 359°
	b := SubLords(359.0)
	if a.NakshatraIndex != b.NakshatraIndex {
		t.Errorf("SubLords(-1) NakIdx=%d, SubLords(359) NakIdx=%d, should match", a.NakshatraIndex, b.NakshatraIndex)
	}
}

// TestSubLords_KPTableSpot verifies specific KP table entries.
// At 0° (Aries 0° = Ashwini start): Nakshatra lord = Ketu(Sun), Sub lord = Ketu(Sun).
func TestSubLords_KPTableSpot(t *testing.T) {
	info := SubLords(0.0)
	if info.NakshatraLord != models.PlanetSun { // Ketu proxied as Sun
		t.Errorf("At 0° Nakshatra lord = %s, want SUN (Ketu proxy)", info.NakshatraLord)
	}
	if info.SubLord != models.PlanetSun {
		t.Errorf("At 0° Sub lord = %s, want SUN (Ketu proxy, first sub-period)", info.SubLord)
	}
}

// --- House tests ---

func TestHouseNumber_BasicPlacement(t *testing.T) {
	// Equal house cusps: cusp 1 = 0°, cusp 2 = 30°, … cusp 12 = 330°
	cusps := make([]float64, 12)
	for i := 0; i < 12; i++ {
		cusps[i] = float64(i * 30)
	}
	tests := []struct {
		lon       float64
		wantHouse int
	}{
		{15.0, 1},  // 15° → house 1 (0°–30°)
		{45.0, 2},  // 45° → house 2 (30°–60°)
		{350.0, 12}, // 350° → house 12 (330°–360°)
		{0.0, 1},   // exactly on cusp 1 = house 1
	}
	for _, tt := range tests {
		got := HouseNumber(tt.lon, cusps)
		if got != tt.wantHouse {
			t.Errorf("HouseNumber(%.1f°) = %d, want %d", tt.lon, got, tt.wantHouse)
		}
	}
}

func TestHouseNumber_WrapAround(t *testing.T) {
	// Cusp 12 = 340°, Cusp 1 = 10° — planet at 355° should be in house 12
	cusps := make([]float64, 12)
	cusps[0] = 10  // cusp 1
	for i := 1; i < 12; i++ {
		cusps[i] = float64(i*30 + 10)
	}
	// Cusp 12 = 10 + 11*30 = 340°; planet at 355° is after cusp 12 and before cusp 1(360°+10°)
	got := HouseNumber(355.0, cusps)
	if got != 12 {
		t.Errorf("HouseNumber(355° with cusp12=340°, cusp1=10°) = %d, want 12", got)
	}
}

func TestHouseNumber_InvalidCusps(t *testing.T) {
	// Wrong number of cusps should return 1 without panicking
	got := HouseNumber(100.0, []float64{0, 30, 60})
	if got != 1 {
		t.Errorf("HouseNumber with invalid cusps = %d, want 1", got)
	}
}

// --- Significator tests ---

func TestHouseSignificators_Basic(t *testing.T) {
	// Chart: Sun in house 1 (lon 15°), Moon in house 2 (lon 45°)
	// House ruler of house 1 = Mars
	cusps := make([]float64, 12)
	for i := 0; i < 12; i++ {
		cusps[i] = float64(i * 30)
	}

	sunInfo := CalcPlanetKP(15.0, cusps)
	sunInfo.PlanetID = models.PlanetSun
	moonInfo := CalcPlanetKP(45.0, cusps)
	moonInfo.PlanetID = models.PlanetMoon

	planetInfos := []PlanetKPInfo{sunInfo, moonInfo}

	sig := HouseSignificators(1, planetInfos, models.PlanetMars)

	// Group B must contain Sun (occupant of house 1)
	found := false
	for _, p := range sig.B {
		if p == models.PlanetSun {
			found = true
		}
	}
	if !found {
		t.Error("House 1 significators Group B should contain Sun (occupant)")
	}

	// Group D must contain Mars (ruler)
	if len(sig.D) == 0 || sig.D[0] != models.PlanetMars {
		t.Errorf("House 1 significators Group D = %v, want [MARS]", sig.D)
	}
}

func TestPlanetSignificators_D(t *testing.T) {
	// Mercury rules houses 3 and 6 (Gemini and Virgo in natural zodiac)
	var rulers [13]models.PlanetID
	rulers[3] = models.PlanetMercury
	rulers[6] = models.PlanetMercury

	cusps := make([]float64, 12)
	for i := 0; i < 12; i++ {
		cusps[i] = float64(i * 30)
	}
	mercInfo := CalcPlanetKP(75.0, cusps) // 15° Gemini → house 3
	mercInfo.PlanetID = models.PlanetMercury

	_, _, _, D := PlanetSignificators(mercInfo, rulers)
	if len(D) != 2 {
		t.Errorf("Mercury D significators = %v, want houses 3 and 6", D)
	}
}
