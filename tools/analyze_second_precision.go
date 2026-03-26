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
	DiffSeconds  float64
	EventType    string
	ChartType    string
	Description  string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run analyze_second_precision.go <compare_output_file>")
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
	
	// Parse MATCH lines with second precision
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
						
						diff := compTime.Sub(sfTime).Seconds()
						
						matches = append(matches, MatchEvent{
							SFTime:       sfTime,
							ComputedTime: compTime,
							DiffSeconds:  diff,
							EventType:    aspect,
							ChartType:    chartType,
							Description:  fmt.Sprintf("%s %s %s", sfFields[3], aspect, sfFields[5]),
						})
					}
				}
			}
		}
	}

	// Print second-precision analysis
	fmt.Printf("=== 秒级精度匹配分析 ===\n")
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
	
	// Display results by chart type with second precision
	for _, ct := range chartTypes {
		group := chartGroups[ct]
		fmt.Printf("%s (%d 个匹配):\n", ct, len(group))
		
		var totalDiff, absTotalDiff float64
		maxDiff := group[0].DiffSeconds
		minDiff := group[0].DiffSeconds
		
		for _, match := range group {
			totalDiff += match.DiffSeconds
			absDiff := abs(match.DiffSeconds)
			absTotalDiff += absDiff
			
			if match.DiffSeconds > maxDiff {
				maxDiff = match.DiffSeconds
			}
			if match.DiffSeconds < minDiff {
				minDiff = match.DiffSeconds
			}
			
			fmt.Printf("  %s: SF %s | 计算 %s | 差异: %+.0f秒 (%+.1f分钟)\n",
				match.Description,
				match.SFTime.Format("15:04:05"),
				match.ComputedTime.Format("15:04:05"),
				match.DiffSeconds,
				match.DiffSeconds/60.0)
		}
		
		avgDiff := totalDiff / float64(len(group))
		avgAbsDiff := absTotalDiff / float64(len(group))
		
		fmt.Printf("  平均差异: %+.0f秒 (%+.1f分钟)\n", avgDiff, avgDiff/60.0)
		fmt.Printf("  绝对平均: %.0f秒 (%.1f分钟)\n", avgAbsDiff, avgAbsDiff/60.0)
		fmt.Printf("  最大差异: %+.0f秒 (%+.1f分钟)\n", maxDiff, maxDiff/60.0)
		fmt.Printf("  最小差异: %+.0f秒 (%+.1f分钟)\n\n", minDiff, minDiff/60.0)
	}
	
	// Overall statistics
	var totalDiff, absTotalDiff float64
	maxDiff := matches[0].DiffSeconds
	minDiff := matches[0].DiffSeconds
	
	for _, match := range matches {
		totalDiff += match.DiffSeconds
		absDiff := abs(match.DiffSeconds)
		absTotalDiff += absDiff
		
		if match.DiffSeconds > maxDiff {
			maxDiff = match.DiffSeconds
		}
		if match.DiffSeconds < minDiff {
			minDiff = match.DiffSeconds
		}
	}
	
	overallAvg := totalDiff / float64(len(matches))
	overallAvgAbs := absTotalDiff / float64(len(matches))
	stdDev := calculateStdDev(matches, overallAvg)
	
	fmt.Printf("=== 总体统计 (秒级精度) ===\n")
	fmt.Printf("平均时间差异: %+.0f秒 (%+.1f分钟)\n", overallAvg, overallAvg/60.0)
	fmt.Printf("平均绝对差异: %.0f秒 (%.1f分钟)\n", overallAvgAbs, overallAvgAbs/60.0)
	fmt.Printf("标准差: %.0f秒 (%.1f分钟)\n", stdDev, stdDev/60.0)
	fmt.Printf("最大差异: %+.0f秒 (%+.1f分钟)\n", maxDiff, maxDiff/60.0)
	fmt.Printf("最小差异: %+.0f秒 (%+.1f分钟)\n", minDiff, minDiff/60.0)
	
	// Precision assessment
	fmt.Println()
	fmt.Println("=== 精度评估 ===")
	if overallAvgAbs <= 1.0 {
		fmt.Println("🎯 秒级精度达标 (<1秒差异)")
	} else if overallAvgAbs <= 5.0 {
		fmt.Println("✅ 高精度 (<5秒差异)")
	} else if overallAvgAbs <= 30.0 {
		fmt.Println("⚠️ 中等精度 (<30秒差异)")
	} else {
		fmt.Println("❌ 精度不足 (>30秒差异)")
	}
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
		diffFromMean := match.DiffSeconds - mean
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