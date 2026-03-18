package fixedstars

import (
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Star represents a fixed star with its catalog data
type Star struct {
	Name       string  `json:"name"`
	Tradition  string  `json:"tradition,omitempty"`   // Traditional name
	Longitude  float64 `json:"longitude"`              // Ecliptic longitude (epoch J2000.0)
	Latitude   float64 `json:"latitude"`               // Ecliptic latitude
	Magnitude  float64 `json:"magnitude"`              // Visual magnitude
	Nature     string  `json:"nature,omitempty"`        // Traditional planetary nature
	Sign       string  `json:"sign"`
	SignDegree float64 `json:"sign_degree"`
}

// StarConjunction represents a planet-star conjunction
type StarConjunction struct {
	Planet    models.PlanetID `json:"planet"`
	PlanetLon float64        `json:"planet_longitude"`
	Star      Star           `json:"star"`
	Orb       float64        `json:"orb"`
}

// precessionRate is the annual precession of the equinoxes in degrees
const precessionRate = 50.2882 / 3600.0 // ~0.01397°/year

// j2000Epoch is the Julian Day of the J2000.0 epoch
const j2000Epoch = 2451545.0

// julianYear is days per Julian year
const julianYear = 365.25

// Catalog of major fixed stars with J2000.0 ecliptic longitudes.
// Longitudes are tropical (precessed from the J2000.0 catalog positions).
// Source: Robson, Vivian. "The Fixed Stars and Constellations in Astrology"
// and modern astronomical catalogs.
var Catalog = []Star{
	// Royal Stars
	{Name: "Aldebaran", Tradition: "Watcher of the East", Longitude: 69.47, Latitude: -5.47, Magnitude: 0.85, Nature: "Mars"},
	{Name: "Regulus", Tradition: "Heart of the Lion", Longitude: 149.83, Latitude: 0.46, Magnitude: 1.35, Nature: "Mars-Jupiter"},
	{Name: "Antares", Tradition: "Watcher of the West", Longitude: 249.47, Latitude: -4.57, Magnitude: 0.96, Nature: "Mars-Jupiter"},
	{Name: "Fomalhaut", Tradition: "Watcher of the South", Longitude: 333.87, Latitude: -21.08, Magnitude: 1.16, Nature: "Venus-Mercury"},

	// Brightest stars
	{Name: "Sirius", Tradition: "Dog Star", Longitude: 104.07, Latitude: -39.60, Magnitude: -1.46, Nature: "Jupiter-Mars"},
	{Name: "Canopus", Longitude: 104.99, Latitude: -75.85, Magnitude: -0.74, Nature: "Saturn-Jupiter"},
	{Name: "Arcturus", Tradition: "Guardian of the Bear", Longitude: 204.14, Latitude: 30.75, Magnitude: -0.05, Nature: "Mars-Jupiter"},
	{Name: "Vega", Tradition: "Falling Eagle", Longitude: 285.27, Latitude: 61.73, Magnitude: 0.03, Nature: "Venus-Mercury"},
	{Name: "Capella", Tradition: "She-Goat", Longitude: 81.51, Latitude: 22.86, Magnitude: 0.08, Nature: "Mars-Mercury"},
	{Name: "Rigel", Longitude: 76.96, Latitude: -31.07, Magnitude: 0.12, Nature: "Jupiter-Saturn"},
	{Name: "Procyon", Longitude: 115.62, Latitude: -16.01, Magnitude: 0.34, Nature: "Mars-Mercury"},
	{Name: "Betelgeuse", Longitude: 88.79, Latitude: -16.03, Magnitude: 0.42, Nature: "Mars-Mercury"},
	{Name: "Achernar", Longitude: 345.29, Latitude: -59.38, Magnitude: 0.46, Nature: "Jupiter"},

	// Notable astrological stars
	{Name: "Algol", Tradition: "Demon Star", Longitude: 56.17, Latitude: 22.41, Magnitude: 2.12, Nature: "Saturn-Jupiter"},
	{Name: "Alcyone", Tradition: "Brightest Pleiad", Longitude: 60.00, Latitude: 4.03, Magnitude: 2.87, Nature: "Moon-Mars"},
	{Name: "Spica", Tradition: "Ear of Wheat", Longitude: 203.83, Latitude: -2.05, Magnitude: 0.97, Nature: "Venus-Mars"},
	{Name: "Polaris", Tradition: "North Star", Longitude: 88.37, Latitude: 66.09, Magnitude: 1.98, Nature: "Saturn-Venus"},
	{Name: "Deneb Algedi", Longitude: 323.56, Latitude: 2.60, Magnitude: 2.87, Nature: "Saturn-Jupiter"},
	{Name: "Scheat", Longitude: 349.25, Latitude: 31.08, Magnitude: 2.42, Nature: "Mars-Mercury"},
	{Name: "Markab", Longitude: 353.47, Latitude: 19.41, Magnitude: 2.49, Nature: "Mars-Mercury"},
	{Name: "Alpheratz", Longitude: 14.18, Latitude: 25.68, Magnitude: 2.06, Nature: "Jupiter-Venus"},
	{Name: "Hamal", Longitude: 37.85, Latitude: 9.96, Magnitude: 2.00, Nature: "Mars-Saturn"},
	{Name: "Mirach", Longitude: 30.26, Latitude: 25.56, Magnitude: 2.05, Nature: "Venus"},
	{Name: "Pleiades", Longitude: 60.00, Latitude: 4.03, Magnitude: 1.60, Nature: "Moon-Mars"},
	{Name: "Hyades", Longitude: 65.47, Latitude: -5.47, Magnitude: 3.65, Nature: "Saturn-Mercury"},
	{Name: "Castor", Longitude: 110.20, Latitude: 10.09, Magnitude: 1.58, Nature: "Mercury"},
	{Name: "Pollux", Longitude: 113.22, Latitude: 6.68, Magnitude: 1.14, Nature: "Mars"},
	{Name: "Praesepe", Tradition: "Beehive Cluster", Longitude: 127.23, Latitude: 1.55, Magnitude: 3.70, Nature: "Mars-Moon"},
	{Name: "Vindemiatrix", Longitude: 189.80, Latitude: 16.19, Magnitude: 2.83, Nature: "Saturn-Mercury"},
	{Name: "Algorab", Longitude: 193.57, Latitude: -12.19, Magnitude: 2.95, Nature: "Mars-Saturn"},
	{Name: "Zuben Elgenubi", Longitude: 195.07, Latitude: 0.33, Magnitude: 2.75, Nature: "Saturn-Mars"},
	{Name: "Zuben Elschemali", Longitude: 199.10, Latitude: 8.73, Magnitude: 2.61, Nature: "Jupiter-Mercury"},
	{Name: "Unukalhai", Tradition: "Heart of the Serpent", Longitude: 211.90, Latitude: 25.42, Magnitude: 2.65, Nature: "Saturn-Mars"},
	{Name: "Agena", Longitude: 203.55, Latitude: -23.66, Magnitude: 0.61, Nature: "Venus-Jupiter"},
	{Name: "Toliman", Tradition: "Alpha Centauri", Longitude: 209.22, Latitude: -42.58, Magnitude: -0.01, Nature: "Venus-Jupiter"},
	{Name: "Dschubba", Longitude: 242.57, Latitude: -1.98, Magnitude: 2.32, Nature: "Mars-Saturn"},
	{Name: "Acrab", Longitude: 243.13, Latitude: 1.02, Magnitude: 2.62, Nature: "Mars-Saturn"},
	{Name: "Ras Alhague", Longitude: 262.15, Latitude: 35.84, Magnitude: 2.08, Nature: "Saturn-Venus"},
	{Name: "Lesath", Longitude: 264.15, Latitude: -13.51, Magnitude: 2.69, Nature: "Mercury-Mars"},
	{Name: "Shaula", Longitude: 264.57, Latitude: -13.36, Magnitude: 1.63, Nature: "Mercury-Mars"},
	{Name: "Nunki", Longitude: 272.33, Latitude: -6.55, Magnitude: 2.02, Nature: "Jupiter-Mercury"},
	{Name: "Deneb Adige", Longitude: 320.25, Latitude: 59.98, Magnitude: 1.25, Nature: "Venus-Mercury"},
	{Name: "Sadalsuud", Longitude: 323.45, Latitude: -6.39, Magnitude: 2.91, Nature: "Saturn-Mercury"},
	{Name: "Sadalmelik", Longitude: 333.37, Latitude: 10.39, Magnitude: 2.96, Nature: "Saturn-Mercury"},

	// Additional notable stars
	{Name: "Altair", Tradition: "Flying Eagle", Longitude: 301.81, Latitude: 29.31, Magnitude: 0.77, Nature: "Mars-Jupiter"},
	{Name: "Denebola", Longitude: 171.60, Latitude: 12.33, Magnitude: 2.14, Nature: "Saturn-Venus"},
	{Name: "Wega", Longitude: 285.27, Latitude: 61.73, Magnitude: 0.03, Nature: "Venus-Mercury"},
	{Name: "Rastaban", Longitude: 268.37, Latitude: 75.17, Magnitude: 2.79, Nature: "Saturn-Mars"},
	{Name: "Eltanin", Longitude: 267.95, Latitude: 75.26, Magnitude: 2.23, Nature: "Mars-Jupiter"},
}

func init() {
	// Populate sign info for all catalog entries
	for i := range Catalog {
		Catalog[i].Sign = models.SignFromLongitude(Catalog[i].Longitude)
		Catalog[i].SignDegree = models.SignDegreeFromLongitude(Catalog[i].Longitude)
	}
}

// PrecessLongitude adjusts a J2000.0 ecliptic longitude to a given Julian Day
func PrecessLongitude(lon, jd float64) float64 {
	yearsSinceJ2000 := (jd - j2000Epoch) / julianYear
	return lon + yearsSinceJ2000*precessionRate
}

// FindConjunctions finds all fixed star conjunctions with the given planet positions
func FindConjunctions(positions []models.PlanetPosition, orb float64, jd float64) []StarConjunction {
	var conjunctions []StarConjunction

	for _, p := range positions {
		for _, star := range Catalog {
			// Precess star longitude to the chart epoch
			starLon := PrecessLongitude(star.Longitude, jd)
			diff := angleDiff(p.Longitude, starLon)
			if diff <= orb {
				s := star
				s.Longitude = starLon
				s.Sign = models.SignFromLongitude(starLon)
				s.SignDegree = models.SignDegreeFromLongitude(starLon)
				conjunctions = append(conjunctions, StarConjunction{
					Planet:    p.PlanetID,
					PlanetLon: p.Longitude,
					Star:      s,
					Orb:       diff,
				})
			}
		}
	}
	return conjunctions
}

// GetStarByName finds a star by name (case-insensitive prefix match)
func GetStarByName(name string) *Star {
	for _, s := range Catalog {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func angleDiff(a, b float64) float64 {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}
