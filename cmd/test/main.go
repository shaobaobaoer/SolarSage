package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anthropic/swisseph-mcp/pkg/chart"
	"github.com/anthropic/swisseph-mcp/pkg/geo"
	"github.com/anthropic/swisseph-mcp/pkg/julian"
	"github.com/anthropic/swisseph-mcp/pkg/models"
	"github.com/anthropic/swisseph-mcp/pkg/sweph"
	"github.com/anthropic/swisseph-mcp/pkg/transit"
)

func main() {
	// Initialize ephemeris
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	if _, err := os.Stat(ephePath); err != nil {
		ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
	}
	sweph.Init(ephePath)
	defer sweph.Close()

	fmt.Println("========================================")
	fmt.Println("  Swisseph MCP 功能测试")
	fmt.Println("========================================")

	// Test 1: Geocode
	fmt.Println("\n--- 测试 2.1: 地点名称 → 经纬度 ---")
	loc, err := geo.Geocode("北京")
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  北京: lat=%.4f, lon=%.4f, tz=%s, name=%s\n",
			loc.Latitude, loc.Longitude, loc.Timezone, loc.DisplayName)
	}

	loc2, err := geo.Geocode("london")
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  London: lat=%.4f, lon=%.4f, tz=%s\n",
			loc2.Latitude, loc2.Longitude, loc2.Timezone)
	}

	// Test 2: DateTime to JD
	fmt.Println("\n--- 测试 2.2: 公历时间 → 儒略日 ---")
	jdResult, err := julian.DateTimeToJD("1990-06-15T08:30:00+08:00", models.CalendarGregorian)
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  1990-06-15T08:30:00+08:00:\n")
		fmt.Printf("    JD(UT) = %.6f\n", jdResult.JDUT)
		fmt.Printf("    JD(TT) = %.6f\n", jdResult.JDTT)
	}

	// Test 2b: JD to DateTime
	fmt.Println("\n--- 测试 2.2b: 儒略日 → 公历时间 ---")
	dt, err := julian.JDToDateTime(jdResult.JDUT, "Asia/Shanghai")
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  JD %.6f → %s\n", jdResult.JDUT, dt)
	}

	// Test 3: Single Chart Calculation
	fmt.Println("\n--- 测试 3.1.1: 单盘计算 ---")
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}
	orbs := models.DefaultOrbConfig()

	chartInfo, err := chart.CalcSingleChart(
		39.9042, 116.4074, // 北京
		jdResult.JDUT,
		planets, orbs, models.HousePlacidus,
	)
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  出生盘 (1990-06-15 08:30 北京)\n")
		fmt.Printf("  ASC: %.2f° (%s)\n", chartInfo.Angles.ASC, models.SignFromLongitude(chartInfo.Angles.ASC))
		fmt.Printf("  MC:  %.2f° (%s)\n", chartInfo.Angles.MC, models.SignFromLongitude(chartInfo.Angles.MC))
		fmt.Println("  行星位置:")
		for _, p := range chartInfo.Planets {
			retro := ""
			if p.IsRetrograde {
				retro = " (R)"
			}
			fmt.Printf("    %-10s %6.2f° %s %5.2f°  宫%d%s\n",
				p.PlanetID, p.Longitude, p.Sign, p.SignDegree, p.House, retro)
		}
		fmt.Printf("  相位数量: %d\n", len(chartInfo.Aspects))
		for i, a := range chartInfo.Aspects {
			if i >= 5 {
				fmt.Printf("    ... (共 %d 个相位)\n", len(chartInfo.Aspects))
				break
			}
			applying := "离相"
			if a.IsApplying {
				applying = "入相"
			}
			fmt.Printf("    %s %s %s (容许度 %.2f°, %s)\n",
				a.PlanetA, a.AspectType, a.PlanetB, a.Orb, applying)
		}
	}

	// Test 4: Double Chart
	fmt.Println("\n--- 测试 3.1.2: 双盘计算 ---")
	// Current transit time: 2024-01-01
	transitJD, _ := julian.DateTimeToJD("2024-01-01T12:00:00+08:00", models.CalendarGregorian)
	innerChart, outerChart, crossAspects, err := chart.CalcDoubleChart(
		39.9042, 116.4074, jdResult.JDUT, planets,
		39.9042, 116.4074, transitJD.JDUT, planets,
		&models.SpecialPointsConfig{
			InnerPoints: []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		orbs, models.HousePlacidus,
	)
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  内圈行星: %d, 外圈行星: %d\n", len(innerChart.Planets), len(outerChart.Planets))
		fmt.Printf("  跨盘相位数量: %d\n", len(crossAspects))
		for i, ca := range crossAspects {
			if i >= 5 {
				fmt.Printf("    ... (共 %d 个跨盘相位)\n", len(crossAspects))
				break
			}
			fmt.Printf("    内圈%s %s 外圈%s (容许度 %.2f°)\n",
				ca.InnerBody, ca.AspectType, ca.OuterBody, ca.Orb)
		}
	}

	// Test 5: Transit calculation
	fmt.Println("\n--- 测试 3.2.1: 推运计算 ---")
	// Search transit events for 30 days starting from 2024-01-01
	startJD := transitJD.JDUT
	endJD := startJD + 30.0 // 30 days

	transitEvents, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalLat:     39.9042,
		NatalLon:     116.4074,
		NatalJD:      jdResult.JDUT,
		NatalPlanets: planets,
		TransitLat:   39.9042,
		TransitLon:   116.4074,
		StartJD:      startJD,
		EndJD:        endJD,
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMercury, models.PlanetVenus, models.PlanetMars,
		},
		EventConfig: models.DefaultEventConfig(),
		OrbConfig:   orbs,
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
	} else {
		fmt.Printf("  搜索范围: 30天\n")
		fmt.Printf("  找到事件: %d 个\n", len(transitEvents))

		// Show some events
		shown := 0
		for _, e := range transitEvents {
			if shown >= 15 {
				fmt.Printf("    ... (共 %d 个事件)\n", len(transitEvents))
				break
			}
			retro := ""
			if e.IsRetrograde {
				retro = "(R)"
			}
			dtStr, _ := julian.JDToDateTime(e.JD, "Asia/Shanghai")
			switch e.EventType {
			case models.EventAspectEnter:
				fmt.Printf("    %s %s 进入 %s %s (容许度 %.2f°) %s %s\n",
					dtStr, e.Planet, e.Target, e.AspectType, e.OrbAtEnter, e.PlanetSign, retro)
			case models.EventAspectExact:
				fmt.Printf("    %s %s 精确 %s %s (第%d击) %s %s\n",
					dtStr, e.Planet, e.Target, e.AspectType, e.ExactCount, e.PlanetSign, retro)
			case models.EventAspectLeave:
				fmt.Printf("    %s %s 离开 %s %s (容许度 %.2f°) %s %s\n",
					dtStr, e.Planet, e.Target, e.AspectType, e.OrbAtLeave, e.PlanetSign, retro)
			case models.EventSignIngress:
				fmt.Printf("    %s %s 换座 %s → %s %s\n",
					dtStr, e.Planet, e.FromSign, e.ToSign, retro)
			case models.EventHouseIngress:
				fmt.Printf("    %s %s 变宫 %d宫 → %d宫 %s\n",
					dtStr, e.Planet, e.FromHouse, e.ToHouse, retro)
			case models.EventStation:
				fmt.Printf("    %s %s 站点 %s %s %s\n",
					dtStr, e.Planet, e.StationType, e.PlanetSign, retro)
			}
			shown++
		}
	}

	// Output a sample event as JSON
	if len(transitEvents) > 0 {
		fmt.Println("\n--- 示例 JSON 输出 ---")
		sample := transitEvents[0]
		j, _ := json.MarshalIndent(sample, "  ", "  ")
		fmt.Printf("  %s\n", string(j))
	}

	fmt.Println("\n========================================")
	fmt.Println("  所有测试完成!")
	fmt.Println("========================================")
}
