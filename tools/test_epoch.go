package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/julian"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

const JulianYear = 365.25

func main() {
	// 测试案例：1997年12月18日 09:36:00 UTC 出生
	birthDatetime := "1997-12-18T09:36:00Z"
	
	// 转换为JD
	jdResult, err := julian.DateTimeToJD(birthDatetime, models.CalendarGregorian)
	if err != nil {
		panic(err)
	}
	
	natalJD := jdResult.JDUT
	fmt.Printf("出生JD (UT): %.6f\n", natalJD)
	fmt.Printf("出生时间: %s\n", birthDatetime)
	
	// 计算出生当天的中午JD
	year, month, day, _ := sweph.RevJul(natalJD, true)
	noonJD := sweph.JulDay(year, month, day, 12.0, true) // 当天中午12:00 UT
	fmt.Printf("出生当天中午JD: %.6f\n", noonJD)
	fmt.Printf("差异: %.6f 天 = %.1f 小时\n", noonJD-natalJD, (noonJD-natalJD)*24)
	
	// 测试不同的进展计算方法
	testJD := natalJD + 10000 // 测试10000天后的进展
	
	fmt.Println("\n=== 不同进展方法对比 ===")
	
	// 方法1: 当前实现 - 使用出生时刻
	method1 := natalJD + (testJD-natalJD)/JulianYear
	fmt.Printf("方法1 (出生时刻): %.6f\n", method1)
	
	// 方法2: 使用出生当天中午
	method2 := noonJD + (testJD-noonJD)/JulianYear
	fmt.Printf("方法2 (当天中午): %.6f\n", method2)
	
	// 方法3: 使用前一天中午
	prevNoonJD := sweph.JulDay(year, month, day-1, 12.0, true)
	method3 := prevNoonJD + (testJD-prevNoonJD)/JulianYear
	fmt.Printf("方法3 (前一天中午): %.6f\n", method3)
	
	// 方法4: 使用后一天中午
	nextNoonJD := sweph.JulDay(year, month, day+1, 12.0, true)
	method4 := nextNoonJD + (testJD-nextNoonJD)/JulianYear
	fmt.Printf("方法4 (后一天中午): %.6f\n", method4)
	
	// 计算各方法之间的差异
	fmt.Println("\n=== 差异分析 ===")
	fmt.Printf("方法2 vs 方法1: %.6f 天 = %.1f 小时\n", method2-method1, (method2-method1)*24)
	fmt.Printf("方法3 vs 方法1: %.6f 天 = %.1f 小时\n", method3-method1, (method3-method1)*24)
	fmt.Printf("方法4 vs 方法1: %.6f 天 = %.1f 小时\n", method4-method1, (method4-method1)*24)
	
	// 看看哪个最接近我们观察到的5.3小时偏差
	targetOffsetHours := 5.32
	targetOffsetDays := targetOffsetHours / 24.0
	
	// 我们需要找到使得进展JD等于目标JD的epoch
	// targetJD = epochJD + (testJD - epochJD)/JulianYear
	// targetJD = epochJD*(1 - 1/JulianYear) + testJD/JulianYear
	// epochJD = (targetJD - testJD/JulianYear) / (1 - 1/JulianYear)
	
	targetJD := natalJD + targetOffsetDays
	epochJD := (targetJD - testJD/JulianYear) / (1 - 1/JulianYear)
	
	fmt.Printf("\n=== 反向计算Solar Fire的epoch ===")
	fmt.Printf("\n目标JD (natalJD + %.2f小时): %.6f", targetOffsetHours, targetJD)
	fmt.Printf("\n计算得出的Solar Fire epoch JD: %.6f", epochJD)
	
	// 转换回时间和日期
	epochYear, epochMonth, epochDay, epochHour := sweph.RevJul(epochJD, true)
	fmt.Printf("\nSolar Fire epoch时间: %04d-%02d-%02d %.1f:00 UTC", 
		epochYear, epochMonth, epochDay, epochHour)
	
	// 计算与出生时间的差异
	hourDiff := epochHour - 9.6 // 9:36 = 9.6小时
	fmt.Printf("\n与出生时间(09:36)的差异: %.1f 小时", hourDiff)
	
	// 验证这个epoch是否正确
	sfMethod := epochJD + (testJD-epochJD)/JulianYear
	fmt.Printf("\n使用SF epoch计算的进展JD: %.6f", sfMethod)
	fmt.Printf("\n与目标JD的差异: %.6f 天 = %.1f 秒", sfMethod-targetJD, (sfMethod-targetJD)*86400)
}