package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type MatchEvent struct {
	SFTime       time.Time
	ComputedTime time.Time
	DiffMinutes  float64
	EventType    string
	ChartType    string
	Description  string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run analyze_matches.go <compare_output_file>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	var matches []MatchEvent
	
	// Parse MATCH lines
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "MATCH #") {
			// Format: MATCH #1: SF 2026-03-25 07:19:29 Chiron Square Sun Tr-Sa | Computed 2026-03-25 07:04:29 Square
			parts := strings.Split(line, " | ")
			if len(parts) == 2 {
				sfPart := strings.TrimSpace(strings.TrimPrefix(parts[0], strings.Split(parts[0], ":")[0]+": "))
				computedPart := strings.TrimSpace(strings.TrimPrefix(parts[1], "Computed "))
				
				// Parse SF time and event
				sfFields := strings.Fields(sfPart)
				if len(sfFields) >= 6 {
					sfDateTime := sfFields[1] + " " + sfFields[2]
					sfTime, err := time.Parse("2006-01-02 15:04:05", sfDateTime)
					if err != nil {
						continue
					}
					
					// Parse computed time
					compFields := strings.Fields(computedPart)
					if len(compFields) >= 3 {
						compDateTime := compFields[0] + " " + compFields[1]
						compTime, err := time.Parse("2006-01-02 15:04:05", compDateTime)
						if err != nil {
							continue
						}
						
						// Extract chart type (last field)
						chartType := sfFields[len(sfFields)-1]
						aspect := sfFields[4]
						
						diff := compTime.Sub(sfTime).Minutes()
						
						matches = append(matches, MatchEvent{
							SFTime:       sfTime,
							ComputedTime: compTime,
							DiffMinutes:  diff,
							EventType:    aspect,
							ChartType:    chartType,
							Description:  fmt.Sprintf("%s %s %s", sfFields[3], aspect, sfFields[5]),
						})
					}
				}
			}
		}
	}

	// Print analysis
	fmt.Printf("=== 时间精度匹配分析 ===\n")
	fmt.Printf("总匹配事件数: %d\n\n", len(matches))
	
	if len(matches) == 0 {
		fmt.Println("没有找到匹配的事件")
		return
	}
	
	// Group by chart type
	chartGroups := make(map[string][]MatchEvent)
	for _, match := range matches {
		chartGroups[match.ChartType] = append(chartGroups[match.ChartType], match)
	}
	
	// Sort chart types
	var chartTypes []string
	for ct := range chartGroups {
		chartTypes = append(chartTypes, ct)
	}
	sort.Strings(chartTypes)
	
	// Display results by chart type
	for _, ct := range chartTypes {
		group := chartGroups[ct]
		fmt.Printf("%s (%d 个匹配):\n", ct, len(group))
		
		var totalDiff, absTotalDiff float64
		maxDiff := group[0].DiffMinutes
		minDiff := group[0].DiffMinutes
		
		for _, match := range group {
			totalDiff += match.DiffMinutes
			absDiff := abs(match.DiffMinutes)
			absTotalDiff += absDiff
			
			if match.DiffMinutes > maxDiff {
				maxDiff = match.DiffMinutes
			}
			if match.DiffMinutes < minDiff {
				minDiff = match.DiffMinutes
			}
			
			fmt.Printf("  %s: SF %s | 计算 %s | 差异: %+.1f分钟\n",
				match.Description,
				match.SFTime.Format("15:04:05"),
				match.ComputedTime.Format("15:04:05"),
				match.DiffMinutes)
		}
		
		avgDiff := totalDiff / float64(len(group))
		avgAbsDiff := absTotalDiff / float64(len(group))
		
		fmt.Printf("  平均差异: %+.1f分钟, 绝对平均: %.1f分钟\n", avgDiff, avgAbsDiff)
		fmt.Printf("  最大差异: %+.1f分钟, 最小差异: %+.1f分钟\n\n", maxDiff, minDiff)
	}
	
	// Overall statistics
	var totalDiff, absTotalDiff float64
	maxDiff := matches[0].DiffMinutes
	minDiff := matches[0].DiffMinutes
	
	for _, match := range matches {
		totalDiff += match.DiffMinutes
		absDiff := abs(match.DiffMinutes)
		absTotalDiff += absDiff
		
		if match.DiffMinutes > maxDiff {
			maxDiff = match.DiffMinutes
		}
		if match.DiffMinutes < minDiff {
			minDiff = match.DiffMinutes
		}
	}
	
	overallAvg := totalDiff / float64(len(matches))
	overallAvgAbs := absTotalDiff / float64(len(matches))
	stdDev := calculateStdDev(matches, overallAvg)
	
	fmt.Printf("=== 总体统计 ===\n")
	fmt.Printf("平均时间差异: %+.1f分钟\n", overallAvg)
	fmt.Printf("平均绝对差异: %.1f分钟\n", overallAvgAbs)
	fmt.Printf("标准差: %.1f分钟\n", stdDev)
	fmt.Printf("最大差异: %+.1f分钟\n", maxDiff)
	fmt.Printf("最小差异: %+.1f分钟\n", minDiff)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func calculateStdDev(matches []MatchEvent, mean float64) float64 {
	if len(matches) <= 1 {
		return 0
	}
	
	var sumSquaredDiffs float64
	for _, match := range matches {
		diffFromMean := match.DiffMinutes - mean
		sumSquaredDiffs += diffFromMean * diffFromMean
	}
	
	variance := sumSquaredDiffs / float64(len(matches)-1)
	return sqrt(variance)
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	
	z := x
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}