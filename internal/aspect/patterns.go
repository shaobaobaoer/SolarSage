package aspect

import (
	"math"
	"sort"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// PatternType identifies an aspect pattern
type PatternType string

const (
	PatternGrandTrine  PatternType = "GRAND_TRINE"
	PatternTSquare     PatternType = "T_SQUARE"
	PatternGrandCross  PatternType = "GRAND_CROSS"
	PatternYod         PatternType = "YOD"
	PatternKite        PatternType = "KITE"
	PatternMysticRect  PatternType = "MYSTIC_RECTANGLE"
	PatternStelllium   PatternType = "STELLIUM"
)

// AspectPattern represents a detected aspect pattern
type AspectPattern struct {
	Type    PatternType `json:"type"`
	Bodies  []string    `json:"bodies"`
	Element string      `json:"element,omitempty"` // for Grand Trine
	Quality string      `json:"quality,omitempty"` // for T-Square/Grand Cross (cardinal/fixed/mutable)
	Apex    string      `json:"apex,omitempty"`    // for T-Square or Yod
}

// FindPatterns detects aspect patterns from a list of aspects and bodies
func FindPatterns(aspects []models.AspectInfo, bodies []Body, orbs models.OrbConfig) []AspectPattern {
	var patterns []AspectPattern

	// Build adjacency for each aspect type
	trines := buildAdjacency(aspects, models.AspectTrine)
	squares := buildAdjacency(aspects, models.AspectSquare)
	opps := buildAdjacency(aspects, models.AspectOpposition)
	sextiles := buildAdjacency(aspects, models.AspectSextile)
	quincunxes := buildAdjacency(aspects, models.AspectQuincunx)

	// Grand Trine: 3 bodies mutually in trine
	patterns = append(patterns, findGrandTrines(trines, bodies)...)

	// T-Square: 2 bodies in opposition, both square to a 3rd (apex)
	patterns = append(patterns, findTSquares(squares, opps, bodies)...)

	// Grand Cross: 4 bodies forming 2 oppositions and 4 squares
	patterns = append(patterns, findGrandCrosses(squares, opps, bodies)...)

	// Yod (Finger of God): 2 bodies in sextile, both quincunx to a 3rd (apex)
	patterns = append(patterns, findYods(sextiles, quincunxes, bodies)...)

	// Kite: Grand Trine + one body in opposition to a trine point and sextile to the other two
	patterns = append(patterns, findKites(trines, opps, sextiles, bodies)...)

	// Mystic Rectangle: 2 oppositions connected by 2 sextiles and 2 trines
	patterns = append(patterns, findMysticRectangles(opps, sextiles, trines, bodies)...)

	// Stellium: 3+ planets within 8 degrees (or same sign)
	patterns = append(patterns, findStelliums(bodies, orbs)...)

	return patterns
}

type adjacency map[string]map[string]bool

func buildAdjacency(aspects []models.AspectInfo, aspType models.AspectType) adjacency {
	adj := make(adjacency)
	for _, a := range aspects {
		if a.AspectType != aspType {
			continue
		}
		if adj[a.PlanetA] == nil {
			adj[a.PlanetA] = make(map[string]bool)
		}
		if adj[a.PlanetB] == nil {
			adj[a.PlanetB] = make(map[string]bool)
		}
		adj[a.PlanetA][a.PlanetB] = true
		adj[a.PlanetB][a.PlanetA] = true
	}
	return adj
}

func hasEdge(adj adjacency, a, b string) bool {
	if m, ok := adj[a]; ok {
		return m[b]
	}
	return false
}

func allNodes(adj adjacency) []string {
	seen := make(map[string]bool)
	var nodes []string
	for k := range adj {
		if !seen[k] {
			seen[k] = true
			nodes = append(nodes, k)
		}
	}
	sort.Strings(nodes)
	return nodes
}

func findGrandTrines(trines adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	nodes := allNodes(trines)

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if !hasEdge(trines, nodes[i], nodes[j]) {
				continue
			}
			for k := j + 1; k < len(nodes); k++ {
				if hasEdge(trines, nodes[i], nodes[k]) && hasEdge(trines, nodes[j], nodes[k]) {
					element := detectTrineElement(nodes[i], nodes[j], nodes[k], bodies)
					patterns = append(patterns, AspectPattern{
						Type:    PatternGrandTrine,
						Bodies:  []string{nodes[i], nodes[j], nodes[k]},
						Element: element,
					})
				}
			}
		}
	}
	return patterns
}

