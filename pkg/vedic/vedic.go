// Package vedic provides sidereal zodiac calculations for Vedic (Jyotish) astrology.
// It supports multiple ayanamsa systems, Nakshatra (lunar mansion) calculations,
// and Vedic-specific chart analysis.
//
// Use TropicalToSidereal to convert tropical longitudes, or compute a full
// sidereal chart with CalcSiderealChart.
package vedic

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// Ayanamsa represents a sidereal zodiac offset system
type Ayanamsa string

const (
	AyanamsaLahiri      Ayanamsa = "LAHIRI"       // Most widely used in India
	AyanamsaRaman       Ayanamsa = "RAMAN"         // B.V. Raman system
	AyanamsaKrishnamurti Ayanamsa = "KRISHNAMURTI" // KP system
	AyanamsaFaganBradley Ayanamsa = "FAGAN_BRADLEY" // Western sidereal
	AyanamsaYukteshwar  Ayanamsa = "YUKTESHWAR"    // Sri Yukteshwar
)

// ayanamsaJ2000 stores the ayanamsa value at J2000.0 for each system
var ayanamsaJ2000 = map[Ayanamsa]float64{
	AyanamsaLahiri:       23.853,  // Lahiri at J2000.0
	AyanamsaRaman:        22.378,  // Raman at J2000.0
	AyanamsaKrishnamurti: 23.795,  // Krishnamurti at J2000.0
	AyanamsaFaganBradley: 24.736,  // Fagan-Bradley at J2000.0
	AyanamsaYukteshwar:   22.277,  // Yukteshwar at J2000.0
}

// precessionRate is approximately 50.2882 arcseconds per year
const precessionRate = 50.2882 / 3600.0
const j2000Epoch = 2451545.0
const julianYear = 365.25

// GetAyanamsa returns the ayanamsa (sidereal offset) at a given JD for the specified system.
func GetAyanamsa(jdUT float64, system Ayanamsa) (float64, error) {
	base, ok := ayanamsaJ2000[system]
	if !ok {
		return 0, fmt.Errorf("unknown ayanamsa system: %s", system)
	}
	yearsSinceJ2000 := (jdUT - j2000Epoch) / julianYear
	return base + yearsSinceJ2000*precessionRate, nil
}

// TropicalToSidereal converts a tropical ecliptic longitude to sidereal.
func TropicalToSidereal(tropicalLon, ayanamsa float64) float64 {
	return sweph.NormalizeDegrees(tropicalLon - ayanamsa)
}

// SiderealPosition extends PlanetPosition with sidereal data
type SiderealPosition struct {
	models.PlanetPosition
	SiderealLon  float64 `json:"sidereal_longitude"`
	SiderealSign string  `json:"sidereal_sign"`
	SiderealDeg  float64 `json:"sidereal_sign_degree"`
	Nakshatra    string  `json:"nakshatra"`
	NakshatraPada int    `json:"nakshatra_pada"` // 1-4
	NakshatraLord models.PlanetID `json:"nakshatra_lord"`
}

// SiderealChart holds a complete sidereal chart
type SiderealChart struct {
	Ayanamsa     Ayanamsa           `json:"ayanamsa"`
	AyanamsaValue float64           `json:"ayanamsa_value"`
	Planets      []SiderealPosition `json:"planets"`
	Houses       []float64          `json:"houses"`
	Angles       models.AnglesInfo  `json:"angles"`
	SiderealAngles models.AnglesInfo `json:"sidereal_angles"`
}

