package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/transit"
)

// TrTrTestCase 专门测试Tr-Tr事件的测试案例
type TrTrTestCase struct {
	Name       string
	NatalJD    float64
	NatalLat   float64
	NatalLon   float64
	Timezone   string
	SFCSVPath  string
}

func main() {
	// 命令行参数
	testRange := flag.String("range", "all", "测试范围: all, early, middle, late")
	debug := flag.Bool("debug", false, "显示详细调试信息")
	flag.Parse()

	// 初始化星历
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	if _, err := os.Stat(ephePath); err != nil {
		ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
	}
	sweph.Init(ephePath)
	defer sweph.Close()

	// 测试案例参数 (JN案例)
	tc := TrTrTestCase{
		Name:      "JN - Tr-Tr Events Only",
		NatalJD:   2450800.900009, // 1997-12-18 09:36 UTC
		NatalLat:  30.9,           // 30°54'N
		NatalLon:  121.15,         // 121°09'E
		Timezone:  "Australia/Perth",
		SFCSVPath: "testdata/solarfire/testcase-1-transit.csv",
	}

	fmt.Printf("=== Tr-Tr Events Step-by-Step Validation ===\n")
	fmt.Printf("Test case: %s\n", tc.Name)
	fmt.Printf("Natal JD: %.6f\n", tc.NatalJD)
	fmt.Printf("Test range: %s\n", *testRange)

	// 获取Tr-Tr事件时间段
	startJD, endJD := getTimeRange(tc, *testRange)
	fmt.Printf("Time range: JD %.6f to %.6f\n", startJD, endJD)

	// 只计算Tr-Tr事件
	fmt.Println("\n=== Running Tr-Tr Events Calculation ===")
	events, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalLat:     tc.NatalLat,
		NatalLon:     tc.NatalLon,
		NatalJD:      tc.NatalJD,
		NatalPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
		},
		TransitLat:   tc.NatalLat,
		TransitLon:   tc.NatalLon,
		StartJD:      startJD,
		EndJD:        endJD,
		TransitPlanets: []models.PlanetID{
			// SF只包含外行星的Tr-Tr事件，不包含内行星和Moon
			models.PlanetJupiter, models.PlanetSaturn, models.PlanetUranus,
			models.PlanetNeptune, models.PlanetPluto, models.PlanetChiron,
			models.PlanetNorthNodeMean,
		},
		EventConfig: models.EventConfig{
			IncludeTrTr: true, // 只包含Tr-Tr事件
		},
		OrbConfigTransit: models.OrbConfig{
			// SF只包含这些相位（不包含Semi-Sextile）
			Conjunction: 1, Opposition: 1, Trine: 1, Square: 1,
			Sextile: 1, Quincunx: 1,
			SemiSextile: -1, // 禁用
		},
		HouseSystem: models.HousePlacidus,
	})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Generated %d Tr-Tr events\n", len(events))

	// 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].JD < events[j].JD
	})

	// 显示前几个和后几个事件
	showSampleEvents(events, 5)

	// 与Solar Fire对比
	if *debug {
		compareWithSolarFire(tc.SFCSVPath, events, tc.Timezone, *debug)
	} else {
		simpleMatchRate(tc.SFCSVPath, events)
	}
}

func getTimeRange(tc TrTrTestCase, rangeType string) (float64, float64) {
	// 从CSV中获取实际的时间范围
	f, err := os.Open(tc.SFCSVPath)
	if err != nil {
		// 默认范围：2026-2027年
		start := sweph.JulDay(2026, 1, 1, 0, true)
		end := sweph.JulDay(2027, 12, 31, 23.99, true)
		return start, end
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		start := sweph.JulDay(2026, 1, 1, 0, true)
		end := sweph.JulDay(2027, 12, 31, 23.99, true)
		return start, end
	}

	// 找到第一个和最后一个Tr-Tr事件
	var firstJD, lastJD float64
	foundFirst := false

	for _, row := range records[1:] { // 跳过标题行
		if len(row) >= 7 && row[6] == "Tr-Tr" {
			// 解析日期时间
			if len(row) >= 9 {
				dateStr := row[7]
				timeStr := row[8]
				
				t, err := time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
				if err != nil {
					continue
				}
				utcTime := t.Add(-8 * time.Hour) // AWST to UTC
				
				year, month, day := utcTime.Date()
				hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60 + float64(utcTime.Second())/3600
				jd := sweph.JulDay(year, int(month), day, hour, true)
				
				if !foundFirst {
					firstJD = jd
					foundFirst = true
				}
				lastJD = jd
			}
		}
	}

	if foundFirst {
		// 添加缓冲时间确保捕获边界事件
		return firstJD - 1.0, lastJD + 1.0
	}

	// 如果没找到Tr-Tr事件，使用默认范围
	start := sweph.JulDay(2026, 1, 1, 0, true)
	end := sweph.JulDay(2027, 12, 31, 23.99, true)
	return start, end
}

