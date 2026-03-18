package models

import "fmt"

// Zodiac sign glyphs (Unicode astrology symbols)
var ZodiacGlyphs = []string{
	"\u2648", // Aries ♈
	"\u2649", // Taurus ♉
	"\u264A", // Gemini ♊
	"\u264B", // Cancer ♋
	"\u264C", // Leo ♌
	"\u264D", // Virgo ♍
	"\u264E", // Libra ♎
	"\u264F", // Scorpio ♏
	"\u2650", // Sagittarius ♐
	"\u2651", // Capricorn ♑
	"\u2652", // Aquarius ♒
	"\u2653", // Pisces ♓
}

// Planet glyphs (Unicode astrology symbols)
var planetGlyphs = map[PlanetID]string{
	PlanetSun:           "\u2609", // ☉
	PlanetMoon:          "\u263D", // ☽
	PlanetMercury:       "\u263F", // ☿
	PlanetVenus:         "\u2640", // ♀
	PlanetMars:          "\u2642", // ♂
	PlanetJupiter:       "\u2643", // ♃
	PlanetSaturn:        "\u2644", // ♄
	PlanetUranus:        "\u2645", // ♅
	PlanetNeptune:       "\u2646", // ♆
	PlanetPluto:         "\u2647", // ♇
	PlanetChiron:        "\u26B7", // ⚷
	PlanetNorthNodeTrue: "\u260A", // ☊
	PlanetNorthNodeMean: "\u260A", // ☊
	PlanetSouthNode:     "\u260B", // ☋
	PlanetLilithMean:    "\u26B8", // ⚸
	PlanetLilithTrue:    "\u26B8", // ⚸
}

// Aspect glyphs (Unicode symbols)
var aspectGlyphs = map[AspectType]string{
	AspectConjunction:    "\u260C", // ☌
	AspectOpposition:     "\u260D", // ☍
	AspectTrine:          "\u25B3", // △
	AspectSquare:         "\u25A1", // □
	AspectSextile:        "\u2731", // ✱
	AspectQuincunx:       "\u26BB", // ⚻
	AspectSemiSextile:    "\u26BA", // ⚺
	AspectSemiSquare:     "\u2220", // ∠
	AspectSesquiquadrate: "\u26BC", // ⚼
}

// Special point glyphs
var specialPointGlyphs = map[SpecialPointID]string{
	PointASC:        "AC",
	PointMC:         "MC",
	PointDSC:        "DC",
	PointIC:         "IC",
	PointVertex:     "Vx",
	PointAntiVertex: "Av",
	PointEastPoint:  "EP",
	PointLotFortune: "\u2297", // ⊗
	PointLotSpirit:  "\u2295", // ⊕
}

// PlanetGlyph returns the Unicode astrology glyph for a planet
func PlanetGlyph(pid PlanetID) string {
	if g, ok := planetGlyphs[pid]; ok {
		return g
	}
	return string(pid)
}

// AspectGlyph returns the Unicode glyph for an aspect type
func AspectGlyph(at AspectType) string {
	if g, ok := aspectGlyphs[at]; ok {
		return g
	}
	return string(at)
}

// SignGlyph returns the Unicode zodiac glyph for a sign name
func SignGlyph(sign string) string {
	for i, s := range ZodiacSigns {
		if s == sign {
			return ZodiacGlyphs[i]
		}
	}
	return sign
}

// SignGlyphFromLongitude returns the zodiac glyph for an ecliptic longitude
func SignGlyphFromLongitude(lon float64) string {
	return ZodiacGlyphs[signIndex(lon)]
}

// SpecialPointGlyph returns the glyph for a special point
func SpecialPointGlyph(sp SpecialPointID) string {
	if g, ok := specialPointGlyphs[sp]; ok {
		return g
	}
	return string(sp)
}

// FormatLonGlyph formats a longitude as "DD°MM' ♈" with sign glyph
func FormatLonGlyph(lon float64) string {
	signDeg := SignDegreeFromLongitude(lon)
	glyph := SignGlyphFromLongitude(lon)
	dms := ToDMS(signDeg)
	return fmt.Sprintf("%d°%02d'%s", dms.Degrees, dms.Minutes, glyph)
}

// FormatPlanetGlyph formats a planet position as "☉ 10°15'♑"
func FormatPlanetGlyph(p PlanetPosition) string {
	glyph := PlanetGlyph(p.PlanetID)
	lonStr := FormatLonGlyph(p.Longitude)
	retro := ""
	if p.IsRetrograde {
		retro = " \u211E" // ℞
	}
	return glyph + " " + lonStr + retro
}

// FormatAspectGlyph formats an aspect as "☉ □ ☽ 2.5°"
func FormatAspectGlyph(a AspectInfo) string {
	pA := PlanetGlyph(PlanetID(a.PlanetA))
	pB := PlanetGlyph(PlanetID(a.PlanetB))
	asp := AspectGlyph(a.AspectType)
	return fmt.Sprintf("%s %s %s %.1f°", pA, asp, pB, a.Orb)
}
