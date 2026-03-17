package julian

import (
	"fmt"
	"time"

	"github.com/anthropic/swisseph-mcp/pkg/models"
	"github.com/anthropic/swisseph-mcp/pkg/sweph"
)

// DateTimeToJD converts an ISO 8601 datetime string to Julian Day (UT and TT)
func DateTimeToJD(datetime string, calendar models.CalendarType) (*models.JulianDayResult, error) {
	t, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return nil, fmt.Errorf("invalid datetime format: %w", err)
	}

	// Convert to UTC
	utc := t.UTC()
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	hour := float64(utc.Hour()) + float64(utc.Minute())/60.0 + float64(utc.Second())/3600.0 + float64(utc.Nanosecond())/3600000000000.0

	gregorian := calendar != models.CalendarJulian
	jdUT := sweph.JulDay(year, month, day, hour, gregorian)
	deltaT := sweph.DeltaT(jdUT)
	jdTT := jdUT + deltaT

	return &models.JulianDayResult{
		JDUT: jdUT,
		JDTT: jdTT,
	}, nil
}

// JDToDateTime converts Julian Day back to ISO 8601 datetime string
func JDToDateTime(jd float64, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", timezone, err)
	}

	year, month, day, hour := sweph.RevJul(jd, true)

	hours := int(hour)
	minutesFrac := (hour - float64(hours)) * 60.0
	minutes := int(minutesFrac)
	secondsFrac := (minutesFrac - float64(minutes)) * 60.0
	seconds := int(secondsFrac)
	nanos := int((secondsFrac - float64(seconds)) * 1e9)

	t := time.Date(year, time.Month(month), day, hours, minutes, seconds, nanos, time.UTC)
	t = t.In(loc)

	return t.Format(time.RFC3339), nil
}
