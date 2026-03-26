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
	// 出生的UTC时间
	birthTimeUTC, _ := time.Parse("2006-01-02T15:04:05Z", "1997-12-18T09:36:00Z")
	
	// 测试各种可能的本地时间参考点
	testTimezones := []struct {
		name string
		tz   string
	}{
		{"洛杉矶", "America/Los_Angeles"},
		{"纽约", "America/New_York"},
		{"伦敦", "Europe/London"},
		{"北京", "Asia/Shanghai"},
		{"东京", "Asia/Tokyo"},
		{"悉尼", "Australia/Sydney"},
	}
	
	fmt.Println("=== 不同时区的标准时间点测试 ===")
	fmt.Printf("%-10s %-15s %-15s %-10s %-10s\n", "时区", "标准时间点", "对应UTC", "时间差异", "匹配度")
	fmt.Println(strings.Repeat("-", 70))
	
	testPoints := []struct {
		name string
		hour int
		min  int
	}{
		{"午夜", 0, 0},
		{"凌晨6点", 6, 0},
		{"上午9点", 9, 0},
		{"中午", 12, 0},
		{"下午3点", 15, 0},
		{"傍晚6点", 18, 0},
		{"晚上9点", 21, 0},
	}
	
	observedDiffs := []float64{1.1, 5.3, 13.6, 18.9}
	
	for _, tz := range testTimezones {
		loc, err := time.LoadLocation(tz.tz)
		if err != nil {
			continue
		}
		
		localBirth := birthTimeUTC.In(loc)
		
		for _, point := range testPoints {
			// 创建该时区该时间点的datetime
			testTime := time.Date(localBirth.Year(), localBirth.Month(), localBirth.Day(),
				point.hour, point.min, 0, 0, loc)
			utcTime := testTime.UTC()
			
			// 转换为JD
			year, month, day := utcTime.Year(), int(utcTime.Month()), utcTime.Day()
			hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60.0
			epochJD := julianDay(year, month, day, hour)
			
			// 计算进展差异
			testJD := natalJD + 10000
			currentMethod := natalJD + (testJD-natalJD)/JulianYear
			newMethod := epochJD + (testJD-epochJD)/JulianYear
			diffHours := (newMethod - currentMethod) * 24
			
			// 检查匹配度
			matchLevel := 0
			for _, obs := range observedDiffs {
				if abs(abs(diffHours)-obs) < 1.0 {
					matchLevel++
				}
			}
			
			matchStr := ""
			switch matchLevel {
			case 4:
				matchStr = "★★★★"
			case 3:
				matchStr = "★★★☆"
			case 2:
				matchStr = "★★☆☆"
			case 1:
				matchStr = "★☆☆☆"
			default:
				matchStr = "☆☆☆☆"
			}
			
			fmt.Printf("%-10s %-15s %-15s %+7.1f   %s\n",
				tz.name,
				point.name,
				utcTime.Format("15:04"),
				diffHours,
				matchStr)
		}
		fmt.Println()
	}
	
	// 特别关注匹配度高的组合
	fmt.Println("=== 高匹配度组合分析 ===")
	for _, tz := range testTimezones {
		loc, _ := time.LoadLocation(tz.tz)
		localBirth := birthTimeUTC.In(loc)
		
		// 测试出生时间前后的几个小时
		for offset := -6; offset <= 6; offset++ {
			testHour := (localBirth.Hour() + offset + 24) % 24
			testTime := time.Date(localBirth.Year(), localBirth.Month(), localBirth.Day(),
				testHour, 0, 0, 0, loc)
			utcTime := testTime.UTC()
			
			year, month, day := utcTime.Year(), int(utcTime.Month()), utcTime.Day()
			hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60.0
			epochJD := julianDay(year, month, day, hour)
			
			testJD := natalJD + 10000
			currentMethod := natalJD + (testJD-natalJD)/JulianYear
			newMethod := epochJD + (testJD-epochJD)/JulianYear
			diffHours := (newMethod - currentMethod) * 24
			
			// 统计匹配数量
			matchCount := 0
			for _, obs := range observedDiffs {
				if abs(abs(diffHours)-obs) < 1.0 {
					matchCount++
				}
			}
			
			if matchCount >= 2 {
				fmt.Printf("%-10s %02d:00本地时间 %+7.1f小时 (%d个匹配)\n",
					tz.name, testHour, diffHours, matchCount)
			}
		}
	}
}

func julianDay(year, month, day int, hour float64) float64 {
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

var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}