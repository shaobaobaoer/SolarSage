package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	defer sweph.Close()

	const natalJD = 2450800.900009
	const transitJD = 2461168.835104 // 2026-05-08 08:02:33 UTC

	fmt.Println("=== Progression Formula Analysis ===")
	
	// 标准Secondary Progression公式验证
	fmt.Println("1. Standard Secondary Progression Formula:")
	fmt.Println("   Progressed JD = Natal JD + (Transit JD - Natal JD) / 365.25")
	
	currentProgJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
	fmt.Printf("   Current result: %.6f\n", currentProgJD)
	
	// 验证数学
	expected := natalJD + (transitJD-natalJD)/365.25
	fmt.Printf("   Manual calculation: %.6f\n", expected)
	fmt.Printf("   Difference: %.9f\n", currentProgJD-expected)
	
	// 分析年龄计算
	age := (transitJD - natalJD) / 365.25
	fmt.Printf("\n2. Age calculation:")
	fmt.Printf("   Age = (Transit JD - Natal JD) / 365.25\n")
	fmt.Printf("   Age = (%.6f - %.6f) / 365.25 = %.4f years\n", transitJD, natalJD, age)
	
	// 关键问题：natal epoch的选择
	fmt.Println("\n3. Natal Epoch Considerations:")
	fmt.Printf("   Birth time (UTC): 1997-12-18 09:36:00 (JD %.6f)\n", natalJD)
	
	// 方案1：使用出生时刻作为natal epoch
	fmt.Println("   Option 1: Use birth moment as natal epoch")
	fmt.Printf("     Natal epoch JD: %.6f\n", natalJD)
	
	// 方案2：使用出生日正午作为natal epoch (一些传统做法)
	noonJD := math.Floor(natalJD) + 0.5 // 当日正午
	fmt.Println("   Option 2: Use birth day noon as natal epoch")
	fmt.Printf("     Noon JD: %.6f (diff: %+.3fd)\n", noonJD, noonJD-natalJD)
	
	// 方案3：使用出生日子夜作为natal epoch
	midnightJD := math.Floor(natalJD) // 当日子夜
	fmt.Println("   Option 3: Use birth day midnight as natal epoch")
	fmt.Printf("     Midnight JD: %.6f (diff: %+.3fd)\n", midnightJD, midnightJD-natalJD)
	
	// 计算各种方案的结果
	fmt.Println("\n4. Comparing Results:")
	
	scenarios := []struct {
		name string
		epochJD float64
	}{
		{"Birth Moment", natalJD},
		{"Birth Day Noon", noonJD},
		{"Birth Day Midnight", midnightJD},
	}
	
	for _, scenario := range scenarios {
		progJD := progressions.SecondaryProgressionJD(scenario.epochJD, transitJD)
		progSun, _ := sweph.CalcUT(progJD, sweph.SE_SUN)
		natalSun, _ := sweph.CalcUT(scenario.epochJD, sweph.SE_SUN)
		solarArc := sweph.NormalizeDegrees(progSun.Longitude - natalSun.Longitude)
		
		age := (transitJD - scenario.epochJD) / 365.25
		
		fmt.Printf("   %s:\n", scenario.name)
		fmt.Printf("     Epoch JD: %.6f\n", scenario.epochJD)
		fmt.Printf("     Progressed JD: %.6f\n", progJD)
		fmt.Printf("     Age: %.2f years\n", age)
		fmt.Printf("     Solar Arc: %.4f°\n", solarArc)
	}
	
	// 最重要的是：检查SF是否使用了不同的year长度
	fmt.Println("\n5. Year Length Considerations:")
	
	yearLengths := []struct {
		name string
		length float64
	}{
		{"Julian Year (current)", 365.25},
		{"Tropical Year", 365.2422},
		{"Sidereal Year", 365.2564},
		{"Anomalistic Year", 365.2596},
	}
	
	for _, yl := range yearLengths {
		// 使用birth day noon场景
		modifiedProgJD := noonJD + (transitJD-noonJD)/yl.length
		progSun, _ := sweph.CalcUT(modifiedProgJD, sweph.SE_SUN)
		natalSun, _ := sweph.CalcUT(noonJD, sweph.SE_SUN)
		solarArc := sweph.NormalizeDegrees(progSun.Longitude - natalSun.Longitude)
		
		fmt.Printf("   %s (%.4f days):\n", yl.name, yl.length)
		fmt.Printf("     Progressed JD: %.6f\n", modifiedProgJD)
		fmt.Printf("     Solar Arc: %.4f°\n", solarArc)
		fmt.Printf("     Diff from Julian: %+.4f°\n", solarArc-28.6162)
	}
}