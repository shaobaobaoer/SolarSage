package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 简化进展计算测试 ===")
	
	// 出生数据: 1997-12-18 09:36:00 UTC
	natalJD := 2450800.900000
	
	// 测试时间: 2026-02-01 00:00:00 (从CSV中看到的第一个时间点)
	// AWST时区是UTC+8，所以00:00:00 AWST = 前一天的16:00:00 UTC
	testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	
	fmt.Printf("出生JD: %.6f\n", natalJD)
	fmt.Printf("测试JD: %.6f (2026-01-31 16:00 UTC)\n", testJD)
	fmt.Printf("时间跨度: %.1f天 = %.2f年\n", testJD-natalJD, (testJD-natalJD)/365.25)
	
	// 计算进展JD
	progJD := progressions.SecondaryProgressionJD(natalJD, testJD)
	fmt.Printf("\n进展JD: %.6f\n", progJD)
	
	// 计算几个关键行星的位置
	planets := []models.PlanetID{
		models.PlanetSun,
		models.PlanetMoon,
		models.PlanetMercury,
		models.PlanetMars,
		models.PlanetSaturn,
	}
	
	fmt.Println("\n行星位置对比:")
	fmt.Printf("%-10s %-15s %-15s\n", "行星", "本命位置", "进展位置")
	fmt.Println("----------------------------------------")
	
	for _, p := range planets {
		// 本命位置
		natalLon, _, err := chart.CalcPlanetLongitude(p, natalJD)
		if err != nil {
			fmt.Printf("%-10s 计算错误: %v\n", p, err)
			continue
		}
		
		// 进展位置
		progLon, _, err := progressions.CalcProgressedLongitude(p, natalJD, testJD)
		if err != nil {
			fmt.Printf("%-10s 进展计算错误: %v\n", p, err)
			continue
		}
		
		// 计算进展移动的距离
		movement := progLon - natalLon
		for movement < 0 {
			movement += 360
		}
		
		fmt.Printf("%-10s %6.2f°         %6.2f° (移动 %.2f°)\n", 
			p, natalLon, progLon, movement)
	}
	
	// 验证进展公式
	fmt.Println("\n=== 进展公式验证 ===")
	
	// 传统方法: progressed_JD = natal_JD + days/365.25
	tradProgJD := natalJD + (testJD-natalJD)/365.25
	fmt.Printf("传统进展JD: %.6f\n", tradProgJD)
	
	// Solar Fire方法
	sfProgJD := progressions.SecondaryProgressionJD(natalJD, testJD)
	fmt.Printf("Solar Fire进展JD: %.6f\n", sfProgJD)
	
	// 差异
	diffDays := sfProgJD - tradProgJD
	diffYears := diffDays / 365.25
	fmt.Printf("JD差异: %.6f天 = %.4f年\n", diffDays, diffYears)
	
	// 测试太阳弧
	fmt.Println("\n=== 太阳弧验证 ===")
	offset, err := progressions.SolarArcOffset(natalJD, testJD)
	if err != nil {
		fmt.Printf("太阳弧计算错误: %v\n", err)
	} else {
		fmt.Printf("太阳弧偏移: %.4f°\n", offset)
		
		// 验证太阳位置
		natalSun, _, _ := chart.CalcPlanetLongitude(models.PlanetSun, natalJD)
		progSun, _, _ := progressions.CalcProgressedLongitude(models.PlanetSun, natalJD, testJD)
		calculatedOffset := progSun - natalSun
		for calculatedOffset < 0 {
			calculatedOffset += 360
		}
		fmt.Printf("太阳进展偏移验证: %.4f° (期望: %.4f°)\n", calculatedOffset, offset)
	}
}