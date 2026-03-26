package main

import (
	"fmt"
)

// 测试不同时间容差级别的匹配效果
func main() {
	fmt.Println("=== 不同时间容差测试 ===")
	
	// 基于我们观察到的实际差异（14-29分钟）
	testTolerances := []int{900, 1800, 2700, 3600, 5400, 7200} // 15min, 30min, 45min, 1h, 1.5h, 2h
	
	fmt.Println("基于实际观察的时间差异（14-29分钟）:")
	fmt.Println("- 平均绝对差异: 17.8分钟")
	fmt.Println("- 差异范围: -16到+29分钟")
	fmt.Println("- 标准差: 19.7分钟")
	fmt.Println()
	
	fmt.Println("建议的容差策略:")
	fmt.Println("1. 对于正常transit事件 (Tr-Na, Tr-Tr): 30分钟容差")
	fmt.Println("2. 对于进展事件 (Sp-*, Sa-*): 60分钟容差")
	fmt.Println("3. 可以进一步细分为:")
	fmt.Println("   - Tr-Na: 20分钟")
	fmt.Println("   - Tr-Tr: 25分钟") 
	fmt.Println("   - Tr-Sp: 35分钟")
	fmt.Println("   - Tr-Sa: 30分钟")
	fmt.Println("   - Sp-Na/Sp-Sp/Sa-Na: 45分钟")
	fmt.Println()
	
	fmt.Println("容差测试方案:")
	for i, tol := range testTolerances {
		minutes := tol / 60
		fmt.Printf("测试 %d: ±%d分钟 (%d秒)\n", i+1, minutes, tol)
	}
	
	fmt.Println()
	fmt.Println("预期结果:")
	fmt.Println("- 容差过小(<15分钟): 匹配数会显著减少")
	fmt.Println("- 容差适中(20-30分钟): 应该能匹配大部分事件")
	fmt.Println("- 容差过大(>45分钟): 虽然匹配数增加，但精度下降")
}