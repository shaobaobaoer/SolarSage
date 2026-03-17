package models

import "github.com/anthropic/swisseph-mcp/pkg/sweph"

// PlanetToSweID maps PlanetID to Swiss Ephemeris planet number
func PlanetToSweID(p PlanetID) (int, bool) {
	m := map[PlanetID]int{
		PlanetSun:           sweph.SE_SUN,
		PlanetMoon:          sweph.SE_MOON,
		PlanetMercury:       sweph.SE_MERCURY,
		PlanetVenus:         sweph.SE_VENUS,
		PlanetMars:          sweph.SE_MARS,
		PlanetJupiter:       sweph.SE_JUPITER,
		PlanetSaturn:        sweph.SE_SATURN,
		PlanetUranus:        sweph.SE_URANUS,
		PlanetNeptune:       sweph.SE_NEPTUNE,
		PlanetPluto:         sweph.SE_PLUTO,
		PlanetChiron:        sweph.SE_CHIRON,
		PlanetNorthNodeTrue: sweph.SE_TRUE_NODE,
		PlanetNorthNodeMean: sweph.SE_MEAN_NODE,
		PlanetLilithMean:    sweph.SE_MEAN_APOG,
		PlanetLilithTrue:    sweph.SE_OSCU_APOG,
	}
	id, ok := m[p]
	return id, ok
}

// HouseSystemToChar maps HouseSystem to Swiss Ephemeris char code
func HouseSystemToChar(hs HouseSystem) int {
	m := map[HouseSystem]int{
		HousePlacidus:      sweph.HousePlacidus,
		HouseKoch:          sweph.HouseKoch,
		HouseEqual:         sweph.HouseEqual,
		HouseWholeSign:     sweph.HouseWholeSign,
		HouseCampanus:      sweph.HouseCampanus,
		HouseRegiomontanus: sweph.HouseRegiomontanus,
		HousePorphyry:      sweph.HousePorphyry,
	}
	if c, ok := m[hs]; ok {
		return c
	}
	return sweph.HousePlacidus // default
}

// AllPlanets returns all standard planets
var AllPlanets = []PlanetID{
	PlanetSun, PlanetMoon, PlanetMercury, PlanetVenus, PlanetMars,
	PlanetJupiter, PlanetSaturn, PlanetUranus, PlanetNeptune, PlanetPluto,
	PlanetChiron, PlanetNorthNodeTrue, PlanetNorthNodeMean,
	PlanetSouthNode, PlanetLilithMean, PlanetLilithTrue,
}
