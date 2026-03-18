package julian

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropic/swisseph-mcp/pkg/models"
	"github.com/anthropic/swisseph-mcp/pkg/sweph"
)

func TestMain(m *testing.M) {
	ephePath := filepath.Join("..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	code := m.Run()
	sweph.Close()
	os.Exit(code)
}

func TestDateTimeToJD(t *testing.T) {
	// 1990-06-15T08:30:00+08:00 = 1990-06-15 00:30 UTC
	r, err := DateTimeToJD("1990-06-15T08:30:00+08:00", models.CalendarGregorian)
	if err != nil {
		t.Fatalf("DateTimeToJD: %v", err)
	}

	// Expected JD_UT for 1990-06-15 00:30 UTC ≈ 2448057.5208
	if math.Abs(r.JDUT-2448057.5208) > 0.001 {
		t.Errorf("JDUT = %f, expected ~2448057.5208", r.JDUT)
	}

	// TT should be slightly larger than UT
	if r.JDTT <= r.JDUT {
		t.Errorf("JDTT (%f) should be > JDUT (%f)", r.JDTT, r.JDUT)
	}
}

func TestDateTimeToJD_InvalidFormat(t *testing.T) {
	_, err := DateTimeToJD("not-a-date", models.CalendarGregorian)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestJDToDateTime(t *testing.T) {
	// Round-trip: convert to JD and back
	r, _ := DateTimeToJD("2024-01-01T12:00:00+00:00", models.CalendarGregorian)
	dt, err := JDToDateTime(r.JDUT, "UTC")
	if err != nil {
		t.Fatalf("JDToDateTime: %v", err)
	}
	if !strings.HasPrefix(dt, "2024-01-01T12:00:0") {
		t.Errorf("Round-trip: got %s, expected 2024-01-01T12:00:0x", dt)
	}
}

func TestJDToDateTime_InvalidTimezone(t *testing.T) {
	_, err := JDToDateTime(2451545.0, "Invalid/TZ")
	if err == nil {
		t.Error("Expected error for invalid timezone")
	}
}

func TestJDToDateTime_UTC(t *testing.T) {
	dt, err := JDToDateTime(2451545.0, "UTC") // J2000.0 = 2000-01-01 12:00 UTC
	if err != nil {
		t.Fatalf("JDToDateTime: %v", err)
	}
	if !strings.HasPrefix(dt, "2000-01-01T12:00:0") {
		t.Errorf("J2000.0 in UTC = %s, expected 2000-01-01T12:00:0x", dt)
	}
}
