package divisional

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

func TestMain(m *testing.M) {
	ephePath, _ := filepath.Abs("../../third_party/swisseph/ephe")
	sweph.Init(ephePath)
	defer sweph.Close()
	os.Exit(m.Run())
}

func TestCalcVargaPosition_D1(t *testing.T) {
	// D1 (division=1) should return the same sign position.
	lon := 45.0 // 15 Taurus sidereal
	result := CalcVargaPosition(lon, 1)
	if math.Abs(result-lon) > 0.001 {
		t.Errorf("D1: got %.4f, want %.4f", result, lon)
	}
}

func TestCalcNavamsaPosition_Aries(t *testing.T) {
	// A planet at 0 Aries sidereal (fire sign) -> Navamsa starts from Aries.
	// First pada (0-3.333) -> sign 0 (Aries).
	result := CalcNavamsaPosition(0.5) // 0d30' Aries
	resultSign := int(result / 30.0)
	if resultSign != 0 { // Aries
		t.Errorf("Navamsa of 0.5 Aries: sign index = %d, want 0 (Aries)", resultSign)
	}
}

func TestCalcNavamsaPosition_AriesLast(t *testing.T) {
	// Last pada of Aries (26.667-30) -> 9th subdivision -> sign index = (0+8)%12 = 8 (Sagittarius)
	result := CalcNavamsaPosition(29.0)
	resultSign := int(result / 30.0)
	if resultSign != 8 { // Sagittarius
		t.Errorf("Navamsa of 29 Aries: sign index = %d, want 8 (Sagittarius)", resultSign)
	}
}

func TestCalcNavamsaPosition_Taurus(t *testing.T) {
	// Taurus (earth sign, index 1) -> Navamsa starts from Capricorn (9).
	// First pada of Taurus (30.0-33.333) -> sign 9 (Capricorn).
	result := CalcNavamsaPosition(31.0)
	resultSign := int(result / 30.0)
	if resultSign != 9 { // Capricorn
		t.Errorf("Navamsa of 31 (1 Taurus): sign index = %d, want 9 (Capricorn)", resultSign)
	}
}

func TestCalcNavamsaPosition_Gemini(t *testing.T) {
	// Gemini (air sign, index 2) -> Navamsa starts from Libra (6).
	// First pada of Gemini (60-63.333) -> sign 6 (Libra).
	result := CalcNavamsaPosition(61.0)
	resultSign := int(result / 30.0)
	if resultSign != 6 { // Libra
		t.Errorf("Navamsa of 61 (1 Gemini): sign index = %d, want 6 (Libra)", resultSign)
	}
}

func TestCalcNavamsaPosition_Cancer(t *testing.T) {
	// Cancer (water sign, index 3) -> Navamsa starts from Cancer (3).
	// First pada of Cancer (90-93.333) -> sign 3 (Cancer).
	result := CalcNavamsaPosition(91.0)
	resultSign := int(result / 30.0)
	if resultSign != 3 { // Cancer
		t.Errorf("Navamsa of 91 (1 Cancer): sign index = %d, want 3 (Cancer)", resultSign)
	}
}

func TestCalcNavamsaPosition_Leo(t *testing.T) {
	// Leo (fire sign, index 4) -> Navamsa starts from Aries (0).
	// First pada of Leo (120-123.333) -> sign 0 (Aries).
	result := CalcNavamsaPosition(121.0)
	resultSign := int(result / 30.0)
	if resultSign != 0 { // Aries
		t.Errorf("Navamsa of 121 (1 Leo): sign index = %d, want 0 (Aries)", resultSign)
	}
}

func TestCalcVargaPosition_D12(t *testing.T) {
	// D12 (Dwadasamsa): 30/12 = 2.5 degrees per part.
	// 0 Aries (signIdx=0, partIdx=0): vargaSign = (0*12 + 0) % 12 = 0 (Aries)
	result := CalcVargaPosition(1.0, 12)
	resultSign := int(result / 30.0)
	if resultSign != 0 {
		t.Errorf("D12 of 1 Aries: sign index = %d, want 0 (Aries)", resultSign)
	}

	// 15 Aries (partIdx=6): vargaSign = (0*12 + 6) % 12 = 6 (Libra)
	result = CalcVargaPosition(15.0, 12)
	resultSign = int(result / 30.0)
	if resultSign != 6 {
		t.Errorf("D12 of 15 Aries: sign index = %d, want 6 (Libra)", resultSign)
	}
}

func TestCalcVargaPosition_D2(t *testing.T) {
	// D2 (Hora): 30/2 = 15 degrees per part.
	// 10 Aries (signIdx=0, partIdx=0): vargaSign = (0*2+0)%12 = 0 (Aries)
	result := CalcVargaPosition(10.0, 2)
	resultSign := int(result / 30.0)
	if resultSign != 0 {
		t.Errorf("D2 of 10 Aries: sign index = %d, want 0", resultSign)
	}

	// 20 Aries (signIdx=0, partIdx=1): vargaSign = (0*2+1)%12 = 1 (Taurus)
	result = CalcVargaPosition(20.0, 2)
	resultSign = int(result / 30.0)
	if resultSign != 1 {
		t.Errorf("D2 of 20 Aries: sign index = %d, want 1 (Taurus)", resultSign)
	}
}

