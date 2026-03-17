package models

// PlanetID represents a celestial body identifier
type PlanetID string

const (
	PlanetSun           PlanetID = "SUN"
	PlanetMoon          PlanetID = "MOON"
	PlanetMercury       PlanetID = "MERCURY"
	PlanetVenus         PlanetID = "VENUS"
	PlanetMars          PlanetID = "MARS"
	PlanetJupiter       PlanetID = "JUPITER"
	PlanetSaturn        PlanetID = "SATURN"
	PlanetUranus        PlanetID = "URANUS"
	PlanetNeptune       PlanetID = "NEPTUNE"
	PlanetPluto         PlanetID = "PLUTO"
	PlanetChiron        PlanetID = "CHIRON"
	PlanetNorthNodeTrue PlanetID = "NORTH_NODE_TRUE"
	PlanetNorthNodeMean PlanetID = "NORTH_NODE_MEAN"
	PlanetSouthNode     PlanetID = "SOUTH_NODE"
	PlanetLilithMean    PlanetID = "LILITH_MEAN"
	PlanetLilithTrue    PlanetID = "LILITH_TRUE"
)

// SpecialPointID represents a derived astrological point
type SpecialPointID string

const (
	PointASC        SpecialPointID = "ASC"
	PointMC         SpecialPointID = "MC"
	PointDSC        SpecialPointID = "DSC"
	PointIC         SpecialPointID = "IC"
	PointVertex     SpecialPointID = "VERTEX"
	PointAntiVertex SpecialPointID = "ANTI_VERTEX"
	PointEastPoint  SpecialPointID = "EAST_POINT"
	PointLotFortune SpecialPointID = "LOT_FORTUNE"
	PointLotSpirit  SpecialPointID = "LOT_SPIRIT"
)

// HouseSystem represents a house system type
type HouseSystem string

const (
	HousePlacidus      HouseSystem = "PLACIDUS"
	HouseKoch          HouseSystem = "KOCH"
	HouseEqual         HouseSystem = "EQUAL"
	HouseWholeSign     HouseSystem = "WHOLE_SIGN"
	HouseCampanus      HouseSystem = "CAMPANUS"
	HouseRegiomontanus HouseSystem = "REGIOMONTANUS"
	HousePorphyry      HouseSystem = "PORPHYRY"
)

// CalendarType represents calendar type
type CalendarType string

const (
	CalendarGregorian CalendarType = "GREGORIAN"
	CalendarJulian    CalendarType = "JULIAN"
)

// AspectType represents a type of aspect
type AspectType string

const (
	AspectConjunction    AspectType = "合相"
	AspectOpposition     AspectType = "对分相"
	AspectTrine          AspectType = "三分相"
	AspectSquare         AspectType = "刑相"
	AspectSextile        AspectType = "六分相"
	AspectQuincunx       AspectType = "补十二分相"
	AspectSemiSextile    AspectType = "十二分相"
	AspectSemiSquare     AspectType = "八分相"
	AspectSesquiquadrate AspectType = "倍半刑"
)

// AspectDef defines a standard aspect with its angle
type AspectDef struct {
	Type  AspectType
	Angle float64
}

// StandardAspects lists all standard aspects
var StandardAspects = []AspectDef{
	{AspectConjunction, 0},
	{AspectOpposition, 180},
	{AspectTrine, 120},
	{AspectSquare, 90},
	{AspectSextile, 60},
	{AspectQuincunx, 150},
	{AspectSemiSextile, 30},
	{AspectSemiSquare, 45},
	{AspectSesquiquadrate, 135},
}

// OrbConfig holds the orb (tolerance) for each aspect type
type OrbConfig struct {
	Conjunction    float64 `json:"conjunction"`
	Opposition     float64 `json:"opposition"`
	Trine          float64 `json:"trine"`
	Square         float64 `json:"square"`
	Sextile        float64 `json:"sextile"`
	Quincunx       float64 `json:"quincunx"`
	SemiSextile    float64 `json:"semi_sextile"`
	SemiSquare     float64 `json:"semi_square"`
	Sesquiquadrate float64 `json:"sesquiquadrate"`
}

// DefaultOrbConfig returns default orb values
func DefaultOrbConfig() OrbConfig {
	return OrbConfig{
		Conjunction:    8,
		Opposition:     8,
		Trine:          7,
		Square:         7,
		Sextile:        5,
		Quincunx:       3,
		SemiSextile:    2,
		SemiSquare:     2,
		Sesquiquadrate: 2,
	}
}

// GetOrb returns the orb for a given aspect type
func (o OrbConfig) GetOrb(at AspectType) float64 {
	switch at {
	case AspectConjunction:
		return o.Conjunction
	case AspectOpposition:
		return o.Opposition
	case AspectTrine:
		return o.Trine
	case AspectSquare:
		return o.Square
	case AspectSextile:
		return o.Sextile
	case AspectQuincunx:
		return o.Quincunx
	case AspectSemiSextile:
		return o.SemiSextile
	case AspectSemiSquare:
		return o.SemiSquare
	case AspectSesquiquadrate:
		return o.Sesquiquadrate
	default:
		return 0
	}
}

// EventType represents a transit event type
type EventType string

const (
	EventAspectEnter  EventType = "ASPECT_ENTER"
	EventAspectExact  EventType = "ASPECT_EXACT"
	EventAspectLeave  EventType = "ASPECT_LEAVE"
	EventSignIngress  EventType = "SIGN_INGRESS"
	EventHouseIngress EventType = "HOUSE_INGRESS"
	EventStation      EventType = "STATION"
)

