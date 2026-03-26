package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// TestCase1 参数 (JN案例)
const (
	natalJD = 2450800.900009 // 1997-12-18 09:36 UTC
)

func main() {
	// 初始化星历
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	defer sweph.Close()

	fmt.Println("=== Progression Epoch Analysis ===")
	fmt.Printf("Natal JD: %.6f (1997-12-18 09:36 UTC)\n", natalJD)

	// 选择一个SF有记录的progression事件进行详细分析
	// 例如: Tr-Sp事件 Uranus Square Moon at 2026-05-08 16:02:33 AWST
	// UTC时间: 2026-05-08 08:02:33
	transitTime := time.Date(2026, 5, 8, 8, 2, 33, 0, time.UTC)
	transitJD := julianDate(transitTime)
	
	fmt.Printf("\nAnalyzing Tr-Sp event at:\n")
	fmt.Printf("  UTC: %s\n", transitTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  JD:  %.6f\n", transitJD)

	// 1. 当前的progression计算方法
	fmt.Println("\n=== Current Progression Calculation ===")
	currentProgJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
	currentAge := (transitJD - natalJD) / 365.25
	fmt.Printf("  Transit JD: %.6f\n", transitJD)
	fmt.Printf("  Current Progressed JD: %.6f\n", currentProgJD)
	fmt.Printf("  Age: %.2f years\n", currentAge)

	// 计算progressed太阳位置
	currentProgSun, err := sweph.CalcUT(currentProgJD, sweph.SE_SUN)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("  Progressed Sun: %.4f°\n", currentProgSun.Longitude)

	// 2. SF期望的太阳位置分析
	// 我们需要反向推导SF使用的progressed JD
	// 假设SF的误差是由于progressed JD偏差造成的
	
	// 让我们尝试不同的natal epoch假设
	fmt.Println("\n=== Testing Different Natal Epochs ===")
	
	// 原始natal JD
	testEpoch("Original Natal JD", natalJD, transitJD)
	
	// 尝试不同的调整
	offsets := []float64{
		0.2218,  // +5.32小时 (之前分析的偏差)
		0.5,     // +12小时
		1.0,     // +1天
		-0.2218, // -5.32小时
	}
	
	for _, offset := range offsets {
		desc := fmt.Sprintf("Natal JD %+0.4f days", offset)
		testEpoch(desc, natalJD+offset, transitJD)
	}

	// 测试使用JDE而不是JD_UT
	fmt.Println("\n=== Testing with JDE (TT) instead of JD_UT ===")
	deltaT := sweph.DeltaT(natalJD)
	natalJDE := natalJD + deltaT
	fmt.Printf("  ΔT: %.3f seconds\n", deltaT*86400)
	fmt.Printf("  Natal JDE: %.6f\n", natalJDE)
	
	progJDE := progressions.SecondaryProgressionJD(natalJDE, transitJD)
	progSunJDE, _ := sweph.CalcUT(progJDE, sweph.SE_SUN)
	natalSunJDE, _ := sweph.CalcUT(natalJDE, sweph.SE_SUN)
	solarArcJDE := sweph.NormalizeDegrees(progSunJDE.Longitude - natalSunJDE.Longitude)
	
	fmt.Printf("  Progressed JDE: %.6f\n", progJDE)
	fmt.Printf("  Progressed Sun (JDE): %.4f°\n", progSunJDE.Longitude)
	fmt.Printf("  Solar Arc (JDE): %.4f°\n", solarArcJDE)
	fmt.Printf("  Difference from UT: %+.4f°\n", solarArcJDE-28.9233)
}

func testEpoch(description string, epochJD, transitJD float64) {
	fmt.Printf("\n--- %s ---\n", description)
	fmt.Printf("  Epoch JD: %.6f\n", epochJD)
	
	progJD := progressions.SecondaryProgressionJD(epochJD, transitJD)
	age := (transitJD - epochJD) / 365.25
	
	progSun, err := sweph.CalcUT(progJD, sweph.SE_SUN)
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
		return
	}
	
	// 计算natal太阳位置
	natalSun, _ := sweph.CalcUT(epochJD, sweph.SE_SUN)
	solarArc := sweph.NormalizeDegrees(progSun.Longitude - natalSun.Longitude)
	
	fmt.Printf("  Progressed JD: %.6f\n", progJD)
	fmt.Printf("  Age: %.2f years\n", age)
	fmt.Printf("  Natal Sun: %.4f°\n", natalSun.Longitude)
	fmt.Printf("  Progressed Sun: %.4f°\n", progSun.Longitude)
	fmt.Printf("  Solar Arc: %.4f°\n", solarArc)
}

