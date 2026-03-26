package main

import (
	"fmt"
)

const (
	JulianYear = 365.25
	natalJD    = 2450800.900000 // 1997-12-18 09:36:00 UTC
)

func main() {
	fmt.Println("=== 直接测试Solar Fire因子 ===")
	
	// 从逆向计算得到的Solar Fire因子
	sfFactors := []float64{
		0.999924, // 从Sp-Sp案例
		1.000054, // 从Tr-Sa案例
	}
	
	testJD := natalJD + 10000
	ourMethod := natalJD + (testJD-natalJD)/JulianYear
	
	fmt.Printf("我们的方法: %.6f\n", ourMethod)
	fmt.Println()
	
	for i, sfFactor := range sfFactors {
		sfMethod := natalJD + (testJD-natalJD)/sfFactor
		diffHours := (sfMethod - ourMethod) * 24
		
		fmt.Printf("Solar Fire因子%d: %.6f\n", i+1, sfFactor)
		fmt.Printf("  计算结果: %.6f\n", sfMethod)
		fmt.Printf("  与我们差异: %.1f小时\n", diffHours)
		
		// 验证能否重现具体案例
		testCases := []struct {
			name    string
			sfJD    float64
			ourJD   float64
		}{
			{
				name:  "Neptune Opposition ASC Sp-Sp",
				sfJD:  2461183.005012,
				ourJD: 2461182.216817,
			},
			{
				name:  "Uranus Square Uranus Tr-Sa",
				sfJD:  2461294.037361,
				ourJD: 2461294.602257,
			},
		}
		
		fmt.Printf("  案例验证:\n")
		for _, tc := range testCases {
			// 使用SF因子重新计算
			eventJD := tc.ourJD
			calculatedSFJD := natalJD + (eventJD-natalJD)/sfFactor
			diffSeconds := (calculatedSFJD - tc.sfJD) * 86400
			
			fmt.Printf("    %s: 差异 %.1f秒\n", tc.name, diffSeconds)
		}
		fmt.Println()
	}
	
	// 测试平均因子
	avgFactor := (sfFactors[0] + sfFactors[1]) / 2
	fmt.Printf("平均因子: %.6f\n", avgFactor)
	avgMethod := natalJD + (testJD-natalJD)/avgFactor
	avgDiff := (avgMethod - ourMethod) * 24
	fmt.Printf("平均方法结果: %.6f (差异: %.1f小时)\n", avgMethod, avgDiff)
	
	// 分析因子的数学含义
	fmt.Println("\n=== 因子数学分析 ===")
	
	// SF因子 = (event_JD - natal_JD) / (SF_time_JD - natal_JD)
	// 这意味着 SF_time_JD = natal_JD + (event_JD - natal_JD) / SF_factor
	
	// 对于第一个案例
	eventJD1 := 2461182.216817
	sfJD1 := 2461183.005012
	sfFactor1 := (eventJD1 - natalJD) / (sfJD1 - natalJD)
	
	fmt.Printf("案例1反向计算因子: %.6f\n", sfFactor1)
	fmt.Printf("与观测因子差异: %.6f\n", sfFactor1-0.999924)
	
	// 这揭示了Solar Fire的真实算法：
	// Solar Fire不是改变年长度，而是改变了进展的映射关系
	// SF_progressed_JD = natal_JD + (event_JD - natal_JD) / SF_factor
	
	fmt.Println("\n=== Solar Fire真实算法推测 ===")
	fmt.Printf("SF进展公式: progressed_JD = natal_JD + (event_JD - natal_JD) / %.6f\n", avgFactor)
	fmt.Printf("我们公式:    progressed_JD = natal_JD + (event_JD - natal_JD) / %.6f\n", JulianYear)
	fmt.Printf("差异比例: %.6f\n", JulianYear/avgFactor)
	
	// 计算这种差异对应的时间影响
	timeRatio := JulianYear / avgFactor
	fmt.Printf("时间压缩比例: %.6f (即Solar Fire的时间流逝%.1f%%快于我们的计算)\n", 
		timeRatio, (timeRatio-1)*100)
}