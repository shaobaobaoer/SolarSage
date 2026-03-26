package progressions

import (
	"math"
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestEclipticToRA(t *testing.T) {
	// At lon=0, RA should be 0
	ra := eclipticToRA(0, 23.44)
	if math.Abs(ra) > 0.01 {
		t.Errorf("eclipticToRA(0, 23.44) = %f, want ~0", ra)
	}

	// At lon=90 (summer solstice), RA should be ~90
	ra90 := eclipticToRA(90, 23.44)
	if math.Abs(ra90-90) > 0.5 {
		t.Errorf("eclipticToRA(90, 23.44) = %f, want ~90", ra90)
	}
}

func TestMcFromRAMC(t *testing.T) {
	// At RAMC=0, MC should be 0
	mc := mcFromRAMC(0, 23.44)
	if math.Abs(mc) > 0.01 {
		t.Errorf("mcFromRAMC(0, 23.44) = %f, want ~0", mc)
	}

	// At RAMC=180, MC should be ~180
	mc180 := mcFromRAMC(180, 23.44)
	if math.Abs(mc180-180) > 0.5 {
		t.Errorf("mcFromRAMC(180, 23.44) = %f, want ~180", mc180)
	}
}

func TestAscFromRAMC(t *testing.T) {
	asc := ascFromRAMC(0, 23.44, 39.9)
	// At RAMC=0, latitude=40°N, ASC should be somewhere around 90°+
	if asc < 0 || asc >= 360 {
		t.Errorf("ascFromRAMC(0, 23.44, 39.9) = %f, want 0-360", asc)
	}
}

func TestCalcProgressedAngles(t *testing.T) {
	natalJD := 2448057.5208  // 1990-07-12
	transitJD := 2460310.667 // ~2024

	asc, mc, err := CalcProgressedAngles(natalJD, transitJD, 39.9, 116.4, models.HousePlacidus, 0, 0, 0)
	if err != nil {
		t.Fatalf("CalcProgressedAngles error: %v", err)
	}
	if asc < 0 || asc >= 360 {
		t.Errorf("ASC = %f, want 0-360", asc)
	}
	if mc < 0 || mc >= 360 {
		t.Errorf("MC = %f, want 0-360", mc)
	}
}

func TestCalcProgressedSpecialPoint(t *testing.T) {
	natalJD := 2448057.5208
	transitJD := 2460310.667

	tests := []struct {
		sp models.SpecialPointID
	}{
		{models.PointASC},
		{models.PointMC},
		{models.PointDSC},
		{models.PointIC},
	}
	for _, tt := range tests {
		lon, err := CalcProgressedSpecialPoint(tt.sp, natalJD, transitJD, 39.9, 116.4, models.HousePlacidus, 0, 0, 0)
		if err != nil {
			t.Errorf("CalcProgressedSpecialPoint(%s) error: %v", tt.sp, err)
			continue
		}
		if lon < 0 || lon >= 360 {
			t.Errorf("CalcProgressedSpecialPoint(%s) = %f, want 0-360", tt.sp, lon)
		}
	}

	// Unsupported point
	_, err := CalcProgressedSpecialPoint(models.PointVertex, natalJD, transitJD, 39.9, 116.4, models.HousePlacidus, 0, 0, 0)
	if err == nil {
		t.Error("Expected error for unsupported special point Vertex")
	}
}
