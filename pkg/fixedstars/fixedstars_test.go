package fixedstars

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCatalogSize(t *testing.T) {
	if len(Catalog) < 40 {
		t.Errorf("Expected at least 40 stars in catalog, got %d", len(Catalog))
	}
}

func TestCatalogSignsPopulated(t *testing.T) {
	for _, s := range Catalog {
		if s.Sign == "" {
			t.Errorf("Star %s has empty sign", s.Name)
		}
		if s.SignDegree < 0 || s.SignDegree >= 30 {
			t.Errorf("Star %s sign degree out of range: %.2f", s.Name, s.SignDegree)
		}
	}
}

func TestGetStarByName(t *testing.T) {
	s := GetStarByName("Regulus")
	if s == nil {
		t.Fatal("Regulus not found in catalog")
	}
	if s.Magnitude < 1 || s.Magnitude > 2 {
		t.Errorf("Regulus magnitude = %.2f, expected ~1.35", s.Magnitude)
	}
}

func TestGetStarByName_NotFound(t *testing.T) {
	s := GetStarByName("NonExistent")
	if s != nil {
		t.Error("Expected nil for nonexistent star")
	}
}

func TestFindConjunctions(t *testing.T) {
	// Create a position near Regulus (~149.83 at J2000)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 150.0},
	}
	conj := FindConjunctions(positions, 1.5, j2000Epoch)

	found := false
	for _, c := range conj {
		if c.Star.Name == "Regulus" {
			found = true
			if c.Orb > 1.5 {
				t.Errorf("Regulus conjunction orb = %.2f, too large", c.Orb)
			}
		}
	}
	if !found {
		t.Error("Expected to find Sun conjunct Regulus")
	}
}

func TestFindConjunctions_NoMatch(t *testing.T) {
	// Position far from any star
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Longitude: 300.0},
	}
	conj := FindConjunctions(positions, 0.5, j2000Epoch)
	// With a very tight orb, may or may not find anything
	for _, c := range conj {
		if c.Orb > 0.5 {
			t.Errorf("Conjunction orb %.2f exceeds limit 0.5", c.Orb)
		}
	}
}

func TestPrecessLongitude(t *testing.T) {
	// After ~72 years, precession should add ~1 degree
	yearsForOneDeg := 3600 / 50.2882 // ~71.6 years
	jd := j2000Epoch + yearsForOneDeg*julianYear
	precessed := PrecessLongitude(0, jd)
	if precessed < 0.95 || precessed > 1.05 {
		t.Errorf("Precession after %.0f years = %.4f, expected ~1.0", yearsForOneDeg, precessed)
	}
}
