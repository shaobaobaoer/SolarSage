package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	if _, err := os.Stat(ephePath); err != nil {
		ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
	}
	sweph.Init(ephePath)
	defer sweph.Close()

	// SF期望: Chiron-NorthNode Semi-Square Exact at 2026-03-01 00:38:36 AWST
	// AWST = UTC+8, so UTC = 2026-02-28 16:38:36
	// SF期望: Jupiter-Chiron Square Exact at 2026-07-02 02:23:46 AWST
	// UTC = 2026-07-01 18:23:46
	// SF期望: Chiron-NorthNode Sextile Exact at 2026-08-02 12:04:18 AWST
	// UTC = 2026-08-02 04:04:18

	fmt.Println("=== 验证缺失的3个Tr-Tr事件 ===")

	type check struct {
		name string
		p1   models.PlanetID
		p2   models.PlanetID
		angle float64
		year, month, day, hour, min, sec int
	}

	checks := []check{
		{"Chiron-NorthNode Semi-Square", models.PlanetChiron, models.PlanetNorthNodeMean, 45,
			2026, 2, 28, 16, 38, 36},
		{"Jupiter-Chiron Square", models.PlanetJupiter, models.PlanetChiron, 90,
			2026, 7, 1, 18, 23, 46},
		{"Chiron-NorthNode Sextile", models.PlanetChiron, models.PlanetNorthNodeMean, 60,
			2026, 8, 2, 4, 4, 18},
	}

	for _, c := range checks {
		fmt.Printf("\n--- %s ---\n", c.name)
		h := float64(c.hour) + float64(c.min)/60.0 + float64(c.sec)/3600.0
		jd := sweph.JulDay(c.year, c.month, c.day, h, true)
		fmt.Printf("SF Exact time (UTC): %04d-%02d-%02d %02d:%02d:%02d  JD=%.6f\n",
			c.year, c.month, c.day, c.hour, c.min, c.sec, jd)

		// 计算两个行星在该时刻的位置
		lon1, _, _ := chart.CalcPlanetLongitude(c.p1, jd)
		lon2, _, _ := chart.CalcPlanetLongitude(c.p2, jd)

		diff := math.Abs(lon1 - lon2)
		if diff > 180 {
			diff = 360 - diff
		}
		orb := diff - c.angle

		fmt.Printf("  %s: %.4f°\n", c.p1, lon1)
		fmt.Printf("  %s: %.4f°\n", c.p2, lon2)
		fmt.Printf("  Angular diff: %.4f°\n", diff)
		fmt.Printf("  Expected aspect angle: %.0f°\n", c.angle)
		fmt.Printf("  Orb at SF exact time: %.4f° (%.1f arcsec)\n", math.Abs(orb), math.Abs(orb)*3600)

		// 扫描附近几天，看看aspect diff的变化趋势
		fmt.Println("  Scanning nearby days:")
		for d := -5.0; d <= 5.0; d += 1.0 {
			testJD := jd + d
			l1, _, _ := chart.CalcPlanetLongitude(c.p1, testJD)
			l2, _, _ := chart.CalcPlanetLongitude(c.p2, testJD)
			dd := math.Abs(l1 - l2)
			if dd > 180 {
				dd = 360 - dd
			}
			o := dd - c.angle
			marker := ""
			if math.Abs(o) < 0.01 {
				marker = " <-- EXACT"
			} else if math.Abs(o) <= 1.0 {
				marker = " (within 1° orb)"
			}
			year, month, day, hour := sweph.RevJul(testJD, true)
			fmt.Printf("    %04d-%02d-%02d %02d:%02d  orb=%+.4f°%s\n",
				year, month, day, int(hour), int((hour-float64(int(hour)))*60), o, marker)
		}
	}
}
