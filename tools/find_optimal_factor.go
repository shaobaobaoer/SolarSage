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
	fmt.Println("=== 重新分析Solar Fire因子 ===")
	
	// 从实际案例中获取的数据
	testCases := []struct {
		name    string
		sfJD    float64  // Solar Fire时间JD
		ourJD   float64  // 我们计算的时间JD
	}{
		{
			name:  "Neptune Opposition ASC Sp-Sp",
			sfJD:  2461183.005012,   // 2026-05-22 12:07:13
			ourJD: 2461182.216817,   // 2026-05-21 17:12:13
		},
		{
			name:  "Uranus Square Uranus Tr-Sa",
			sfJD:  2461294.037361,   // 2026-09-10 12:53:48
			ourJD: 2461294.602257,   // 2026-09-11 02:27:15
		},
	}
	
	// 对每个案例计算Solar Fire使用的因子
	fmt.Println("案例分析:")
	for _, tc := range testCases {
		// SF_progressed_JD = natal_JD + (event_JD - natal_JD) / SF_factor
		// 因此: SF_factor = (event_JD - natal_JD) / (SF_progressed_JD - natal_JD)
		
		eventJD := tc.ourJD
		sfFactor := (eventJD - natalJD) / (tc.sfJD - natalJD)
		
		fmt.Printf("\n%s:\n", tc.name)
		fmt.Printf("  SF时间JD: %.6f\n", tc.sfJD)
		fmt.Printf("  我们时间JD: %.6f\n", tc.ourJD)
		fmt.Printf("  JD差异: %.6f\n", tc.sfJD-tc.ourJD)
		fmt.Printf("  推导出的SF因子: %.6f\n", sfFactor)
		
		// 验证这个因子是否合理
		calculatedSFJD := natalJD + (eventJD-natalJD)/sfFactor
		verificationDiff := (calculatedSFJD - tc.sfJD) * 86400 // 秒
		fmt.Printf("  验证差异: %.1f秒\n", verificationDiff)
	}
	
	// 计算平均因子
	var sumFactor float64
	for _, tc := range testCases {
		eventJD := tc.ourJD
		sfFactor := (eventJD - natalJD) / (tc.sfJD - natalJD)
		sumFactor += sfFactor
	}
	avgFactor := sumFactor / float64(len(testCases))
	fmt.Printf("\n平均SF因子: %.6f\n", avgFactor)
	
	// 测试这个因子在不同时间跨度下的表现
	fmt.Println("\n=== 因子性能测试 ===")
	testPeriods := []float64{1000, 5000, 10000, 15000} // 天数
	
	for _, days := range testPeriods {
		testJD := natalJD + days
		ourProgJD := natalJD + days/JulianYear
		sfProgJD := natalJD + days/avgFactor
		
		timeDiffHours := (sfProgJD - ourProgJD) * 24
		fmt.Printf("时间跨度 %.0f天: 差异 %.1f小时\n", days, timeDiffHours)
	}
	
	// 寻找使差异在合理范围内的因子
	fmt.Println("\n=== 寻找最优因子 ===")
	targetRanges := []struct {
		min, max float64
		desc     string
	}{
		{1.0, 2.0, "轻微差异"},
		{5.0, 6.0, "中等差异"},
		{13.0, 14.0, "较大差异"},
		{18.0, 19.0, "显著差异"},
	}
	
	bestFactors := make(map[string]float64)
	
	for _, target := range targetRanges {
		bestFactor := 0.0
		bestDiff := math.MaxFloat64
		
		// 在合理范围内搜索因子
		for factor := 360.0; factor <= 370.0; factor += 0.000001 {
			testJD := natalJD + 10000 // 固定测试时间
			ourProgJD := natalJD + 10000/JulianYear
			sfProgJD := natalJD + 10000/factor
			
			diffHours := math.Abs((sfProgJD - ourProgJD) * 24)
			
			if diffHours >= target.min && diffHours <= target.max {
				if math.Abs(diffHours-(target.min+target.max)/2) < bestDiff {
					bestDiff = math.Abs(diffHours - (target.min+target.max)/2)
					bestFactor = factor
				}
			}
		}
		
		if bestFactor > 0 {
			bestFactors[target.desc] = bestFactor
			fmt.Printf("%s (%.1f-%.1f小时): 因子=%.6f\n", 
				target.desc, target.min, target.max, bestFactor)
		}
	}
	
	// 选择最适合整体匹配的因子
	if len(bestFactors) > 0 {
		// 使用中等差异的因子作为折中方案
		chosenFactor := bestFactors["中等差异"]
		fmt.Printf("\n推荐使用因子: %.6f\n", chosenFactor)
		
		// 验证在所有案例上的表现
		fmt.Println("\n=== 最终验证 ===")
		for _, tc := range testCases {
			eventJD := tc.ourJD
			sfProgJD := natalJD + (eventJD-natalJD)/chosenFactor
			timeDiffHours := (sfProgJD - tc.sfJD) * 24
			
			fmt.Printf("%s: 时间差异 %.1f小时\n", tc.name, timeDiffHours)
		}
	}
}