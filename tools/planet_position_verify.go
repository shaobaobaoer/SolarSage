package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// 行星名称映射
var planetMap = map[string]models.PlanetID{
	"Sun":   models.PlanetSun,
	"Moon":  models.PlanetMoon,
	"Mercury": models.PlanetMercury,
	"Venus": models.PlanetVenus,
	"Mars":  models.PlanetMars,
	"Jupiter": models.PlanetJupiter,
	"Saturn": models.PlanetSaturn,
	"Uranus": models.PlanetUranus,
	"Neptune": models.PlanetNeptune,
	"Pluto": models.PlanetPluto,
	"Chiron": models.PlanetChiron,
	"NorthNode": models.PlanetNorthNodeMean,
}

// 星座到度数的映射
var signToDegree = map[string]float64{
	"Aries": 0, "Taurus": 30, "Gemini": 60, "Cancer": 90,
	"Leo": 120, "Virgo": 150, "Libra": 180, "Scorpio": 210,
	"Sagittarius": 240, "Capricorn": 270, "Aquarius": 300, "Pisces": 330,
}

func main() {
	fmt.Println("=== 行星位置精确验证 ===")
	
	// 出生数据
	natalJD := 2450800.900000 // 1997-12-18 09:36:00 UTC
	fmt.Printf("出生JD: %.6f\n\n", natalJD)
	
	// 读取Solar Fire数据
	file, err := os.Open("testdata/solarfire/testcase-1-transit.csv")
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("读取CSV错误: %v\n", err)
		return
	}
	
	// 跳过标题行
	if len(records) < 2 {
		fmt.Println("CSV文件为空")
		return
	}
	
	// 统计信息
	var totalDiff, maxDiff float64
	var count int
	var maxDiffEvent string
	
	fmt.Println("逐行星位置对比:")
	fmt.Printf("%-25s %-10s %-8s %-15s %-15s %-10s %-10s\n", 
		"事件", "行星", "类型", "SF位置", "计算位置", "差异(°)", "状态")
	fmt.Println(strings.Repeat("-", 100))
	
	// 处理每个记录
	for i, record := range records[1:] {
		if len(record) < 17 {
			continue
		}
		
		p1Name := record[0]
		p2Name := record[3]
		eventType := record[6] // Tr-Na, Tr-Sp, etc.
		dateStr := record[7]
		timeStr := record[8]
		
		// 解析时间
		t, err := time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
		if err != nil {
			continue
		}
		
		// 转换为JD
		utcTime := t.UTC()
		transitJD := sweph.JulDay(utcTime.Year(), int(utcTime.Month()), utcTime.Day(),
			float64(utcTime.Hour())+float64(utcTime.Minute())/60.0+float64(utcTime.Second())/3600.0, true)
		
		// 获取SF报告的位置
		// CSV格式: ..., Age, Pos1_Deg, Pos1_Sign, Pos1_Dir, Pos2_Deg, Pos2_Sign, Pos2_Dir
		// Index:    10   11        12         13        14        15         16
		sfPos1Deg, _ := strconv.ParseFloat(record[11], 64)
		sfPos1Sign := record[12]
		sfLon1 := signToDegree[sfPos1Sign] + sfPos1Deg
		
		sfPos2Deg, _ := strconv.ParseFloat(record[14], 64)
		sfPos2Sign := record[15]
		sfLon2 := signToDegree[sfPos2Sign] + sfPos2Deg
		
		// 根据事件类型计算我们的位置
		var ourLon1, ourLon2 float64
		var calcErr error
		
		// 解析P1
		if p1, ok := planetMap[p1Name]; ok {
			ourLon1, calcErr = calcPosition(p1, eventType, natalJD, transitJD, true)
			if calcErr != nil {
				continue
			}
		} else if p1Name == "ASC" {
			ourLon1, calcErr = calcAngle(models.PointASC, natalJD, transitJD, eventType)
			if calcErr != nil {
				continue
			}
		} else if p1Name == "MC" {
			ourLon1, calcErr = calcAngle(models.PointMC, natalJD, transitJD, eventType)
			if calcErr != nil {
				continue
			}
		}
		
		// 解析P2
		if p2, ok := planetMap[p2Name]; ok {
			ourLon2, calcErr = calcPosition(p2, eventType, natalJD, transitJD, false)
			if calcErr != nil {
				continue
			}
		} else if p2Name == "ASC" {
			ourLon2, calcErr = calcAngle(models.PointASC, natalJD, transitJD, eventType)
			if calcErr != nil {
				continue
			}
		} else if p2Name == "MC" {
			ourLon2, calcErr = calcAngle(models.PointMC, natalJD, transitJD, eventType)
			if calcErr != nil {
				continue
			}
		}
		
		// 计算差异
		diff1 := normalizeDiff(ourLon1 - sfLon1)
		diff2 := normalizeDiff(ourLon2 - sfLon2)
		
		avgDiff := (abs(diff1) + abs(diff2)) / 2
		
		status := "✓"
		if avgDiff > 0.5 {
			status = "⚠"
		}
		if avgDiff > 1.0 {
			status = "✗"
		}
		
		if i < 20 || p1Name == "Moon" || p2Name == "Moon" { // 显示前20个或Moon事件
			// 计算我们的星座位置
			ourSign1 := int(ourLon1 / 30.0)
			_ = ourLon1 - float64(ourSign1)*30.0 // ourDeg1 not used
			signNames := []string{"Ari", "Tau", "Gem", "Can", "Leo", "Vir",
				"Lib", "Sco", "Sag", "Cap", "Aqu", "Pis"}
			
			fmt.Printf("%-25s %-10s %-8s %6.2f° %-8s %6.2f° %-8s %8.3f   %s\n",
				record[4]+" "+record[2]+" "+record[5],
				p1Name,
				eventType,
				sfLon1, sfPos1Sign[:3],
				ourLon1, signNames[ourSign1],
				avgDiff,
				status)
		}
		
		totalDiff += avgDiff
		count++
		
		if avgDiff > maxDiff {
			maxDiff = avgDiff
			maxDiffEvent = p1Name + " " + record[2] + " " + p2Name
		}
	}
	
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("\n统计摘要:\n")
	fmt.Printf("总事件数: %d\n", count)
	fmt.Printf("平均位置差异: %.3f°\n", totalDiff/float64(count))
	fmt.Printf("最大位置差异: %.3f° (%s)\n", maxDiff, maxDiffEvent)
	
	if maxDiff < 0.1 {
		fmt.Println("✓ 位置精度优秀 (< 0.1°)")
	} else if maxDiff < 0.5 {
		fmt.Println("✓ 位置精度良好 (< 0.5°)")
	} else if maxDiff < 1.0 {
		fmt.Println("⚠ 位置精度一般 (< 1.0°)")
	} else {
		fmt.Println("✗ 位置精度较差 (> 1.0°)")
	}
}

