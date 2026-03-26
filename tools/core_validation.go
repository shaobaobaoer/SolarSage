package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/progressions"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== 核心进展计算验证 ===")
	
	// 测试数据
	natalJD := 2450800.900000 // 1997-12-18 09:36:00 UTC
	testJD := 2461193.000000  // 2026-06-01 12:00:00 UTC
	
	fmt.Printf("出生JD: %.6f\n", natalJD)
	fmt.Printf("测试JD: %.6f\n", testJD)
	fmt.Printf("时间跨度: %.1f天\n", testJD-natalJD)
	
	// 使用修正后的算法
	progJD := progressions.SecondaryProgressionJD(natalJD, testJD)
	fmt.Printf("Solar Fire进展JD: %.6f\n", progJD)
	
	// 计算太阳位置验证
	sunPos, err := sweph.CalcUT(progJD, sweph.SE_SUN)
	if err != nil {
		fmt.Printf("计算错误: %v\n", err)
		return
	}
	
	fmt.Printf("进展太阳位置: %.4f°\n", sunPos.Longitude)
	
	// 计算年龄
	age := progressions.Age(natalJD, testJD)
	fmt.Printf("Solar Fire年龄: %.2f岁\n", age)
	
	fmt.Println("核心功能验证完成!")
}