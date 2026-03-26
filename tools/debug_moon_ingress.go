package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 调试Moon SignIngress ===")
	
	file, _ := os.Open("testdata/solarfire/testcase-1-transit.csv")
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	
	signToDegree := map[string]float64{
		"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
		"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
		"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
	}
	
	// 找到一个Moon Conjunction SignIngress事件
	for _, record := range records[1:] {
		if len(record) < 17 {
			continue
		}
		
		p1Name := record[0]
		eventType := record[5] // EventType: SignIngress, Exact, etc.
		aspect := record[2]
		
		if p1Name == "Moon" && eventType == "SignIngress" && aspect == "Conjunction" {
			fmt.Printf("\n找到Moon Conjunction SignIngress事件:\n")
			fmt.Printf("CSV行: %v\n", record[:13])
			
			// 解析SF位置
			sfDeg, _ := strconv.ParseFloat(record[11], 64)
			sfSign := record[12]
			sfLon := signToDegree[sfSign] + sfDeg
			
			fmt.Printf("SF位置: %.3f° (%.3f° %s)\n", sfLon, sfDeg, sfSign)
			
			// 解析时间
			dateStr := record[7]
			timeStr := record[8]
			fmt.Printf("日期: %s %s\n", dateStr, timeStr)
			
			// 解析准确时间: 2026-02-01 08:08:52 AWST
			// AWST = UTC+8, 所以 08:08:52 AWST = 00:08:52 UTC
			testJD := sweph.JulDay(2026, 2, 1, 0.0 + 8.0/60.0 + 52.0/3600.0, true)
			fmt.Printf("事件JD: %.6f\n", testJD)
			
			// 计算月亮位置
			moonLon, _, err := chart.CalcPlanetLongitude(models.PlanetMoon, testJD)
			if err != nil {
				fmt.Printf("月亮计算错误: %v\n", err)
				continue
			}
			
			fmt.Printf("计算月亮位置: %.4f°\n", moonLon)
			
			// 计算差异
			diff := moonLon - sfLon
			for diff > 180 {
				diff -= 360
			}
			for diff < -180 {
				diff += 360
			}
			if diff < 0 {
				diff = -diff
			}
			
			fmt.Printf("差异: %.4f°\n", diff)
			
			// 检查星座边界
			signIndex := int(moonLon / 30.0)
			signStart := float64(signIndex) * 30.0
			signNames := []string{"Ari", "Tau", "Gem", "Can", "Leo", "Vir",
				"Lib", "Sco", "Sag", "Cap", "Aqu", "Pis"}
			
			fmt.Printf("计算星座: %s (%.2f° - %.2f°)\n", 
				signNames[signIndex], signStart, signStart+30)
			fmt.Printf("SF期望星座: %s\n", sfSign)
			
			break
		}
	}
}