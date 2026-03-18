package export

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anthropic/swisseph-mcp/pkg/julian"
	"github.com/anthropic/swisseph-mcp/pkg/models"
)

// CSVRow represents one row of transit CSV output
type CSVRow struct {
	P1        string
	P1House   int
	Aspect    string
	P2        string
	P2House   int
	EventType string
	Type      string
	Date      string
	Time      string
	Timezone  string
	Age       float64
	Pos1Deg   float64
	Pos1Sign  string
	Pos1Dir   string
	Pos2Deg   float64
	Pos2Sign  string
	Pos2Dir   string
}

// CSVHeader returns the CSV header line
func CSVHeader() string {
	return "P1,P1_House,Aspect,P2,P2_House,EventType,Type,Date,Time,Timezone,Age,Pos1_Deg,Pos1_Sign,Pos1_Dir,Pos2_Deg,Pos2_Sign,Pos2_Dir"
}

// direction returns "Dir" or "Rx" based on retrograde status
func direction(isRetro bool) string {
	if isRetro {
		return "Rx"
	}
	return "Dir"
}

// chartPairType returns the Type column value (e.g. "Tr-Na", "Sp-Sp")
func chartPairType(evt models.TransitEvent) string {
	switch evt.EventType {
	case models.EventStation:
		return models.ChartTypeShort(evt.ChartType)
	case models.EventSignIngress:
		// Sign ingress: use chart type pair
		ct := models.ChartTypeShort(evt.ChartType)
		switch evt.ChartType {
		case models.ChartTransit:
			return ct + "-" + ct
		default:
			return ct + "-Na"
		}
	case models.EventVoidOfCourse:
		return "Tr-Tr"
	default:
		ct1 := models.ChartTypeShort(evt.ChartType)
		ct2 := models.ChartTypeShort(evt.TargetChartType)
		return ct1 + "-" + ct2
	}
}

// EventToCSVRow converts a TransitEvent to a CSVRow
func EventToCSVRow(evt models.TransitEvent, tz string) CSVRow {
	row := CSVRow{
		Timezone: tz,
		Age:      evt.Age,
	}

	// Format date/time in the specified timezone
	dtStr, err := julian.JDToDateTime(evt.JD, tz)
	if err == nil {
		t, err2 := time.Parse(time.RFC3339, dtStr)
		if err2 == nil {
			row.Date = t.Format("2006-01-02")
			row.Time = t.Format("15:04:05")
		}
	}

	row.Type = chartPairType(evt)

	switch evt.EventType {
	case models.EventStation:
		row.P1 = models.BodyDisplayName(string(evt.Planet))
		row.P1House = evt.PlanetHouse
		row.Aspect = "Station"
		row.EventType = models.EventTypeCSV(evt.EventType, evt.StationType)
		row.Pos1Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.PlanetLongitude))
		row.Pos1Sign = models.SignFromLongitude(evt.PlanetLongitude)
		row.Pos1Dir = direction(evt.IsRetrograde)

	case models.EventSignIngress:
		row.P1 = models.BodyDisplayName(string(evt.Planet))
		row.P1House = evt.PlanetHouse
		row.Aspect = "Conjunction"
		row.P2 = evt.ToSign
		row.P2House = evt.PlanetHouse
		row.EventType = "SignIngress"
		row.Pos1Deg = 0.0
		row.Pos1Sign = evt.ToSign
		row.Pos1Dir = direction(evt.IsRetrograde)
		row.Pos2Deg = 0.0
		row.Pos2Sign = evt.ToSign
		row.Pos2Dir = direction(evt.IsRetrograde)

	case models.EventVoidOfCourse:
		row.P1 = models.BodyDisplayName(string(evt.Planet))
		row.P1House = evt.PlanetHouse
		row.Aspect = evt.LastAspectType
		row.P2 = models.BodyDisplayName(evt.LastAspectTarget)
		row.P2House = evt.TargetHouse
		row.EventType = "Void"
		row.Pos1Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.PlanetLongitude))
		row.Pos1Sign = evt.PlanetSign
		row.Pos1Dir = direction(evt.IsRetrograde)
		row.Pos2Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.TargetLongitude))
		row.Pos2Sign = evt.TargetSign
		row.Pos2Dir = direction(evt.TargetIsRetrograde)

	case models.EventHouseIngress:
		row.P1 = models.BodyDisplayName(string(evt.Planet))
		row.P1House = evt.ToHouse
		row.Aspect = "HouseIngress"
		row.P2 = fmt.Sprintf("House%d", evt.ToHouse)
		row.P2House = evt.ToHouse
		row.EventType = "HouseIngress"
		row.Pos1Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.PlanetLongitude))
		row.Pos1Sign = evt.PlanetSign
		row.Pos1Dir = direction(evt.IsRetrograde)

	default:
		// Aspect events (Begin, Enter, Exact, Leave)
		row.P1 = models.BodyDisplayName(string(evt.Planet))
		row.P1House = evt.PlanetHouse
		row.Aspect = models.AspectCSVName(evt.AspectType)
		row.P2 = models.BodyDisplayName(evt.Target)
		row.P2House = evt.TargetHouse
		row.EventType = models.EventTypeCSV(evt.EventType, "")
		row.Pos1Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.PlanetLongitude))
		row.Pos1Sign = evt.PlanetSign
		row.Pos1Dir = direction(evt.IsRetrograde)
		row.Pos2Deg = models.FormatSignDegreeCSV(models.SignDegreeFromLongitude(evt.TargetLongitude))
		row.Pos2Sign = evt.TargetSign
		row.Pos2Dir = direction(evt.TargetIsRetrograde)
	}

	return row
}

