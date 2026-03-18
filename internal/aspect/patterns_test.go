package aspect

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestFindPatterns_GrandTrine(t *testing.T) {
	// Three bodies at 0°, 120°, 240° (perfect grand trine)
	bodies := []Body{
		{ID: "SUN", Longitude: 5},   // Aries (Fire)
		{ID: "MARS", Longitude: 125}, // Leo (Fire)
		{ID: "JUPITER", Longitude: 245}, // Sagittarius (Fire)
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternGrandTrine {
			found = true
			if len(p.Bodies) != 3 {
				t.Errorf("Grand Trine should have 3 bodies, got %d", len(p.Bodies))
			}
		}
	}
	if !found {
		t.Error("Expected to find a Grand Trine pattern")
	}
}

func TestFindPatterns_TSquare(t *testing.T) {
	// Opposition at 0° and 180°, with square at 90°
	bodies := []Body{
		{ID: "SUN", Longitude: 0},
		{ID: "MOON", Longitude: 180},
		{ID: "MARS", Longitude: 90}, // apex
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternTSquare {
			found = true
			if p.Apex != "MARS" {
				t.Errorf("T-Square apex = %s, want MARS", p.Apex)
			}
		}
	}
	if !found {
		t.Error("Expected to find a T-Square pattern")
	}
}

func TestFindPatterns_GrandCross(t *testing.T) {
	// Four bodies at 0°, 90°, 180°, 270°
	bodies := []Body{
		{ID: "SUN", Longitude: 0},
		{ID: "MOON", Longitude: 90},
		{ID: "MARS", Longitude: 180},
		{ID: "JUPITER", Longitude: 270},
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternGrandCross {
			found = true
			if len(p.Bodies) != 4 {
				t.Errorf("Grand Cross should have 4 bodies, got %d", len(p.Bodies))
			}
		}
	}
	if !found {
		t.Error("Expected to find a Grand Cross pattern")
	}
}

func TestFindPatterns_Yod(t *testing.T) {
	// Two bodies in sextile (0° and 60°), both quincunx to apex at 210°
	// SUN(0)->MARS(210) = 150° quincunx, MOON(60)->MARS(210) = 150° quincunx
	bodies := []Body{
		{ID: "SUN", Longitude: 0},
		{ID: "MOON", Longitude: 60},
		{ID: "MARS", Longitude: 210}, // apex: quincunx to both
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternYod {
			found = true
			if p.Apex != "MARS" {
				t.Errorf("Yod apex = %s, want MARS", p.Apex)
			}
		}
	}
	if !found {
		t.Error("Expected to find a Yod pattern")
	}
}

func TestFindPatterns_Stellium(t *testing.T) {
	// 4 planets in the same sign (Aries: 0°-30°)
	bodies := []Body{
		{ID: "SUN", Longitude: 5},
		{ID: "MERCURY", Longitude: 10},
		{ID: "VENUS", Longitude: 15},
		{ID: "MARS", Longitude: 25},
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternStelllium {
			found = true
			if len(p.Bodies) != 4 {
				t.Errorf("Stellium should have 4 bodies, got %d", len(p.Bodies))
			}
		}
	}
	if !found {
		t.Error("Expected to find a Stellium pattern")
	}
}

func TestFindPatterns_NoPatterns(t *testing.T) {
	// Two bodies with no significant pattern
	bodies := []Body{
		{ID: "SUN", Longitude: 0},
		{ID: "MOON", Longitude: 75}, // not a standard aspect
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	for _, p := range patterns {
		if p.Type == PatternGrandTrine || p.Type == PatternTSquare || p.Type == PatternGrandCross || p.Type == PatternYod {
			t.Errorf("Unexpected pattern found: %s", p.Type)
		}
	}
}

func TestFindPatterns_Kite(t *testing.T) {
	// Grand Trine at 0°, 120°, 240° + a 4th body opposite one of them
	bodies := []Body{
		{ID: "SUN", Longitude: 0},
		{ID: "MARS", Longitude: 120},
		{ID: "JUPITER", Longitude: 240},
		{ID: "SATURN", Longitude: 180}, // opposite SUN, sextile MARS and JUPITER
	}
	orbs := models.DefaultOrbConfig()
	aspects := FindAspects(bodies, bodies, orbs, true)
	patterns := FindPatterns(aspects, bodies, orbs)

	found := false
	for _, p := range patterns {
		if p.Type == PatternKite {
			found = true
		}
	}
	if !found {
		t.Error("Expected to find a Kite pattern")
	}
}
