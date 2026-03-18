package synastry

import (
	"math"
	"sort"

	"github.com/shaobaobaoer/solarsage-mcp/internal/aspect"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// SynastryScore holds the relationship compatibility analysis
type SynastryScore struct {
	TotalScore     float64         `json:"total_score"`
	Harmony        float64         `json:"harmony"`        // positive aspects score
	Tension        float64         `json:"tension"`         // challenging aspects score
	Compatibility  float64         `json:"compatibility"`   // 0-100 percentage
	TopAspects     []ScoredAspect  `json:"top_aspects"`
	CategoryScores []CategoryScore `json:"category_scores"`
}

// ScoredAspect is an aspect with its contribution score
type ScoredAspect struct {
	PersonA     string             `json:"person_a"`
	PersonB     string             `json:"person_b"`
	AspectType  models.AspectType  `json:"aspect_type"`
	Orb         float64            `json:"orb"`
	Score       float64            `json:"score"`
	Category    string             `json:"category"`
}

// CategoryScore groups scores by relationship category
type CategoryScore struct {
	Category string  `json:"category"`
	Score    float64 `json:"score"`
	Count    int     `json:"count"`
}

// planetWeight assigns importance weights to planets in synastry
var planetWeight = map[string]float64{
	"SUN":             1.0,
	"MOON":            1.0,
	"MERCURY":         0.6,
	"VENUS":           0.9,
	"MARS":            0.8,
	"JUPITER":         0.7,
	"SATURN":          0.7,
	"URANUS":          0.5,
	"NEPTUNE":         0.5,
	"PLUTO":           0.5,
	"CHIRON":          0.4,
	"NORTH_NODE_TRUE": 0.4,
	"ASC":             0.8,
	"MC":              0.5,
}

// aspectScore returns a base score for each aspect type
// Positive = harmonious, Negative = challenging
var aspectBaseScore = map[models.AspectType]float64{
	models.AspectConjunction:    4.0,  // powerful, can be either
	models.AspectTrine:          3.0,  // harmonious
	models.AspectSextile:        2.0,  // mildly harmonious
	models.AspectOpposition:     -2.0, // tension, awareness
	models.AspectSquare:         -3.0, // friction, growth
	models.AspectQuincunx:       -1.0, // adjustment needed
	models.AspectSemiSextile:    0.5,  // subtle
	models.AspectSemiSquare:     -1.0, // irritation
	models.AspectSesquiquadrate: -1.0, // tension
}

// synastryCategory determines the relationship area an aspect activates
func synastryCategory(planetA, planetB string) string {
	lights := map[string]bool{"SUN": true, "MOON": true}
	personal := map[string]bool{"MERCURY": true, "VENUS": true, "MARS": true}

	if lights[planetA] && lights[planetB] {
		return "Core Identity"
	}
	if (planetA == "VENUS" || planetB == "VENUS") && (planetA == "MARS" || planetB == "MARS") {
		return "Passion & Attraction"
	}
	if planetA == "VENUS" || planetB == "VENUS" {
		return "Love & Affection"
	}
	if planetA == "MARS" || planetB == "MARS" {
		return "Drive & Desire"
	}
	if lights[planetA] || lights[planetB] {
		if personal[planetA] || personal[planetB] {
			return "Communication & Connection"
		}
		return "Core Identity"
	}
	if planetA == "SATURN" || planetB == "SATURN" {
		return "Commitment & Structure"
	}
	if planetA == "JUPITER" || planetB == "JUPITER" {
		return "Growth & Expansion"
	}
	return "Transpersonal"
}

// CalcSynastryScore computes a compatibility score from cross-aspects
func CalcSynastryScore(crossAspects []models.CrossAspectInfo) *SynastryScore {
	var scored []ScoredAspect
	categoryMap := make(map[string]*CategoryScore)

	var harmony, tension float64

	for _, ca := range crossAspects {
		baseScore, ok := aspectBaseScore[ca.AspectType]
		if !ok {
			continue
		}

		// Weight by planet importance
		wA := planetWeight[ca.InnerBody]
		if wA == 0 {
			wA = 0.3
		}
		wB := planetWeight[ca.OuterBody]
		if wB == 0 {
			wB = 0.3
		}

		// Orb penalty: tighter orb = stronger
		orbFactor := 1.0 - (ca.Orb / 10.0)
		if orbFactor < 0.1 {
			orbFactor = 0.1
		}

		score := baseScore * wA * wB * orbFactor
		category := synastryCategory(ca.InnerBody, ca.OuterBody)

		sa := ScoredAspect{
			PersonA:    ca.InnerBody,
			PersonB:    ca.OuterBody,
			AspectType: ca.AspectType,
			Orb:        ca.Orb,
			Score:      math.Round(score*100) / 100,
			Category:   category,
		}
		scored = append(scored, sa)

		if score > 0 {
			harmony += score
		} else {
			tension += math.Abs(score)
		}

		if categoryMap[category] == nil {
			categoryMap[category] = &CategoryScore{Category: category}
		}
		categoryMap[category].Score += score
		categoryMap[category].Count++
	}

	// Sort by absolute score descending
	sort.Slice(scored, func(i, j int) bool {
		return math.Abs(scored[i].Score) > math.Abs(scored[j].Score)
	})

	// Top 10 aspects
	top := scored
	if len(top) > 10 {
		top = top[:10]
	}

	// Compatibility: ratio of harmony to total, scaled to 0-100
	total := harmony + tension
	compatibility := 50.0 // neutral baseline
	if total > 0 {
		compatibility = (harmony / total) * 100
	}

	var categories []CategoryScore
	for _, v := range categoryMap {
		v.Score = math.Round(v.Score*100) / 100
		categories = append(categories, *v)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Score > categories[j].Score
	})

	return &SynastryScore{
		TotalScore:     math.Round((harmony-tension)*100) / 100,
		Harmony:        math.Round(harmony*100) / 100,
		Tension:        math.Round(tension*100) / 100,
		Compatibility:  math.Round(compatibility*100) / 100,
		TopAspects:     top,
		CategoryScores: categories,
	}
}

// CalcSynastryFromCharts computes synastry from two charts' planet positions
func CalcSynastryFromCharts(chart1, chart2 []models.PlanetPosition, orbs models.OrbConfig) *SynastryScore {
	var innerBodies, outerBodies []aspect.Body
	for _, p := range chart1 {
		innerBodies = append(innerBodies, aspect.Body{
			ID: string(p.PlanetID), Longitude: p.Longitude, Speed: p.Speed,
		})
	}
	for _, p := range chart2 {
		outerBodies = append(outerBodies, aspect.Body{
			ID: string(p.PlanetID), Longitude: p.Longitude, Speed: p.Speed,
		})
	}

	crossAspects := aspect.FindCrossAspects(innerBodies, outerBodies, orbs)
	return CalcSynastryScore(crossAspects)
}