// CalcSiderealChart computes a sidereal natal chart with Nakshatra data.
func CalcSiderealChart(lat, lon, jdUT float64, system Ayanamsa) (*SiderealChart, error) {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}

	tropicalChart, err := chart.CalcSingleChart(lat, lon, jdUT, planets,
		models.DefaultOrbConfig(), models.HousePlacidus)
	if err != nil {
		return nil, err
	}

	ayanamsaVal, err := GetAyanamsa(jdUT, system)
	if err != nil {
		return nil, err
	}

	siderealPlanets := make([]SiderealPosition, len(tropicalChart.Planets))
	for i, p := range tropicalChart.Planets {
		sidLon := TropicalToSidereal(p.Longitude, ayanamsaVal)
		nak, pada, lord := CalcNakshatra(sidLon)
		siderealPlanets[i] = SiderealPosition{
			PlanetPosition: p,
			SiderealLon:    sidLon,
			SiderealSign:   models.SignFromLongitude(sidLon),
			SiderealDeg:    models.SignDegreeFromLongitude(sidLon),
			Nakshatra:      nak,
			NakshatraPada:  pada,
			NakshatraLord:  lord,
		}
	}

	// Convert house cusps to sidereal
	sidHouses := make([]float64, len(tropicalChart.Houses))
	for i, h := range tropicalChart.Houses {
		sidHouses[i] = TropicalToSidereal(h, ayanamsaVal)
	}

	sidAngles := models.AnglesInfo{
		ASC: TropicalToSidereal(tropicalChart.Angles.ASC, ayanamsaVal),
		MC:  TropicalToSidereal(tropicalChart.Angles.MC, ayanamsaVal),
		DSC: TropicalToSidereal(tropicalChart.Angles.DSC, ayanamsaVal),
		IC:  TropicalToSidereal(tropicalChart.Angles.IC, ayanamsaVal),
	}

	return &SiderealChart{
		Ayanamsa:       system,
		AyanamsaValue:  ayanamsaVal,
		Planets:        siderealPlanets,
		Houses:         sidHouses,
		Angles:         tropicalChart.Angles,
		SiderealAngles: sidAngles,
	}, nil
}

// Nakshatra represents one of the 27 lunar mansions
type NakshatraInfo struct {
	Number int             `json:"number"` // 1-27
	Name   string          `json:"name"`
	Start  float64         `json:"start_degree"`
	Lord   models.PlanetID `json:"lord"`
}

// Nakshatras is the ordered list of 27 lunar mansions with their Vimshottari lords
var Nakshatras = []NakshatraInfo{
	{1, "Ashwini", 0, models.PlanetSun},       // Ketu -> map to Sun for simplicity
	{2, "Bharani", 13.333, models.PlanetVenus},
	{3, "Krittika", 26.667, models.PlanetSun},
	{4, "Rohini", 40.0, models.PlanetMoon},
	{5, "Mrigashirsha", 53.333, models.PlanetMars},
	{6, "Ardra", 66.667, models.PlanetSaturn},    // Rahu -> Saturn
	{7, "Punarvasu", 80.0, models.PlanetJupiter},
	{8, "Pushya", 93.333, models.PlanetSaturn},
	{9, "Ashlesha", 106.667, models.PlanetMercury},
	{10, "Magha", 120.0, models.PlanetSun},        // Ketu -> Sun
	{11, "Purva Phalguni", 133.333, models.PlanetVenus},
	{12, "Uttara Phalguni", 146.667, models.PlanetSun},
	{13, "Hasta", 160.0, models.PlanetMoon},
	{14, "Chitra", 173.333, models.PlanetMars},
	{15, "Swati", 186.667, models.PlanetSaturn},   // Rahu -> Saturn
	{16, "Vishakha", 200.0, models.PlanetJupiter},
	{17, "Anuradha", 213.333, models.PlanetSaturn},
	{18, "Jyeshtha", 226.667, models.PlanetMercury},
	{19, "Mula", 240.0, models.PlanetSun},          // Ketu -> Sun
	{20, "Purva Ashadha", 253.333, models.PlanetVenus},
	{21, "Uttara Ashadha", 266.667, models.PlanetSun},
	{22, "Shravana", 280.0, models.PlanetMoon},
	{23, "Dhanishtha", 293.333, models.PlanetMars},
	{24, "Shatabhisha", 306.667, models.PlanetSaturn}, // Rahu -> Saturn
	{25, "Purva Bhadrapada", 320.0, models.PlanetJupiter},
	{26, "Uttara Bhadrapada", 333.333, models.PlanetSaturn},
	{27, "Revati", 346.667, models.PlanetMercury},
}

