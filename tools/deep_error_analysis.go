package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	defer sweph.Close()

	const natalJD = 2450800.900009
	const natalLat = 30.9
	const natalLon = 121.15

	fmt.Println("=== Detailed Progression Error Analysis ===")

	// 分析一个典型的大误差事件
	// ASC Semi-Square NorthNode Sp-Na: SF diff -42424s (-707min)
	awstTime := "2026-04-29 17:18:12"
	awstLoc, _ := time.LoadLocation("Australia/Perth")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", awstTime, awstLoc)
	utcTime := t.UTC()
	transitJD := julianDate(utcTime)

	fmt.Printf("Analyzing event at:\n")
	fmt.Printf("  AWST: %s\n", awstTime)
	fmt.Printf("  UTC:  %s\n", utcTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  JD:   %.6f\n", transitJD)
	fmt.Printf("  SF reported error: -42424s (-707min)\n")

	// 1. 当前我们的计算结果
	currentProgJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
	currentASC, _ := progressions.CalcProgressedSpecialPoint(models.PointASC, natalJD, transitJD, natalLat, natalLon, models.HousePlacidus)
	currentNorthNode, _ := progressions.CalcProgressedSpecialPoint(models.PointNorthNode, natalJD, transitJD, natalLat, natalLon, models.HousePlacidus)
	
	diff := currentASC - currentNorthNode
	if diff < 0 { diff += 360 }
	if diff > 180 { diff = 360 - diff }
	
	fmt.Printf("\nCurrent calculation:\n")
	fmt.Printf("  Progressed JD: %.6f\n", currentProgJD)
	fmt.Printf("  Progressed ASC: %.4f°\n", currentASC)
	fmt.Printf("  Progressed NorthNode: %.4f°\n", currentNorthNode)
	fmt.Printf("  Angular separation: %.4f°\n", diff)
	fmt.Printf("  Aspect orb: %.4f°\n", diff-45) // Semi-Square = 45°

	// 2. SF期望的结果（反向推导）
	// SF的时间比我们早707分钟 = 11.78小时
	// 对应的角度差：11.78 * 15° = 176.7°
	// 但这不合理，因为Semi-Square只需要45°分离
	
	// 更合理的解释：SF使用了不同的progressed JD
	targetTimeDiff := -42424.0 // 秒
	targetJDOffset := targetTimeDiff / 86400.0 // 天
	targetProgJD := currentProgJD + targetJDOffset
	
	targetASC, _ := calcASCAtJD(targetProgJD, natalLat, natalLon)
	targetNorthNode, _ := sweph.CalcUT(targetProgJD, sweph.SE_MEAN_NODE)
	
	targetDiff := targetASC - targetNorthNode
	if targetDiff < 0 { targetDiff += 360 }
	if targetDiff > 180 { targetDiff = 360 - targetDiff }
	
	fmt.Printf("\nTarget calculation (to match SF):\n")
	fmt.Printf("  Target JD offset: %+.4fd (%+.1fh)\n", targetJDOffset, targetJDOffset*24)
	fmt.Printf("  Target Progressed JD: %.6f\n", targetProgJD)
	fmt.Printf("  Target ASC: %.4f°\n", targetASC)
	fmt.Printf("  Target NorthNode: %.4f°\n", targetNorthNode)
	fmt.Printf("  Target angular separation: %.4f°\n", targetDiff)
	fmt.Printf("  Target orb: %.4f°\n", targetDiff-45)

	// 3. 分析差异来源
	fmt.Printf("\n=== Error Source Analysis ===\n")
	
	// 检查是否是natal epoch问题
	natalOffsets := []float64{-0.5, -0.25, -0.1, -0.05, 0, 0.05, 0.1, 0.25, 0.5}
	fmt.Println("Testing different natal epoch offsets:")
	
	bestOffset := 0.0
	bestDiff := 1e9
	
	for _, offset := range natalOffsets {
		testNatalJD := natalJD + offset
		testProgJD := testNatalJD + (transitJD-testNatalJD)/365.25
		testASC, _ := calcASCAtJD(testProgJD, natalLat, natalLon)
		testNorthNode, _ := sweph.CalcUT(testProgJD, sweph.SE_MEAN_NODE)
		
		testDiff := testASC - testNorthNode
		if testDiff < 0 { testDiff += 360 }
		if testDiff > 180 { testDiff = 360 - testDiff }
		orb := testDiff - 45
		
		improvement := abs(orb) - abs(targetDiff-45)
		
		status := ""
		if abs(improvement) < 0.1 {
			status = "★"
			if abs(orb) < bestDiff {
				bestDiff = abs(orb)
				bestOffset = offset
			}
		}
		
		fmt.Printf("  %+5.2fd: orb=%+7.4f° improvement=%+7.4f° %s\n", 
			offset, orb, improvement, status)
	}
	
	fmt.Printf("Best offset: %+5.2fd\n", bestOffset)

	// 4. 检查是否是house系统或计算方法问题
	fmt.Println("\nTesting different calculation approaches:")
	
	// 方法A: 当前方法 (Solar Arc in RA)
	ascA, mcA, _ := progressions.CalcProgressedAngles(natalJD, transitJD, natalLat, natalLon, models.HousePlacidus)
	fmt.Printf("  Method A (Solar Arc RA): ASC=%.4f° MC=%.4f°\n", ascA, mcA)
	
	// 方法B: 直接progress MC然后推导ASC
	// 这是一些传统做法
	fmt.Println("  Method B (Direct MC progression): [需要实现]")
	
	// 方法C: 使用不同的house系统
	ascC, mcC, _ := progressions.CalcProgressedAngles(natalJD, transitJD, natalLat, natalLon, models.HouseKoch)
	fmt.Printf("  Method C (Koch houses): ASC=%.4f° MC=%.4f°\n", ascC, mcC)
}

func calcASCAtJD(jd float64, lat, lon float64) (float64, error) {
	hsysChar := models.HouseSystemToChar(models.HousePlacidus)
	hr, err := sweph.Houses(jd, lat, lon, hsysChar)
	if err != nil {
		return 0, err
	}
	return hr.ASC, nil
}

func julianDate(t time.Time) float64 {
	year, month, day := t.Date()
	hour := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
	return sweph.JulDay(year, int(month), day, hour, true)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}