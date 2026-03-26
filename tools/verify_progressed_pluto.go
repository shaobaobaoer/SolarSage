package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 进展冥王星位置验证 ===")
	
	// 出生: 1997-12-18 09:36:00 UTC
	natalJD := 2450800.900000
	
	// 事件时间: 2026-02-01 00:00:00 AWST = 2026-01-31 16:00:00 UTC
	eventJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	
	fmt.Printf("出生JD: %.6f\n", natalJD)
	fmt.Printf("事件JD: %.6f (2026-01-31 16:00 UTC)\n", eventJD)
	fmt.Printf("时间跨度: %.1f天 = %.2f年\n", eventJD-natalJD, (eventJD-natalJD)/365.25)
	
	// 计算本命冥王星
	natalPluto, _, err := chart.CalcPlanetLongitude(models.PlanetPluto, natalJD)
	if err != nil {
		fmt.Printf("本命计算错误: %v\n", err)
		return
	}
	
	sign := int(natalPluto / 30.0)
	degInSign := natalPluto - float64(sign)*30.0
	signNames := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	
	fmt.Printf("\n本命冥王星: %.4f° (%s %.2f°)\n", natalPluto, signNames[sign], degInSign)
	
	// 计算进展冥王星
	progPluto, _, err := progressions.CalcProgressedLongitude(models.PlanetPluto, natalJD, eventJD)
	if err != nil {
		fmt.Printf("进展计算错误: %v\n", err)
		return
	}
	
	progSign := int(progPluto / 30.0)
	progDegInSign := progPluto - float64(progSign)*30.0
	
	fmt.Printf("进展冥王星: %.4f° (%s %.2f°)\n", progPluto, signNames[progSign], progDegInSign)
	
	// 计算进展移动
	movement := progPluto - natalPluto
	for movement < 0 {
		movement += 360
	}
	fmt.Printf("进展移动: %.2f°\n", movement)
	
	// 对比Solar Fire
	// SF显示: 7.217° Sagittarius = 240 + 7.217 = 247.217°
	sfProgPluto := 240.0 + 7.217
	fmt.Printf("\nSolar Fire进展冥王星: %.3f° (7.217° Sagittarius)\n", sfProgPluto)
	fmt.Printf("计算差异: %.3f°\n", progPluto-sfProgPluto)
	
	// 进展JD
	progJD := progressions.SecondaryProgressionJD(natalJD, eventJD)
	fmt.Printf("\n进展JD: %.6f\n", progJD)
	fmt.Printf("进展日期: 出生+%.1f天\n", progJD-natalJD)
	
	// 验证进展JD的冥王星位置
	progJDPluto, _, _ := chart.CalcPlanetLongitude(models.PlanetPluto, progJD)
	progJDSign := int(progJDPluto / 30.0)
	progJDDegInSign := progJDPluto - float64(progJDSign)*30.0
	
	fmt.Printf("进展JD冥王星位置: %.4f° (%s %.2f°)\n", 
		progJDPluto, signNames[progJDSign], progJDDegInSign)
}