package mcp

import (
	"encoding/json"
	"testing"

	"github.com/anthropic/swisseph-mcp/pkg/models"
)

func TestHandleCalcPlanetPosition_Direct(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"planet":"SUN","jd_ut":2451545.0}`)
	result, err := s.handleCalcPlanetPosition(args)
	if err != nil {
		t.Fatalf("handleCalcPlanetPosition error: %v", err)
	}
	m := result.(map[string]interface{})
	if m["planet"] != models.PlanetID("SUN") {
		t.Errorf("expected planet SUN, got %v", m["planet"])
	}
	lon := m["longitude"].(float64)
	if lon < 279 || lon > 281 {
		t.Errorf("Sun longitude at J2000.0 expected ~280, got %f", lon)
	}
}

func TestHandleCalcPlanetPosition_InvalidPlanet(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"planet":"INVALID","jd_ut":2451545.0}`)
	_, err := s.handleCalcPlanetPosition(args)
	if err == nil {
		t.Error("expected error for invalid planet")
	}
}

func TestHandleJDToDatetime_Direct(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"jd":2451545.0,"timezone":"UTC"}`)
	result, err := s.handleJDToDatetime(args)
	if err != nil {
		t.Fatalf("handleJDToDatetime error: %v", err)
	}
	m := result.(map[string]string)
	if m["datetime"] == "" {
		t.Error("expected non-empty datetime")
	}
}

func TestHandleJDToDatetime_NoTimezone(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"jd":2451545.0}`)
	result, err := s.handleJDToDatetime(args)
	if err != nil {
		t.Fatalf("handleJDToDatetime error: %v", err)
	}
	m := result.(map[string]string)
	if m["datetime"] == "" {
		t.Error("expected non-empty datetime with default timezone")
	}
}

func TestHandleJDToDatetime_InvalidTimezone(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"jd":2451545.0,"timezone":"Invalid/Timezone"}`)
	_, err := s.handleJDToDatetime(args)
	if err == nil {
		t.Error("expected error for invalid timezone")
	}
}

func TestHandleCalcSolarArc_Direct(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")

	args := json.RawMessage(`{"natal_jd_ut":2448000.5,"transit_jd_ut":2451545.0}`)
	result, err := s.handleCalcSolarArc(args)
	if err != nil {
		t.Fatalf("handleCalcSolarArc error: %v", err)
	}
	m := result.(map[string]interface{})
	if m["age"] == nil {
		t.Error("expected age field")
	}
	if m["solar_arc_offset"] == nil {
		t.Error("expected solar_arc_offset field")
	}
}

func TestOrbOrDefault_Nil(t *testing.T) {
	fallback := models.DefaultOrbConfig()
	result := orbOrDefault(nil, fallback)
	if result.Conjunction != fallback.Conjunction {
		t.Error("nil custom should return fallback")
	}
}

func TestOrbOrDefault_Custom(t *testing.T) {
	custom := models.OrbConfig{Conjunction: 2.0}
	fallback := models.DefaultOrbConfig()
	result := orbOrDefault(&custom, fallback)
	if result.Conjunction != 2.0 {
		t.Errorf("custom should return 2.0, got %f", result.Conjunction)
	}
}

func TestHandleGeocode_InvalidJSON(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")
	_, err := s.handleGeocode(json.RawMessage(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestHandleCalcPlanetPosition_InvalidJSON(t *testing.T) {
	s := NewServer("../../third_party/swisseph/ephe")
	_, err := s.handleCalcPlanetPosition(json.RawMessage(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
