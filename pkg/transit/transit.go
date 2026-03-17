package transit

import (
	"math"
	"sort"

	"github.com/anthropic/swisseph-mcp/pkg/chart"
	"github.com/anthropic/swisseph-mcp/pkg/models"
)

const (
	dayStep       = 1.0         // 1 day coarse scan step
	bisectEps     = 1.0 / 86400 // ~1 second precision
	fineStep      = 0.5         // fine scan step for RQ2
)

// StationInfo represents a retrograde/direct station
type StationInfo struct {
	JD          float64
	IsDirecting bool // true = station direct (retro -> direct)
}

// MonoInterval represents a monotonic longitude interval (no station inside)
type MonoInterval struct {
	Start float64
	End   float64
}

// TransitCalcInput holds all inputs for transit calculation
type TransitCalcInput struct {
	NatalLat     float64
	NatalLon     float64
	NatalJD      float64
	NatalPlanets []models.PlanetID

	TransitLat float64
	TransitLon float64

	StartJD float64
	EndJD   float64

	TransitPlanets []models.PlanetID

	SpecialPoints *models.SpecialPointsConfig
	EventConfig   models.EventConfig
	OrbConfig     models.OrbConfig
	HouseSystem   models.HouseSystem
}

// CalcTransitEvents computes all transit events in the given time range
func CalcTransitEvents(input TransitCalcInput) ([]models.TransitEvent, error) {
	var allEvents []models.TransitEvent

	// Pre-calculate natal chart data (fixed)
	natalHouses, err := chart.CalcNatalFixedHouses(input.NatalLat, input.NatalLon, input.NatalJD, input.HouseSystem)
	if err != nil {
		return nil, err
	}

	// Collect natal reference points (planets + special points)
	type refPoint struct {
		ID        string
		Longitude float64
	}
	var natalRefs []refPoint

	if input.EventConfig.IncludeAspects {
		for _, pid := range input.NatalPlanets {
			lon, _, err := chart.CalcPlanetLongitude(pid, input.NatalJD)
			if err != nil {
				continue
			}
			natalRefs = append(natalRefs, refPoint{ID: string(pid), Longitude: lon})
		}

		// Add natal special points
		if input.SpecialPoints != nil {
			for _, sp := range input.SpecialPoints.NatalPoints {
				lon, err := chart.CalcSpecialPointLongitude(sp, input.NatalLat, input.NatalLon, input.NatalJD, input.HouseSystem)
				if err != nil {
					continue
				}
				natalRefs = append(natalRefs, refPoint{ID: string(sp), Longitude: lon})
			}
		}
	}

	// Process each transit planet
	for _, tPlanet := range input.TransitPlanets {
		// Step 1: Find all stations (retrograde/direct) in the time range
		stations := findStations(tPlanet, input.StartJD, input.EndJD)

		// Build monotonic intervals from stations
		intervals := buildMonoIntervals(input.StartJD, input.EndJD, stations)

		// Station events
		if input.EventConfig.IncludeStation {
			for _, st := range stations {
				lon, _, _ := chart.CalcPlanetLongitude(tPlanet, st.JD)
				stType := models.StationRetrograde
				if st.IsDirecting {
					stType = models.StationDirect
				}
				allEvents = append(allEvents, models.TransitEvent{
					EventType:       models.EventStation,
					Planet:          tPlanet,
					JD:              st.JD,
					PlanetLongitude: lon,
					PlanetSign:      models.SignFromLongitude(lon),
					PlanetHouse:     chart.FindHouseForLongitude(lon, natalHouses),
					IsRetrograde:    stType == models.StationRetrograde,
					StationType:     stType,
				})
			}
		}

		// Aspect events (RQ1: transit planet vs fixed natal point)
		if input.EventConfig.IncludeAspects {
			exactCounters := make(map[string]int) // track exact_count per target+aspect
			for _, ref := range natalRefs {
				events := findAspectEventsRQ1(tPlanet, ref.ID, ref.Longitude,
					intervals, input.OrbConfig, natalHouses, exactCounters)
				allEvents = append(allEvents, events...)
			}
		}

		// Sign ingress events (RQ1 special case: cross 30° boundaries)
		if input.EventConfig.IncludeSignIngress {
			events := findSignIngressEvents(tPlanet, intervals, natalHouses)
			allEvents = append(allEvents, events...)
		}

		// House ingress events (RQ1 special case: cross natal house cusps)
		if input.EventConfig.IncludeHouseIngress {
			events := findHouseIngressEvents(tPlanet, intervals, natalHouses)
			allEvents = append(allEvents, events...)
		}
	}

	// Sort all events by JD
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].JD < allEvents[j].JD
	})

	return allEvents, nil
}

