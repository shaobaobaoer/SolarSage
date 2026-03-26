package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TimeDiff struct {
	sfTime    string
	compTime  string
	diffSec   int
	eventDesc string
	chartType string
}

func main() {
	// 从比较输出中提取时间差异数据
	file, err := os.Open("/tmp/comparison_output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var diffs []TimeDiff
	scanner := bufio.NewScanner(file)
	
	inWorstSection := false
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.Contains(line, "=== Worst 20") {
			inWorstSection = true
			continue
		}
		
		if strings.Contains(line, "=== Best 20") {
			break
		}
		
		if inWorstSection && strings.Contains(line, "| SF ") {
			// 解析格式: "  -68100s (-1135.0min) | SF 2026-05-22 12:07:13 | Comp 2026-05-21 17:12:13 | Neptune Opposition ASC Sp-Sp"
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				// 提取时间差
				diffPart := strings.TrimSpace(parts[0])
				diffSec := 0
				if strings.Contains(diffPart, "s") {
					diffStr := strings.TrimSuffix(strings.Fields(diffPart)[0], "s")
					diffSec, _ = strconv.Atoi(diffStr)
				}
				
				// 提取SF时间和Comp时间
				sfTime := strings.TrimSpace(parts[1])
				compTime := strings.TrimSpace(parts[2])
				
				// 提取事件描述和图表类型
				eventAndType := strings.TrimSpace(parts[3])
				eventParts := strings.Split(eventAndType, " ")
				chartType := eventParts[len(eventParts)-1]
				eventDesc := strings.Join(eventParts[:len(eventParts)-1], " ")
				
				diffs = append(diffs, TimeDiff{
					sfTime:    sfTime,
					compTime:  compTime,
					diffSec:   diffSec,
					eventDesc: eventDesc,
					chartType: chartType,
				})
			}
		}
	}

	// 按图表类型分类统计
	typeStats := make(map[string][]TimeDiff)
	for _, diff := range diffs {
		typeStats[diff.chartType] = append(typeStats[diff.chartType], diff)
	}

	fmt.Println("=== 各类图表时间差异统计 ===")
	
	for chartType, diffs := range typeStats {
		if len(diffs) == 0 {
			continue
		}
		
		// 计算统计信息
		var totalDiff int
		var absDiffs []int
		for _, d := range diffs {
			totalDiff += d.diffSec
			absDiffs = append(absDiffs, abs(d.diffSec))
		}
		
		avgDiff := float64(totalDiff) / float64(len(diffs))
		sort.Ints(absDiffs)
		medianAbs := absDiffs[len(absDiffs)/2]
		
		minAbs := absDiffs[0]
		maxAbs := absDiffs[len(absDiffs)-1]
		
		fmt.Printf("\n图表类型: %s (样本数: %d)\n", chartType, len(diffs))
		fmt.Printf("  平均差异: %.1f 秒 (%+.1f小时)\n", avgDiff, avgDiff/3600)
		fmt.Printf("  平均绝对差异: %.1f 秒 (%.1f小时)\n", float64(sum(absDiffs))/float64(len(absDiffs)), float64(sum(absDiffs))/float64(len(absDiffs))/3600)
		fmt.Printf("  中位数绝对差异: %d 秒 (%.1f小时)\n", medianAbs, float64(medianAbs)/3600)
		fmt.Printf("  最小绝对差异: %d 秒 (%.1f小时)\n", minAbs, float64(minAbs)/3600)
		fmt.Printf("  最大绝对差异: %d 秒 (%.1f小时)\n", maxAbs, float64(maxAbs)/3600)
		
		// 显示几个典型例子
		fmt.Printf("  典型例子:\n")
		count := 0
		for _, diff := range diffs {
			if count >= 3 {
				break
			}
			fmt.Printf("    %s: %s vs %s (%+ds)\n", 
				diff.eventDesc, diff.sfTime, diff.compTime, diff.diffSec)
			count++
		}
	}
	
	// 分析所有进展相关事件的一致性
	fmt.Println("\n=== 进展事件一致性分析 ===")
	progressTypes := []string{"Tr-Sp", "Tr-Sa", "Sp-Na", "Sp-Sp", "Sa-Na"}
	
	for _, pt := range progressTypes {
		if diffs, exists := typeStats[pt]; exists && len(diffs) > 0 {
			fmt.Printf("\n%s 事件的一致性检验:\n", pt)
			
			// 检查是否有系统性偏差
			var positiveCount, negativeCount int
			var totalPositive, totalNegative int
			
			for _, diff := range diffs {
				if diff.diffSec > 0 {
					positiveCount++
					totalPositive += diff.diffSec
				} else {
					negativeCount++
					totalNegative += diff.diffSec
				}
			}
			
			fmt.Printf("  正偏差事件: %d 个 (平均 %+d秒)\n", positiveCount, totalPositive/positiveCount)
			fmt.Printf("  负偏差事件: %d 个 (平均 %+d秒)\n", negativeCount, totalNegative/negativeCount)
			
			// 如果正负偏差数量相近，说明是系统性偏差而非随机误差
			if len(diffs) > 1 && abs(positiveCount-negativeCount) <= len(diffs)/3 {
				fmt.Printf("  → 存在明显的系统性偏差模式\n")
			} else if len(diffs) > 1 {
				fmt.Printf("  → 偏差方向相对随机\n")
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sum(slice []int) int {
	total := 0
	for _, v := range slice {
		total += v
	}
	return total
}