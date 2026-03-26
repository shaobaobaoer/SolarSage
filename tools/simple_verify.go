package main

import (
	"fmt"
	"math"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 验证Solar Fire进展修正 ===")
	
	// 使用相同的出生数据
	natalDatetime := "1997-12-18T09:36:00Z"
	year, month, day := 1997, 12, 18
	hour := 9.6 // 09:36
	
	natalJD := sweph.JulDay(year, month, day, hour, true)
	fmt.Printf("出生JD: %.6f (%s)\n", natalJD, natalDatetime)
	
	// 测试时间点：2026年的某个时间
	testYear, testMonth, testDay := 2026, 6, 1
	testHour := 12.0
	testJD := sweph.JulDay(testYear, testMonth, testDay, testHour, true)
	fmt.Printf("测试JD: %.6f (%04d-%02d-%02d %.1f:00)\n", testJD, testYear, testMonth, testDay, testHour)
	
	// 计算时间差
	daysDiff := testJD - natalJD
	yearsDiff := daysDiff / 365.25
	fmt.Printf("时间跨度: %.1f天 = %.2f年\n", daysDiff, yearsDiff)
	
	// 比较两种方法
	fmt.Println("\n=== 进展计算对比 ===")
	
	// 传统方法 (Julian Year = 365.25)
	traditionalProgJD := natalJD + daysDiff/365.25
	fmt.Printf("传统方法进展JD: %.6f\n", traditionalProgJD)
	
	// Solar Fire方法 (修正因子 = 0.999989)
	sfProgJD := progressions.SecondaryProgressionJD(natalJD, testJD)
	fmt.Printf("Solar Fire方法进展JD: %.6f\n", sfProgJD)
	
	// 差异分析
	jdDiff := sfProgJD - traditionalProgJD
	timeDiffHours := jdDiff * 24
	fmt.Printf("JD差异: %.6f\n", jdDiff)
	fmt.Printf("时间差异: %.1f小时\n", timeDiffHours)
	
	// 验证行星位置计算
	fmt.Println("\n=== 行星位置验证 ===")
	
	planet := models.PlanetSun
	fmt.Printf("计算 %s 在不同进展时间的位置:\n", planet)
	
	// 传统方法位置
	tradLon, tradSpeed, err := chart.CalcPlanetLongitude(planet, traditionalProgJD)
	if err != nil {
		fmt.Printf("传统方法计算错误: %v\n", err)
	} else {
		fmt.Printf("  传统方法: %.4f° (速度: %.6f°/day)\n", tradLon, tradSpeed)
	}
	
	// Solar Fire方法位置
	sfLon, sfSpeed, err := progressions.CalcProgressedLongitude(planet, natalJD, testJD)
	if err != nil {
		fmt.Printf("Solar Fire方法计算错误: %v\n", err)
	} else {
		fmt.Printf("  Solar Fire: %.4f° (速度: %.6f°/day)\n", sfLon, sfSpeed)
		fmt.Printf("  位置差异: %.4f°\n", math.Abs(sfLon-tradLon))
	}
	
	// 年龄计算验证
	fmt.Println("\n=== 年龄计算验证 ===")
	
	tradAge := (testJD - natalJD) / 365.25
	sfAge := progressions.Age(natalJD, testJD)
	
	fmt.Printf("传统年龄计算: %.2f岁\n", tradAge)
	fmt.Printf("Solar Fire年龄计算: %.2f岁\n", sfAge)
	fmt.Printf("年龄差异: %.2f岁\n", sfAge-tradAge)
	
	fmt.Println("\n=== 总结 ===")
	fmt.Printf("修正后时间差异: %.1f小时\n", timeDiffHours)
	if math.Abs(timeDiffHours) < 1.0 {
		fmt.Println("✓ 时间差异在可接受范围内 (< 1小时)")
	} else if math.Abs(timeDiffHours) < 5.0 {
		fmt.Println("⚠ 时间差异中等 (1-5小时)")
	} else {
		fmt.Println("✗ 时间差异较大 (> 5小时)")
	}
}