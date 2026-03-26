package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 进展ASC验证 ===")
	
	natalJD := 2450800.900000
	eventJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	
	geoLat, geoLon := 39.9042, 116.4074
	hs := models.HousePlacidus
	
	// 本命ASC
	natalASC, err := chart.CalcSpecialPointLongitude(models.PointASC, geoLat, geoLon, natalJD, hs)
	if err != nil {
		fmt.Printf("本命ASC错误: %v\n", err)
		return
	}
	
	natalSign := int(natalASC / 30.0)
	natalDeg := natalASC - float64(natalSign)*30.0
	signNames := []string{"Ari", "Tau", "Gem", "Can", "Leo", "Vir",
		"Lib", "Sco", "Sag", "Cap", "Aqu", "Pis"}
	
	fmt.Printf("本命ASC: %.2f° (%s %.2f°)\n", natalASC, signNames[natalSign], natalDeg)
	
	// 进展ASC (使用CalcProgressedSpecialPoint)
	progASC, err := progressions.CalcProgressedSpecialPoint(models.PointASC, natalJD, eventJD, geoLat, geoLon, hs)
	if err != nil {
		fmt.Printf("进展ASC错误: %v\n", err)
		return
	}
	
	progSign := int(progASC / 30.0)
	progDeg := progASC - float64(progSign)*30.0
	
	fmt.Printf("进展ASC (CalcProgressedSpecialPoint): %.2f° (%s %.2f°)\n", 
		progASC, signNames[progSign], progDeg)
	
	// 太阳弧偏移
	offset, err := progressions.SolarArcOffset(natalJD, eventJD)
	if err != nil {
		fmt.Printf("太阳弧偏移错误: %v\n", err)
		return
	}
	
	fmt.Printf("太阳弧偏移: %.2f°\n", offset)
	
	// 太阳弧ASC
	saASC := sweph.NormalizeDegrees(natalASC + offset)
	saSign := int(saASC / 30.0)
	saDeg := saASC - float64(saSign)*30.0
	
	fmt.Printf("太阳弧ASC: %.2f° (%s %.2f°)\n", saASC, signNames[saSign], saDeg)
	
	// 对比Solar Fire
	// SF显示: 29.25° Cancer = 90 + 29.25 = 119.25°
	sfASC := 90.0 + 29.25
	fmt.Printf("\nSolar Fire ASC: %.2f° (29.25° Cancer)\n", sfASC)
	
	fmt.Printf("与本命ASC差异: %.2f°\n", natalASC-sfASC)
	fmt.Printf("与进展ASC差异: %.2f°\n", progASC-sfASC)
	fmt.Printf("与太阳弧ASC差异: %.2f°\n", saASC-sfASC)
}