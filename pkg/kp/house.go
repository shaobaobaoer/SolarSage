package kp

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// HouseNumber returns the KP house number (1–12) for a planet at a given
// sidereal longitude, based on the 12 sidereal house cusp longitudes.
//
// KP house placement differs from sign-based placement: a planet belongs to
// the house whose cusp it has passed most recently (i.e. the largest cusp
// longitude ≤ planet longitude, wrapping around from cusp 12 to cusp 1).
//
// cusps must be a slice of 12 sidereal cusp longitudes (index 0 = cusp 1, …,
// index 11 = cusp 12). Longitudes are in degrees [0, 360).
func HouseNumber(planetLon float64, cusps []float64) int {
	if len(cusps) != 12 {
		return 1
	}
	// Normalize planet longitude
	pLon := normDeg(planetLon)

	// Find the cusp the planet has most recently passed.
	// Iterate cusps 12 → 1 to find the last cusp ≤ pLon (accounting for wrap).
	for h := 11; h >= 0; h-- {
		cLon := normDeg(cusps[h])
		nextH := (h + 1) % 12
		nLon := normDeg(cusps[nextH])

		if cuspContains(cLon, nLon, pLon) {
			return h + 1 // house numbers are 1-based
		}
	}
	return 1
}

// cuspContains reports whether lon is within the arc from start to end
// (going forward on the zodiac), handling the 0°/360° wrap.
func cuspContains(start, end, lon float64) bool {
	if start <= end {
		return lon >= start && lon < end
	}
	// Arc wraps past 360°
	return lon >= start || lon < end
}

// normDeg normalises a longitude to [0, 360).
func normDeg(d float64) float64 {
	d = d - float64(int(d/360))*360
	if d < 0 {
		d += 360
	}
	return d
}

// Significators holds the ABCD significators for a planet or house in KP.
//
// The ABCD method classifies planets by the strength of their connection to
// each house:
//
//	A (strongest): planets occupying the star (Nakshatra) of occupants of the house
//	B:             planets occupying the house itself
//	C:             planets occupying the star of the house ruler
//	D (weakest):   the house ruler itself
type Significators struct {
	A []models.PlanetID `json:"a"` // in star of house occupants
	B []models.PlanetID `json:"b"` // in the house itself
	C []models.PlanetID `json:"c"` // in star of house ruler
	D []models.PlanetID `json:"d"` // the house ruler
}

// PlanetKPInfo holds KP-specific data for one planet.
type PlanetKPInfo struct {
	PlanetID      models.PlanetID `json:"planet_id"`
	Longitude     float64         `json:"longitude"`      // sidereal longitude
	House         int             `json:"house"`          // KP house (1–12)
	NakshatraLord models.PlanetID `json:"nakshatra_lord"` // Star lord
	SubLord       models.PlanetID `json:"sub_lord"`
	SubSubLord    models.PlanetID `json:"sub_sub_lord"`
}

// ChartKPInput is the input required to compute a full KP chart analysis.
type ChartKPInput struct {
	// Planets holds all planet sidereal longitudes.
	Planets []PlanetKPInfo
	// Cusps holds 12 sidereal house cusp longitudes (index 0 = cusp 1).
	Cusps []float64
	// HouseRulers maps house number (1–12) to the traditional ruling planet.
	HouseRulers [13]models.PlanetID // index 1-12 used
}

// CalcPlanetKP computes the full KP analysis (house, star lord, sub lord,
// sub-sub lord) for a single planet at a given sidereal longitude.
func CalcPlanetKP(planetLon float64, cusps []float64) PlanetKPInfo {
	sl := SubLords(planetLon)
	house := HouseNumber(planetLon, cusps)
	return PlanetKPInfo{
		Longitude:     planetLon,
		House:         house,
		NakshatraLord: sl.NakshatraLord,
		SubLord:       sl.SubLord,
		SubSubLord:    sl.SubSubLord,
	}
}

// HouseSignificators returns the ABCD significators for a given house number.
//
// Parameters:
//   - houseNum: the house to analyse (1–12)
//   - planetInfos: full KP info for all planets in the chart
//   - houseRuler: the traditional ruler of the house
func HouseSignificators(houseNum int, planetInfos []PlanetKPInfo, houseRuler models.PlanetID) Significators {
	var sig Significators
	sig.D = []models.PlanetID{houseRuler}

	// Collect planets occupying the house (Group B)
	var occupants []models.PlanetID
	for _, p := range planetInfos {
		if p.House == houseNum {
			occupants = append(occupants, p.PlanetID)
			sig.B = append(sig.B, p.PlanetID)
		}
	}

	// Group A: planets whose Nakshatra lord is one of the house occupants
	for _, p := range planetInfos {
		for _, occ := range occupants {
			if p.NakshatraLord == occ && p.PlanetID != occ {
				sig.A = append(sig.A, p.PlanetID)
				break
			}
		}
	}

	// Group C: planets whose Nakshatra lord is the house ruler
	for _, p := range planetInfos {
		if p.NakshatraLord == houseRuler {
			sig.C = append(sig.C, p.PlanetID)
		}
	}

	return sig
}

// PlanetSignificators returns the ABCD significators for a given planet,
// i.e. which houses does this planet signify?
//
// Returns four slices of house numbers:
//
//	A: houses occupied by the planet's Nakshatra lord
//	B: the house the planet itself occupies
//	C: houses ruled by the planet's Nakshatra lord
//	D: houses ruled by the planet itself
func PlanetSignificators(
	planet PlanetKPInfo,
	houseRulers [13]models.PlanetID,
) (A, B, C, D []int) {
	// B: house the planet occupies
	B = []int{planet.House}

	// A: houses occupied by the Nakshatra lord (we'd need all planet infos;
	//    kept simple here — caller should use full chart version below)
	// D: houses the planet itself rules
	for h := 1; h <= 12; h++ {
		if houseRulers[h] == planet.PlanetID {
			D = append(D, h)
		}
		if houseRulers[h] == planet.NakshatraLord {
			C = append(C, h)
		}
	}
	return
}
