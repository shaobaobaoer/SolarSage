package avastha

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcAvasthas_Garvita(t *testing.T) {
	// Sun in Leo (sign 4) = own sign → Garvita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 130, SignIndex: 4},
	}
	results := CalcAvasthas(planets, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !hasAvastha(results[0], Garvita) {
		t.Errorf("Sun in Leo should be Garvita, got %v", results[0].Avasthas)
	}
}

func TestCalcAvasthas_Exaltation(t *testing.T) {
	// Sun in Aries (sign 0) = exalted → Garvita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 10, SignIndex: 0},
	}
	results := CalcAvasthas(planets, nil)
	if !hasAvastha(results[0], Garvita) {
		t.Errorf("Sun in Aries (exalted) should be Garvita, got %v", results[0].Avasthas)
	}
}

func TestCalcAvasthas_Lajjita(t *testing.T) {
	// Sun in Libra (sign 6) = debilitated → Lajjita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 190, SignIndex: 6},
	}
	results := CalcAvasthas(planets, nil)
	if !hasAvastha(results[0], Lajjita) {
		t.Errorf("Sun in Libra (debilitated) should be Lajjita, got %v", results[0].Avasthas)
	}
}

func TestCalcAvasthas_Kshudhita(t *testing.T) {
	// Sun in Aquarius (sign 10, ruler Saturn → enemy) + aspected by malefic → Kshudhita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 310, SignIndex: 10},
		{PlanetID: models.PlanetSaturn, Longitude: 130, SignIndex: 4},
	}
	aspects := []AspectRef{
		{PlanetA: models.PlanetSaturn, PlanetB: models.PlanetSun, Type: "opposition"},
	}
	results := CalcAvasthas(planets, aspects)
	sunResult := findResult(results, models.PlanetSun)
	if sunResult == nil {
		t.Fatal("no Sun result")
	}
	if !hasAvastha(*sunResult, Kshudhita) {
		t.Errorf("Sun in enemy sign + malefic aspect should be Kshudhita, got %v", sunResult.Avasthas)
	}
}

func TestCalcAvasthas_Mudita(t *testing.T) {
	// Sun in Sagittarius (sign 8, ruler Jupiter → friend) + benefic aspect → Mudita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 250, SignIndex: 8},
		{PlanetID: models.PlanetJupiter, Longitude: 130, SignIndex: 4},
	}
	aspects := []AspectRef{
		{PlanetA: models.PlanetJupiter, PlanetB: models.PlanetSun, Type: "trine"},
	}
	results := CalcAvasthas(planets, aspects)
	sunResult := findResult(results, models.PlanetSun)
	if sunResult == nil {
		t.Fatal("no Sun result")
	}
	if !hasAvastha(*sunResult, Mudita) {
		t.Errorf("Sun in friend's sign + benefic aspect should be Mudita, got %v", sunResult.Avasthas)
	}
}

func TestCalcAvasthas_Trishita(t *testing.T) {
	// Moon in Scorpio (sign 7, water) + malefic aspect, no benefic → Trishita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetMoon, Longitude: 220, SignIndex: 7},
		{PlanetID: models.PlanetSaturn, Longitude: 310, SignIndex: 10},
	}
	aspects := []AspectRef{
		{PlanetA: models.PlanetSaturn, PlanetB: models.PlanetMoon, Type: "square"},
	}
	results := CalcAvasthas(planets, aspects)
	moonResult := findResult(results, models.PlanetMoon)
	if moonResult == nil {
		t.Fatal("no Moon result")
	}
	if !hasAvastha(*moonResult, Trishita) {
		t.Errorf("Moon in Scorpio + malefic aspect should be Trishita, got %v", moonResult.Avasthas)
	}
}

func TestCalcAvasthas_Kshobhita(t *testing.T) {
	// Mercury conjunct Sun + aspected by Mars → Kshobhita
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetMercury, Longitude: 130, SignIndex: 4},
		{PlanetID: models.PlanetSun, Longitude: 132, SignIndex: 4},
		{PlanetID: models.PlanetMars, Longitude: 220, SignIndex: 7},
	}
	aspects := []AspectRef{
		{PlanetA: models.PlanetMercury, PlanetB: models.PlanetSun, Type: "conjunction"},
		{PlanetA: models.PlanetMars, PlanetB: models.PlanetMercury, Type: "square"},
	}
	results := CalcAvasthas(planets, aspects)
	mercResult := findResult(results, models.PlanetMercury)
	if mercResult == nil {
		t.Fatal("no Mercury result")
	}
	if !hasAvastha(*mercResult, Kshobhita) {
		t.Errorf("Mercury conjunct Sun + malefic aspect should be Kshobhita, got %v", mercResult.Avasthas)
	}
}

func TestCalcAvasthas_EmptyInput(t *testing.T) {
	results := CalcAvasthas(nil, nil)
	if len(results) != 0 {
		t.Errorf("empty input should return empty results, got %d", len(results))
	}
}

func TestCalcAvasthas_SkipsOuterPlanets(t *testing.T) {
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetUranus, Longitude: 10, SignIndex: 0},
		{PlanetID: models.PlanetNeptune, Longitude: 40, SignIndex: 1},
	}
	results := CalcAvasthas(planets, nil)
	if len(results) != 0 {
		t.Errorf("outer planets should be skipped, got %d results", len(results))
	}
}

func TestCalcAvasthas_MultipleStates(t *testing.T) {
	// Sun in Libra (debilitated → Lajjita) + aspected by malefic Saturn → also Kshudhita (enemy sign + malefic)
	// Sun's enemy is Venus (ruler of Libra) → enemy sign confirmed
	planets := []AvasthaPlanetInfo{
		{PlanetID: models.PlanetSun, Longitude: 190, SignIndex: 6},
		{PlanetID: models.PlanetSaturn, Longitude: 280, SignIndex: 9},
	}
	aspects := []AspectRef{
		{PlanetA: models.PlanetSaturn, PlanetB: models.PlanetSun, Type: "square"},
	}
	results := CalcAvasthas(planets, aspects)
	sunResult := findResult(results, models.PlanetSun)
	if sunResult == nil {
		t.Fatal("no Sun result")
	}
	// Should have both Lajjita and Kshudhita
	if !hasAvastha(*sunResult, Lajjita) {
		t.Errorf("expected Lajjita, got %v", sunResult.Avasthas)
	}
	if !hasAvastha(*sunResult, Kshudhita) {
		t.Errorf("expected Kshudhita, got %v", sunResult.Avasthas)
	}
}

func hasAvastha(r AvasthResult, a Avastha) bool {
	for _, av := range r.Avasthas {
		if av == a {
			return true
		}
	}
	return false
}

func findResult(results []AvasthResult, id models.PlanetID) *AvasthResult {
	for _, r := range results {
		if r.PlanetID == id {
			return &r
		}
	}
	return nil
}
