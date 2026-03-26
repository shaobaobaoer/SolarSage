//go:build ignore

package main

import (
	"fmt"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/julian"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
	"github.com/shaobaobaoer/solarsage-mcp/pkg/transit"
)

func main() {
	sweph.Init("/home/ecs-user/SolarSage/ephe")

	natalASC := 96.529167// 06°Cancer31'45''
	natalMC := 351.4975  // 21°Pisces29'51''

	events, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalJD:  2450800.900009,
		NatalLat: 30.9,
		NatalLon: 121.15,
		NatalASC: natalASC,
		NatalMC:  natalMC,
		NatalPlanets: []models.PlanetID{
			models.PlanetPluto,
		},
		StartJD: sweph.JulDay(2026, 1, 1, 0, true),
		EndJD:   sweph.JulDay(2026, 12, 31, 0, true),
		TransitPlanets: []models.PlanetID{
			models.PlanetPluto,
		},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{models.PlanetPluto},
		},
		SpecialPoints: &models.SpecialPointsConfig{
			NatalPoints:[]models.SpecialPointID{models.PointASC},
			SolarArcPoints: []models.SpecialPointID{models.PointASC},
		},
		EventConfig: models.EventConfig{
			IncludeTrNa: true,
			IncludeTrSa: true,
		},
		OrbConfigSolarArc: models.OrbConfig{
			Opposition: 1,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total events: %d\n\n", len(events))

	fmt.Println("=== All events ===")
	for _, e := range events {
		dt, _ := julian.JDToDateTime(e.JD, "Australia/Perth")
		fmt.Printf("%-12s %-10s %-10s %-8s JD=%.4f %s AWST\n",
			e.EventType, e.Planet, e.AspectType, e.ChartType, e.JD, dt)
	}
}