// StationType represents retrograde/direct station
type StationType string

const (
	StationRetrograde StationType = "RETROGRADE"
	StationDirect     StationType = "DIRECT"
)

// PlanetPosition holds calculated position data for a planet
type PlanetPosition struct {
	PlanetID    PlanetID `json:"planet_id"`
	Longitude   float64  `json:"longitude"`
	Latitude    float64  `json:"latitude"`
	Speed       float64  `json:"speed"`
	IsRetrograde bool    `json:"is_retrograde"`
	Sign        string   `json:"sign"`
	SignDegree  float64  `json:"sign_degree"`
	House       int      `json:"house"`
}

// AnglesInfo holds the four angles
type AnglesInfo struct {
	ASC float64 `json:"asc"`
	MC  float64 `json:"mc"`
	DSC float64 `json:"dsc"`
	IC  float64 `json:"ic"`
}

// AspectInfo holds aspect data between two bodies
type AspectInfo struct {
	PlanetA     string     `json:"planet_a"`
	PlanetB     string     `json:"planet_b"`
	AspectType  AspectType `json:"aspect_type"`
	AspectAngle float64    `json:"aspect_angle"`
	ActualAngle float64    `json:"actual_angle"`
	Orb         float64    `json:"orb"`
	IsApplying  bool       `json:"is_applying"`
}

// ChartInfo holds complete chart data
type ChartInfo struct {
	Planets []PlanetPosition `json:"planets"`
	Houses  []float64        `json:"houses"`
	Angles  AnglesInfo       `json:"angles"`
	Aspects []AspectInfo     `json:"aspects"`
}

// CrossAspectInfo holds aspect data between two charts
type CrossAspectInfo struct {
	InnerBody   string     `json:"inner_body"`
	OuterBody   string     `json:"outer_body"`
	AspectType  AspectType `json:"aspect_type"`
	AspectAngle float64    `json:"aspect_angle"`
	ActualAngle float64    `json:"actual_angle"`
	Orb         float64    `json:"orb"`
	IsApplying  bool       `json:"is_applying"`
}

// GeoLocation holds geographic coordinates
type GeoLocation struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	DisplayName string  `json:"display_name"`
}

// JulianDayResult holds Julian Day conversion result
type JulianDayResult struct {
	JDUT float64 `json:"jd_ut"`
	JDTT float64 `json:"jd_tt"`
}

// SpecialPointsConfig configures which special points to include
type SpecialPointsConfig struct {
	InnerPoints []SpecialPointID `json:"inner_points,omitempty"`
	OuterPoints []SpecialPointID `json:"outer_points,omitempty"`
	NatalPoints   []SpecialPointID `json:"natal_points,omitempty"`
	TransitPoints []SpecialPointID `json:"transit_points,omitempty"`
}

// EventConfig configures which event types to include
type EventConfig struct {
	IncludeAspects      bool `json:"include_aspects"`
	IncludeSignIngress  bool `json:"include_sign_ingress"`
	IncludeHouseIngress bool `json:"include_house_ingress"`
	IncludeStation      bool `json:"include_station"`
}

// DefaultEventConfig returns config with all events enabled
func DefaultEventConfig() EventConfig {
	return EventConfig{
		IncludeAspects:      true,
		IncludeSignIngress:  true,
		IncludeHouseIngress: true,
		IncludeStation:      true,
	}
}

// TransitEvent represents an astrological transit event
type TransitEvent struct {
	EventType      EventType  `json:"event_type"`
	Planet         PlanetID   `json:"planet"`
	JD             float64    `json:"jd"`
	PlanetLongitude float64   `json:"planet_longitude"`
	PlanetSign     string     `json:"planet_sign"`
	PlanetHouse    int        `json:"planet_house"`
	IsRetrograde   bool       `json:"is_retrograde"`

	// Aspect events
	Target      string     `json:"target,omitempty"`
	AspectType  AspectType `json:"aspect_type,omitempty"`
	AspectAngle float64    `json:"aspect_angle,omitempty"`
	OrbAtEnter  float64    `json:"orb_at_enter,omitempty"`
	OrbAtLeave  float64    `json:"orb_at_leave,omitempty"`
	ExactCount  int        `json:"exact_count,omitempty"`

	// Sign ingress
	FromSign string `json:"from_sign,omitempty"`
	ToSign   string `json:"to_sign,omitempty"`

	// House ingress
	FromHouse int `json:"from_house,omitempty"`
	ToHouse   int `json:"to_house,omitempty"`

	// Station
	StationType StationType `json:"station_type,omitempty"`
}

// ZodiacSigns maps sign index (0-11) to Chinese name
var ZodiacSigns = []string{
	"白羊座", "金牛座", "双子座", "巨蟹座",
	"狮子座", "处女座", "天秤座", "天蝎座",
	"射手座", "摩羯座", "水瓶座", "双鱼座",
}

// SignFromLongitude returns the zodiac sign name for a given ecliptic longitude
func SignFromLongitude(lon float64) string {
	idx := int(lon / 30.0)
	if idx < 0 {
		idx = 0
	}
	if idx > 11 {
		idx = 11
	}
	return ZodiacSigns[idx]
}

// SignDegreeFromLongitude returns the degree within the sign (0-30)
func SignDegreeFromLongitude(lon float64) float64 {
	return lon - float64(int(lon/30.0))*30.0
}