func analyzeSpecificEvents() {
	// 分析几个典型的错误事件
	events := []struct {
		name      string
		sfTime    string  // AWST时间
		chartType string
		sfDiffSec float64 // SF报告的时间差(秒)
	}{
		{"Uranus Square Moon Tr-Sp", "2026-05-08 16:02:33", "Tr-Sp", -26061},
		{"Saturn Quincunx Moon Tr-Sp", "2027-01-15 00:10:46", "Tr-Sp", -23174},
		{"ASC Semi-Square NorthNode Sp-Na", "2026-04-29 17:18:12", "Sp-Na", -42424},
	}

	awstLoc, _ := time.LoadLocation("Australia/Perth")

	for _, evt := range events {
		fmt.Printf("\n--- %s ---\n", evt.name)
		
		// 解析AWST时间
		t, err := time.ParseInLocation("2006-01-02 15:04:05", evt.sfTime, awstLoc)
		if err != nil {
			fmt.Printf("  Time parse error: %v\n", err)
			continue
		}
		utcTime := t.UTC()
		transitJD := julianDate(utcTime)
		
		fmt.Printf("  AWST: %s\n", evt.sfTime)
		fmt.Printf("  UTC:  %s\n", utcTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  JD:   %.6f\n", transitJD)
		fmt.Printf("  SF Diff: %+.0fs (%+.1fmin)\n", evt.sfDiffSec, evt.sfDiffSec/60)
		
		// 计算当前方法的progressed位置
		progJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
		progSun, _ := sweph.CalcUT(progJD, sweph.SE_SUN)
		natalSun, _ := sweph.CalcUT(natalJD, sweph.SE_SUN)
		solarArc := sweph.NormalizeDegrees(progSun.Longitude - natalSun.Longitude)
		
		fmt.Printf("  Current Progressed Sun: %.4f°\n", progSun.Longitude)
		fmt.Printf("  Current Solar Arc: %.4f°\n", solarArc)
		
		// 反向计算SF可能使用的progressed JD
		// 如果SF的时间比我们早，说明他们用的progressed JD更小
		targetDiffDays := evt.sfDiffSec / 86400.0
		targetProgJD := progJD + targetDiffDays
		targetProgSun, _ := sweph.CalcUT(targetProgJD, sweph.SE_SUN)
		targetSolarArc := sweph.NormalizeDegrees(targetProgSun.Longitude - natalSun.Longitude)
		
		fmt.Printf("  Target Progressed JD: %.6f (diff: %+.1fd)\n", targetProgJD, targetProgJD-progJD)
		fmt.Printf("  Target Progressed Sun: %.4f°\n", targetProgSun.Longitude)
		fmt.Printf("  Target Solar Arc: %.4f°\n", targetSolarArc)
		fmt.Printf("  Solar Arc Diff: %+.4f°\n", targetSolarArc-solarArc)
		
		// 推导所需的natal epoch偏移
		requiredEpochOffset := deriveRequiredEpochOffset(natalJD, transitJD, targetProgJD)
		fmt.Printf("  Required Epoch Offset: %+.4f days (%+.1f hours)\n", requiredEpochOffset, requiredEpochOffset*24)
	}
}

func deriveRequiredEpochOffset(originalNatalJD, transitJD, targetProgJD float64) float64 {
	// SecondaryProgressionJD公式: progJD = natalJD + (transitJD - natalJD) / 365.25
	// 要使 progJD = targetProgJD
	// targetProgJD = natalJD + (transitJD - natalJD) / 365.25
	// targetProgJD - natalJD = (transitJD - natalJD) / 365.25
	// (targetProgJD - natalJD) * 365.25 = transitJD - natalJD
	// natalJD = transitJD - (targetProgJD - natalJD) * 365.25
	
	// 但这形成了自引用。实际上我们要找的是:
	// targetProgJD = newNatalJD + (transitJD - newNatalJD) / 365.25
	// 解这个方程求 newNatalJD
	
	// targetProgJD = newNatalJD + (transitJD - newNatalJD) / 365.25
	// targetProgJD = newNatalJD + transitJD/365.25 - newNatalJD/365.25
	// targetProgJD = newNatalJD*(1 - 1/365.25) + transitJD/365.25
	// targetProgJD - transitJD/365.25 = newNatalJD*(1 - 1/365.25)
	// newNatalJD = (targetProgJD - transitJD/365.25) / (1 - 1/365.25)
	
	factor := 1.0 - 1.0/365.25
	newNatalJD := (targetProgJD - transitJD/365.25) / factor
	return newNatalJD - originalNatalJD
}

func julianDate(t time.Time) float64 {
	year, month, day := t.Date()
	hour := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
	return sweph.JulDay(year, int(month), day, hour, true)
}