func TestCalcVargaPosition_ZeroDivision(t *testing.T) {
	// Division 0 should return the original longitude.
	result := CalcVargaPosition(100.0, 0)
	if math.Abs(result-100.0) > 0.001 {
		t.Errorf("Division 0: got %.4f, want 100.0", result)
	}
}

func TestCalcVargaPosition_Normalization(t *testing.T) {
	// Negative longitude should be normalized.
	result := CalcVargaPosition(-10.0, 9)
	if result < 0 || result >= 360 {
		t.Errorf("Normalization: got %.4f, expected 0-360 range", result)
	}
}

func TestCalcDivisionalChart_Navamsa(t *testing.T) {
	// London, J2000.0
	dc, err := CalcDivisionalChart(51.5074, -0.1278, 2451545.0, VargaNavamsa, vedic.AyanamsaLahiri)
	if err != nil {
		t.Fatalf("CalcDivisionalChart: %v", err)
	}

	if dc.Varga != VargaNavamsa {
		t.Errorf("Varga = %s, want D9", dc.Varga)
	}
	if len(dc.Positions) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(dc.Positions))
	}

	for _, p := range dc.Positions {
		if p.VargaSign == "" {
			t.Errorf("%s has empty VargaSign", p.PlanetID)
		}
		if p.VargaLon < 0 || p.VargaLon >= 360 {
			t.Errorf("%s VargaLon = %.4f, out of range", p.PlanetID, p.VargaLon)
		}
		if p.VargaDegree < 0 || p.VargaDegree >= 30 {
			t.Errorf("%s VargaDegree = %.4f, out of range", p.PlanetID, p.VargaDegree)
		}
	}
}

func TestCalcDivisionalChart_AllVargas(t *testing.T) {
	vargas := []VargaType{
		VargaRasi, VargaHora, VargaDrekkana, VargaChaturthamsa,
		VargaSaptamsa, VargaNavamsa, VargaDasamsa, VargaDwadasamsa,
		VargaShodasamsa, VargaVimsamsa, VargaSiddhamsa, VargaSaptavimsamsa,
		VargaTrimsamsa, VargaKhavedamsa, VargaAkshavedamsa, VargaShashtiamsa,
	}
	for _, v := range vargas {
		dc, err := CalcDivisionalChart(51.5074, -0.1278, 2451545.0, v, vedic.AyanamsaLahiri)
		if err != nil {
			t.Errorf("%s: %v", v, err)
			continue
		}
		if len(dc.Positions) != 10 {
			t.Errorf("%s: expected 10 planets, got %d", v, len(dc.Positions))
		}
	}
}

func TestCalcDivisionalChart_InvalidVarga(t *testing.T) {
	_, err := CalcDivisionalChart(51.5074, -0.1278, 2451545.0, "D99", vedic.AyanamsaLahiri)
	if err == nil {
		t.Error("Expected error for invalid varga type")
	}
}

func TestCalcDivisionalChart_InvalidAyanamsa(t *testing.T) {
	_, err := CalcDivisionalChart(51.5074, -0.1278, 2451545.0, VargaNavamsa, "INVALID")
	if err == nil {
		t.Error("Expected error for invalid ayanamsa")
	}
}

func TestVargaDescription(t *testing.T) {
	if VargaDescription[VargaNavamsa] == "" {
		t.Error("VargaNavamsa should have a description")
	}
	if VargaDescription[VargaShashtiamsa] == "" {
		t.Error("VargaShashtiamsa should have a description")
	}
}

func TestNavamsaStartOffset(t *testing.T) {
	tests := []struct {
		signIdx int
		want    int
	}{
		{0, 0},  // Aries -> Aries
		{4, 0},  // Leo -> Aries
		{8, 0},  // Sagittarius -> Aries
		{1, 9},  // Taurus -> Capricorn
		{5, 9},  // Virgo -> Capricorn
		{9, 9},  // Capricorn -> Capricorn
		{2, 6},  // Gemini -> Libra
		{6, 6},  // Libra -> Libra
		{10, 6}, // Aquarius -> Libra
		{3, 3},  // Cancer -> Cancer
		{7, 3},  // Scorpio -> Cancer
		{11, 3}, // Pisces -> Cancer
	}
	for _, tt := range tests {
		got := navamsaStartOffset(tt.signIdx)
		if got != tt.want {
			t.Errorf("navamsaStartOffset(%d) = %d, want %d", tt.signIdx, got, tt.want)
		}
	}
}