func calcPosition(planet models.PlanetID, eventType string, natalJD, transitJD float64, isP1 bool) (float64, error) {
	// 根据事件类型决定如何计算位置
	if strings.HasPrefix(eventType, "Tr-") {
		if isP1 {
			// 过境行星 - 使用当前时间
			lon, _, err := chart.CalcPlanetLongitude(planet, transitJD)
			return lon, err
		} else {
			// 本命/进展/太阳弧行星
			if strings.Contains(eventType, "-Na") {
				// 本命位置
				lon, _, err := chart.CalcPlanetLongitude(planet, natalJD)
				return lon, err
			} else if strings.Contains(eventType, "-Sp") {
				// 进展位置
				lon, _, err := progressions.CalcProgressedLongitude(planet, natalJD, transitJD)
				return lon, err
			} else if strings.Contains(eventType, "-Sa") {
				// 太阳弧位置
				lon, _, err := progressions.CalcSolarArcLongitude(planet, natalJD, transitJD)
				return lon, err
			}
		}
	} else if strings.HasPrefix(eventType, "Sp-") || strings.HasPrefix(eventType, "Sa-") {
		// 进展或太阳弧事件
		if strings.HasPrefix(eventType, "Sp-") {
			if isP1 {
				lon, _, err := progressions.CalcProgressedLongitude(planet, natalJD, transitJD)
				return lon, err
			} else {
				if strings.HasSuffix(eventType, "-Na") {
					lon, _, err := chart.CalcPlanetLongitude(planet, natalJD)
					return lon, err
				} else if strings.HasSuffix(eventType, "-Sp") {
					lon, _, err := progressions.CalcProgressedLongitude(planet, natalJD, transitJD)
					return lon, err
				}
			}
		} else if strings.HasPrefix(eventType, "Sa-") {
			if isP1 {
				lon, _, err := progressions.CalcSolarArcLongitude(planet, natalJD, transitJD)
				return lon, err
			} else {
				lon, _, err := chart.CalcPlanetLongitude(planet, natalJD)
				return lon, err
			}
		}
	}
	
	// 默认使用当前时间
	lon, _, err := chart.CalcPlanetLongitude(planet, transitJD)
	return lon, err
}

func calcAngle(angle models.SpecialPointID, natalJD, transitJD float64, eventType string) (float64, error) {
	// 简化处理：使用本命盘的角度
	// 实际应该根据事件类型计算进展/太阳弧角度
	hs := models.HousePlacidus
	geoLat, geoLon := 39.9042, 116.4074 // 北京坐标
	
	// 检查事件类型 (Sp-Na, Tr-Sp, etc.)
	if strings.HasPrefix(eventType, "Sp-") || strings.Contains(eventType, "-Sp-") || strings.HasSuffix(eventType, "-Sp") {
		return progressions.CalcProgressedSpecialPoint(angle, natalJD, transitJD, geoLat, geoLon, hs)
	} else if strings.HasPrefix(eventType, "Sa-") || strings.Contains(eventType, "-Sa-") || strings.HasSuffix(eventType, "-Sa") {
		// 太阳弧角度 = 本命角度 + 太阳弧偏移
		offset, err := progressions.SolarArcOffset(natalJD, transitJD)
		if err != nil {
			return 0, err
		}
		natalAngle, err := chart.CalcSpecialPointLongitude(angle, geoLat, geoLon, natalJD, hs)
		if err != nil {
			return 0, err
		}
		return sweph.NormalizeDegrees(natalAngle + offset), nil
	}
	
	// 默认本命角度
	return chart.CalcSpecialPointLongitude(angle, geoLat, geoLon, natalJD, hs)
}

func normalizeDiff(diff float64) float64 {
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	return diff
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}