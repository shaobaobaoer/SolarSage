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
	AyanamsaLahiri       Ayanamsa = "LAHIRI"        // Most widely used in India
	AyanamsaRaman        Ayanamsa = "RAMAN"          // B.V. Raman system
	AyanamsaKrishnamurti Ayanamsa = "KRISHNAMURTI"   // KP system
	AyanamsaFaganBradley Ayanamsa = "FAGAN_BRADLEY"  // Western sidereal
	AyanamsaYukteshwar   Ayanamsa = "YUKTESHWAR"     // Sri Yukteshwar
	AyanamsaTrueCitra    Ayanamsa = "TRUE_CITRA"     // True Chitrapaksha
	AyanamsaTrueRevati   Ayanamsa = "TRUE_REVATI"    // True Revati
	AyanamsaTruePushya   Ayanamsa = "TRUE_PUSHYA"    // True Pushya
	AyanamsaTrueMula     Ayanamsa = "TRUE_MULA"      // True Mula
)

// ayanamsaSidMode maps each Ayanamsa to the Swiss Ephemeris SIDM_* constant.
// Values match the swe_set_sid_mode() constants in swephexp.h.
var ayanamsaSidMode = map[Ayanamsa]int{
	AyanamsaFaganBradley: sweph.SidmFaganBradley,
	AyanamsaLahiri:       sweph.SidmLahiri,
	AyanamsaRaman:        sweph.SidmRaman,
	AyanamsaKrishnamurti: sweph.SidmKrishnamurti,
	AyanamsaYukteshwar:   sweph.SidmYukteshwar,
	AyanamsaTrueCitra:    sweph.SidmTrueCitra,
	AyanamsaTrueRevati:   sweph.SidmTrueRevati,
	AyanamsaTruePushya:   sweph.SidmTruePushya,
	AyanamsaTrueMula:     sweph.SidmTrueMula,
}

