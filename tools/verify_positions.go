package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// PositionCheck 验证特定时刻的行星位置
type PositionCheck struct {
	DateTime    string  // YYYY-MM-DD HH:MM:SS
	JD          float64
	Description string // 描述这个检查点的意义
}

// PlanetPosition 行星位置数据
type PlanetPosition struct {
	Planet    models.PlanetID
	Longitude float64
	Latitude  float64
	Speed     float64
	IsRetro   bool
}

// PositionComparison 位置对比结果
type PositionComparison struct {
	DateTime      string
	Planet        models.PlanetID
	SF_Lon        float64
	SS_Lon        float64
	Diff_ArcSec   float64
	Diff_Degrees  float64
	Within_Tolerance bool
}

func main() {
	// 初始化星历
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	sweph.Init(ephePath)
	defer sweph.Close()

	// 测试案例1参数
	natalJD := 2450800.900009 // 1997-12-18 09:36 UTC
	natalLat := 30.9
	natalLon := 121.15

	fmt.Println("=== 验证行星位置计算准确性 ===")
	fmt.Printf("出生时间: JD=%.6f (1997-12-18 09:36 UTC)\n", natalJD)
	fmt.Printf("出生地点: %.2f°N, %.2f°E\n", natalLat, natalLon)

	// 从Solar Fire CSV中提取一些关键时间点进行验证
	checkPoints := extractCheckPointsFromSF()
	
	fmt.Printf("\n找到 %d 个检查时间点\n", len(checkPoints))
	
	// 验证每个时间点的行星位置
	var allComparisons []PositionComparison
	
	for i, cp := range checkPoints {
		fmt.Printf("\n[%d/%d] 检查时间点: %s (JD=%.6f)\n", i+1, len(checkPoints), cp.DateTime, cp.JD)
		
		// 分析这个时间点涉及的事件类型
		eventTypes := getEventTypesAtTime(cp.DateTime)
		fmt.Printf("  相关事件类型: %s\n", strings.Join(eventTypes, ", "))
		
		// 验证不同图表类型的行星位置
		comparisons := verifyPositionsAtTime(cp, natalJD, natalLat, natalLon)
		allComparisons = append(allComparisons, comparisons...)
	}
	
	// 统计分析
	analyzeResults(allComparisons)
}

func extractCheckPointsFromSF() []PositionCheck {
	// 读取Solar Fire CSV，提取一些代表性的时间点
	// 特别关注那些时间差异较大的事件
	f, err := os.Open("testdata/solarfire/testcase-1-transit.csv")
	if err != nil {
		fmt.Printf("无法打开SF CSV: %v\n", err)
		return nil
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return nil
	}
	
	var checks []PositionCheck
	seenTimes := make(map[string]bool) // 避免重复时间点
	
	// 提取前几个和后几个事件作为检查点
	sampleIndices := []int{1, 2, 3, 4, 5, len(records)-3, len(records)-2, len(records)-1}
	
	for _, idx := range sampleIndices {
		if idx >= len(records) {
			continue
		}
		
		row := records[idx]
		if len(row) < 9 {
			continue
		}
		
		dateStr := row[7]
		timeStr := row[8]
		dtStr := dateStr + " " + timeStr
		
		if seenTimes[dtStr] {
			continue
		}
		seenTimes[dtStr] = true
		
		// 解析日期时间 (AWST = UTC+8，所以需要转换为UTC)
		t, err := time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
		if err != nil {
			continue
		}
		utcTime := t.Add(-8 * time.Hour) // AWST to UTC
		
		// 转换为JD
		year, month, day := utcTime.Date()
		hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60 + float64(utcTime.Second())/3600
		jdResult := sweph.JulDay(year, int(month), day, hour, true) // true = Gregorian
		
		checks = append(checks, PositionCheck{
			DateTime:    dtStr,
			JD:          jdResult,
			Description: fmt.Sprintf("%s-%s %s", row[0], row[3], row[2]), // P1-P2 aspect
		})
	}
	
	// 添加一些中间时间点
	midIndex := len(records) / 2
	for i := midIndex - 2; i <= midIndex + 2; i++ {
		if i >= 0 && i < len(records) {
			row := records[i]
			if len(row) < 9 {
				continue
			}
			
			dateStr := row[7]
			timeStr := row[8]
			dtStr := dateStr + " " + timeStr
			
			if seenTimes[dtStr] {
				continue
			}
			seenTimes[dtStr] = true
			
			// 解析日期时间 (AWST = UTC+8，所以需要转换为UTC)
			t, err := time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
			if err != nil {
				continue
			}
			utcTime := t.Add(-8 * time.Hour) // AWST to UTC
			
			// 转换为JD
			year, month, day := utcTime.Date()
			hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60 + float64(utcTime.Second())/3600
			jdResult := sweph.JulDay(year, int(month), day, hour, true) // true = Gregorian
			
			checks = append(checks, PositionCheck{
				DateTime:    dtStr,
				JD:          jdResult,
				Description: fmt.Sprintf("%s-%s %s", row[0], row[3], row[2]),
			})
		}
	}
	
	// 按时间排序
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].JD < checks[j].JD
	})
	
	return checks
}

