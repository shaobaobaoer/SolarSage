package kp

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Vimshottari dasha sequence and durations (years) — 120-year cycle.
// Ketu is represented as PlanetSun proxy following the convention in pkg/vedic.
var dashaLords = []models.PlanetID{
	models.PlanetSun,     // Ketu   — 7 years
	models.PlanetVenus,   // Venus  — 20 years
	models.PlanetSun,     // Sun    — 6 years  (second entry = actual Sun)
	models.PlanetMoon,    // Moon   — 10 years
	models.PlanetMars,    // Mars   — 7 years
	models.PlanetSaturn,  // Rahu   — 18 years (Saturn proxy)
	models.PlanetJupiter, // Jupiter— 16 years
	models.PlanetSaturn,  // Saturn — 19 years
	models.PlanetMercury, // Mercury— 17 years
}

// dashaYears stores the dasha duration for each position in the sequence.
var dashaYears = []float64{7, 20, 6, 10, 7, 18, 16, 19, 17} // total = 120

// nakshatraSpan is the span of one Nakshatra in degrees (360°/27).
const nakshatraSpan = 360.0 / 27.0 // 13.3333...°

// SubLordInfo holds the Krishnamurti sub-lord analysis for a sidereal longitude.
type SubLordInfo struct {
	// Longitude is the input sidereal longitude (0–360°).
	Longitude float64 `json:"longitude"`

	// NakshatraIndex is the 0-based Nakshatra index (0=Ashwini … 26=Revati).
	NakshatraIndex int `json:"nakshatra_index"`

	// NakshatraLord is the ruler of the Nakshatra (Star lord in KP terminology).
	NakshatraLord models.PlanetID `json:"nakshatra_lord"`

	// SubLord is the ruler of the Vimshottari sub-division within the Nakshatra.
	SubLord models.PlanetID `json:"sub_lord"`

	// SubSubLord is the ruler of the further sub-division within the sub-period.
	SubSubLord models.PlanetID `json:"sub_sub_lord"`

	// CuspDegree is the degree within the Nakshatra (0–13.333°).
	CuspDegree float64 `json:"cusp_degree"`
}

// nakshatraLords maps Nakshatra index (0-26) to the Vimshottari dasha lord.
// The sequence repeats from index 0: Ketu(0), Venus(1), Sun(2), …, Mercury(8),
// Ketu(9), …
var nakshatraLords = [27]models.PlanetID{
	models.PlanetSun,     // 0 Ashwini     — Ketu (Sun proxy)
	models.PlanetVenus,   // 1 Bharani     — Venus
	models.PlanetSun,     // 2 Krittika    — Sun
	models.PlanetMoon,    // 3 Rohini      — Moon
	models.PlanetMars,    // 4 Mrigashirsha— Mars
	models.PlanetSaturn,  // 5 Ardra       — Rahu (Saturn proxy)
	models.PlanetJupiter, // 6 Punarvasu   — Jupiter
	models.PlanetSaturn,  // 7 Pushya      — Saturn
	models.PlanetMercury, // 8 Ashlesha    — Mercury
	models.PlanetSun,     // 9 Magha       — Ketu (Sun proxy)
	models.PlanetVenus,   // 10 Purva Phalguni — Venus
	models.PlanetSun,     // 11 Uttara Phalguni— Sun
	models.PlanetMoon,    // 12 Hasta      — Moon
	models.PlanetMars,    // 13 Chitra     — Mars
	models.PlanetSaturn,  // 14 Swati      — Rahu (Saturn proxy)
	models.PlanetJupiter, // 15 Vishakha   — Jupiter
	models.PlanetSaturn,  // 16 Anuradha   — Saturn
	models.PlanetMercury, // 17 Jyeshtha   — Mercury
	models.PlanetSun,     // 18 Mula       — Ketu (Sun proxy)
	models.PlanetVenus,   // 19 Purva Ashadha  — Venus
	models.PlanetSun,     // 20 Uttara Ashadha — Sun
	models.PlanetMoon,    // 21 Shravana   — Moon
	models.PlanetMars,    // 22 Dhanishtha — Mars
	models.PlanetSaturn,  // 23 Shatabhisha— Rahu (Saturn proxy)
	models.PlanetJupiter, // 24 Purva Bhadrapada  — Jupiter
	models.PlanetSaturn,  // 25 Uttara Bhadrapada — Saturn
	models.PlanetMercury, // 26 Revati     — Mercury
}

