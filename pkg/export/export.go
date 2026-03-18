package export

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/julian"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
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

// formatDeg formats a degree value for CSV output (trim trailing zeros, keep at least one decimal)
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

// === Chart Export (CSV/JSON for any chart type) ===

// ChartToCSV exports a natal chart as CSV with planet positions
func ChartToCSV(chart *models.ChartInfo) string {
	var sb strings.Builder
	sb.WriteString("Planet,Longitude,Sign,SignDegree,House,Speed,Retrograde,Glyph\n")
	for _, p := range chart.Planets {
		sb.WriteString(fmt.Sprintf("%s,%.4f,%s,%.4f,%d,%.6f,%t,%s\n",
			models.BodyDisplayName(string(p.PlanetID)),
			p.Longitude, p.Sign, p.SignDegree,
			p.House, p.Speed, p.IsRetrograde,
			models.PlanetGlyph(p.PlanetID),
		))
	}
	return sb.String()
}

// ChartToJSON exports a natal chart as pretty JSON
func ChartToJSON(chart *models.ChartInfo) (string, error) {
	data, err := json.MarshalIndent(chart, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// AspectsToCSV exports aspects as CSV
func AspectsToCSV(aspects []models.AspectInfo) string {
	var sb strings.Builder
	sb.WriteString("PlanetA,PlanetB,Aspect,Angle,ActualAngle,Orb,Applying,Glyph\n")
	for _, a := range aspects {
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%.1f,%.4f,%.4f,%t,%s\n",
			models.BodyDisplayName(a.PlanetA),
			models.BodyDisplayName(a.PlanetB),
			a.AspectType, a.AspectAngle, a.ActualAngle,
			a.Orb, a.IsApplying,
			models.AspectGlyph(a.AspectType),
		))
	}
	return sb.String()
}

// CrossAspectsToCSV exports cross-aspects as CSV
func CrossAspectsToCSV(aspects []models.CrossAspectInfo) string {
	var sb strings.Builder
	sb.WriteString("InnerBody,OuterBody,Aspect,Angle,ActualAngle,Orb,Applying,Glyph\n")
	for _, a := range aspects {
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%.1f,%.4f,%.4f,%t,%s\n",
			models.BodyDisplayName(a.InnerBody),
			models.BodyDisplayName(a.OuterBody),
			a.AspectType, a.AspectAngle, a.ActualAngle,
			a.Orb, a.IsApplying,
			models.AspectGlyph(a.AspectType),
		))
	}
	return sb.String()
}

// HousesToCSV exports house cusps as CSV
func HousesToCSV(houses []float64, angles models.AnglesInfo) string {
	var sb strings.Builder
	sb.WriteString("House,Longitude,Sign,SignDegree\n")
	for i, lon := range houses {
		sb.WriteString(fmt.Sprintf("%d,%.4f,%s,%.4f\n",
			i+1, lon,
			models.SignFromLongitude(lon),
			models.SignDegreeFromLongitude(lon),
		))
	}
	sb.WriteString(fmt.Sprintf("ASC,%.4f,%s,%.4f\n", angles.ASC, models.SignFromLongitude(angles.ASC), models.SignDegreeFromLongitude(angles.ASC)))
	sb.WriteString(fmt.Sprintf("MC,%.4f,%s,%.4f\n", angles.MC, models.SignFromLongitude(angles.MC), models.SignDegreeFromLongitude(angles.MC)))
	return sb.String()
}

// PositionsToCSV exports a slice of planet positions as CSV (generic, works for progressions/solar arc/etc.)
func PositionsToCSV(positions []models.PlanetPosition) string {
	var sb strings.Builder
	sb.WriteString("Planet,Longitude,Sign,SignDegree,Speed,Retrograde,Glyph\n")
	for _, p := range positions {
		sb.WriteString(fmt.Sprintf("%s,%.4f,%s,%.4f,%.6f,%t,%s\n",
			models.BodyDisplayName(string(p.PlanetID)),
			p.Longitude, p.Sign, p.SignDegree,
			p.Speed, p.IsRetrograde,
			models.PlanetGlyph(p.PlanetID),
		))
	}
	return sb.String()
}