// findStations finds all retrograde/direct station points for a planet in [startJD, endJD]
func findStations(planet models.PlanetID, startJD, endJD float64) []StationInfo {
	var stations []StationInfo

	// Sun and Moon never retrograde
	if planet == models.PlanetSun || planet == models.PlanetMoon {
		return stations
	}

	// Coarse scan: check speed sign changes at 1-day intervals
	prevSpeed := getPlanetSpeed(planet, startJD)

	for jd := startJD + dayStep; jd <= endJD; jd += dayStep {
		curSpeed := getPlanetSpeed(planet, jd)

		if prevSpeed*curSpeed < 0 {
			// Speed sign changed, bisect to find exact station
			stationJD := bisectStation(planet, jd-dayStep, jd)
			isDirecting := prevSpeed < 0 && curSpeed > 0
			stations = append(stations, StationInfo{
				JD:          stationJD,
				IsDirecting: isDirecting,
			})
		}
		prevSpeed = curSpeed
	}

	return stations
}

func getPlanetSpeed(planet models.PlanetID, jd float64) float64 {
	_, speed, err := chart.CalcPlanetLongitude(planet, jd)
	if err != nil {
		return 0
	}
	return speed
}

// bisectStation uses binary search to find the exact moment of station (speed = 0)
func bisectStation(planet models.PlanetID, lo, hi float64) float64 {
	speedLo := getPlanetSpeed(planet, lo)

	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		speedMid := getPlanetSpeed(planet, mid)
		if speedLo*speedMid <= 0 {
			hi = mid
		} else {
			lo = mid
			speedLo = speedMid
		}
	}
	return (lo + hi) / 2
}

// buildMonoIntervals creates monotonic intervals from stations
func buildMonoIntervals(startJD, endJD float64, stations []StationInfo) []MonoInterval {
	var intervals []MonoInterval
	prev := startJD
	for _, st := range stations {
		if st.JD > prev && st.JD <= endJD {
			intervals = append(intervals, MonoInterval{Start: prev, End: st.JD})
			prev = st.JD
		}
	}
	if prev < endJD {
		intervals = append(intervals, MonoInterval{Start: prev, End: endJD})
	}
	return intervals
}

