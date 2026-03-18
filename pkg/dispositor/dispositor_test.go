package dispositor

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcDispositors_FinalDispositor(t *testing.T) {
	// Sun in Leo = Sun disposes itself (final dispositor)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Leo"},
		{PlanetID: models.PlanetMoon, Sign: "Cancer"},
		{PlanetID: models.PlanetMercury, Sign: "Virgo"},
	}
	result := CalcDispositors(positions, false)

	if result.FinalDispositor == nil {
		t.Fatal("Expected final dispositor")
	}
	// Sun in Leo rules itself, Moon in Cancer rules itself too
	// At least one should be found
	fd := *result.FinalDispositor
	if fd != models.PlanetSun && fd != models.PlanetMoon && fd != models.PlanetMercury {
		t.Errorf("Unexpected final dispositor: %s", fd)
	}
}

func TestCalcDispositors_MutualDispositors(t *testing.T) {
	// Sun in Cancer (ruled by Moon), Moon in Leo (ruled by Sun) = mutual
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Cancer"},
		{PlanetID: models.PlanetMoon, Sign: "Leo"},
	}
	result := CalcDispositors(positions, false)

	if len(result.MutualDispositors) != 1 {
		t.Fatalf("Expected 1 mutual dispositor pair, got %d", len(result.MutualDispositors))
	}
}

func TestCalcDispositors_Traditional(t *testing.T) {
	// Scorpio: modern ruler = Pluto, traditional ruler = Mars
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Scorpio"},
	}
	modern := CalcDispositors(positions, false)
	trad := CalcDispositors(positions, true)

	if modern.Dispositors[models.PlanetSun] != models.PlanetPluto {
		t.Errorf("Modern: Sun in Scorpio dispositor = %s, want PLUTO", modern.Dispositors[models.PlanetSun])
	}
	if trad.Dispositors[models.PlanetSun] != models.PlanetMars {
		t.Errorf("Traditional: Sun in Scorpio dispositor = %s, want MARS", trad.Dispositors[models.PlanetSun])
	}
}

func TestCalcDispositors_Chains(t *testing.T) {
	// Sun in Aries (Mars), Mars in Capricorn (Saturn), Saturn in Aquarius (Uranus)
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Aries"},
		{PlanetID: models.PlanetMars, Sign: "Capricorn"},
		{PlanetID: models.PlanetSaturn, Sign: "Aquarius"},
	}
	result := CalcDispositors(positions, false)

	if len(result.Chains) == 0 {
		t.Error("Expected at least one chain")
	}
	// Sun -> Mars -> Saturn chain
	found := false
	for _, c := range result.Chains {
		if len(c.Links) >= 2 {
			found = true
		}
	}
	if !found {
		t.Error("Expected a chain with at least 2 links")
	}
}

func TestCalcDispositors_NoFinalDispositor(t *testing.T) {
	// No planet in its own sign
	positions := []models.PlanetPosition{
		{PlanetID: models.PlanetSun, Sign: "Gemini"},
		{PlanetID: models.PlanetMoon, Sign: "Aries"},
	}
	result := CalcDispositors(positions, false)

	if result.FinalDispositor != nil {
		t.Errorf("Expected no final dispositor, got %s", *result.FinalDispositor)
	}
}