// GetAyanamsa returns the precise ayanamsa value at a given JD UT using the
// Swiss Ephemeris native computation. Replaces the previous linear approximation
// which could drift up to 0.3° from the true value.
func GetAyanamsa(jdUT float64, system Ayanamsa) (float64, error) {
	sidMode, ok := ayanamsaSidMode[system]
	if !ok {
		return 0, fmt.Errorf("unknown ayanamsa system: %s", system)
	}
	return sweph.GetAyanamsaUT(jdUT, sidMode), nil
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

// ---- Ashtottari Dasha (108-year cycle) ----

// ashtottariYears maps planets to their Ashtottari Maha Dasha durations (total = 108 years).
// Sequence: Sun(6) Venus(20) Moon(15) Mars(8) Mercury(17) Jupiter(10) Saturn(19) Rahu(12)->Saturn proxy.
// Lord sequence used for the 8-planet system (Rahu represented as extra Saturn slot).
var ashtottariYears = map[models.PlanetID]float64{
	models.PlanetSun:     6,
	models.PlanetVenus:   20,
	models.PlanetMoon:    15,
	models.PlanetMars:    8,
	models.PlanetMercury: 17,
	models.PlanetJupiter: 10,
	models.PlanetSaturn:  19, // Rahu (12) + Saturn (19) merged; canonical Ashtottari uses Rahu separately
}

// ashtottariSequence is the 8-planet Ashtottari sequence.
// Ketu is omitted in this system; the sequence is Sun-Venus-Moon-Mars-Mercury-Jupiter-Saturn-Rahu(->Saturn).
var ashtottariSequence = []models.PlanetID{
	models.PlanetSun,
	models.PlanetVenus,
	models.PlanetMoon,
	models.PlanetMars,
	models.PlanetMercury,
	models.PlanetJupiter,
	models.PlanetSaturn,
	models.PlanetSaturn, // Rahu slot (proxy)
}

// ashtottariFullYears gives the canonical year per slot in ashtottariSequence
var ashtottariFullYears = []float64{6, 20, 15, 8, 17, 10, 19, 12}

// ashtottariNakshatra maps Nakshatra index (0-26) to the starting dasha planet index
// in ashtottariSequence. The 27 Nakshatras distribute evenly over the 8 lords.
var ashtottariNakshatraLord = []int{
	// Ashwini(0)-Revati(26): repeating pattern Sun,Venus,Moon,Mars,Mercury,Jupiter,Saturn,Rahu
	0, 1, 2, 3, 4, 5, 6, 7, // 0-7
	0, 1, 2, 3, 4, 5, 6, 7, // 8-15
	0, 1, 2, 3, 4, 5, 6, 7, // 16-23
	0, 1, 2, // 24-26
}

// CalcAshtottariDasha calculates the Ashtottari (108-year) Maha Dasha sequence.
// Valid for charts where the Moon occupies a Rahu-associated Nakshatra or
// for charts with an Aquarius Lagna; widely used in Kerala tradition.
func CalcAshtottariDasha(moonSiderealLon float64) []DashaPeriod {
	moonSiderealLon = math.Mod(moonSiderealLon, 360)
	if moonSiderealLon < 0 {
		moonSiderealLon += 360
	}

	nakshatraSpan := 360.0 / 27.0
	nakIdx := int(moonSiderealLon / nakshatraSpan)
	if nakIdx >= 27 {
		nakIdx = 26
	}

	startSlot := ashtottariNakshatraLord[nakIdx]
	totalYears := 0.0
	for _, y := range ashtottariFullYears {
		totalYears += y
	}

	// Fraction elapsed in the first nakshatra
	posInNak := moonSiderealLon - float64(nakIdx)*nakshatraSpan
	fractionUsed := posInNak / nakshatraSpan
	firstDashaFull := ashtottariFullYears[startSlot]
	firstDashaRemaining := firstDashaFull * (1 - fractionUsed)

	var periods []DashaPeriod
	age := 0.0

	periods = append(periods, DashaPeriod{
		Lord:     ashtottariSequence[startSlot],
		Years:    math.Round(firstDashaRemaining*100) / 100,
		StartAge: age,
	})
	age += firstDashaRemaining

	for i := 1; i < len(ashtottariSequence); i++ {
		slot := (startSlot + i) % len(ashtottariSequence)
		years := ashtottariFullYears[slot]
		periods = append(periods, DashaPeriod{
			Lord:     ashtottariSequence[slot],
			Years:    years,
			StartAge: math.Round(age*100) / 100,
		})
		age += years
	}

	return periods
}

// ---- Yogini Dasha (36-year cycle) ----

// YoginiName represents one of the 8 Yogini Dasha periods
type YoginiName string

const (
	YoginiMangala  YoginiName = "Mangala" // Moon (1 year)
	YoginiPingala  YoginiName = "Pingala" // Sun (2 years)
	YoginiBhramari  YoginiName = "Bhramari"  // Jupiter (3 years)
	YoginiBhadrika  YoginiName = "Bhadrika"  // Mercury (4 years)
	YoginiUlka      YoginiName = "Ulka"      // Saturn (5 years)
	YoginiSiddha    YoginiName = "Siddha"    // Venus (6 years)
	YoginiSankata   YoginiName = "Sankata"   // Rahu/Saturn (7 years)
	YoginaDhanya YoginiName = "Dhanya" // Mars (8 years)
)

// YoginiDashaPeriod extends DashaPeriod with the Yogini name
type YoginiDashaPeriod struct {
	DashaPeriod
	Yogini YoginiName `json:"yogini"`
}

// yoginiData pairs each of the 8 Yoginis with its ruling planet and duration.
// Sequence starts from Mangala and repeats over the 36-year total cycle.
var yoginiData = []struct {
	Yogini YoginiName
	Lord   models.PlanetID
	Years  float64
}{
	{YoginiMangala, models.PlanetMoon, 1},
	{YoginiPingala, models.PlanetSun, 2},
	{YoginiBhramari, models.PlanetJupiter, 3},
	{YoginiBhadrika, models.PlanetMercury, 4},
	{YoginiUlka, models.PlanetSaturn, 5},
	{YoginiSiddha, models.PlanetVenus, 6},
	{YoginiSankata, models.PlanetSaturn, 7}, // Rahu -> Saturn proxy
	{YoginaDhanya, models.PlanetMars, 8},
}

// yoginiNakshatraStart maps Nakshatra index (0-26) to Yogini slot index (0-7).
// Each of the 27 Nakshatras is assigned to one of the 8 Yoginis in a repeating cycle.
var yoginiNakshatraStart = [27]int{
	0, 1, 2, 3, 4, 5, 6, 7, // Naks 0-7
	0, 1, 2, 3, 4, 5, 6, 7, // Naks 8-15
	0, 1, 2, 3, 4, 5, 6, 7, // Naks 16-23
	0, 1, 2, // Naks 24-26
}

// CalcYoginiDasha calculates the Yogini Dasha sequence (36-year cycle) from birth.
// Each cycle repeats every 36 years (sum of 1+2+3+4+5+6+7+8).
func CalcYoginiDasha(moonSiderealLon float64) []YoginiDashaPeriod {
	moonSiderealLon = math.Mod(moonSiderealLon, 360)
	if moonSiderealLon < 0 {
		moonSiderealLon += 360
	}

	nakshatraSpan := 360.0 / 27.0
	nakIdx := int(moonSiderealLon / nakshatraSpan)
	if nakIdx >= 27 {
		nakIdx = 26
	}

	startSlot := yoginiNakshatraStart[nakIdx]

	posInNak := moonSiderealLon - float64(nakIdx)*nakshatraSpan
	fractionUsed := posInNak / nakshatraSpan
	firstFull := yoginiData[startSlot].Years
	firstRemaining := firstFull * (1 - fractionUsed)

	var periods []YoginiDashaPeriod
	age := 0.0

	periods = append(periods, YoginiDashaPeriod{
		DashaPeriod: DashaPeriod{
			Lord:     yoginiData[startSlot].Lord,
			Years:    math.Round(firstRemaining*100) / 100,
			StartAge: age,
		},
		Yogini: yoginiData[startSlot].Yogini,
	})
	age += firstRemaining

	for i := 1; i < len(yoginiData); i++ {
		slot := (startSlot + i) % len(yoginiData)
		d := yoginiData[slot]
		periods = append(periods, YoginiDashaPeriod{
			DashaPeriod: DashaPeriod{
				Lord:     d.Lord,
				Years:    d.Years,
				StartAge: math.Round(age*100) / 100,
			},
			Yogini: d.Yogini,
		})
		age += d.Years
	}

	return periods
}

// ---- Chara (Jaimini) Dasha ----

// CharaDashaPeriod represents one sign-based dasha period in the Jaimini system.
type CharaDashaPeriod struct {
	Sign     string  `json:"sign"`      // Zodiac sign name
	SignIdx  int     `json:"sign_index"` // 0-11
	Years    int     `json:"years"`     // Duration in years (1-12)
	StartAge float64 `json:"start_age"`
}

// CalcCharaDasha computes the Chara (Jaimini) Dasha sequence based on the
// Lagna (Ascendant) sign. Chara Dasha is sign-based rather than planet-based.
//
// Rules (BPHS Jaimini Sutras):
//   - The sequence starts from the Lagna sign.
//   - For odd signs (Aries, Gemini, Leo, Libra, Sagittarius, Aquarius),
//     the sequence proceeds in direct order (zodiacal).
//   - For even signs (Taurus, Cancer, Virgo, Scorpio, Capricorn, Pisces),
//     the sequence proceeds in reverse order (anti-zodiacal).
//   - Duration of each sign's dasha = distance from the sign's ruler to the sign,
//     counting inclusively. If the ruler is in its own sign, the duration is 12 years.
//
// Parameters:
//   - lagnaSignIdx: sidereal sign index of the Ascendant (0=Aries, 11=Pisces)
//   - planetSignPositions: map of planet → sign index (0-11), sidereal
func CalcCharaDasha(lagnaSignIdx int, planetSignPositions map[models.PlanetID]int) []CharaDashaPeriod {
	isOdd := lagnaSignIdx%2 == 0 // 0=Aries(odd), 1=Taurus(even), etc.

	var periods []CharaDashaPeriod
	age := 0.0

	for i := 0; i < 12; i++ {
		var signIdx int
		if isOdd {
			signIdx = (lagnaSignIdx + i) % 12
		} else {
			signIdx = (lagnaSignIdx - i + 12) % 12
		}

		ruler := charaSignRuler[signIdx]
		rulerSign, ok := planetSignPositions[ruler]
		if !ok {
			rulerSign = signIdx // fallback: ruler in own sign
		}

		years := charaDuration(signIdx, rulerSign, signIdx%2 == 0)

		periods = append(periods, CharaDashaPeriod{
			Sign:     models.ZodiacSigns[signIdx],
			SignIdx:  signIdx,
			Years:    years,
			StartAge: math.Round(age*100) / 100,
		})
		age += float64(years)
	}

	return periods
}

// charaSignRuler maps each sign to its Jaimini ruler.
// Jaimini uses the same rulers as Parashari for most signs, but some
// traditions assign co-rulers. We use the standard BPHS assignment.
var charaSignRuler = [12]models.PlanetID{
	models.PlanetMars,    // 0 Aries
	models.PlanetVenus,   // 1 Taurus
	models.PlanetMercury, // 2 Gemini
	models.PlanetMoon,    // 3 Cancer
	models.PlanetSun,     // 4 Leo
	models.PlanetMercury, // 5 Virgo
	models.PlanetVenus,   // 6 Libra
	models.PlanetMars,    // 7 Scorpio
	models.PlanetJupiter, // 8 Sagittarius
	models.PlanetSaturn,  // 9 Capricorn
	models.PlanetSaturn,  // 10 Aquarius
	models.PlanetJupiter, // 11 Pisces
}

// charaDuration calculates the Chara Dasha period length for a sign.
// For odd signs: count from sign to ruler's sign (direct, inclusive).
// For even signs: count from sign to ruler's sign (reverse, inclusive).
// If ruler is in its own sign, the duration is 12 years.
func charaDuration(signIdx, rulerSignIdx int, isOddSign bool) int {
	if signIdx == rulerSignIdx {
		return 12
	}

	var dist int
	if isOddSign {
		// Direct counting: how many signs from signIdx to rulerSignIdx, going forward
		dist = (rulerSignIdx - signIdx + 12) % 12
	} else {
		// Reverse counting: how many signs from signIdx to rulerSignIdx, going backward
		dist = (signIdx - rulerSignIdx + 12) % 12
	}

	if dist == 0 {
		dist = 12
	}

	return dist
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