// findAspectEventsRQ1 finds aspect ENTER/EXACT/LEAVE events for transit planet vs fixed natal point
func findAspectEventsRQ1(
	tPlanet models.PlanetID, targetID string, targetLon float64,
	intervals []MonoInterval, orbs models.OrbConfig,
	natalHouses []float64, exactCounters map[string]int,
) []models.TransitEvent {
	var events []models.TransitEvent

	for _, asp := range models.StandardAspects {
		orb := orbs.GetOrb(asp.Type)
		if orb == 0 {
			continue
		}

		counterKey := targetID + ":" + string(asp.Type)
		inAspect := false // tracking whether we're currently in the orb

		// Check initial state
		initLon, _, _ := chart.CalcPlanetLongitude(tPlanet, intervals[0].Start)
		initDiff := angleDiffToAspect(initLon, targetLon, asp.Angle)
		if math.Abs(initDiff) <= orb {
			inAspect = true
		}

		for _, interval := range intervals {
			// Scan this interval with fine steps to find orb boundary crossings and exact hits
			step := fineStep
			prevJD := interval.Start
			prevLon, _, _ := chart.CalcPlanetLongitude(tPlanet, prevJD)
			prevDiff := angleDiffToAspect(prevLon, targetLon, asp.Angle)

			for jd := interval.Start + step; jd <= interval.End+step*0.5; jd += step {
				if jd > interval.End {
					jd = interval.End
				}
				curLon, _, _ := chart.CalcPlanetLongitude(tPlanet, jd)
				curDiff := angleDiffToAspect(curLon, targetLon, asp.Angle)

				// Check for ENTER (crossing into orb)
				if !inAspect && math.Abs(curDiff) <= orb && math.Abs(prevDiff) > orb {
					enterJD := bisectThreshold(tPlanet, targetLon, asp.Angle, orb, prevJD, jd, true)
					eLon, _, _ := chart.CalcPlanetLongitude(tPlanet, enterJD)
					eRetro := getPlanetSpeed(tPlanet, enterJD) < 0
					events = append(events, models.TransitEvent{
						EventType:       models.EventAspectEnter,
						Planet:          tPlanet,
						JD:              enterJD,
						PlanetLongitude: eLon,
						PlanetSign:      models.SignFromLongitude(eLon),
						PlanetHouse:     chart.FindHouseForLongitude(eLon, natalHouses),
						IsRetrograde:    eRetro,
						Target:          targetID,
						AspectType:      asp.Type,
						AspectAngle:     asp.Angle,
						OrbAtEnter:      math.Abs(angleDiffToAspect(eLon, targetLon, asp.Angle)),
					})
					inAspect = true
				}

				// Check for EXACT (diff crosses zero)
				if prevDiff*curDiff < 0 || (math.Abs(curDiff) < 0.01 && inAspect) {
					if prevDiff*curDiff < 0 {
						exactJD := bisectExact(tPlanet, targetLon, asp.Angle, prevJD, jd)
						exactCounters[counterKey]++
						eLon, _, _ := chart.CalcPlanetLongitude(tPlanet, exactJD)
						eRetro := getPlanetSpeed(tPlanet, exactJD) < 0
						events = append(events, models.TransitEvent{
							EventType:       models.EventAspectExact,
							Planet:          tPlanet,
							JD:              exactJD,
							PlanetLongitude: eLon,
							PlanetSign:      models.SignFromLongitude(eLon),
							PlanetHouse:     chart.FindHouseForLongitude(eLon, natalHouses),
							IsRetrograde:    eRetro,
							Target:          targetID,
							AspectType:      asp.Type,
							AspectAngle:     asp.Angle,
							ExactCount:      exactCounters[counterKey],
						})
					}
				}

				// Check for LEAVE (crossing out of orb)
				if inAspect && math.Abs(curDiff) > orb && math.Abs(prevDiff) <= orb {
					leaveJD := bisectThreshold(tPlanet, targetLon, asp.Angle, orb, prevJD, jd, false)
					eLon, _, _ := chart.CalcPlanetLongitude(tPlanet, leaveJD)
					eRetro := getPlanetSpeed(tPlanet, leaveJD) < 0
					events = append(events, models.TransitEvent{
						EventType:       models.EventAspectLeave,
						Planet:          tPlanet,
						JD:              leaveJD,
						PlanetLongitude: eLon,
						PlanetSign:      models.SignFromLongitude(eLon),
						PlanetHouse:     chart.FindHouseForLongitude(eLon, natalHouses),
						IsRetrograde:    eRetro,
						Target:          targetID,
						AspectType:      asp.Type,
						AspectAngle:     asp.Angle,
						OrbAtLeave:      math.Abs(angleDiffToAspect(eLon, targetLon, asp.Angle)),
					})
					inAspect = false
				}

				prevJD = jd
				prevLon = curLon
				prevDiff = curDiff

				if jd >= interval.End {
					break
				}
			}
		}
	}

	return events
}

// angleDiffToAspect returns the signed difference between actual angle and target aspect angle
// Result is in [-180, 180), where 0 means exact aspect
func angleDiffToAspect(transitLon, natalLon, aspectAngle float64) float64 {
	actualAngle := shortestAngle(transitLon, natalLon)
	return wrapAngle(actualAngle - aspectAngle)
}

func shortestAngle(lon1, lon2 float64) float64 {
	diff := math.Abs(lon1 - lon2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

func wrapAngle(a float64) float64 {
	a = math.Mod(a+180, 360)
	if a < 0 {
		a += 360
	}
	return a - 180
}

// bisectThreshold finds the exact JD where |angleDiff| crosses the orb threshold
func bisectThreshold(planet models.PlanetID, targetLon, aspectAngle, orb, lo, hi float64, entering bool) float64 {
	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		midLon, _, _ := chart.CalcPlanetLongitude(planet, mid)
		midDiff := math.Abs(angleDiffToAspect(midLon, targetLon, aspectAngle))

		if entering {
			// Looking for transition from outside to inside orb
			if midDiff > orb {
				lo = mid
			} else {
				hi = mid
			}
		} else {
			// Looking for transition from inside to outside orb
			if midDiff <= orb {
				lo = mid
			} else {
				hi = mid
			}
		}
	}
	return (lo + hi) / 2
}

// bisectExact finds the exact JD where angleDiff crosses zero
func bisectExact(planet models.PlanetID, targetLon, aspectAngle, lo, hi float64) float64 {
	loLon, _, _ := chart.CalcPlanetLongitude(planet, lo)
	loDiff := angleDiffToAspect(loLon, targetLon, aspectAngle)

	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		midLon, _, _ := chart.CalcPlanetLongitude(planet, mid)
		midDiff := angleDiffToAspect(midLon, targetLon, aspectAngle)

		if loDiff*midDiff <= 0 {
			hi = mid
		} else {
			lo = mid
			loDiff = midDiff
		}
	}
	return (lo + hi) / 2
}

