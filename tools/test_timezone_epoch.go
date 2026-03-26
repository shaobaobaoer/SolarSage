package main

import (
	"fmt"
	"time"

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
	
	// 获取出生地点的时区信息
	// 假设出生在北京 (UTC+8)
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai")
	utcTime, _ := time.Parse(time.RFC3339, birthDatetime)
	beijingTime := utcTime.In(beijingLoc)
	
	fmt.Printf("北京时间: %s\n", beijingTime.Format("2006-01-02 15:04:05 MST"))
	
	// 计算北京时区的午夜JD
	beijingMidnight := time.Date(beijingTime.Year(), beijingTime.Month(), beijingTime.Day(), 
		0, 0, 0, 0, beijingLoc)
	utcMidnight := beijingMidnight.UTC()
	
	year, month, day := utcMidnight.Year(), int(utcMidnight.Month()), utcMidnight.Day()
	hour := float64(utcMidnight.Hour()) + float64(utcMidnight.Minute())/60.0
	midnightJD := sweph.JulDay(year, month, day, hour, true)
	
	fmt.Printf("北京午夜UTC时间: %s\n", utcMidnight.Format("2006-01-02 15:04:05"))
	fmt.Printf("北京午夜JD: %.6f\n", midnightJD)
	fmt.Printf("与出生JD差异: %.6f 天 = %.1f 小时\n", midnightJD-natalJD, (midnightJD-natalJD)*24)
	
	// 测试不同的时间点作为进展epoch
	testPoints := []struct {
		name string
		hour float64
	}{
		{"出生时刻", 9.6},           // 09:36
		{"UTC午夜", 0.0},            // 00:00 UTC
		{"UTC中午", 12.0},           // 12:00 UTC
		{"北京午夜", hour},          // 北京时间00:00对应的UTC时间
		{"出生后5.32小时", 9.6 + 5.32}, // 我们观察到的目标偏移
	}
	
	testJD := natalJD + 10000 // 测试10000天后的进展
	
	fmt.Println("\n=== 不同时点作为进展epoch的对比 ===")
	currentMethod := natalJD + (testJD-natalJD)/JulianYear
	fmt.Printf("当前方法 (出生时刻): %.6f\n", currentMethod)
	
	for _, point := range testPoints {
		epochJD := sweph.JulDay(year, month, day, point.hour, true)
		progJD := epochJD + (testJD-epochJD)/JulianYear
		diffHours := (progJD - currentMethod) * 24
		fmt.Printf("%s (%.1f:00 UTC): %.6f (差异: %+.1f小时)\n", 
			point.name, point.hour, progJD, diffHours)
	}
	
	// 分析我们观察到的最大差异案例
	fmt.Println("\n=== 分析最大差异案例 ===")
	
	// 选取几个典型的最大差异事件进行反向计算
	cases := []struct {
		description string
		sfDateTime  string    // Solar Fire时间
		compDateTime string   // 我们的计算时间
		expectedDiffHours float64 // 期望的小时差异
	}{
		{
			description: "Neptune Opposition ASC Sp-Sp",
			sfDateTime: "2026-05-22 12:07:13",
			compDateTime: "2026-05-21 17:12:13",
			expectedDiffHours: -18.9, // 实际差异约-18.9小时
		},
		{
			description: "Uranus Square Uranus Tr-Sa", 
			sfDateTime: "2026-09-10 12:53:48",
			compDateTime: "2026-09-11 02:27:15",
			expectedDiffHours: +13.6, // 实际差异约+13.6小时
		},
	}
	
	for _, c := range cases {
		fmt.Printf("\n案例: %s\n", c.description)
		
		// 转换时间为JD
		sfJD := dateTimeToJD(c.sfDateTime)
		compJD := dateTimeToJD(c.compDateTime)
		
		actualDiffHours := (sfJD - compJD) * 24
		fmt.Printf("  SF时间: %s (JD: %.6f)\n", c.sfDateTime, sfJD)
		fmt.Printf("  我们时间: %s (JD: %.6f)\n", c.compDateTime, compJD)
		fmt.Printf("  实际差异: %.1f 小时\n", actualDiffHours)
		fmt.Printf("  期望差异: %.1f 小时\n", c.expectedDiffHours)
		
		// 反向计算Solar Fire使用的进展因子
		// 如果 SF_time = epoch + (event_jd - epoch)/factor
		// 则 factor = (event_jd - epoch)/(SF_time - epoch)
		
		eventJD := compJD // 使用我们的计算时间作为基准事件JD
		epochJD := natalJD // 假设使用natalJD作为epoch
		
		// 计算Solar Fire使用的进展因子
		sfFactor := (eventJD - epochJD) / (sfJD - epochJD)
		currentFactor := JulianYear
		
		fmt.Printf("  Solar Fire进展因子: %.6f\n", sfFactor)
		fmt.Printf("  我们使用的因子: %.6f\n", currentFactor)
		fmt.Printf("  因子差异: %.6f\n", sfFactor-currentFactor)
	}
}

func dateTimeToJD(dateTimeStr string) float64 {
	// 假设输入格式为 "YYYY-MM-DD HH:MM:SS"
	t, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		panic(err)
	}
	
	utc := t.UTC()
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	hour := float64(utc.Hour()) + float64(utc.Minute())/60.0 + float64(utc.Second())/3600.0
	
	return sweph.JulDay(year, month, day, hour, true)
}