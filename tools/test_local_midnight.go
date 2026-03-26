package main

import (
	"fmt"
	"time"
)

const (
	JulianYear = 365.25
	natalJD    = 2450800.900000 // 1997-12-18 09:36:00 UTC
)

func main() {
	// 假设出生在北京 (东八区)
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai")
	
	// 出生的UTC时间
	birthTimeUTC, _ := time.Parse("2006-01-02T15:04:05Z", "1997-12-18T09:36:00Z")
	
	// 转换为北京时间
	birthTimeBeijing := birthTimeUTC.In(beijingLoc)
	fmt.Printf("出生UTC时间: %s\n", birthTimeUTC.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("出生北京时间: %s\n", birthTimeBeijing.Format("2006-01-02 15:04:05 MST"))
	
	// 计算出生当天北京时间午夜
	beijingMidnight := time.Date(birthTimeBeijing.Year(), birthTimeBeijing.Month(), birthTimeBeijing.Day(),
		0, 0, 0, 0, beijingLoc)
	
	// 转换为UTC
	utcMidnight := beijingMidnight.UTC()
	fmt.Printf("北京午夜UTC时间: %s\n", utcMidnight.Format("2006-01-02 15:04:05 MST"))
	
	// 计算JD
	year, month, day := utcMidnight.Year(), int(utcMidnight.Month()), utcMidnight.Day()
	hour := float64(utcMidnight.Hour()) + float64(utcMidnight.Minute())/60.0 + float64(utcMidnight.Second())/3600.0
	
	midnightJD := julianDay(year, month, day, hour)
	fmt.Printf("北京午夜JD: %.6f\n", midnightJD)
	fmt.Printf("与出生JD差异: %.6f 天 = %.1f 小时\n", midnightJD-natalJD, (midnightJD-natalJD)*24)
	
	// 测试使用这个午夜作为进展epoch的效果
	testJD := natalJD + 10000
	
	fmt.Println("\n=== 使用北京午夜作为进展epoch ===")
	
	// 当前方法 (出生时刻)
	currentMethod := natalJD + (testJD-natalJD)/JulianYear
	fmt.Printf("当前方法 (出生时刻): %.6f\n", currentMethod)
	
	// 新方法 (北京午夜)
	newMethod := midnightJD + (testJD-midnightJD)/JulianYear
	fmt.Printf("新方法 (北京午夜): %.6f\n", newMethod)
	
	diffHours := (newMethod - currentMethod) * 24
	fmt.Printf("时间差异: %.1f 小时\n", diffHours)
	
	// 检查这是否接近我们观察到的差异
	observedDiffs := []float64{1.1, 5.3, 13.6, 18.9}
	fmt.Println("\n=== 与观察差异对比 ===")
	for _, obs := range observedDiffs {
		match := "✗"
		if abs(abs(diffHours)-obs) < 1.0 { // 1小时内认为匹配
			match = "✓"
		}
		fmt.Printf("%.1f小时: %s\n", obs, match)
	}
	
	// 尝试其他可能的时区
	fmt.Println("\n=== 测试其他时区午夜 ===")
	
	timezones := []struct {
		name string
		tz   string
	}{
		{"UTC", "UTC"},
		{"洛杉矶", "America/Los_Angeles"},
		{"纽约", "America/New_York"},
		{"伦敦", "Europe/London"},
		{"东京", "Asia/Tokyo"},
		{"悉尼", "Australia/Sydney"},
	}
	
	for _, tz := range timezones {
		loc, err := time.LoadLocation(tz.tz)
		if err != nil {
			continue
		}
		
		localBirth := birthTimeUTC.In(loc)
		localMidnight := time.Date(localBirth.Year(), localBirth.Month(), localBirth.Day(),
			0, 0, 0, 0, loc)
		utcMidnight := localMidnight.UTC()
		
		year, month, day := utcMidnight.Year(), int(utcMidnight.Month()), utcMidnight.Day()
		hour := float64(utcMidnight.Hour()) + float64(utcMidnight.Minute())/60.0
		localMidnightJD := julianDay(year, month, day, hour)
		
		localMethod := localMidnightJD + (testJD-localMidnightJD)/JulianYear
		localDiffHours := (localMethod - currentMethod) * 24
		
		matchSymbol := "✗"
		for _, obs := range observedDiffs {
			if abs(abs(localDiffHours)-obs) < 1.0 {
				matchSymbol = "✓"
				break
			}
		}
		
		fmt.Printf("%-10s: %+6.1f小时 %s\n", tz.name, localDiffHours, matchSymbol)
	}
}

func julianDay(year, month, day int, hour float64) float64 {
	// 简化的儒略日计算（适用于格里高利历）
	if month <= 2 {
		year--
		month += 12
	}
	
	a := year / 100
	b := 2 - a + a/4
	
	jd := float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		float64(day) + float64(b) - 1524.5 +
		hour/24.0
	
	return jd
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}