package main

import (
	"fmt"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/export"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/transit"
)

func main() {
	fmt.Println("=== 快速进展修正验证 ===")
	
	// 创建测试输入
	input := &transit.Input{
		NatalChart: &models.Chart{
			JD:      2450800.900000, // 1997-12-18 09:36:00 UTC
			Lat:     39.9042,
			Lon:     116.4074,
			Hsys:    models.HousePlacidus,
			Calendar: models.CalendarGregorian,
		},
		StartJD: 2461100.0, // 2026-04-01
		EndJD:   2461200.0, // 2026-07-10
		Charts: transit.ChartConfig{
			Transits:    true,
			Progressions: true,
			SolarArc:    true,
		},
		Aspects: transit.DefaultAspects(),
		Orbs:    transit.DefaultOrbs(),
	}
	
	// 运行计算
	fmt.Println("开始计算...")
	events, err := transit.Calculate(input)
	if err != nil {
		fmt.Printf("计算错误: %v\n", err)
		return
	}
	
	fmt.Printf("计算完成，共 %d 个事件\n", len(events))
	
	// 分析事件类型
	typeCounts := make(map[string]int)
	for _, e := range events {
		if e.EventType == models.EventAspectExact {
			key := fmt.Sprintf("%s-%s", e.ChartType1, e.ChartType2)
			typeCounts[key]++
		}
	}
	
	fmt.Println("事件类型统计:")
	for k, v := range typeCounts {
		fmt.Printf("  %s: %d个\n", k, v)
	}
	
	// 导出几个事件查看时间格式
	if len(events) > 0 {
		fmt.Println("\n示例事件:")
		count := 0
		for _, e := range events {
			if e.EventType == models.EventAspectExact && count < 3 {
				row := export.EventToCSVRow(e, "Asia/Shanghai")
				fmt.Printf("  %s %s %s %s %s | %s %s\n", 
					row.Date, row.Time, row.P1, row.Aspect, row.P2, row.Type, row.PairType)
				count++
			}
		}
	}
	
	fmt.Println("验证完成!")
}