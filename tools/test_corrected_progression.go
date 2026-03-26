package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	fmt.Println("=== Testing Corrected Progression Formula ===")
	
	// SF的JDE值
	sfJDE := 2450800.900729
	deltaT := sfJDE - natalJD
	fmt.Printf("SF JDE: %.6f\n", sfJDE)
	fmt.Printf("ΔT from SF: %.3f seconds\n", deltaT*86400)
	
	// 我们计算的ΔT
	ourDeltaT := sweph.DeltaT(natalJD)
	fmt.Printf("Our ΔT calculation: %.3f seconds\n", ourDeltaT*86400)
	fmt.Printf("Difference: %.3f seconds\n", (deltaT-ourDeltaT)*86400)
	
	// 使用旧方法（JD_UT）
	oldProgJD := natalJD + (transitJD-natalJD)/365.25
	oldProgSun, _ := sweph.CalcUT(oldProgJD, sweph.SE_SUN)
	oldNatalSun, _ := sweph.CalcUT(natalJD, sweph.SE_SUN)
	oldSolarArc := sweph.NormalizeDegrees(oldProgSun.Longitude - oldNatalSun.Longitude)
	
	fmt.Printf("\nOld method (JD_UT):\n")
	fmt.Printf("  Progressed JD: %.6f\n", oldProgJD)
	fmt.Printf("  Solar Arc: %.4f°\n", oldSolarArc)
	
	// 使用新方法（JDE）
	newProgJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
	newProgSun, _ := sweph.CalcUT(newProgJD, sweph.SE_SUN)
	newNatalSun, _ := sweph.CalcUT(sfJDE, sweph.SE_SUN) // 用SF的JDE作为natal
	newSolarArc := sweph.NormalizeDegrees(newProgSun.Longitude - newNatalSun.Longitude)
	
	fmt.Printf("\nNew method (JDE):\n")
	fmt.Printf("  Progressed JD: %.6f\n", newProgJD)
	fmt.Printf("  Solar Arc: %.4f°\n", newSolarArc)
	
	fmt.Printf("\nImprovement:\n")
	fmt.Printf("  Solar Arc diff: %+.4f°\n", newSolarArc-oldSolarArc)
	fmt.Printf("  Time equivalent: %+.1f hours\n", (newSolarArc-oldSolarArc)/15.0) // 15°/小时
	
	// 验证几个具体事件
	testEvents := []struct {
		name string
		awstTime string
		sfDiffSec float64
	}{
		{"Uranus Square Moon Tr-Sp", "2026-05-08 16:02:33", -26061},
		{"ASC Semi-Square NorthNode Sp-Na", "2026-04-29 17:18:12", -42424},
	}
	
	awstLoc, _ := time.LoadLocation("Australia/Perth")
	
	for _, evt := range testEvents {
		fmt.Printf("\n--- %s ---\n", evt.name)
		
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", evt.awstTime, awstLoc)
		utcTime := t.UTC()
		testTransitJD := julianDate(utcTime)
		
		// 旧方法
		oldProg := natalJD + (testTransitJD-natalJD)/365.25
		oldSun, _ := sweph.CalcUT(oldProg, sweph.SE_SUN)
		oldArc := sweph.NormalizeDegrees(oldSun.Longitude - oldNatalSun.Longitude)
		
		// 新方法
		newProg := progressions.SecondaryProgressionJD(natalJD, testTransitJD)
		newSun, _ := sweph.CalcUT(newProg, sweph.SE_SUN)
		newArc := sweph.NormalizeDegrees(newSun.Longitude - newNatalSun.Longitude)
		
		fmt.Printf("  Old Solar Arc: %.4f°\n", oldArc)
		fmt.Printf("  New Solar Arc: %.4f°\n", newArc)
		fmt.Printf("  Improvement: %+.4f° (%+.1f hours)\n", newArc-oldArc, (newArc-oldArc)/15.0)
		
		// 估算时间修正
		arcDiff := newArc - oldArc
		timeCorrectionHours := arcDiff / 15.0 // 15°/小时
		timeCorrectionSeconds := timeCorrectionHours * 3600
		fmt.Printf("  Estimated time correction: %+.0f seconds\n", timeCorrectionSeconds)
		fmt.Printf("  SF reported diff: %+.0f seconds\n", evt.sfDiffSec)
		fmt.Printf("  Remaining error: %+.0f seconds\n", evt.sfDiffSec-timeCorrectionSeconds)
	}
}

func julianDate(t time.Time) float64 {
	year, month, day := t.Date()
	hour := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
	return sweph.JulDay(year, int(month), day, hour, true)
}