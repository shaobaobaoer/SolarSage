package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 调试验证 ===")
	
	natalJD := 2450800.900000
	
	// 读取CSV
	file, _ := os.Open("testdata/solarfire/testcase-1-transit.csv")
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	
	signToDegree := map[string]float64{
		"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
		"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
		"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
	}
	
	// 找到Pluto Sp-Na事件
	for _, record := range records[1:] {
		if len(record) < 17 {
			continue
		}
		
		p1Name := record[0]
		eventType := record[6]
		
		if p1Name == "Pluto" && eventType == "Sp-Na" {
			fmt.Printf("\n找到Pluto Sp-Na事件:\n")
			fmt.Printf("CSV行: %v\n", record[:12])
			
			// 解析SF位置 (根据CSV格式: Age, Pos1_Deg, Pos1_Sign, Pos1_Dir, Pos2_Deg, Pos2_Sign, Pos2_Dir)
			// Index: 10=Age, 11=Pos1_Deg, 12=Pos1_Sign, 13=Pos1_Dir, 14=Pos2_Deg, 15=Pos2_Sign, 16=Pos2_Dir
			sfDeg, _ := strconv.ParseFloat(record[11], 64)
			sfSign := record[12]
			sfLon := signToDegree[sfSign] + sfDeg
			
			fmt.Printf("SF位置: %.3f° (%.3f° %s)\n", sfLon, sfDeg, sfSign)
			
			// 解析时间
			dateStr := record[7]
			timeStr := record[8]
			fmt.Printf("日期: %s %s\n", dateStr, timeStr)
			
			// 计算JD (AWST = UTC+8, 所以00:00:00 AWST = 16:00:00 UTC前一天)
			testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
			fmt.Printf("事件JD: %.6f\n", testJD)
			
			// 计算我们的位置
			ourLon, _, err := progressions.CalcProgressedLongitude(models.PlanetPluto, natalJD, testJD)
			if err != nil {
				fmt.Printf("计算错误: %v\n", err)
				continue
			}
			
			fmt.Printf("我们的位置: %.4f°\n", ourLon)
			
			// 计算差异
			diff := ourLon - sfLon
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
			
			// 直接验证
			progJD := progressions.SecondaryProgressionJD(natalJD, testJD)
			fmt.Printf("进展JD: %.6f\n", progJD)
			
			directPluto, _, _ := chart.CalcPlanetLongitude(models.PlanetPluto, progJD)
			fmt.Printf("直接计算Pluto: %.4f°\n", directPluto)
			
			break
		}
	}
}