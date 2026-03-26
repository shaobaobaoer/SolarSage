package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 土星位置验证 ===")
	
	// 出生JD: 1997-12-18 09:36:00 UTC
	natalJD := 2450800.900000
	
	// 2026-02-01 00:00:00 AWST = 2026-01-31 16:00:00 UTC
	testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	
	fmt.Printf("出生JD: %.6f\n", natalJD)
	fmt.Printf("测试JD: %.6f (2026-01-31 16:00 UTC)\n", testJD)
	
	// 计算土星位置
	fmt.Println("\n土星位置:")
	
	// 本命位置
	natalSaturn, _, err := chart.CalcPlanetLongitude(models.PlanetSaturn, natalJD)
	if err != nil {
		fmt.Printf("本命计算错误: %v\n", err)
		return
	}
	fmt.Printf("  本命: %.4f° (%.2f° Aries)\n", natalSaturn, natalSaturn)
	
	// 过境位置
	transitSaturn, _, err := chart.CalcPlanetLongitude(models.PlanetSaturn, testJD)
	if err != nil {
		fmt.Printf("过境计算错误: %v\n", err)
		return
	}
	fmt.Printf("  过境: %.4f° (%.2f° Pisces)\n", transitSaturn, transitSaturn)
	
	// 转换为星座度数
	sign := int(transitSaturn / 30.0)
	degInSign := transitSaturn - float64(sign)*30.0
	signNames := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	fmt.Printf("  过境位置: %.2f° %s\n", degInSign, signNames[sign])
	
	// 对比Solar Fire数据
	// SF显示: Saturn at 28.6° Pisces = 330 + 28.6 = 358.6°
	sfSaturn := 330.0 + 28.6
	fmt.Printf("\nSolar Fire土星位置: %.2f° (28.6° Pisces)\n", sfSaturn)
	fmt.Printf("计算差异: %.4f°\n", transitSaturn-sfSaturn)
	
	// 显示土星详细信息
	fmt.Printf("\n土星详细信息:\n")
	fmt.Printf("  计算经度: %.4f°\n", transitSaturn)
	saturnSign := int(transitSaturn / 30.0)
	saturnDeg := transitSaturn - float64(saturnSign)*30.0
	fmt.Printf("  星座位置: %s %.2f°\n", signNames[saturnSign], saturnDeg)
	
	// 检查其他行星
	fmt.Println("\n=== 其他行星位置对比 ===")
	planets := []struct {
		name   string
		planet models.PlanetID
		sfDeg  float64
		sfSign string
	}{
		{"Uranus", models.PlanetUranus, 27.45, "Taurus"},
		{"Neptune", models.PlanetNeptune, 0.117, "Aries"},
		{"Pluto", models.PlanetPluto, 7.217, "Sagittarius"},
		{"Chiron", models.PlanetChiron, 22.983, "Aries"},
	}
	
	signDegrees := map[string]float64{
		"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
		"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
		"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
	}
	
	for _, p := range planets {
		lon, _, err := chart.CalcPlanetLongitude(p.planet, testJD)
		if err != nil {
			continue
		}
		
		sfLon := signDegrees[p.sfSign] + p.sfDeg
		diff := lon - sfLon
		
		fmt.Printf("%-10s: 计算=%.2f°, SF=%.2f° (%.2f° %s), 差异=%.3f°\n",
			p.name, lon, sfLon, p.sfDeg, p.sfSign, diff)
	}
}