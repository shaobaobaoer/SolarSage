package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Event struct {
	Date     string
	Time     string
	Planet1  string
	Aspect   string
	Planet2  string
	Chart    string
	FullLine string
}

type TimeDiff struct {
	SFTime     time.Time
	ComputedTime time.Time
	DiffMinutes float64
	Event      Event
}

func parseEvent(line string) (*Event, error) {
	// Parse lines like: "SF 2026-03-25 07:19:29 Chiron Square Sun Tr-Sa"
	parts := strings.Fields(line)
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid line format")
	}
	
	event := &Event{
		Date: parts[1],
		Time: parts[2],
		Planet1: parts[3],
		Aspect: parts[4],
		Planet2: parts[5],
		FullLine: line,
	}
	
	// Extract chart type (last part)
	if len(parts) > 6 {
		event.Chart = parts[len(parts)-1]
	}
	
	return event, nil
}

func parseDateTime(date, timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", date+" "+timeStr)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run analyze_time_diffs.go <compare_output_file>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	var sfEvents []Event
	var computedEvents []Event
	
	// Parse SF and computed events
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MATCH #") {
			// Extract SF and computed times from match lines
			// Format: MATCH #1: SF 2026-03-25 07:19:29 Chiron Square Sun Tr-Sa | Computed 2026-03-25 07:04:29 Square
			parts := strings.Split(line, " | ")
			if len(parts) == 2 {
				sfPart := strings.TrimPrefix(parts[0], "MATCH #"+strings.Fields(parts[0])[1]+": ")
				computedPart := strings.TrimPrefix(parts[1], "Computed ")
				
				if sfEvent, err := parseEvent(sfPart); err == nil {
					sfEvents = append(sfEvents, *sfEvent)
				}
				if computedEvent, err := parseEvent(computedPart); err == nil {
					computedEvents = append(computedEvents, *computedEvent)
				}
			}
		} else if strings.HasPrefix(line, "SF unmatched:") {
			// Parse SF unmatched events
			sfLine := strings.TrimPrefix(line, "SF unmatched: ")
			if event, err := parseEvent(sfLine); err == nil {
				sfEvents = append(sfEvents, *event)
			}
		} else if strings.HasPrefix(line, "Extra computed #") {
			// Parse extra computed events
			computedLine := strings.TrimPrefix(line, "Extra computed #"+strings.Fields(line)[2]+": ")
			if event, err := parseEvent(computedLine); err == nil {
				computedEvents = append(computedEvents, *event)
			}
		}
	}

	// Calculate time differences for matched events
	var diffs []TimeDiff
	matchedCount := 0
	
	for i := 0; i < len(sfEvents) && i < len(computedEvents); i++ {
		sfTime, err1 := parseDateTime(sfEvents[i].Date, sfEvents[i].Time)
		compTime, err2 := parseDateTime(computedEvents[i].Date, computedEvents[i].Time)
		
		if err1 == nil && err2 == nil {
			diff := compTime.Sub(sfTime).Minutes()
			diffs = append(diffs, TimeDiff{
				SFTime: sfTime,
				ComputedTime: compTime,
				DiffMinutes: diff,
				Event: sfEvents[i],
			})
			matchedCount++
		}
	}

	// Group by chart type
	chartTypeStats := make(map[string][]TimeDiff)
	for _, diff := range diffs {
		chartType := diff.Event.Chart
		chartTypeStats[chartType] = append(chartTypeStats[chartType], diff)
	}

	// Print statistics
	fmt.Printf("=== 时间差异分析报告 ===\n")
	fmt.Printf("总匹配事件数: %d\n", matchedCount)
	fmt.Printf("\n按图表类型分组:\n")
	
	// Sort chart types for consistent output
	var chartTypes []string
	for ct := range chartTypeStats {
		chartTypes = append(chartTypes, ct)
	}
	sort.Strings(chartTypes)
	
	for _, ct := range chartTypes {
		stats := chartTypeStats[ct]
		if len(stats) == 0 {
			continue
		}
		
		var totalDiff float64
		var maxDiff, minDiff float64
		maxDiff = stats[0].DiffMinutes
		minDiff = stats[0].DiffMinutes
		
		fmt.Printf("\n%s (%d events):\n", ct, len(stats))
		for _, stat := range stats {
			totalDiff += stat.DiffMinutes
			if stat.DiffMinutes > maxDiff {
				maxDiff = stat.DiffMinutes
			}
			if stat.DiffMinutes < minDiff {
				minDiff = stat.DiffMinutes
			}
			fmt.Printf("  %s: SF %s | Computed %s | 差异: %.1f分钟\n",
				stat.Event.Planet1+" "+stat.Event.Aspect+" "+stat.Event.Planet2,
				stat.SFTime.Format("15:04:05"),
				stat.ComputedTime.Format("15:04:05"),
				stat.DiffMinutes)
		}
		
		avgDiff := totalDiff / float64(len(stats))
		fmt.Printf("  平均差异: %.1f分钟, 最大: %.1f分钟, 最小: %.1f分钟\n", avgDiff, maxDiff, minDiff)
	}

	// Overall statistics
	if len(diffs) > 0 {
		var totalDiff float64
		var absTotalDiff float64
		for _, diff := range diffs {
			totalDiff += diff.DiffMinutes
			absTotalDiff += abs(diff.DiffMinutes)
		}
		avgDiff := totalDiff / float64(len(diffs))
		avgAbsDiff := absTotalDiff / float64(len(diffs))
		
		fmt.Printf("\n=== 总体统计 ===\n")
		fmt.Printf("平均时间差异: %.1f分钟\n", avgDiff)
		fmt.Printf("平均绝对差异: %.1f分钟\n", avgAbsDiff)
		fmt.Printf("标准差: %.1f分钟\n", calculateStdDev(diffs, avgDiff))
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func calculateStdDev(diffs []TimeDiff, mean float64) float64 {
	if len(diffs) <= 1 {
		return 0
	}
	
	var sumSquaredDiffs float64
	for _, diff := range diffs {
		diffFromMean := diff.DiffMinutes - mean
		sumSquaredDiffs += diffFromMean * diffFromMean
	}
	
	variance := sumSquaredDiffs / float64(len(diffs)-1)
	return sqrt(variance)
}

func sqrt(x float64) float64 {
	// Simple Newton-Raphson method for square root
	if x <= 0 {
		return 0
	}
	
	z := x
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}