// findSignIngressEvents finds when a transit planet crosses sign boundaries (0°, 30°, 60°, ...)
func findSignIngressEvents(planet models.PlanetID, intervals []MonoInterval, natalHouses []float64) []models.TransitEvent {
	var events []models.TransitEvent

	for _, interval := range intervals {
		prevJD := interval.Start
		prevLon, _, _ := chart.CalcPlanetLongitude(planet, prevJD)
		prevSign := int(prevLon / 30.0)

		step := fineStep
		for jd := interval.Start + step; jd <= interval.End+step*0.5; jd += step {
			if jd > interval.End {
				jd = interval.End
			}
			curLon, _, _ := chart.CalcPlanetLongitude(planet, jd)
			curSign := int(curLon / 30.0)

			if curSign != prevSign {
				// Bisect to find exact crossing
				crossJD := bisectSignBoundary(planet, prevJD, jd, prevSign, curSign)
				cLon, _, _ := chart.CalcPlanetLongitude(planet, crossJD)
				cRetro := getPlanetSpeed(planet, crossJD) < 0

				fromSign := models.ZodiacSigns[prevSign%12]
				toSign := models.ZodiacSigns[curSign%12]

				events = append(events, models.TransitEvent{
					EventType:       models.EventSignIngress,
					Planet:          planet,
					JD:              crossJD,
					PlanetLongitude: cLon,
					PlanetSign:      toSign,
					PlanetHouse:     chart.FindHouseForLongitude(cLon, natalHouses),
					IsRetrograde:    cRetro,
					FromSign:        fromSign,
					ToSign:          toSign,
				})
			}

			prevJD = jd
			prevLon = curLon
			prevSign = curSign

			if jd >= interval.End {
				break
			}
		}
	}
	return events
}

func bisectSignBoundary(planet models.PlanetID, lo, hi float64, prevSign, curSign int) float64 {
	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		midLon, _, _ := chart.CalcPlanetLongitude(planet, mid)
		midSign := int(midLon / 30.0)
		if midSign == prevSign {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

// findHouseIngressEvents finds when a transit planet crosses natal house cusps
func findHouseIngressEvents(planet models.PlanetID, intervals []MonoInterval, natalHouses []float64) []models.TransitEvent {
	var events []models.TransitEvent
	if len(natalHouses) < 12 {
		return events
	}

	for _, interval := range intervals {
		prevJD := interval.Start
		prevLon, _, _ := chart.CalcPlanetLongitude(planet, prevJD)
		prevHouse := chart.FindHouseForLongitude(prevLon, natalHouses)

		step := fineStep
		for jd := interval.Start + step; jd <= interval.End+step*0.5; jd += step {
			if jd > interval.End {
				jd = interval.End
			}
			curLon, _, _ := chart.CalcPlanetLongitude(planet, jd)
			curHouse := chart.FindHouseForLongitude(curLon, natalHouses)

			if curHouse != prevHouse {
				crossJD := bisectHouseBoundary(planet, prevJD, jd, natalHouses, prevHouse)
				cLon, _, _ := chart.CalcPlanetLongitude(planet, crossJD)
				cRetro := getPlanetSpeed(planet, crossJD) < 0

				events = append(events, models.TransitEvent{
					EventType:       models.EventHouseIngress,
					Planet:          planet,
					JD:              crossJD,
					PlanetLongitude: cLon,
					PlanetSign:      models.SignFromLongitude(cLon),
					PlanetHouse:     curHouse,
					IsRetrograde:    cRetro,
					FromHouse:       prevHouse,
					ToHouse:         curHouse,
				})
			}

			prevJD = jd
			prevLon = curLon
			prevHouse = curHouse

			if jd >= interval.End {
				break
			}
		}
	}
	return events
}

func bisectHouseBoundary(planet models.PlanetID, lo, hi float64, cusps []float64, prevHouse int) float64 {
	for hi-lo > bisectEps {
		mid := (lo + hi) / 2
		midLon, _, _ := chart.CalcPlanetLongitude(planet, mid)
		midHouse := chart.FindHouseForLongitude(midLon, cusps)
		if midHouse == prevHouse {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}