func getEventTypesAtTime(dateTime string) []string {
	// 从CSV中查找指定时间的所有事件类型
	f, err := os.Open("testdata/solarfire/testcase-1-transit.csv")
	if err != nil {
		return nil
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	
	types := make(map[string]bool)
	
	for _, row := range records {
		if len(row) >= 9 && row[7]+" "+row[8] == dateTime {
			types[row[6]] = true // chart type列
		}
	}
	
	var result []string
	for t := range types {
		result = append(result, t)
	}
	sort.Strings(result)
	
	return result
}

func verifyPositionsAtTime(cp PositionCheck, natalJD, natalLat, natalLon float64) []PositionComparison {
	var comparisons []PositionComparison
	
	// 1. 验证Transit行星位置 (使用transit JD)
	fmt.Println("  验证Transit行星位置...")
	transitPositions := getTransitPositions(cp.JD)
	_ = transitPositions // 保留变量供后续使用
	
	// 2. 验证Progressed行星位置 (使用progressed JD)
	fmt.Println("  验证Progressed行星位置...")
	progressedPositions := getProgressedPositions(cp.JD, natalJD)
	
	// 3. 验证Solar Arc行星位置
	fmt.Println("  验证Solar Arc行星位置...")
	solarArcPositions := getSolarArcPositions(cp.JD, natalJD)
	_ = solarArcPositions // 保留变量供后续使用
	
	// 这里需要从Solar Fire获取对应的真实位置来进行对比
	// 由于我们没有SF的内部行星位置数据，我们只能验证计算的一致性
	// 但可以检查是否存在明显的逻辑错误
	
	// 检查Progressed和Solar Arc计算是否合理
	fmt.Printf("    Progressed Sun位置: %.4f°\n", findPosition(progressedPositions, models.PlanetSun))
	fmt.Printf("    Solar Arc偏移量: %.4f°\n", getSolarArcOffset(cp.JD, natalJD))
	
	// 验证Progressed JD计算
	pJD := progressions.SecondaryProgressionJD(natalJD, cp.JD)
	age := (cp.JD - natalJD) / 365.25
	fmt.Printf("    Progressed JD: %.6f (年龄: %.1f年)\n", pJD, age)
	
	// 验证natal太阳位置
	natalSun, _ := sweph.CalcUT(natalJD, sweph.SE_SUN)
	fmt.Printf("    Natal Sun: %.4f°\n", natalSun.Longitude)
	
	// 验证progressed太阳位置
	progressedSun, _ := sweph.CalcUT(pJD, sweph.SE_SUN)
	fmt.Printf("    Progressed Sun: %.4f°\n", progressedSun.Longitude)
	
	// 太阳弧偏移应该是这两个的差值
	expectedOffset := sweph.NormalizeDegrees(progressedSun.Longitude - natalSun.Longitude)
	actualOffset := getSolarArcOffset(cp.JD, natalJD)
	fmt.Printf("    预期Solar Arc偏移: %.4f°\n", expectedOffset)
	fmt.Printf("    实际Solar Arc偏移: %.4f°\n", actualOffset)
	fmt.Printf("    差异: %.4f°\n", sweph.NormalizeDegrees(expectedOffset-actualOffset))
	
	return comparisons
}

func getTransitPositions(jd float64) []PlanetPosition {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
	}
	
	var positions []PlanetPosition
	
	for _, planet := range planets {
		pos, err := sweph.CalcUT(jd, getSwephPlanetID(planet))
		if err != nil {
			continue
		}
		
		positions = append(positions, PlanetPosition{
			Planet:    planet,
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
			Speed:     pos.SpeedLong,
			IsRetro:   pos.IsRetrograde,
		})
	}
	
	return positions
}

