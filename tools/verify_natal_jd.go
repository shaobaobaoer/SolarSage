package main

import (
	"fmt"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	fmt.Println("=== Natal JD 计算验证 ===")
	
	// 原始出生信息
	// Birth: 1997-12-18 17:36:00 AWST (UTC+8)
	// 即 UTC 09:36:00
	
	birthUTC := time.Date(1997, 12, 18, 9, 36, 0, 0, time.UTC)
	fmt.Printf("Birth UTC: %s\n", birthUTC.Format("2006-01-02 15:04:05"))
	
	// 方法1: 直接用sweph.JulDay
	year, month, day := birthUTC.Date()
	hour := float64(birthUTC.Hour()) + float64(birthUTC.Minute())/60 + float64(birthUTC.Second())/3600
	jd1 := sweph.JulDay(year, int(month), day, hour, true)
	fmt.Printf("Method 1 (sweph.JulDay): %.6f\n", jd1)
	
	// 方法2: 用time.Unix转为儒略日
	// Unix时间戳从1970-01-01 00:00:00 UTC开始
	// J2000.0 = 2451545.0 = 2000-01-01 12:00:00 UTC
	// 1970-01-01 00:00:00 UTC = JD 2440587.5
	
	unixEpochJD := 2440587.5
	secondsPerDay := 86400.0
	daysSinceUnixEpoch := float64(birthUTC.Unix()) / secondsPerDay
	jd2 := unixEpochJD + daysSinceUnixEpoch
	fmt.Printf("Method 2 (Unix timestamp): %.6f\n", jd2)
	
	// 差异
	diff := jd1 - jd2
	fmt.Printf("Difference: %.9f days (%.3f seconds)\n", diff, diff*86400)
	
	// SF使用的JD值
	sfJD := 2450800.900009
	fmt.Printf("SF JD: %.6f\n", sfJD)
	fmt.Printf("Diff from Method 1: %.9f days (%.3f seconds)\n", sfJD-jd1, (sfJD-jd1)*86400)
	fmt.Printf("Diff from Method 2: %.9f days (%.3f seconds)\n", sfJD-jd2, (sfJD-jd2)*86400)
	
	// 反向验证SF的JD
	year2, month2, day2, hour2 := sweph.RevJul(sfJD, true)
	fmt.Printf("\nReverse calculation from SF JD:\n")
	fmt.Printf("  Date: %04d-%02d-%02d %02d:%02d:%02d\n", 
		year2, month2, day2, 
		int(hour2), int((hour2-float64(int(hour2)))*60), 
		int(((hour2-float64(int(hour2)))*60-float64(int((hour2-float64(int(hour2)))*60)))*60))
	
	// 检查是否考虑了ΔT
	// ΔT = TT - UT1
	deltaT := sweph.DeltaT(sfJD)
	fmt.Printf("ΔT at natal JD: %.3f seconds\n", deltaT*86400)
	
	// JDE = JD_UT + ΔT/86400
	jde := sfJD + deltaT
	fmt.Printf("JDE (Ephemeris Time): %.6f\n", jde)
	
	// 用JDE重新计算
	year3, month3, day3, hour3 := sweph.RevJul(jde, true)
	fmt.Printf("Date from JDE: %04d-%02d-%02d %02d:%02d:%02d\n",
		year3, month3, day3,
		int(hour3), int((hour3-float64(int(hour3)))*60),
		int(((hour3-float64(int(hour3)))*60-float64(int((hour3-float64(int(hour3)))*60)))*60))
}