func findTSquares(squares, opps adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	oppNodes := allNodes(opps)

	// Find all oppositions, then check if any body squares both
	for i := 0; i < len(oppNodes); i++ {
		for j := i + 1; j < len(oppNodes); j++ {
			if !hasEdge(opps, oppNodes[i], oppNodes[j]) {
				continue
			}
			// Find apex: body that squares both opposition bodies
			sqNodes := allNodes(squares)
			for _, apex := range sqNodes {
				if apex == oppNodes[i] || apex == oppNodes[j] {
					continue
				}
				if hasEdge(squares, apex, oppNodes[i]) && hasEdge(squares, apex, oppNodes[j]) {
					quality := detectSquareQuality(apex, bodies)
					patterns = append(patterns, AspectPattern{
						Type:    PatternTSquare,
						Bodies:  []string{oppNodes[i], oppNodes[j], apex},
						Apex:    apex,
						Quality: quality,
					})
				}
			}
		}
	}
	return patterns
}

func findGrandCrosses(squares, opps adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	seen := make(map[string]bool)

	// Grand Cross: two opposition pairs where all 4 bodies are mutually connected by squares
	oppPairs := findAllEdges(opps)
	for i := 0; i < len(oppPairs); i++ {
		for j := i + 1; j < len(oppPairs); j++ {
			a, b := oppPairs[i][0], oppPairs[i][1]
			c, d := oppPairs[j][0], oppPairs[j][1]
			// All 4 must be distinct
			if a == c || a == d || b == c || b == d {
				continue
			}
			// Check all 4 square connections
			if hasEdge(squares, a, c) && hasEdge(squares, a, d) &&
				hasEdge(squares, b, c) && hasEdge(squares, b, d) {
				key := sortedKey([]string{a, b, c, d})
				if seen[key] {
					continue
				}
				seen[key] = true
				quality := detectSquareQuality(a, bodies)
				patterns = append(patterns, AspectPattern{
					Type:    PatternGrandCross,
					Bodies:  []string{a, b, c, d},
					Quality: quality,
				})
			}
		}
	}
	return patterns
}

func findAllEdges(adj adjacency) [][2]string {
	seen := make(map[string]bool)
	var edges [][2]string
	for a, neighbors := range adj {
		for b := range neighbors {
			key := a + "|" + b
			rev := b + "|" + a
			if !seen[key] && !seen[rev] {
				seen[key] = true
				edges = append(edges, [2]string{a, b})
			}
		}
	}
	return edges
}

func sortedKey(strs []string) string {
	sort.Strings(strs)
	key := ""
	for _, s := range strs {
		key += s + "|"
	}
	return key
}

func findYods(sextiles, quincunxes adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	sxNodes := allNodes(sextiles)

	for i := 0; i < len(sxNodes); i++ {
		for j := i + 1; j < len(sxNodes); j++ {
			if !hasEdge(sextiles, sxNodes[i], sxNodes[j]) {
				continue
			}
			// Find apex: body quincunx to both sextile bodies
			qxNodes := allNodes(quincunxes)
			for _, apex := range qxNodes {
				if apex == sxNodes[i] || apex == sxNodes[j] {
					continue
				}
				if hasEdge(quincunxes, apex, sxNodes[i]) && hasEdge(quincunxes, apex, sxNodes[j]) {
					patterns = append(patterns, AspectPattern{
						Type:   PatternYod,
						Bodies: []string{sxNodes[i], sxNodes[j], apex},
						Apex:   apex,
					})
				}
			}
		}
	}
	return patterns
}

