package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 理解Secondary Progression ===")
	
	// 出生: 1997-12-18 09:36:00 UTC
	natalJD := 2450800.900000
	
	// 28年后: 2025-12-18 09:36:00 UTC (约28.1岁)
	age28JD := sweph.JulDay(2025, 12, 18, 9.6, true)
	
	// 在secondary progression中:
	// 28岁对应出生后的第28天
	day28JD := natalJD + 28.0
	
	fmt.Printf("出生JD: %.6f (1997-12-18)\n", natalJD)
	fmt.Printf("28岁时JD: %.6f (2025-12-18)\n", age28JD)
	fmt.Printf("出生后第28天JD: %.6f (1998-01-15)\n", day28JD)
	fmt.Printf("时间差: %.1f天\n", age28JD-natalJD)
	
	// 计算太阳在这些时间点的位置
	fmt.Println("\n太阳位置:")
	
	natalSun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, natalJD)
	fmt.Printf("  出生时: %.2f°\n", natalSun)
	
	day28Sun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, day28JD)
	fmt.Printf("  出生后第28天: %.2f° (进展太阳位置)\n", day28Sun)
	
	age28Sun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, age28JD)
	fmt.Printf("  28岁时实际: %.2f°\n", age28Sun)
	
	// 进展偏移
	progOffset := day28Sun - natalSun
	for progOffset < 0 {
		progOffset += 360
	}
	fmt.Printf("\n进展偏移 (28岁): %.2f°\n", progOffset)
	
	// 太阳弧偏移 (实际太阳移动)
	arcOffset := age28Sun - natalSun
	for arcOffset < 0 {
		arcOffset += 360
	}
	fmt.Printf("太阳弧偏移 (28岁): %.2f°\n", arcOffset)
	
	// 测试2026-02-01 (约28.13岁)
	testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	daysSinceBirth := testJD - natalJD
	progDays := daysSinceBirth / 365.25
	progJD := natalJD + progDays
	
	fmt.Printf("\n2026-02-01测试:\n")
	fmt.Printf("  出生后天数: %.1f天\n", daysSinceBirth)
	fmt.Printf("  进展天数: %.1f天 (%.2f年)\n", progDays, progDays/365.25*1)
	fmt.Printf("  进展JD: %.6f\n", progJD)
	
	progSun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, progJD)
	fmt.Printf("  进展太阳位置: %.2f°\n", progSun)
	
	// 对比Solar Fire的算法
	// SF使用: progressed_JD = natal_JD + (transit_JD - natal_JD) * factor
	sfFactor := 0.008066
	sfProgJD := natalJD + (testJD-natalJD)*sfFactor
	fmt.Printf("\nSolar Fire进展JD: %.6f\n", sfProgJD)
	
	sfProgSun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, sfProgJD)
	fmt.Printf("Solar Fire进展太阳: %.2f°\n", sfProgSun)
}