// formatDeg formats a degree value for CSV output (Solar Fire style: trim trailing zeros, keep at least one decimal)
func formatDeg(d float64) string {
	s := fmt.Sprintf("%.3f", d)
	// Trim trailing zeros after decimal point, but keep at least X.Y format
	s = strings.TrimRight(s, "0")
	if strings.HasSuffix(s, ".") {
		s += "0"
	}
	return s
}

// CSVRowToString formats a CSVRow as a comma-separated string
func CSVRowToString(r CSVRow) string {
	// For station events, P2 fields are empty
	if r.Aspect == "Station" {
		return fmt.Sprintf("%s,%d,%s,,,%s,%s,%s,%s,%s,%.3f,%s,%s,%s,,,",
			r.P1, r.P1House, r.Aspect,
			r.EventType, r.Type, r.Date, r.Time, r.Timezone,
			r.Age, formatDeg(r.Pos1Deg), r.Pos1Sign, r.Pos1Dir)
	}

	return fmt.Sprintf("%s,%d,%s,%s,%d,%s,%s,%s,%s,%s,%.3f,%s,%s,%s,%s,%s,%s",
		r.P1, r.P1House, r.Aspect, r.P2, r.P2House,
		r.EventType, r.Type, r.Date, r.Time, r.Timezone,
		r.Age, formatDeg(r.Pos1Deg), r.Pos1Sign, r.Pos1Dir,
		formatDeg(r.Pos2Deg), r.Pos2Sign, r.Pos2Dir)
}

// EventsToCSV converts a list of transit events to CSV string.
// tz is the Go timezone location for date/time conversion.
// tzLabel is the display label (e.g. "AWST"); if empty, uses tz.
func EventsToCSV(events []models.TransitEvent, tz string, tzLabel string) string {
	if tzLabel == "" {
		tzLabel = tz
	}
	var sb strings.Builder
	sb.WriteString(CSVHeader())
	sb.WriteString("\n")
	for _, evt := range events {
		row := EventToCSVRow(evt, tz)
		row.Timezone = tzLabel
		sb.WriteString(CSVRowToString(row))
		sb.WriteString("\n")
	}
	return sb.String()
}

// EventsToJSON converts a list of transit events to pretty JSON string
func EventsToJSON(events []models.TransitEvent) (string, error) {
	data, err := json.MarshalIndent(map[string]interface{}{
		"events": events,
	}, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
