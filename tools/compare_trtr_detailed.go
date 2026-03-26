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
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/transit"
)

func main() {
	// 初始化星历
	exe, _ := os.Executable()
	ephePath := filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
	if _, err := os.Stat(ephePath); err != nil {
		ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
	}
	sweph.Init(ephePath)
	defer sweph.Close()

	natalJD := 2450800.900009
	natalLat := 30.9
	natalLon := 121.15

	fmt.Println("=== Tr-Tr Events Detailed Comparison ===")

	// 读取SF事件以确定时间范围
	sfEvents := readSFTrTrEvents()
	if len(sfEvents) == 0 {
		fmt.Println("No SF events found")
		return
	}
	
	// 使用SF事件的时间范围
	startJD := sfEvents[0].JD - 1.0
	endJD := sfEvents[len(sfEvents)-1].JD + 1.0
	fmt.Printf("Time range: %.6f to %.6f\n", startJD, endJD)

	events, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalLat:     natalLat,
		NatalLon:     natalLon,
		NatalJD:      natalJD,
		NatalPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto, models.PlanetChiron, models.PlanetNorthNodeMean,
		},
		TransitLat: natalLat,
		TransitLon: natalLon,
		StartJD:    startJD,
		EndJD:      endJD,
		TransitPlanets: []models.PlanetID{
			// 只使用外行星
			models.PlanetJupiter, models.PlanetSaturn, models.PlanetUranus,
			models.PlanetNeptune, models.PlanetPluto, models.PlanetChiron,
			models.PlanetNorthNodeMean,
		},
		EventConfig: models.EventConfig{
			IncludeTrTr: true,
		},
		OrbConfigTransit: models.OrbConfig{
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

	// 提取计算的Exact事件
	var computedExactEvents []models.TransitEvent
	for _, e := range events {
		if e.EventType == models.EventAspectExact {
			computedExactEvents = append(computedExactEvents, e)
		}
	}

	fmt.Printf("SF Tr-Tr Exact events: %d\n", len(sfEvents))
	fmt.Printf("Computed Tr-Tr Exact events: %d\n", len(computedExactEvents))

	// 显示计算的事件
	fmt.Println("\n=== Computed Events ===")
	for _, comp := range computedExactEvents {
		fmt.Printf("  %s %s %s at %s\n", comp.Planet, comp.AspectType, comp.Target, jdToString(comp.JD))
	}

	// 逐个匹配
	fmt.Println("\n=== Event-by-Event Comparison ===")
	fmt.Printf("%-25s | %-15s | %-20s | %-20s | %-10s\n",
		"Planet Pair", "Aspect", "SF Time", "Computed Time", "Diff")
	fmt.Println("----------------------------------------------------------------------------------------")

	matched := 0
	unmatchedSF := 0

	for _, sf := range sfEvents {
		// 在计算事件中查找匹配
		bestMatch := -1
		bestDiff := float64(1e9)

		for i, comp := range computedExactEvents {
			// 匹配行星对（忽略大小写，并处理NorthNode命名差异）
			compP1 := normalizePlanetName(string(comp.Planet))
			compP2 := normalizePlanetName(comp.Target)
			sfP1 := normalizePlanetName(sf.P1)
			sfP2 := normalizePlanetName(sf.P2)
			
			p1Match := (compP1 == sfP1 && compP2 == sfP2) ||
				(compP1 == sfP2 && compP2 == sfP1)
			aspectMatch := strings.ToUpper(string(comp.AspectType)) == strings.ToUpper(sf.Aspect)

			if p1Match && aspectMatch {
				diff := abs(sf.JD - comp.JD) * 24 * 3600 // 转换为秒
				if diff < bestDiff {
					bestDiff = diff
					bestMatch = i
				}
			}
		}

		pairStr := fmt.Sprintf("%s-%s", sf.P1, sf.P2)
		if bestMatch >= 0 {
			comp := computedExactEvents[bestMatch]
			compTime := jdToString(comp.JD)
			sfTime := jdToString(sf.JD)
			fmt.Printf("%-25s | %-15s | %-20s | %-20s | %.1fs\n",
				pairStr, sf.Aspect, sfTime, compTime, bestDiff)
			matched++
		} else {
			fmt.Printf("%-25s | %-15s | %-20s | %-20s | NO MATCH\n",
				pairStr, sf.Aspect, jdToString(sf.JD), "-")
			unmatchedSF++
		}
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Matched: %d/%d (%.1f%%)\n", matched, len(sfEvents), float64(matched)/float64(len(sfEvents))*100)
	fmt.Printf("Unmatched SF: %d\n", unmatchedSF)

	// 检查是否有额外的计算事件
	extraComputed := 0
	for _, comp := range computedExactEvents {
		found := false
		for _, sf := range sfEvents {
			p1Match := (string(comp.Planet) == sf.P1 && comp.Target == sf.P2) ||
				(string(comp.Planet) == sf.P2 && comp.Target == sf.P1)
			aspectMatch := string(comp.AspectType) == sf.Aspect
			if p1Match && aspectMatch {
				found = true
				break
			}
		}
		if !found {
			extraComputed++
			fmt.Printf("Extra computed: %s %s %s at %s\n",
				comp.Planet, comp.AspectType, comp.Target, jdToString(comp.JD))
		}
	}
	fmt.Printf("Extra computed events: %d\n", extraComputed)
}

type SFEvent struct {
	P1    string
	P2    string
	Aspect string
	JD    float64
}

func readSFTrTrEvents() []SFEvent {
	f, err := os.Open("testdata/solarfire/testcase-1-transit.csv")
	if err != nil {
		return nil
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	var events []SFEvent
	for _, row := range records[1:] {
		if len(row) < 9 {
			continue
		}
		if row[6] != "Tr-Tr" || row[5] != "Exact" {
			continue
		}

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

		events = append(events, SFEvent{
			P1:     row[0],
			P2:     row[3],
			Aspect: row[2],
			JD:     jd,
		})
	}

	// 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].JD < events[j].JD
	})

	return events
}

func jdToString(jd float64) string {
	year, month, day, hour := sweph.RevJul(jd, true)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		year, month, day,
		int(hour), int((hour-float64(int(hour)))*60),
		int(((hour-float64(int(hour)))*60-float64(int((hour-float64(int(hour)))*60)))*60))
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// normalizePlanetName 统一行星名称格式
func normalizePlanetName(name string) string {
	name = strings.ToUpper(name)
	// 处理NorthNode的命名差异
	if name == "NORTHNODE" || name == "NORTH_NODE" || name == "NORTH_NODE_MEAN" {
		return "NORTHNODE"
	}
	return name
}