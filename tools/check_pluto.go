package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 冥王星位置详细检查 ===")
	
	// 2026-02-01 00:00:00 AWST = 2026-01-31 16:00:00 UTC
	testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
	
	fmt.Printf("测试JD: %.6f\n", testJD)
	
	// 计算冥王星
	result, err := sweph.CalcUT(testJD, 9) // 9 = SE_PLUTO
	if err != nil {
		fmt.Printf("错误 %v\n", err)
		return
	}
	
	sign := int(result.Longitude / 30.0)
	degInSign := result.Longitude - float64(sign)*30.0
	
	signNames := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	
	fmt.Printf("冥王星位置:\n")
	fmt.Printf("  经度: %.4f° (%s %.2f°)\n", result.Longitude, signNames[sign], degInSign)
	fmt.Printf("  纬度: %.4f°\n", result.Latitude)
	fmt.Printf("  距离: %.6f AU\n", result.Distance)
	
	// 检查2026年冥王星的实际位置变化
	fmt.Println("\n=== 冥王星2026年位置变化 ===")
	for month := 1; month <= 12; month++ {
		jd := sweph.JulDay(2026, month, 15, 12.0, true)
		result, _ := sweph.CalcUT(jd, 9)
		
		sign := int(result.Longitude / 30.0)
		degInSign := result.Longitude - float64(sign)*30.0
		signNames := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
			"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
		
		fmt.Printf("2026-%02d-15: %.2f° (%s %.2f°)\n", month, result.Longitude, signNames[sign], degInSign)
	}
}