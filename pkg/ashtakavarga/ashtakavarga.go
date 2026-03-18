package ashtakavarga

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/vedic"
)

// BinduTable holds the per-sign bindu count for one planet.
type BinduTable struct {
	Planet models.PlanetID `json:"planet"`
	Bindus [12]int         `json:"bindus"` // Points per sign (0-11 = Aries-Pisces sidereal)
	Total  int             `json:"total"`  // Sum of all bindus (max 48)
}

// AshtakavargaResult holds the complete Ashtakavarga analysis.
type AshtakavargaResult struct {
	PlanetTables []BinduTable `json:"planet_tables"` // 7 planets (Sun-Saturn)
	SAV          [12]int      `json:"sarvashtakavarga"`
	SAVTotal     int          `json:"sav_total"`
}

// traditionalPlanets are the seven traditional planets used in Ashtakavarga.
var traditionalPlanets = []models.PlanetID{
	models.PlanetSun, models.PlanetMoon, models.PlanetMars,
	models.PlanetMercury, models.PlanetJupiter, models.PlanetVenus,
	models.PlanetSaturn,
}

// contributor represents either a planet or the ascendant.
type contributor struct {
	id    string // planet ID or "ASC"
	isASC bool
}

// binduHouses maps [planet receiving bindus][contributing planet/ASC] -> benefic houses (1-based from contributor).
// Source: Brihat Parashara Hora Shastra (BPHS).
var binduHouses = map[models.PlanetID]map[string][]int{
	models.PlanetSun: {
		string(models.PlanetSun):     {1, 2, 4, 7, 8, 9, 10, 11},
		string(models.PlanetMoon):    {3, 6, 10, 11},
		string(models.PlanetMars):    {1, 2, 4, 7, 8, 9, 10, 11},
		string(models.PlanetMercury): {3, 5, 6, 9, 10, 11, 12},
		string(models.PlanetJupiter): {5, 6, 9, 11},
		string(models.PlanetVenus):   {6, 7, 12},
		string(models.PlanetSaturn):  {1, 2, 4, 7, 8, 9, 10, 11},
		"ASC":                        {3, 4, 6, 10, 11, 12},
	},
	models.PlanetMoon: {
		string(models.PlanetSun):     {3, 6, 7, 8, 10, 11},
		string(models.PlanetMoon):    {1, 3, 6, 7, 10, 11},
		string(models.PlanetMars):    {2, 3, 5, 6, 9, 10, 11},
		string(models.PlanetMercury): {1, 3, 4, 5, 7, 8, 10, 11},
		string(models.PlanetJupiter): {1, 4, 7, 8, 10, 11, 12},
		string(models.PlanetVenus):   {3, 4, 5, 7, 9, 10, 11},
		string(models.PlanetSaturn):  {3, 5, 6, 11},
		"ASC":                        {3, 6, 10, 11},
	},
	models.PlanetMars: {
		string(models.PlanetSun):     {3, 5, 6, 10, 11},
		string(models.PlanetMoon):    {3, 6, 11},
		string(models.PlanetMars):    {1, 2, 4, 7, 8, 10, 11},
		string(models.PlanetMercury): {3, 5, 6, 11},
		string(models.PlanetJupiter): {6, 10, 11, 12},
		string(models.PlanetVenus):   {6, 8, 11, 12},
		string(models.PlanetSaturn):  {1, 4, 7, 8, 9, 10, 11},
		"ASC":                        {1, 3, 6, 10, 11},
	},
	models.PlanetMercury: {
		string(models.PlanetSun):     {5, 6, 9, 11, 12},
		string(models.PlanetMoon):    {2, 4, 6, 8, 10, 11},
		string(models.PlanetMars):    {1, 2, 4, 7, 8, 9, 10, 11},
		string(models.PlanetMercury): {1, 3, 5, 6, 9, 10, 11, 12},
		string(models.PlanetJupiter): {6, 8, 11, 12},
		string(models.PlanetVenus):   {1, 2, 3, 4, 5, 8, 9, 11},
		string(models.PlanetSaturn):  {1, 2, 4, 7, 8, 9, 10, 11},
		"ASC":                        {1, 2, 4, 6, 8, 10, 11},
	},
	models.PlanetJupiter: {
		string(models.PlanetSun):     {1, 2, 3, 4, 7, 8, 9, 10, 11},
		string(models.PlanetMoon):    {2, 5, 7, 9, 11},
		string(models.PlanetMars):    {1, 2, 4, 7, 8, 10, 11},
		string(models.PlanetMercury): {1, 2, 4, 5, 6, 9, 10, 11},
		string(models.PlanetJupiter): {1, 2, 3, 4, 7, 8, 10, 11},
		string(models.PlanetVenus):   {2, 5, 6, 9, 10, 11},
		string(models.PlanetSaturn):  {3, 5, 6, 12},
		"ASC":                        {1, 2, 4, 5, 6, 7, 9, 10, 11},
	},
	models.PlanetVenus: {
		string(models.PlanetSun):     {8, 11, 12},
		string(models.PlanetMoon):    {1, 2, 3, 4, 5, 8, 9, 11, 12},
		string(models.PlanetMars):    {3, 4, 6, 8, 9, 11, 12},
		string(models.PlanetMercury): {3, 5, 6, 9, 11},
		string(models.PlanetJupiter): {5, 8, 9, 10, 11},
		string(models.PlanetVenus):   {1, 2, 3, 4, 5, 8, 9, 10, 11},
		string(models.PlanetSaturn):  {3, 4, 5, 8, 9, 10, 11},
		"ASC":                        {1, 2, 3, 4, 5, 8, 9, 11},
	},
	models.PlanetSaturn: {
		string(models.PlanetSun):     {1, 2, 4, 7, 8, 9, 10, 11},
		string(models.PlanetMoon):    {3, 6, 11},
		string(models.PlanetMars):    {3, 5, 6, 10, 11, 12},
		string(models.PlanetMercury): {6, 8, 9, 10, 11, 12},
		string(models.PlanetJupiter): {5, 6, 11, 12},
		string(models.PlanetVenus):   {6, 11, 12},
		string(models.PlanetSaturn):  {3, 5, 6, 11},
		"ASC":                        {1, 3, 4, 6, 10, 11},
	},
}

