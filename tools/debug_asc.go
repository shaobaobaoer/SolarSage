package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 调试ASC计算 ===")
	
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
	
	// 找到ASC Sp-Na事件
	for _, record := range records[1:] {
		if len(record) < 17 {
			continue
		}
		
		p1Name := record[0]
		eventType := record[6]
		
		if p1Name == "ASC" && eventType == "Sp-Na" {
			fmt.Printf("\n找到ASC Sp-Na事件:\n")
			fmt.Printf("CSV行: %v\n", record[:13])
			
			// 解析SF位置
			sfDeg, _ := strconv.ParseFloat(record[11], 64)
			sfSign := record[12]
			sfLon := signToDegree[sfSign] + sfDeg
			
			fmt.Printf("SF位置: %.3f° (%.3f° %s)\n", sfLon, sfDeg, sfSign)
			fmt.Printf("事件类型: %s\n", eventType)
			
			// 解析时间
			dateStr := record[7]
			timeStr := record[8]
			fmt.Printf("日期: %s %s\n", dateStr, timeStr)
			
			// 计算JD
			testJD := sweph.JulDay(2026, 1, 31, 16.0, true)
			fmt.Printf("事件JD: %.6f\n", testJD)
			
			// 检查事件类型是否包含"-Sp"
			fmt.Printf("包含'-Sp': %v\n", strings.Contains(eventType, "-Sp"))
			
			// 计算进展ASC
			geoLat, geoLon := 39.9042, 116.4074
			hs := models.HousePlacidus
			
			progASC, err := progressions.CalcProgressedSpecialPoint(models.PointASC, natalJD, testJD, geoLat, geoLon, hs)
			if err != nil {
				fmt.Printf("进展ASC错误: %v\n", err)
				continue
			}
			
			fmt.Printf("进展ASC: %.4f°\n", progASC)
			
			// 计算差异
			diff := progASC - sfLon
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
			
			break
		}
	}
}