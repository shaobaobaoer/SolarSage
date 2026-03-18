package models

import (
	"strings"
	"testing"
)

func TestPlanetGlyph(t *testing.T) {
	tests := map[PlanetID]string{
		PlanetSun:     "\u2609", // ☉
		PlanetMoon:    "\u263D", // ☽
		PlanetVenus:   "\u2640", // ♀
		PlanetMars:    "\u2642", // ♂
		PlanetJupiter: "\u2643", // ♃
	}
	for pid, expected := range tests {
		got := PlanetGlyph(pid)
		if got != expected {
			t.Errorf("PlanetGlyph(%s) = %q, want %q", pid, got, expected)
		}
	}
}

func TestPlanetGlyph_Unknown(t *testing.T) {
	got := PlanetGlyph("UNKNOWN")
	if got != "UNKNOWN" {
		t.Errorf("PlanetGlyph(UNKNOWN) = %q, want UNKNOWN", got)
	}
}

func TestAspectGlyph(t *testing.T) {
	tests := map[AspectType]string{
		AspectConjunction: "\u260C", // ☌
		AspectOpposition:  "\u260D", // ☍
		AspectTrine:       "\u25B3", // △
		AspectSquare:      "\u25A1", // □
		AspectSextile:     "\u2731", // ✱
	}
	for at, expected := range tests {
		got := AspectGlyph(at)
		if got != expected {
			t.Errorf("AspectGlyph(%s) = %q, want %q", at, got, expected)
		}
	}
}

func TestAspectGlyph_Unknown(t *testing.T) {
	got := AspectGlyph("Unknown")
	if got != "Unknown" {
		t.Errorf("AspectGlyph(Unknown) = %q", got)
	}
}

func TestSignGlyph(t *testing.T) {
	if SignGlyph("Aries") != "\u2648" {
		t.Errorf("SignGlyph(Aries) = %q", SignGlyph("Aries"))
	}
	if SignGlyph("Pisces") != "\u2653" {
		t.Errorf("SignGlyph(Pisces) = %q", SignGlyph("Pisces"))
	}
	// Unknown sign returns the input
	if SignGlyph("Unknown") != "Unknown" {
		t.Errorf("SignGlyph(Unknown) = %q", SignGlyph("Unknown"))
	}
}

func TestSignGlyphFromLongitude(t *testing.T) {
	// 0° = Aries
	if SignGlyphFromLongitude(0) != "\u2648" {
		t.Error("0° should be Aries glyph")
	}
	// 90° = Cancer
	if SignGlyphFromLongitude(90) != "\u264B" {
		t.Error("90° should be Cancer glyph")
	}
}

func TestFormatLonGlyph(t *testing.T) {
	s := FormatLonGlyph(280.5)
	// Should contain Capricorn glyph ♑ and degree
	if !strings.Contains(s, "\u2651") {
		t.Errorf("FormatLonGlyph(280.5) = %q, missing Capricorn glyph", s)
	}
	if !strings.Contains(s, "10°") {
		t.Errorf("FormatLonGlyph(280.5) = %q, missing degree", s)
	}
}

func TestFormatPlanetGlyph(t *testing.T) {
	p := PlanetPosition{
		PlanetID:     PlanetSun,
		Longitude:    280.5,
		IsRetrograde: false,
	}
	s := FormatPlanetGlyph(p)
	if !strings.Contains(s, "\u2609") { // ☉
		t.Errorf("FormatPlanetGlyph missing Sun glyph: %q", s)
	}

	// Test retrograde
	p.IsRetrograde = true
	s = FormatPlanetGlyph(p)
	if !strings.Contains(s, "\u211E") { // ℞
		t.Errorf("FormatPlanetGlyph missing retrograde glyph: %q", s)
	}
}

func TestFormatAspectGlyph(t *testing.T) {
	a := AspectInfo{
		PlanetA:    "SUN",
		PlanetB:    "MOON",
		AspectType: AspectSquare,
		Orb:        2.5,
	}
	s := FormatAspectGlyph(a)
	if !strings.Contains(s, "\u25A1") { // □
		t.Errorf("FormatAspectGlyph missing square glyph: %q", s)
	}
	if !strings.Contains(s, "2.5") {
		t.Errorf("FormatAspectGlyph missing orb: %q", s)
	}
}

func TestSpecialPointGlyph(t *testing.T) {
	if SpecialPointGlyph(PointASC) != "AC" {
		t.Errorf("ASC glyph = %q", SpecialPointGlyph(PointASC))
	}
	if SpecialPointGlyph(PointLotFortune) != "\u2297" {
		t.Errorf("Fortune glyph = %q", SpecialPointGlyph(PointLotFortune))
	}
	// Unknown
	if SpecialPointGlyph("UNKNOWN") != "UNKNOWN" {
		t.Errorf("Unknown glyph = %q", SpecialPointGlyph("UNKNOWN"))
	}
}

func TestZodiacGlyphsLength(t *testing.T) {
	if len(ZodiacGlyphs) != 12 {
		t.Errorf("Expected 12 zodiac glyphs, got %d", len(ZodiacGlyphs))
	}
}