// signIndexFromLon returns the sidereal sign index (0=Aries .. 11=Pisces).
func signIndexFromLon(lon float64) int {
	idx := int(lon / 30.0)
	if idx < 0 {
		idx += 12
	}
	if idx > 11 {
		idx = 11
	}
	return idx
}

// CalcAshtakavarga computes the Ashtakavarga tables for the seven traditional
// planets. positions must contain sidereal longitudes; siderealASC is the
// sidereal ascendant longitude.
func CalcAshtakavarga(positions []vedic.SiderealPosition, siderealASC float64) *AshtakavargaResult {
	// Build a sign lookup: contributor ID -> sign index (0-11)
	contribSign := make(map[string]int)
	for _, p := range positions {
		contribSign[string(p.PlanetID)] = signIndexFromLon(p.SiderealLon)
	}
	contribSign["ASC"] = signIndexFromLon(siderealASC)

	result := &AshtakavargaResult{}

	for _, planet := range traditionalPlanets {
		table := BinduTable{Planet: planet}
		houses, ok := binduHouses[planet]
		if !ok {
			continue
		}

		for contribID, beneficHouses := range houses {
			fromSign, exists := contribSign[contribID]
			if !exists {
				continue
			}
			for _, h := range beneficHouses {
				// House 1 = contributor's own sign, house 2 = next sign, etc.
				targetSign := (fromSign + h - 1) % 12
				table.Bindus[targetSign]++
			}
		}

		total := 0
		for _, b := range table.Bindus {
			total += b
		}
		table.Total = total

		result.PlanetTables = append(result.PlanetTables, table)
	}

	// Compute Sarvashtakavarga (SAV): sum of all planet tables per sign
	for _, t := range result.PlanetTables {
		for i := 0; i < 12; i++ {
			result.SAV[i] += t.Bindus[i]
		}
	}
	savTotal := 0
	for _, v := range result.SAV {
		savTotal += v
	}
	result.SAVTotal = savTotal

	return result
}
