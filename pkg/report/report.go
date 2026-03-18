// Package report generates comprehensive chart analysis reports that combine
// multiple calculation modules into a single structured result.
package report

import (
	"github.com/shaobaobaoer/solarsage-mcp/internal/aspect"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/antiscia"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/bounds"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/dignity"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/dispositor"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/fixedstars"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/lots"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/lunar"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/midpoint"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// ChartReport is a comprehensive natal chart analysis
type ChartReport struct {
	// Core chart data
	Chart *models.ChartInfo `json:"chart"`

	// Lunar phase at birth
	LunarPhase *lunar.PhaseInfo `json:"lunar_phase"`

	// Essential dignities for all planets
	Dignities []dignity.DignityInfo `json:"dignities"`

	// Mutual receptions
	MutualReceptions []dignity.MutualReceptionInfo `json:"mutual_receptions"`

	// Sect (day/night chart)
	IsDayChart bool              `json:"is_day_chart"`
	Sect       []dignity.SectInfo `json:"sect"`

	// Dispositorship
	Dispositors *dispositor.DispositorResult `json:"dispositors"`

	// Aspect patterns
	Patterns []aspect.AspectPattern `json:"aspect_patterns"`

	// Decans and terms
	Faces []bounds.FaceInfo `json:"faces"`

	// Arabic lots
	Lots []lots.LotResult `json:"lots"`

	// Antiscia pairs
	AntisciaPairs []antiscia.AntisciaPair `json:"antiscia_pairs"`

	// Fixed star conjunctions
	FixedStarConjunctions []fixedstars.StarConjunction `json:"fixed_star_conjunctions"`

	// Midpoint activations
	MidpointActivations []midpoint.MidpointAspect `json:"midpoint_activations"`

	// Element/modality distribution
	ElementBalance  map[string]int `json:"element_balance"`
	ModalityBalance map[string]int `json:"modality_balance"`

	// Hemisphere emphasis
	HemisphereBalance map[string]int `json:"hemisphere_balance"`
}

// DefaultPlanets for the report
var defaultPlanets = []models.PlanetID{
	models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
	models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
	models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
	models.PlanetPluto,
}

// GenerateNatalReport creates a comprehensive natal chart analysis.
func GenerateNatalReport(lat, lon, jdUT float64) (*ChartReport, error) {
	orbs := models.DefaultOrbConfig()
	hsys := models.HousePlacidus

	chartInfo, err := chart.CalcSingleChart(lat, lon, jdUT, defaultPlanets, orbs, hsys)
	if err != nil {
		return nil, err
	}

	report := &ChartReport{Chart: chartInfo}

	// Lunar phase
	report.LunarPhase, _ = lunar.CalcLunarPhase(jdUT)

	// Dignities
	report.Dignities = dignity.CalcChartDignities(chartInfo.Planets)
	report.MutualReceptions = dignity.FindMutualReceptions(chartInfo.Planets)

	// Sect
	report.IsDayChart = chart.IsDayChart(jdUT, chartInfo.Angles.ASC)
	for _, p := range chartInfo.Planets {
		report.Sect = append(report.Sect, dignity.CalcSect(p.PlanetID, report.IsDayChart))
	}

	// Dispositors
	report.Dispositors = dispositor.CalcDispositors(chartInfo.Planets, false)

	// Aspect patterns
	var bodies []aspect.Body
	for _, p := range chartInfo.Planets {
		bodies = append(bodies, aspect.Body{
			ID: string(p.PlanetID), Longitude: p.Longitude, Speed: p.Speed,
		})
	}
	report.Patterns = aspect.FindPatterns(chartInfo.Aspects, bodies, orbs)

	// Faces (decans + terms)
	report.Faces = bounds.CalcChartFaces(chartInfo.Planets)

	// Arabic lots
	report.Lots = lots.CalcStandardLots(chartInfo.Planets, chartInfo.Angles.ASC, report.IsDayChart)

	// Antiscia
	report.AntisciaPairs = antiscia.FindAntisciaPairs(chartInfo.Planets, 2.0)

	// Fixed stars
	report.FixedStarConjunctions = fixedstars.FindConjunctions(chartInfo.Planets, 1.5, jdUT)

	// Midpoints
	tree := midpoint.CalcMidpoints(chartInfo.Planets, 1.5)
	report.MidpointActivations = tree.Activations

	// Element balance
	report.ElementBalance = calcElementBalance(chartInfo.Planets)
	report.ModalityBalance = calcModalityBalance(chartInfo.Planets)
	report.HemisphereBalance = calcHemisphereBalance(chartInfo.Planets, chartInfo.Angles)

	return report, nil
}

func calcElementBalance(planets []models.PlanetPosition) map[string]int {
	elements := map[string]int{"Fire": 0, "Earth": 0, "Air": 0, "Water": 0}
	elementMap := map[string]string{
		"Aries": "Fire", "Leo": "Fire", "Sagittarius": "Fire",
		"Taurus": "Earth", "Virgo": "Earth", "Capricorn": "Earth",
		"Gemini": "Air", "Libra": "Air", "Aquarius": "Air",
		"Cancer": "Water", "Scorpio": "Water", "Pisces": "Water",
	}
	for _, p := range planets {
		if e, ok := elementMap[p.Sign]; ok {
			elements[e]++
		}
	}
	return elements
}

func calcModalityBalance(planets []models.PlanetPosition) map[string]int {
	modalities := map[string]int{"Cardinal": 0, "Fixed": 0, "Mutable": 0}
	modalityMap := map[string]string{
		"Aries": "Cardinal", "Cancer": "Cardinal", "Libra": "Cardinal", "Capricorn": "Cardinal",
		"Taurus": "Fixed", "Leo": "Fixed", "Scorpio": "Fixed", "Aquarius": "Fixed",
		"Gemini": "Mutable", "Virgo": "Mutable", "Sagittarius": "Mutable", "Pisces": "Mutable",
	}
	for _, p := range planets {
		if m, ok := modalityMap[p.Sign]; ok {
			modalities[m]++
		}
	}
	return modalities
}

func calcHemisphereBalance(planets []models.PlanetPosition, angles models.AnglesInfo) map[string]int {
	balance := map[string]int{"Eastern": 0, "Western": 0, "Northern": 0, "Southern": 0}
	for _, p := range planets {
		// East/West: relative to MC-IC axis
		if isInHemisphere(p.Longitude, angles.ASC, 180) {
			balance["Eastern"]++
		} else {
			balance["Western"]++
		}
		// North/South: relative to ASC-DSC axis
		if isInHemisphere(p.Longitude, angles.MC, 180) {
			balance["Southern"]++
		} else {
			balance["Northern"]++
		}
	}
	return balance
}

func isInHemisphere(lon, start, span float64) bool {
	diff := lon - start
	if diff < 0 {
		diff += 360
	}
	return diff < span
}
