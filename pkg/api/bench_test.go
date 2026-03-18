package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkNatalChart(b *testing.B) {
	srv := newTestServer("")
	body, _ := json.Marshal(map[string]interface{}{
		"latitude":  51.5074,
		"longitude": -0.1278,
		"jd_ut":     2451545.0,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chart/natal", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
	}
}

func BenchmarkPlanetPosition(b *testing.B) {
	srv := newTestServer("")
	body, _ := json.Marshal(map[string]interface{}{
		"planet": "SUN",
		"jd_ut":  2451545.0,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/planet/position", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
	}
}

func BenchmarkNatalReport(b *testing.B) {
	srv := newTestServer("")
	body, _ := json.Marshal(map[string]interface{}{
		"latitude":  51.5074,
		"longitude": -0.1278,
		"jd_ut":     2451545.0,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/report/natal", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
	}
}
