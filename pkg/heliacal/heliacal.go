// Package heliacal computes heliacal rising and setting events for the
// classical visible planets (Mercury through Saturn) using the Swiss
// Ephemeris swe_heliacal_ut function.
package heliacal

import (
	"fmt"
	"strings"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

// EventType represents a heliacal visibility event.
type EventType string

const (
	HeliacalRising  EventType = "HELIACAL_RISING"  // Star/planet rises before Sun
	HeliacalSetting EventType = "HELIACAL_SETTING" // Star/planet sets after Sun
	EveningFirst    EventType = "EVENING_FIRST"    // First visible in evening
	MorningLast     EventType = "MORNING_LAST"     // Last visible in morning
)

// HeliacalEvent represents a single heliacal visibility event.
type HeliacalEvent struct {
	Planet    models.PlanetID `json:"planet"`
	EventType EventType       `json:"event_type"`
	JDStart   float64         `json:"jd_start"`   // Start of event window
	JDOptimum float64         `json:"jd_optimum"` // Optimal observation time
	JDEnd     float64         `json:"jd_end"`     // End of event window
}

// HeliacalResult holds all heliacal events found in a search.
type HeliacalResult struct {
	Events []HeliacalEvent `json:"events"`
}

// DefaultPlanets are the five classical visible planets applicable to
// heliacal phenomena.
var DefaultPlanets = []models.PlanetID{
	models.PlanetMercury,
	models.PlanetVenus,
	models.PlanetMars,
	models.PlanetJupiter,
	models.PlanetSaturn,
}

// allEventTypes are the four heliacal event types to search for.
var allEventTypes = []struct {
	et   EventType
	code int
}{
	{HeliacalRising, sweph.SE_HELIACAL_RISING},
	{HeliacalSetting, sweph.SE_HELIACAL_SETTING},
	{EveningFirst, sweph.SE_EVENING_FIRST},
	{MorningLast, sweph.SE_MORNING_LAST},
}

// planetNameMap maps PlanetID to the object name string expected by
// swe_heliacal_ut.
var planetNameMap = map[models.PlanetID]string{
	models.PlanetMercury: "mercury",
	models.PlanetVenus:   "venus",
	models.PlanetMars:    "mars",
	models.PlanetJupiter: "jupiter",
	models.PlanetSaturn:  "saturn",
}

// planetCName returns the Swiss Ephemeris object name for a planet, or an
// error if the planet is not supported for heliacal calculations.
func planetCName(p models.PlanetID) (string, error) {
	name, ok := planetNameMap[p]
	if !ok {
		return "", fmt.Errorf("planet %s is not supported for heliacal calculations (only Mercury-Saturn)", p)
	}
	return name, nil
}

// CalcHeliacalEvents finds heliacal rising and setting events for the given
// planets within a Julian Day range [startJD, endJD]. The observer is at
// (lat, lon) in degrees and alt in metres above sea level.
func CalcHeliacalEvents(lat, lon, alt float64, startJD, endJD float64, planets []models.PlanetID) (*HeliacalResult, error) {
	if startJD >= endJD {
		return nil, fmt.Errorf("startJD (%.2f) must be before endJD (%.2f)", startJD, endJD)
	}
	if len(planets) == 0 {
		planets = DefaultPlanets
	}

	result := &HeliacalResult{}

	for _, planet := range planets {
		cName, err := planetCName(planet)
		if err != nil {
			continue // skip unsupported planets silently
		}

		for _, evtDef := range allEventTypes {
			searchJD := startJD
			// Safety limit: no more than 50 events per planet per event type
			// in any reasonable date range.
			for i := 0; i < 50; i++ {
				hr, err := sweph.HeliacalUT(searchJD, lon, lat, alt, cName, evtDef.code)
				if err != nil {
					// Some event types may not be found; stop searching
					// for this event type.
					break
				}

				// If the result is beyond our range, stop.
				if hr.JDStart > endJD {
					break
				}

				// Only include if within range.
				if hr.JDStart >= startJD {
					result.Events = append(result.Events, HeliacalEvent{
						Planet:    planet,
						EventType: evtDef.et,
						JDStart:   hr.JDStart,
						JDOptimum: hr.JDOptimum,
						JDEnd:     hr.JDEnd,
					})
				}

				// Advance past this event to search for the next one.
				// Move forward at least 30 days (heliacal events don't
				// repeat more often than that for a given planet).
				searchJD = hr.JDStart + 30
				if searchJD > endJD {
					break
				}
			}
		}
	}

	return result, nil
}

// NextHeliacalRising finds the next heliacal rising of a planet after the
// given Julian Day. Returns an error if the planet is not supported.
func NextHeliacalRising(planet models.PlanetID, lat, lon float64, startJD float64) (*HeliacalEvent, error) {
	cName, err := planetCName(planet)
	if err != nil {
		return nil, err
	}

	hr, err := sweph.HeliacalUT(startJD, lon, lat, 0, cName, sweph.SE_HELIACAL_RISING)
	if err != nil {
		return nil, fmt.Errorf("heliacal rising for %s: %w", planet, err)
	}

	return &HeliacalEvent{
		Planet:    planet,
		EventType: HeliacalRising,
		JDStart:   hr.JDStart,
		JDOptimum: hr.JDOptimum,
		JDEnd:     hr.JDEnd,
	}, nil
}

// EventTypeName returns a human-readable name for the event type.
func EventTypeName(et EventType) string {
	switch et {
	case HeliacalRising:
		return "Heliacal Rising"
	case HeliacalSetting:
		return "Heliacal Setting"
	case EveningFirst:
		return "Evening First"
	case MorningLast:
		return "Morning Last"
	default:
		return strings.ReplaceAll(string(et), "_", " ")
	}
}
