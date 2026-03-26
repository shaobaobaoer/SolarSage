package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type EventRecord struct {
	Date      string
	Time      string
	P1        string
	Aspect    string
	P2        string
	Type      string
	DiffSec   int
}

func main() {
	// 读取比较结果
	file, err := os.Open("/tmp/sf_comparison.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var events []EventRecord
	
	// 解析记录
	for _, record := range records {
		if len(record) < 7 {
			continue
		}
		
		diff, err := strconv.Atoi(record[6])
		if err != nil {
			continue
		}
		
		events = append(events, EventRecord{
			Date:    record[0],
			Time:    record[1],
			P1:      record[2],
			Aspect:  record[3],
			P2:      record[4],
			Type:    record[5],
			DiffSec: diff,
		})
	}

	// 按类型分类统计
	typeStats := make(map[string][]int)
	for _, event := range events {
		typeStats[event.Type] = append(typeStats[event.Type], event.DiffSec)
	}

	fmt.Println("=== 时间差异分析报告 ===")
	fmt.Printf("总事件数: %d\n\n", len(events))

	// 分析每种类型的统计信息
	for chartType, diffs := range typeStats {
		sort.Ints(diffs)
		
		count := len(diffs)
		var sum int
		for _, d := range diffs {
			sum += abs(d)
		}
		avg := float64(sum) / float64(count)
		
		median := diffs[count/2]
		min := diffs[0]
		max := diffs[count-1]
		
		// 计算标准差
		var variance float64
		for _, d := range diffs {
			deviation := float64(abs(d)) - avg
			variance += deviation * deviation
		}
		stdDev := 0.0
		if count > 1 {
			stdDev = variance / float64(count-1)
		}
		
		fmt.Printf("图表类型: %s\n", chartType)
		fmt.Printf("  事件数量: %d\n", count)
		fmt.Printf("  平均差异: %.1f 秒 (%.1f 分钟)\n", avg, avg/60)
		fmt.Printf("  中位数: %d 秒 (%.1f 分钟)\n", median, float64(median)/60)
		fmt.Printf("  最小值: %d 秒 (%.1f 分钟)\n", min, float64(min)/60)
		fmt.Printf("  最大值: %d 秒 (%.1f 分钟)\n", max, float64(max)/60)
		fmt.Printf("  标准差: %.1f 秒\n", stdDev)
		fmt.Println()
	}

	// 分析进展相关事件的时间模式
	fmt.Println("=== 进展事件详细分析 ===")
	
	progressEvents := make([]EventRecord, 0)
	for _, event := range events {
		if strings.Contains(event.Type, "Sp-") || strings.Contains(event.Type, "-Sp") ||
		   strings.Contains(event.Type, "Sa-") || strings.Contains(event.Type, "-Sa") {
			progressEvents = append(progressEvents, event)
		}
	}
	
	if len(progressEvents) > 0 {
		// 按差异大小排序
		sort.Slice(progressEvents, func(i, j int) bool {
			return abs(progressEvents[i].DiffSec) < abs(progressEvents[j].DiffSec)
		})
		
		fmt.Printf("进展相关事件总数: %d\n", len(progressEvents))
		fmt.Println("前10个最大差异事件:")
		
		for i := len(progressEvents) - 10; i < len(progressEvents); i++ {
			if i >= 0 {
				event := progressEvents[i]
				dt, _ := time.Parse("2006-01-02 15:04:05", event.Date+" "+event.Time)
				fmt.Printf("  [%s] %s %s %s (%s): %+d秒 (%+.1f小时)\n",
					dt.Format("2006-01-02 15:04"),
					event.P1, event.Aspect, event.P2, event.Type,
					event.DiffSec, float64(event.DiffSec)/3600)
			}
		}
		
		// 统计差异分布
		buckets := map[string]int{
			"< 1分钟": 0, "1-5分钟": 0, "5-30分钟": 0, "30分钟-1小时": 0, "> 1小时": 0,
		}
		
		for _, event := range progressEvents {
			absDiff := abs(event.DiffSec)
			switch {
			case absDiff < 60:
				buckets["< 1分钟"]++
			case absDiff < 300:
				buckets["1-5分钟"]++
			case absDiff < 1800:
				buckets["5-30分钟"]++
			case absDiff < 3600:
				buckets["30分钟-1小时"]++
			default:
				buckets["> 1小时"]++
			}
		}
		
		fmt.Println("\n差异分布:")
		for bucket, count := range buckets {
			percentage := float64(count) / float64(len(progressEvents)) * 100
			fmt.Printf("  %s: %d (%.1f%%)\n", bucket, count, percentage)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}