func showSampleEvents(events []models.TransitEvent, count int) {
	fmt.Printf("\n=== Sample Tr-Tr Events ===\n")
	
	// 显示前几个
	limit := count
	if limit > len(events) {
		limit = len(events)
	}
	
	fmt.Println("First events:")
	for i := 0; i < limit; i++ {
		e := events[i]
		year, month, day, hour := sweph.RevJul(e.JD, true)
		fmt.Printf("  [%d] %04d-%02d-%02d %02d:%02d | %s %s %s | %.4f°\n",
			i+1, year, month, day, int(hour), int((hour-float64(int(hour)))*60),
			e.Planet, e.AspectType, e.Target, e.PlanetLongitude)
	}

	// 显示后几个
	if len(events) > count*2 {
		fmt.Println("... (skipping middle events) ...")
		fmt.Println("Last events:")
		for i := len(events) - count; i < len(events); i++ {
			e := events[i]
			year, month, day, hour := sweph.RevJul(e.JD, true)
			fmt.Printf("  [%d] %04d-%02d-%02d %02d:%02d | %s %s %s | %.4f°\n",
				i+1, year, month, day, int(hour), int((hour-float64(int(hour)))*60),
				e.Planet, e.AspectType, e.Target, e.PlanetLongitude)
		}
	}
}

func simpleMatchRate(sfPath string, computedEvents []models.TransitEvent) {
	// 简单匹配率统计
	f, err := os.Open(sfPath)
	if err != nil {
		fmt.Printf("Cannot open SF CSV: %v\n", err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	// 统计SF中的Tr-Tr事件
	sfTrTrCount := 0
	for _, row := range records[1:] {
		if len(row) >= 7 && row[6] == "Tr-Tr" && len(row) >= 5 && row[5] == "Exact" {
			sfTrTrCount++
		}
	}

	// 统计计算出的Tr-Tr Exact事件
	compTrTrExactCount := 0
	compTrTrEnterCount := 0
	compTrTrLeaveCount := 0
	aspectTypeCounts := make(map[string]int)
	planetPairCounts := make(map[string]int)
	for _, e := range computedEvents {
		switch e.EventType {
		case models.EventAspectExact:
			compTrTrExactCount++
			aspectTypeCounts[string(e.AspectType)]++
			pairKey := string(e.Planet) + "-" + e.Target
			if e.Target < string(e.Planet) {
				pairKey = e.Target + "-" + string(e.Planet)
			}
			planetPairCounts[pairKey]++
		case models.EventAspectEnter:
			compTrTrEnterCount++
		case models.EventAspectLeave:
			compTrTrLeaveCount++
		}
	}

	fmt.Printf("\n=== Simple Match Statistics ===\n")
	fmt.Printf("SF Tr-Tr Exact events: %d\n", sfTrTrCount)
	fmt.Printf("Computed Tr-Tr events:\n")
	fmt.Printf("  Exact: %d\n", compTrTrExactCount)
	fmt.Printf("  Enter: %d\n", compTrTrEnterCount)
	fmt.Printf("  Leave: %d\n", compTrTrLeaveCount)
	fmt.Printf("  Total: %d\n", compTrTrExactCount+compTrTrEnterCount+compTrTrLeaveCount)
	
	if sfTrTrCount > 0 {
		matchRate := float64(compTrTrExactCount) / float64(sfTrTrCount) * 100
		fmt.Printf("Exact match rate: %.1f%%\n", matchRate)
	}

	// 显示按相位类型统计
	fmt.Printf("\n=== Aspect Type Distribution (Computed Exact events) ===\n")
	var aspectTypes []string
	for t := range aspectTypeCounts {
		aspectTypes = append(aspectTypes, t)
	}
	sort.Strings(aspectTypes)
	for _, t := range aspectTypes {
		fmt.Printf("  %s: %d\n", t, aspectTypeCounts[t])
	}

	// 显示按行星对统计
	fmt.Printf("\n=== Planet Pair Distribution (Top 20) ===\n")
	type pairCount struct {
		pair  string
		count int
	}
	var pairs []pairCount
	for p, c := range planetPairCounts {
		pairs = append(pairs, pairCount{p, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	limit := 20
	if limit > len(pairs) {
		limit = len(pairs)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("  %s: %d\n", pairs[i].pair, pairs[i].count)
	}
}

func compareWithSolarFire(sfPath string, computedEvents []models.TransitEvent, tz string, debug bool) {
	// 详细的对比分析
	fmt.Println("\n=== Detailed Comparison ===")
	// 这里可以实现更详细的逐个事件对比
	// 暂时保持简单版本
	simpleMatchRate(sfPath, computedEvents)
}