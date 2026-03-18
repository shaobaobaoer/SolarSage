// Package dispositor computes dispositorship chains and final dispositors
// for a natal chart. The dispositor of a planet is the ruler of the sign
// the planet occupies. Chains reveal the flow of planetary energy.
package dispositor

import (
	"github.com/shaobaobaoer/solarsage-mcp/pkg/dignity"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// Chain represents a single dispositorship chain
type Chain struct {
	Links []models.PlanetID `json:"links"`
}

// DispositorResult holds the full dispositorship analysis
type DispositorResult struct {
	// Dispositors maps each planet to its dispositor (ruler of the sign it's in)
	Dispositors map[models.PlanetID]models.PlanetID `json:"dispositors"`
	// Chains are the dispositorship chains from each planet to the final dispositor
	Chains []Chain `json:"chains"`
	// FinalDispositor is the planet that disposes of itself (in its own sign), if any
	FinalDispositor *models.PlanetID `json:"final_dispositor"`
	// MutualDispositors are pairs that dispose of each other
	MutualDispositors [][2]models.PlanetID `json:"mutual_dispositors"`
	// UsesTraditionalRulers indicates which ruler system is used
	UsesTraditionalRulers bool `json:"uses_traditional_rulers"`
}

// CalcDispositors computes the dispositorship map and chains for a chart.
// If traditional is true, uses traditional rulers (Mars for Scorpio, etc.).
func CalcDispositors(positions []models.PlanetPosition, traditional bool) *DispositorResult {
	result := &DispositorResult{
		Dispositors:       make(map[models.PlanetID]models.PlanetID),
		UsesTraditionalRulers: traditional,
	}

	// Build planet -> sign map
	planetSign := make(map[models.PlanetID]string)
	for _, p := range positions {
		planetSign[p.PlanetID] = p.Sign
	}

	// Build dispositor map
	for _, p := range positions {
		var ruler models.PlanetID
		if traditional {
			ruler = dignity.SignTraditionalRuler(p.Sign)
		} else {
			ruler = dignity.SignRuler(p.Sign)
		}
		result.Dispositors[p.PlanetID] = ruler
	}

	// Find final dispositor: a planet that is its own dispositor (in its own sign)
	for planet, dispositor := range result.Dispositors {
		if planet == dispositor {
			p := planet
			result.FinalDispositor = &p
			break
		}
	}

	// Find mutual dispositors: A disposes B and B disposes A
	seen := make(map[[2]models.PlanetID]bool)
	for a, dispA := range result.Dispositors {
		if dispB, ok := result.Dispositors[dispA]; ok && dispB == a && a != dispA {
			pair := [2]models.PlanetID{a, dispA}
			if a > dispA {
				pair = [2]models.PlanetID{dispA, a}
			}
			if !seen[pair] {
				seen[pair] = true
				result.MutualDispositors = append(result.MutualDispositors, pair)
			}
		}
	}

	// Build chains: follow each planet's dispositor chain until we hit a loop or final
	visited := make(map[models.PlanetID]bool)
	for _, p := range positions {
		if visited[p.PlanetID] {
			continue
		}
		chain := followChain(p.PlanetID, result.Dispositors)
		if len(chain.Links) > 1 {
			result.Chains = append(result.Chains, chain)
		}
		for _, link := range chain.Links {
			visited[link] = true
		}
	}

	return result
}

// followChain follows the dispositor chain from a starting planet
func followChain(start models.PlanetID, dispositors map[models.PlanetID]models.PlanetID) Chain {
	chain := Chain{Links: []models.PlanetID{start}}
	seen := map[models.PlanetID]bool{start: true}
	current := start

	for {
		next, ok := dispositors[current]
		if !ok || seen[next] {
			break
		}
		seen[next] = true
		chain.Links = append(chain.Links, next)
		current = next
	}

	return chain
}
