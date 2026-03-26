package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	// 使用绝对路径确保能找到文件
	ephePath := "/home/ecs-user/SolarSage/third_party/swisseph/ephe"
	sweph.Init(ephePath)
	defer sweph.Close()

	// 测试多个时间点的Chiron位置
	testDates := []struct {
		year, month, day int
		hour float64
	}{
		{2026, 3, 1, 0},
		{2026, 6, 1, 0},
		{2026, 9, 1, 0},
		{2026, 12, 1, 0},
		{2027, 3, 1, 0},
	}

	fmt.Println("=== Chiron位置测试 ===")
	for _, d := range testDates {
		jd := sweph.JulDay(d.year, d.month, d.day, d.hour, true)
		res, err := sweph.CalcUT(jd, sweph.SE_CHIRON)
		if err != nil {
			fmt.Printf("%04d-%02d-%02d: ERROR %v\n", d.year, d.month, d.day, err)
		} else {
			fmt.Printf("%04d-%02d-%02d: %.4f° (speed: %.6f°/d)\n", 
				d.year, d.month, d.day, res.Longitude, res.SpeedLong)
		}
	}

	// 对比Jupiter在同一时间点的位置
	fmt.Println("\n=== Jupiter位置对照 ===")
	for _, d := range testDates {
		jd := sweph.JulDay(d.year, d.month, d.day, d.hour, true)
		res, err := sweph.CalcUT(jd, sweph.SE_JUPITER)
		if err != nil {
			fmt.Printf("%04d-%02d-%02d: ERROR %v\n", d.year, d.month, d.day, err)
		} else {
			fmt.Printf("%04d-%02d-%02d: %.4f° (speed: %.6f°/d)\n", 
				d.year, d.month, d.day, res.Longitude, res.SpeedLong)
		}
	}
}