// CalcNakshatra returns the Nakshatra name, pada (1-4), and Vimshottari lord
// for a given sidereal longitude.
func CalcNakshatra(siderealLon float64) (name string, pada int, lord models.PlanetID) {
	siderealLon = sweph.NormalizeDegrees(siderealLon)

	// Each Nakshatra spans 13.333° (360/27)
	nakshatraSpan := 360.0 / 27.0
	idx := int(siderealLon / nakshatraSpan)
	if idx >= 27 {
		idx = 26
	}

	nak := Nakshatras[idx]

	// Each Nakshatra has 4 padas of 3.333° each
	posInNakshatra := siderealLon - float64(idx)*nakshatraSpan
	pada = int(posInNakshatra/(nakshatraSpan/4.0)) + 1
	if pada > 4 {
		pada = 4
	}

	return nak.Name, pada, nak.Lord
}

// VimshottariDasha calculates the Vimshottari Maha Dasha period active at birth.
// Returns the dasha lord, years remaining in the period, and the full dasha sequence.
type DashaPeriod struct {
	Lord     models.PlanetID `json:"lord"`
	Years    float64         `json:"years"`    // Total duration
	StartAge float64         `json:"start_age"` // Age when this period starts
}

// dashaYears maps each planet to its Vimshottari Maha Dasha duration
var dashaYears = map[models.PlanetID]float64{
	models.PlanetSun:     6,
	models.PlanetMoon:    10,
	models.PlanetMars:    7,
	models.PlanetSaturn:  19, // Rahu mapped to Saturn
	models.PlanetJupiter: 16,
	models.PlanetMercury: 17,
	models.PlanetVenus:   20,
}

// dashaSequence is the Vimshottari sequence starting from Ketu (Sun as proxy)
var dashaSequence = []models.PlanetID{
	models.PlanetSun,     // Ketu (6 years)
	models.PlanetVenus,   // Venus (20 years)
	models.PlanetSun,     // Sun (6 years) - actual Sun
	models.PlanetMoon,    // Moon (10 years)
	models.PlanetMars,    // Mars (7 years)
	models.PlanetSaturn,  // Rahu (18->19 years, mapped to Saturn)
	models.PlanetJupiter, // Jupiter (16 years)
	models.PlanetSaturn,  // Saturn (19 years)
	models.PlanetMercury, // Mercury (17 years)
}

// CalcVimshottariDasha calculates the Maha Dasha sequence starting from birth.
func CalcVimshottariDasha(moonSiderealLon float64) []DashaPeriod {
	_, _, lord := CalcNakshatra(moonSiderealLon)

	// Find starting index in dasha sequence
	startIdx := 0
	for i, l := range dashaSequence {
		if l == lord {
			startIdx = i
			break
		}
	}

	// Calculate remaining balance of first dasha
	nakshatraSpan := 360.0 / 27.0
	idx := int(moonSiderealLon / nakshatraSpan)
	posInNakshatra := moonSiderealLon - float64(idx)*nakshatraSpan
	fractionUsed := posInNakshatra / nakshatraSpan
	firstDashaTotal := dashaYears[lord]
	firstDashaRemaining := firstDashaTotal * (1 - fractionUsed)

	// Build the full sequence
	var periods []DashaPeriod
	age := 0.0

	// First period (partial)
	periods = append(periods, DashaPeriod{
		Lord:     dashaSequence[startIdx],
		Years:    math.Round(firstDashaRemaining*100) / 100,
		StartAge: age,
	})
	age += firstDashaRemaining

	// Remaining periods (full)
	for i := 1; i < len(dashaSequence); i++ {
		idx := (startIdx + i) % len(dashaSequence)
		lord := dashaSequence[idx]
		years := dashaYears[lord]
		periods = append(periods, DashaPeriod{
			Lord:     lord,
			Years:    years,
			StartAge: math.Round(age*100) / 100,
		})
		age += years
	}

	return periods
}
