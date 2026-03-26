package main

import (
	"fmt"
)

const (
	JulianYear = 365.25
	natalJD    = 2450800.900000 // 1997-12-18 09:36:00 UTC
)

func main() {
	fmt.Println("=== 快速因子分析 ===")
	
	// 从实际案例计算Solar Fire因子
	cases := []struct {
		name  string
		sfJD  float64
		ourJD float64
	}{
		{"案例1 Sp-Sp", 2461183.005012, 2461182.216817},
		{"案例2 Tr-Sa", 2461294.037361, 2461294.602257},
	}
	
	var factors []float64
	
	for _, c := range cases {
		sfFactor := (c.ourJD - natalJD) / (c.sfJD - natalJD)
		factors = append(factors, sfFactor)
		fmt.Printf("%s: SF因子 = %.6f\n", c.name, sfFactor)
	}
	
	// 计算平均因子
	avgFactor := (factors[0] + factors[1]) / 2
	fmt.Printf("平均因子: %.6f\n", avgFactor)
	
	// 测试这个因子的表现
	testJD := natalJD + 10000
	ourMethod := natalJD + (testJD-natalJD)/JulianYear
	sfMethod := natalJD + (testJD-natalJD)/avgFactor
	
	timeDiffHours := (sfMethod - ourMethod) * 24
	fmt.Printf("时间差异: %.1f小时\n", timeDiffHours)
	
	// 这个差异太大了，让我尝试另一种理解
	// 也许Solar Fire使用的是: progressed_JD = natal_JD + (transit_JD - natal_JD) * factor
	// 其中factor是一个很小的数
	
	alternativeFactor := avgFactor / JulianYear
	fmt.Printf("替代因子 (乘法形式): %.6f\n", alternativeFactor)
	
	altMethod := natalJD + (testJD-natalJD)*alternativeFactor
	altTimeDiff := (altMethod - ourMethod) * 24
	fmt.Printf("替代方法时间差异: %.1f小时\n", altTimeDiff)
	
	// 寻找合适的乘法因子使差异在合理范围
	targetHours := 5.3 // 我们观察到的典型差异
	requiredFactor := targetHours / 24.0 / (testJD - natalJD) * JulianYear
	
	fmt.Printf("目标差异 %.1f小时所需的因子: %.6f\n", targetHours, requiredFactor)
	
	// 验证这个因子
	finalMethod := natalJD + (testJD-natalJD)*requiredFactor
	finalDiff := (finalMethod - ourMethod) * 24
	fmt.Printf("最终方法时间差异: %.1f小时\n", finalDiff)
	
	fmt.Printf("\n结论: Solar Fire可能使用因子 %.6f\n", requiredFactor)
}