func getProgressedPositions(transitJD, natalJD float64) []PlanetPosition {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
	}
	
	var positions []PlanetPosition
	
	for _, planet := range planets {
		lon, speed, err := progressions.CalcProgressedLongitude(planet, natalJD, transitJD)
		if err != nil {
			continue
		}
		
		positions = append(positions, PlanetPosition{
			Planet:    planet,
			Longitude: lon,
			Speed:     speed,
			IsRetro:   speed < 0,
		})
	}
	
	return positions
}

func getSolarArcPositions(transitJD, natalJD float64) []PlanetPosition {
	planets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
	}
	
	var positions []PlanetPosition
	
	for _, planet := range planets {
		lon, speed, err := progressions.CalcSolarArcLongitude(planet, natalJD, transitJD)
		if err != nil {
			continue
		}
		
		positions = append(positions, PlanetPosition{
			Planet:    planet,
			Longitude: lon,
			Speed:     speed,
			IsRetro:   speed < 0,
		})
	}
	
	return positions
}

func getSolarArcOffset(transitJD, natalJD float64) float64 {
	offset, _ := progressions.SolarArcOffset(natalJD, transitJD)
	return offset
}

func findPosition(positions []PlanetPosition, planet models.PlanetID) float64 {
	for _, pos := range positions {
		if pos.Planet == planet {
			return pos.Longitude
		}
	}
	return 0
}

func getSwephPlanetID(planet models.PlanetID) int {
	switch planet {
	case models.PlanetSun:
		return sweph.SE_SUN
	case models.PlanetMoon:
		return sweph.SE_MOON
	case models.PlanetMercury:
		return sweph.SE_MERCURY
	case models.PlanetVenus:
		return sweph.SE_VENUS
	case models.PlanetMars:
		return sweph.SE_MARS
	case models.PlanetJupiter:
		return sweph.SE_JUPITER
	case models.PlanetSaturn:
		return sweph.SE_SATURN
	case models.PlanetUranus:
		return sweph.SE_URANUS
	case models.PlanetNeptune:
		return sweph.SE_NEPTUNE
	case models.PlanetPluto:
		return sweph.SE_PLUTO
	case models.PlanetChiron:
		return sweph.SE_CHIRON
	case models.PlanetNorthNodeMean:
		return sweph.SE_MEAN_NODE
	default:
		return sweph.SE_SUN
	}
}

func analyzeResults(comparisons []PositionComparison) {
	fmt.Println("\n=== 位置验证结果分析 ===")
	
	if len(comparisons) == 0 {
		fmt.Println("没有可分析的数据")
		return
	}
	
	// 按差异大小排序
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].Diff_ArcSec < comparisons[j].Diff_ArcSec
	})
	
	// 统计在容差范围内的比例
	toleranceArcSec := 30.0 // 30角秒容差
	inTolerance := 0
	outTolerance := 0
	
	for _, comp := range comparisons {
		if comp.Within_Tolerance {
			inTolerance++
		} else {
			outTolerance++
		}
	}
	
	total := len(comparisons)
	fmt.Printf("总比较次数: %d\n", total)
	fmt.Printf("在%.0f角秒容差内: %d (%.1f%%)\n", toleranceArcSec, inTolerance, float64(inTolerance)/float64(total)*100)
	fmt.Printf("超出容差: %d (%.1f%%)\n", outTolerance, float64(outTolerance)/float64(total)*100)
	
	// 显示最大差异
	if len(comparisons) > 0 {
		worst := comparisons[len(comparisons)-1]
		fmt.Printf("\n最大差异: %.2f角秒 (%.4f度) - %s 在 %s\n", 
			worst.Diff_ArcSec, worst.Diff_Degrees, worst.Planet, worst.DateTime)
	}
	
	// 显示最小差异
	if len(comparisons) > 0 {
		best := comparisons[0]
		fmt.Printf("最小差异: %.2f角秒 (%.4f度) - %s 在 %s\n",
			best.Diff_ArcSec, best.Diff_Degrees, best.Planet, best.DateTime)
	}
}