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

	fmt.Println("=== Core Progression Error Analysis ===")

	// 分析典型事件: ASC Semi-Square NorthNode Sp-Na
	// SF时间: 2026-04-29 17:18:12 AWST
	// 我们时间: 比SF晚707分钟
	awstTime := "2026-04-29 17:18:12"
	awstLoc, _ := time.LoadLocation("Australia/Perth")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", awstTime, awstLoc)
	utcTime := t.UTC()
	transitJD := julianDate(utcTime)

	fmt.Printf("Event: ASC Semi-Square NorthNode Sp-Na\n")
	fmt.Printf("SF time (AWST): %s\n", awstTime)
	fmt.Printf("UTC time: %s\n", utcTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Transit JD: %.6f\n", transitJD)
	fmt.Printf("SF error: -42424s (-707min)\n")

	// 当前计算
	currentProgJD := progressions.SecondaryProgressionJD(natalJD, transitJD)
	currentASC, _ := progressions.CalcProgressedSpecialPoint(models.PointASC, natalJD, transitJD, natalLat, natalLon, models.HousePlacidus)
	northNode, _ := sweph.CalcUT(transitJD, sweph.SE_MEAN_NODE)
	
	angularSep := angleSeparation(currentASC, northNode.Longitude)
	orb := angularSep - 45 // Semi-Square
	
	fmt.Printf("\nCurrent calculation:\n")
	fmt.Printf("  Progressed JD: %.6f\n", currentProgJD)
	fmt.Printf("  Progressed ASC: %.4f°\n", currentASC)
	fmt.Printf("  Transit NorthNode: %.4f°\n", northNode.Longitude)
	fmt.Printf("  Angular separation: %.4f°\n", angularSep)
	fmt.Printf("  Aspect orb: %.4f°\n", orb)

	// 要匹配SF，需要什么条件？
	// SF的时间比我们早707分钟 = 11.78小时 = 176.7° (15°/小时)
	// 但这不合理，因为只需要45°分离
	
	// 更可能是progressed JD不同
	targetTimeDiff := -42424.0 // 秒
	targetJDOffset := targetTimeDiff / 86400.0
	targetProgJD := currentProgJD + targetJDOffset
	
	targetASC, _ := calcASCAtJD(targetProgJD, natalLat, natalLon)
	targetNode, _ := sweph.CalcUT(targetProgJD, sweph.SE_MEAN_NODE)
	targetSep := angleSeparation(targetASC, targetNode.Longitude)
	targetOrb := targetSep - 45
	
	fmt.Printf("\nTarget calculation:\n")
	fmt.Printf("  JD offset needed: %+.4fd (%+.1fh)\n", targetJDOffset, targetJDOffset*24)
	fmt.Printf("  Target Progressed JD: %.6f\n", targetProgJD)
	fmt.Printf("  Target ASC: %.4f°\n", targetASC)
	fmt.Printf("  Target NorthNode: %.4f°\n", targetNode.Longitude)
	fmt.Printf("  Target separation: %.4f°\n", targetSep)
	fmt.Printf("  Target orb: %.4f°\n", targetOrb)

	// 测试不同的natal epoch偏移
	fmt.Printf("\n=== Testing Natal Epoch Offsets ===\n")
	offsets := []float64{-0.3, -0.2, -0.1, -0.05, 0, 0.05, 0.1, 0.2, 0.3}
	
	for _, offset := range offsets {
		testNatalJD := natalJD + offset
		// 使用修正后的公式
		testProgJD := testNatalJD + (transitJD-testNatalJD)/365.25
		testASC, _ := calcASCAtJD(testProgJD, natalLat, natalLon)
		testNode, _ := sweph.CalcUT(testProgJD, sweph.SE_MEAN_NODE)
		testSep := angleSeparation(testASC, testNode.Longitude)
		testOrb := testSep - 45
		
		improvement := abs(testOrb) - abs(targetOrb)
		status := ""
		if abs(improvement) < 0.1 {
			status = "★"
		}
		
		fmt.Printf("  %+4.1fd: orb=%+7.4f° improvement=%+7.4f° %s\n", 
			offset, testOrb, improvement, status)
	}

	// 检查是否是house系统问题
	fmt.Printf("\n=== Testing House Systems ===\n")
	houseSystems := []models.HouseSystem{
		models.HousePlacidus,
		models.HouseKoch,
		models.HouseEqual,
		models.HouseWholeSign,
	}
	
	for _, hs := range houseSystems {
		asc, _, _ := progressions.CalcProgressedAngles(natalJD, transitJD, natalLat, natalLon, hs)
		sep := angleSeparation(asc, northNode.Longitude)
		orb := sep - 45
		fmt.Printf("  %-12s: ASC=%.4f° orb=%+7.4f°\n", hs, asc, orb)
	}
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

func angleSeparation(a, b float64) float64 {
	diff := a - b
	if diff < 0 { diff += 360 }
	if diff > 180 { diff = 360 - diff }
	return diff
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}