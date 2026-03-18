package bounds

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// DecanInfo holds decan information for a position
type DecanInfo struct {
	PlanetID    models.PlanetID `json:"planet_id,omitempty"`
	Sign        string          `json:"sign"`
	Decan       int             `json:"decan"` // 1, 2, or 3
	DecanRuler  models.PlanetID `json:"decan_ruler"`
	DecanDegree string          `json:"decan_degrees"` // e.g. "0-10"
}

// TermInfo holds Egyptian/Ptolemaic term information
type TermInfo struct {
	PlanetID  models.PlanetID `json:"planet_id,omitempty"`
	Sign      string          `json:"sign"`
	TermRuler models.PlanetID `json:"term_ruler"`
	TermStart float64         `json:"term_start"`
	TermEnd   float64         `json:"term_end"`
}

// FaceInfo combines decan and term for a position
type FaceInfo struct {
	PlanetID models.PlanetID `json:"planet_id"`
	Sign     string          `json:"sign"`
	SignDeg  float64         `json:"sign_degree"`
	Decan    DecanInfo       `json:"decan"`
	Term     TermInfo        `json:"term"`
}

// Chaldean decan rulers (each sign divided into 3 x 10°)
// Follows the Chaldean order: Mars, Sun, Venus, Mercury, Moon, Saturn, Jupiter
// Starting from Aries 1st decan = Mars
var chaldeanDecans = [12][3]models.PlanetID{
	// Aries
	{models.PlanetMars, models.PlanetSun, models.PlanetVenus},
	// Taurus
	{models.PlanetMercury, models.PlanetMoon, models.PlanetSaturn},
	// Gemini
	{models.PlanetJupiter, models.PlanetMars, models.PlanetSun},
	// Cancer
	{models.PlanetVenus, models.PlanetMercury, models.PlanetMoon},
	// Leo
	{models.PlanetSaturn, models.PlanetJupiter, models.PlanetMars},
	// Virgo
	{models.PlanetSun, models.PlanetVenus, models.PlanetMercury},
	// Libra
	{models.PlanetMoon, models.PlanetSaturn, models.PlanetJupiter},
	// Scorpio
	{models.PlanetMars, models.PlanetSun, models.PlanetVenus},
	// Sagittarius
	{models.PlanetMercury, models.PlanetMoon, models.PlanetSaturn},
	// Capricorn
	{models.PlanetJupiter, models.PlanetMars, models.PlanetSun},
	// Aquarius
	{models.PlanetVenus, models.PlanetMercury, models.PlanetMoon},
	// Pisces
	{models.PlanetSaturn, models.PlanetJupiter, models.PlanetMars},
}

// termBound represents a term boundary
type termBound struct {
	End   float64
	Ruler models.PlanetID
}

// Egyptian (Ptolemaic) terms - the most widely used system
// Each sign has 5 terms ruled by the 5 traditional planets
var egyptianTerms = [12][]termBound{
	// Aries
	{{6, models.PlanetJupiter}, {12, models.PlanetVenus}, {20, models.PlanetMercury}, {25, models.PlanetMars}, {30, models.PlanetSaturn}},
	// Taurus
	{{8, models.PlanetVenus}, {14, models.PlanetMercury}, {22, models.PlanetJupiter}, {27, models.PlanetSaturn}, {30, models.PlanetMars}},
	// Gemini
	{{6, models.PlanetMercury}, {12, models.PlanetJupiter}, {17, models.PlanetVenus}, {24, models.PlanetMars}, {30, models.PlanetSaturn}},
	// Cancer
	{{7, models.PlanetMars}, {13, models.PlanetVenus}, {19, models.PlanetMercury}, {26, models.PlanetJupiter}, {30, models.PlanetSaturn}},
	// Leo
	{{6, models.PlanetJupiter}, {11, models.PlanetVenus}, {18, models.PlanetSaturn}, {24, models.PlanetMercury}, {30, models.PlanetMars}},
	// Virgo
	{{7, models.PlanetMercury}, {17, models.PlanetVenus}, {21, models.PlanetJupiter}, {28, models.PlanetMars}, {30, models.PlanetSaturn}},
	// Libra
	{{6, models.PlanetSaturn}, {14, models.PlanetMercury}, {21, models.PlanetJupiter}, {28, models.PlanetVenus}, {30, models.PlanetMars}},
	// Scorpio
	{{7, models.PlanetMars}, {11, models.PlanetVenus}, {19, models.PlanetMercury}, {24, models.PlanetJupiter}, {30, models.PlanetSaturn}},
	// Sagittarius
	{{12, models.PlanetJupiter}, {17, models.PlanetVenus}, {21, models.PlanetMercury}, {26, models.PlanetMars}, {30, models.PlanetSaturn}},
	// Capricorn
	{{7, models.PlanetMercury}, {14, models.PlanetJupiter}, {22, models.PlanetVenus}, {26, models.PlanetSaturn}, {30, models.PlanetMars}},
	// Aquarius
	{{7, models.PlanetMercury}, {13, models.PlanetVenus}, {20, models.PlanetJupiter}, {25, models.PlanetMars}, {30, models.PlanetSaturn}},
	// Pisces
	{{12, models.PlanetVenus}, {16, models.PlanetJupiter}, {19, models.PlanetMercury}, {28, models.PlanetMars}, {30, models.PlanetSaturn}},
}

// CalcDecan returns the decan information for a given ecliptic longitude
func CalcDecan(lon float64) DecanInfo {
	signIdx := int(lon / 30.0)
	if signIdx > 11 {
		signIdx = 11
	}
	if signIdx < 0 {
		signIdx = 0
	}
	signDeg := lon - float64(signIdx)*30.0

	decanNum := int(signDeg/10.0) + 1
	if decanNum > 3 {
		decanNum = 3
	}

	decanDeg := []string{"0-10", "10-20", "20-30"}

	return DecanInfo{
		Sign:        models.ZodiacSigns[signIdx],
		Decan:       decanNum,
		DecanRuler:  chaldeanDecans[signIdx][decanNum-1],
		DecanDegree: decanDeg[decanNum-1],
	}
}

// CalcTerm returns the Egyptian term information for a given ecliptic longitude
func CalcTerm(lon float64) TermInfo {
	signIdx := int(lon / 30.0)
	if signIdx > 11 {
		signIdx = 11
	}
	if signIdx < 0 {
		signIdx = 0
	}
	signDeg := lon - float64(signIdx)*30.0

	terms := egyptianTerms[signIdx]
	start := 0.0
	for _, t := range terms {
		if signDeg < t.End {
			return TermInfo{
				Sign:      models.ZodiacSigns[signIdx],
				TermRuler: t.Ruler,
				TermStart: start,
				TermEnd:   t.End,
			}
		}
		start = t.End
	}

	// Fallback to last term
	last := terms[len(terms)-1]
	return TermInfo{
		Sign:      models.ZodiacSigns[signIdx],
		TermRuler: last.Ruler,
		TermStart: start,
		TermEnd:   last.End,
	}
}

// CalcChartFaces computes decans and terms for all planet positions
func CalcChartFaces(positions []models.PlanetPosition) []FaceInfo {
	faces := make([]FaceInfo, len(positions))
	for i, p := range positions {
		decan := CalcDecan(p.Longitude)
		decan.PlanetID = p.PlanetID
		term := CalcTerm(p.Longitude)
		term.PlanetID = p.PlanetID
		faces[i] = FaceInfo{
			PlanetID: p.PlanetID,
			Sign:     p.Sign,
			SignDeg:  p.SignDegree,
			Decan:    decan,
			Term:     term,
		}
	}
	return faces
}