func findKites(trines, opps, sextiles adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	nodes := allNodes(trines)

	// First find grand trines
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if !hasEdge(trines, nodes[i], nodes[j]) {
				continue
			}
			for k := j + 1; k < len(nodes); k++ {
				if !hasEdge(trines, nodes[i], nodes[k]) || !hasEdge(trines, nodes[j], nodes[k]) {
					continue
				}
				// Grand trine found: nodes[i], j, k
				// Look for a 4th body that opposes one and sextiles the other two
				oppNodes := allNodes(opps)
				for _, d := range oppNodes {
					if d == nodes[i] || d == nodes[j] || d == nodes[k] {
						continue
					}
					trio := []string{nodes[i], nodes[j], nodes[k]}
					for _, t := range trio {
						others := otherTwo(trio, t)
						if hasEdge(opps, d, t) && hasEdge(sextiles, d, others[0]) && hasEdge(sextiles, d, others[1]) {
							patterns = append(patterns, AspectPattern{
								Type:   PatternKite,
								Bodies: []string{nodes[i], nodes[j], nodes[k], d},
								Apex:   d,
							})
						}
					}
				}
			}
		}
	}
	return patterns
}

func findMysticRectangles(opps, sextiles, trines adjacency, bodies []Body) []AspectPattern {
	var patterns []AspectPattern
	oppNodes := allNodes(opps)

	// Find two pairs of oppositions where the cross-connections are sextiles and trines
	for i := 0; i < len(oppNodes); i++ {
		for j := i + 1; j < len(oppNodes); j++ {
			if !hasEdge(opps, oppNodes[i], oppNodes[j]) {
				continue
			}
			for k := j + 1; k < len(oppNodes); k++ {
				for l := k + 1; l < len(oppNodes); l++ {
					if !hasEdge(opps, oppNodes[k], oppNodes[l]) {
						continue
					}
					a, b := oppNodes[i], oppNodes[j]
					c, d := oppNodes[k], oppNodes[l]
					// Check if a-c, b-d are sextiles and a-d, b-c are trines (or vice versa)
					if (hasEdge(sextiles, a, c) && hasEdge(sextiles, b, d) && hasEdge(trines, a, d) && hasEdge(trines, b, c)) ||
						(hasEdge(trines, a, c) && hasEdge(trines, b, d) && hasEdge(sextiles, a, d) && hasEdge(sextiles, b, c)) {
						patterns = append(patterns, AspectPattern{
							Type:   PatternMysticRect,
							Bodies: []string{a, b, c, d},
						})
					}
				}
			}
		}
	}
	return patterns
}

func findStelliums(bodies []Body, orbs models.OrbConfig) []AspectPattern {
	var patterns []AspectPattern
	if len(bodies) < 3 {
		return nil
	}

	// Group bodies by sign
	signBodies := make(map[int][]string)
	for _, b := range bodies {
		sign := int(b.Longitude / 30.0)
		signBodies[sign] = append(signBodies[sign], b.ID)
	}

	for _, group := range signBodies {
		if len(group) >= 3 {
			patterns = append(patterns, AspectPattern{
				Type:   PatternStelllium,
				Bodies: group,
			})
		}
	}

	return patterns
}

func otherTwo(trio []string, exclude string) []string {
	var result []string
	for _, s := range trio {
		if s != exclude {
			result = append(result, s)
		}
	}
	return result
}

// detectTrineElement returns the element of a grand trine based on body positions
func detectTrineElement(a, b, c string, bodies []Body) string {
	elements := []string{"Fire", "Earth", "Air", "Water"}
	for _, body := range bodies {
		if body.ID == a {
			signIdx := int(math.Mod(body.Longitude, 360) / 30.0)
			return elements[signIdx%4]
		}
	}
	return ""
}

// detectSquareQuality returns the quality (cardinal/fixed/mutable) based on body position
func detectSquareQuality(id string, bodies []Body) string {
	qualities := []string{"Cardinal", "Fixed", "Mutable"}
	for _, body := range bodies {
		if body.ID == id {
			signIdx := int(math.Mod(body.Longitude, 360) / 30.0)
			return qualities[signIdx%3]
		}
	}
	return ""
}
