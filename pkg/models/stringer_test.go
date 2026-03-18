package models

import (
	"strings"
	"testing"
)

func TestPlanetPosition_String(t *testing.T) {
	p := PlanetPosition{
		PlanetID:     PlanetSun,
		Longitude:    280.5,
		Sign:         "Capricorn",
		SignDegree:   10.5,
		House:        10,
		IsRetrograde: false,
	}
	s := p.String()
	// Should contain Sun glyph ☉ and Capricorn glyph ♑
	if !strings.Contains(s, "\u2609") {
		t.Errorf("Expected Sun glyph in %q", s)
	}
	if !strings.Contains(s, "\u2651") {
		t.Errorf("Expected Capricorn glyph in %q", s)
	}
	// Should NOT contain retrograde glyph
	if strings.Contains(s, "\u211E") {
		t.Errorf("Should not contain retrograde glyph: %q", s)
	}
}

func TestPlanetPosition_String_Retrograde(t *testing.T) {
	p := PlanetPosition{
		PlanetID:     PlanetMercury,
		Longitude:    150.0,
		House:        5,
		IsRetrograde: true,
	}
	s := p.String()
	// Should contain retrograde glyph ℞
	if !strings.Contains(s, "\u211E") {
		t.Errorf("Expected retrograde glyph in %q", s)
	}
}

func TestAspectInfo_String(t *testing.T) {
	a := AspectInfo{
		PlanetA:    "SUN",
		PlanetB:    "MOON",
		AspectType: AspectTrine,
		Orb:        2.5,
		IsApplying: true,
	}
	s := a.String()
	// Should contain trine glyph △
	if !strings.Contains(s, "\u25B3") {
		t.Errorf("Expected trine glyph in %q", s)
	}
	if !strings.Contains(s, "applying") {
		t.Errorf("Expected 'applying' in %q", s)
	}
}

func TestCrossAspectInfo_String(t *testing.T) {
	ca := CrossAspectInfo{
		InnerBody:  "SUN",
		OuterBody:  "VENUS",
		AspectType: AspectConjunction,
		Orb:        1.0,
		IsApplying: false,
	}
	s := ca.String()
	if !strings.Contains(s, "separating") {
		t.Errorf("Expected 'separating' in %q", s)
	}
	// Should contain conjunction glyph ☌
	if !strings.Contains(s, "\u260C") {
		t.Errorf("Expected conjunction glyph in %q", s)
	}
}

func TestTransitEvent_String_Aspect(t *testing.T) {
	te := TransitEvent{
		EventType:       EventAspectExact,
		ChartType:       ChartTransit,
		Planet:          PlanetMars,
		AspectType:      AspectSquare,
		Target:          "SUN",
		TargetChartType: ChartNatal,
	}
	s := te.String()
	// Should contain Mars glyph ♂ and square glyph □
	if !strings.Contains(s, "\u2642") {
		t.Errorf("Expected Mars glyph in %q", s)
	}
	if !strings.Contains(s, "\u25A1") {
		t.Errorf("Expected square glyph in %q", s)
	}
}

func TestTransitEvent_String_SignIngress(t *testing.T) {
	te := TransitEvent{
		EventType: EventSignIngress,
		Planet:    PlanetVenus,
		ToSign:    "Aries",
	}
	s := te.String()
	if !strings.Contains(s, "Aries") {
		t.Errorf("Expected 'Aries' in %q", s)
	}
	// Should contain Aries glyph ♈
	if !strings.Contains(s, "\u2648") {
		t.Errorf("Expected Aries glyph in %q", s)
	}
}

func TestTransitEvent_String_Station(t *testing.T) {
	te := TransitEvent{
		EventType:   EventStation,
		Planet:      PlanetMercury,
		StationType: StationRetrograde,
	}
	s := te.String()
	// Should contain retrograde glyph ℞
	if !strings.Contains(s, "\u211E") {
		t.Errorf("Expected retrograde glyph in %q", s)
	}
}

func TestTransitEvent_String_StationDirect(t *testing.T) {
	te := TransitEvent{
		EventType:   EventStation,
		Planet:      PlanetMercury,
		StationType: StationDirect,
	}
	s := te.String()
	if !strings.Contains(s, "D") {
		t.Errorf("Expected 'D' for direct station in %q", s)
	}
}

func TestTransitEvent_String_VOC(t *testing.T) {
	te := TransitEvent{
		EventType:      EventVoidOfCourse,
		Planet:         PlanetMoon,
		LastAspectType: "Trine",
		NextSign:       "Leo",
	}
	s := te.String()
	if !strings.Contains(s, "VOC") {
		t.Errorf("Expected 'VOC' in %q", s)
	}
	// Should contain Moon glyph ☽ and Leo glyph ♌
	if !strings.Contains(s, "\u263D") {
		t.Errorf("Expected Moon glyph in %q", s)
	}
}

func TestTransitEvent_String_HouseIngress(t *testing.T) {
	te := TransitEvent{
		EventType: EventHouseIngress,
		Planet:    PlanetJupiter,
		ToHouse:   7,
	}
	s := te.String()
	if !strings.Contains(s, "H7") {
		t.Errorf("Expected 'H7' in %q", s)
	}
}

func TestTransitEvent_String_Default(t *testing.T) {
	te := TransitEvent{
		EventType: "UNKNOWN",
		Planet:    PlanetSun,
	}
	s := te.String()
	if s == "" {
		t.Error("Expected non-empty string")
	}
}

func TestAnglesInfo_String(t *testing.T) {
	a := AnglesInfo{ASC: 15.5, MC: 280.3}
	s := a.String()
	if !strings.Contains(s, "AC") || !strings.Contains(s, "MC") {
		t.Errorf("Expected 'AC' and 'MC' in %q", s)
	}
}