// dashaIndexOf returns the position (0–8) in dashaLords for a given planet.
// Since Ketu and Rahu are proxied to Sun and Saturn respectively, we match by value.
func dashaIndexOf(lord models.PlanetID) int {
	// The lord of a Nakshatra is the first occurrence in dashaLords.
	// Nakshatras cycle through groups of 9 starting from Ketu=index 0.
	for i, l := range dashaLords {
		if l == lord {
			return i
		}
	}
	return 0
}

// SubLords computes the Star lord, Sub lord, and Sub-sub lord for a given
// sidereal longitude using the Krishnamurti Paddhati triple-subdivision algorithm.
//
// The algorithm:
//  1. Each Nakshatra (13.333°) is subdivided proportionally by the 9 dasha lords
//     in the Vimshottari sequence starting from the Nakshatra lord itself.
//  2. Each sub-period is further subdivided in the same proportional manner.
//  3. The zodiac is periodic in 120° cycles (one full Vimshottari cycle per 3
//     Nakshatras × 9 lords). Reduce longitude modulo 120° before iterating.
func SubLords(siderealLon float64) SubLordInfo {
	// Normalize to [0, 360)
	lon := siderealLon
	for lon < 0 {
		lon += 360
	}
	lon = lon - float64(int(lon/360))*360

	// Nakshatra index and degree within nakshatra
	nakIdx := int(lon / nakshatraSpan)
	if nakIdx > 26 {
		nakIdx = 26
	}
	cuspDeg := lon - float64(nakIdx)*nakshatraSpan

	nakLord := nakshatraLords[nakIdx]
	startDashaIdx := dashaIndexOf(nakLord)

	// Reduce to 120° periodicity for the sub-lord iteration.
	// The KP table repeats every 120° (9 lords × 13.333° = 120°).
	reducedDeg := cuspDeg
	// reducedDeg is already within one Nakshatra (0–13.333°) — no further reduction needed.

	subLord, subSubLord := calcSubLords(reducedDeg, nakshatraSpan, startDashaIdx)

	return SubLordInfo{
		Longitude:      siderealLon,
		NakshatraIndex: nakIdx,
		NakshatraLord:  nakLord,
		SubLord:        subLord,
		SubSubLord:     subSubLord,
		CuspDegree:     cuspDeg,
	}
}

// calcSubLords performs the two-level KP subdivision.
// degInNakshatra is the degree within the current Nakshatra (0–nakshatraSpan).
// nakSpan is the total span of the Nakshatra.
// startIdx is the index into dashaLords for the Nakshatra lord.
func calcSubLords(degInNakshatra, nakSpan float64, startIdx int) (subLord, subSubLord models.PlanetID) {
	// Sub-lord: divide nakshatra into 9 segments proportional to dashaYears
	subStart := 0.0
	subIdx := startIdx
	for i := 0; i < 9; i++ {
		segLen := nakSpan * dashaYears[subIdx%9] / 120.0
		subEnd := subStart + segLen
		if degInNakshatra < subEnd || i == 8 {
			subLord = dashaLords[subIdx%9]
			// Sub-sub lord: divide this sub-segment into 9 further segments
			subSubLord = calcSubSubLord(degInNakshatra-subStart, segLen, subIdx%9)
			return
		}
		subStart = subEnd
		subIdx++
	}
	return dashaLords[startIdx], dashaLords[startIdx]
}

// calcSubSubLord computes the sub-sub lord within a sub-period segment.
// degInSub is the degree within the sub-period (0–subSpan).
// subSpan is the total span of the sub-period.
// subIdx is the index into dashaLords for the sub lord.
func calcSubSubLord(degInSub, subSpan float64, subIdx int) models.PlanetID {
	ssStart := 0.0
	ssIdx := subIdx
	for i := 0; i < 9; i++ {
		segLen := subSpan * dashaYears[ssIdx%9] / 120.0
		ssEnd := ssStart + segLen
		if degInSub < ssEnd || i == 8 {
			return dashaLords[ssIdx%9]
		}
		ssStart = ssEnd
		ssIdx++
	}
	return dashaLords[subIdx]
}