// DignityToCSV exports dignity analysis as CSV
func DignityToCSV(dignities interface{}) string {
	data, _ := json.Marshal(dignities)
	var items []struct {
		PlanetID    string   `json:"planet_id"`
		Sign        string   `json:"sign"`
		Score       int      `json:"score"`
		Ruler       string   `json:"ruler"`
		Exalted     bool     `json:"exalted"`
		InDetriment bool     `json:"in_detriment"`
		InFall      bool     `json:"in_fall"`
		Dignities   []string `json:"dignities"`
	}
	json.Unmarshal(data, &items)

	var sb strings.Builder
	sb.WriteString("Planet,Sign,Score,Ruler,Exalted,Detriment,Fall,Dignities,Glyph\n")
	for _, d := range items {
		sb.WriteString(fmt.Sprintf("%s,%s,%d,%s,%t,%t,%t,%s,%s\n",
			models.BodyDisplayName(d.PlanetID), d.Sign, d.Score,
			models.BodyDisplayName(d.Ruler),
			d.Exalted, d.InDetriment, d.InFall,
			strings.Join(d.Dignities, ";"),
			models.PlanetGlyph(models.PlanetID(d.PlanetID)),
		))
	}
	return sb.String()
}

// LotsToCSV exports Arabic lots as CSV
func LotsToCSV(lotsData interface{}) string {
	data, _ := json.Marshal(lotsData)
	var items []struct {
		Name      string  `json:"name"`
		Longitude float64 `json:"longitude"`
		Sign      string  `json:"sign"`
		SignDeg   float64 `json:"sign_degree"`
		Formula   string  `json:"formula"`
	}
	json.Unmarshal(data, &items)

	var sb strings.Builder
	sb.WriteString("Lot,Longitude,Sign,SignDegree,Formula\n")
	for _, l := range items {
		sb.WriteString(fmt.Sprintf("%s,%.4f,%s,%.4f,%s\n",
			l.Name, l.Longitude, l.Sign, l.SignDeg, l.Formula))
	}
	return sb.String()
}

// EclipsesToCSV exports eclipses as CSV
func EclipsesToCSV(eclipses interface{}) string {
	data, _ := json.Marshal(eclipses)
	var items []struct {
		Type    string  `json:"type"`
		JD      float64 `json:"jd"`
		MoonLon float64 `json:"moon_longitude"`
		MoonSign string `json:"moon_sign"`
		SunSign  string `json:"sun_sign"`
		MoonLat  float64 `json:"moon_latitude"`
	}
	json.Unmarshal(data, &items)

	var sb strings.Builder
	sb.WriteString("Type,JD,MoonLongitude,MoonSign,SunSign,MoonLatitude\n")
	for _, e := range items {
		sb.WriteString(fmt.Sprintf("%s,%.6f,%.4f,%s,%s,%.4f\n",
			e.Type, e.JD, e.MoonLon, e.MoonSign, e.SunSign, e.MoonLat))
	}
	return sb.String()
}

// LunarPhasesToCSV exports lunar phases as CSV
func LunarPhasesToCSV(phases interface{}) string {
	data, _ := json.Marshal(phases)
	var items []struct {
		Phase    string  `json:"phase"`
		JD       float64 `json:"jd"`
		MoonLon  float64 `json:"moon_longitude"`
		MoonSign string  `json:"moon_sign"`
		SunSign  string  `json:"sun_sign"`
	}
	json.Unmarshal(data, &items)

	var sb strings.Builder
	sb.WriteString("Phase,JD,MoonLongitude,MoonSign,SunSign\n")
	for _, p := range items {
		sb.WriteString(fmt.Sprintf("%s,%.6f,%.4f,%s,%s\n",
			p.Phase, p.JD, p.MoonLon, p.MoonSign, p.SunSign))
	}
	return sb.String()
}

// ToJSON is a generic helper that exports any value as pretty JSON
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
