package aspect

import (
	"math"
	"testing"

	"github.com/anthropic/swisseph-mcp/pkg/models"
)

func TestAngleDiff(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{0, 0, 0},
		{0, 180, 180},
		{10, 350, 20},
		{350, 10, 20},
		{90, 270, 180},
		{0, 90, 90},
		{180, 0, 180},
		{1, 359, 2},
	}
	for _, tt := range tests {
		got := AngleDiff(tt.a, tt.b)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("AngleDiff(%f, %f) = %f, want %f", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSignedAngleDiff(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{10, 350, 20},
		{350, 10, -20},
		{0, 180, -180},
		{90, 0, 90},
		{0, 90, -90},
	}
	for _, tt := range tests {
		got := SignedAngleDiff(tt.a, tt.b)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("SignedAngleDiff(%f, %f) = %f, want %f", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFindAspects_Conjunction(t *testing.T) {
	orbs := models.OrbConfig{Conjunction: 8}
	a := []Body{{ID: "A", Longitude: 100, Speed: 1}}
	b := []Body{{ID: "B", Longitude: 105, Speed: 0.5}}

	aspects := FindAspects(a, b, orbs, false)
	if len(aspects) != 1 {
		t.Fatalf("Expected 1 aspect, got %d", len(aspects))
	}
	if aspects[0].AspectType != models.AspectConjunction {
		t.Errorf("Expected conjunction, got %s", aspects[0].AspectType)
	}
	if math.Abs(aspects[0].Orb-5) > 0.01 {
		t.Errorf("Orb = %f, want 5", aspects[0].Orb)
	}
}

func TestFindAspects_Opposition(t *testing.T) {
	orbs := models.OrbConfig{Opposition: 8}
	a := []Body{{ID: "A", Longitude: 0, Speed: 1}}
	b := []Body{{ID: "B", Longitude: 183, Speed: 0.5}}

	aspects := FindAspects(a, b, orbs, false)
	if len(aspects) != 1 {
		t.Fatalf("Expected 1 aspect, got %d", len(aspects))
	}
	if aspects[0].AspectType != models.AspectOpposition {
		t.Errorf("Expected opposition, got %s", aspects[0].AspectType)
	}
}

func TestFindAspects_SameSet_NoDuplicates(t *testing.T) {
	orbs := models.OrbConfig{Conjunction: 10}
	bodies := []Body{
		{ID: "A", Longitude: 100, Speed: 1},
		{ID: "B", Longitude: 105, Speed: 0.5},
	}

	aspects := FindAspects(bodies, bodies, orbs, true)
	// Should only have A-B, not B-A
	if len(aspects) != 1 {
		t.Errorf("Expected 1 aspect (no duplicates), got %d", len(aspects))
	}
}

func TestFindAspects_OutOfOrb(t *testing.T) {
	orbs := models.OrbConfig{Conjunction: 2}
	a := []Body{{ID: "A", Longitude: 100, Speed: 1}}
	b := []Body{{ID: "B", Longitude: 110, Speed: 0.5}}

	aspects := FindAspects(a, b, orbs, false)
	if len(aspects) != 0 {
		t.Errorf("Expected 0 aspects (out of orb), got %d", len(aspects))
	}
}

func TestFindCrossAspects(t *testing.T) {
	orbs := models.OrbConfig{Trine: 7}
	inner := []Body{{ID: "Sun", Longitude: 0, Speed: 1}}
	outer := []Body{{ID: "Moon", Longitude: 122, Speed: 13}}

	ca := FindCrossAspects(inner, outer, orbs)
	if len(ca) != 1 {
		t.Fatalf("Expected 1 cross aspect, got %d", len(ca))
	}
	if ca[0].AspectType != models.AspectTrine {
		t.Errorf("Expected trine, got %s", ca[0].AspectType)
	}
	if ca[0].InnerBody != "Sun" || ca[0].OuterBody != "Moon" {
		t.Errorf("Wrong bodies: %s-%s", ca[0].InnerBody, ca[0].OuterBody)
	}
}

func TestComputeApplying(t *testing.T) {
	// A approaching B for conjunction
	a := Body{ID: "A", Longitude: 95, Speed: 1.5}
	b := Body{ID: "B", Longitude: 100, Speed: 0.5}
	// A is approaching B (faster), aspect is applying
	applying := computeApplying(a, b, 0)
	if !applying {
		t.Error("A approaching B for conjunction should be applying")
	}

	// A separating from B
	a2 := Body{ID: "A", Longitude: 105, Speed: 1.5}
	b2 := Body{ID: "B", Longitude: 100, Speed: 0.5}
	separating := computeApplying(a2, b2, 0)
	if separating {
		t.Error("A separating from B for conjunction should not be applying")
	}
}
