package main

import (
	"fmt"
	"math"
)

const (
	JulianYear = 365.25
	natalJD    = 2450800.900000 // 1997-12-18 09:36:00 UTC
)

func main() {
	testJD := natalJD + 10000 // 测试10000天后的进展
	
	fmt.Printf("出生JD: %.6f\n", natalJD)
	fmt.Printf("测试JD: %.6f\n", testJD)
	
	// 我们的当前计算
	ourProgJD := natalJD + (testJD-natalJD)/JulianYear
	fmt.Printf("\n我们的计算: %.6f\n", ourProgJD)
	
	// 测试各种可能的Solar Fire因子
	fmt.Println("\n=== 测试不同的年长度 ===")
	
	factors := []struct {
		name   string
		factor float64
	}{
		{"Julian Year", 365.25},
		{"Tropical Year", 365.2422},      // 回归年
		{"Sidereal Year", 365.2564},      // 恒星年
		{"Anomalistic Year", 365.2596},   // 近点年
		{"Calendar Year", 365.0},         // 日历年
		{"Solar Fire推断1", 365.249924},   // 从Sp-Sp案例推断
		{"Solar Fire推断2", 365.250054},   // 从Tr-Sa案例推断
		{"Solar Fire平均", 365.249989},    // 两个案例的平均
	}
	
	bestMatch := ""
	minDiff := math.MaxFloat64
	
	for _, f := range factors {
		sfProgJD := natalJD + (testJD-natalJD)/f.factor
		diff := math.Abs(sfProgJD - ourProgJD) * 24 // 差异小时数
		
		fmt.Printf("%-20s: %.6f (差异: %.1f小时)\n", f.name, sfProgJD, diff)
		
		if diff < minDiff && f.factor != JulianYear {
			minDiff = diff
			bestMatch = f.name
		}
	}
	
	fmt.Printf("\n最佳匹配: %s (差异: %.1f小时)\n", bestMatch, minDiff)
	
	// 更精细地搜索Solar Fire因子
	fmt.Println("\n=== 精细搜索Solar Fire因子 ===")
	
	// 根据我们观察到的差异范围(1-19小时)，计算对应的因子范围
	targetHours := []float64{1.1, 5.3, 13.6, 18.9} // 观察到的典型差异
	
	for _, targetHour := range targetHours {
		// 我们需要找到使得差异为targetHour的因子
		timeDiff := targetHour / 24.0 // 转换为天数
		sfFactor := 1.0 / (1.0/JulianYear + timeDiff/(testJD-natalJD))
		
		fmt.Printf("目标差异 %.1f小时 对应的因子: %.6f\n", targetHour, sfFactor)
	}
	
	// 验证逆向计算
	fmt.Println("\n=== 验证逆向计算 ===")
	
	// 使用我们观察到的实际事件时间来验证
	testCases := []struct {
		description string
		sfJD        float64  // Solar Fire计算的时间JD
		ourJD       float64  // 我们计算的时间JD
	}{
		{
			description: "Neptune Opposition ASC Sp-Sp",
			sfJD: 2461183.005012,   // 2026-05-22 12:07:13
			ourJD: 2461182.216817,  // 2026-05-21 17:12:13
		},
		{
			description: "Uranus Square Uranus Tr-Sa",
			sfJD: 2461294.037361,   // 2026-09-10 12:53:48
			ourJD: 2461294.602257,  // 2026-09-11 02:27:15
		},
	}
	
	for _, tc := range testCases {
		fmt.Printf("\n案例: %s\n", tc.description)
		fmt.Printf("  SF JD: %.6f\n", tc.sfJD)
		fmt.Printf("  我们 JD: %.6f\n", tc.ourJD)
		fmt.Printf("  JD差异: %.6f\n", tc.sfJD-tc.ourJD)
		
		// 反向计算Solar Fire使用的因子
		eventJD := tc.ourJD // 使用我们的计算作为事件JD
		sfFactor := (eventJD - natalJD) / (tc.sfJD - natalJD)
		
		fmt.Printf("  Solar Fire因子: %.6f\n", sfFactor)
		fmt.Printf("  与JulianYear差异: %.6f\n", sfFactor-JulianYear)
		
		// 验证这个因子是否能重现SF结果
		calculatedSFJD := natalJD + (eventJD-natalJD)/sfFactor
		fmt.Printf("  验证计算: %.6f (差异: %.1f秒)\n", calculatedSFJD, (calculatedSFJD-tc.sfJD)*86400)
	}
}