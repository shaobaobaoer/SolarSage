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

	// SF exact times for Pluto Opposition ASC Tr-Sa
	// 1. 2026-04-13 16:37:54 AWST = 2026-04-13 08:37:54 UTC
	// 2. 2026-05-18 05:51:23 AWST = 2026-05-17 21:51:23 UTC
	sfJD1 := sweph.JulDay(2026, 4, 13, 8+37.0/60+54.0/3600, true)
	sfJD2 := sweph.JulDay(2026, 5, 17, 21+51.0/60+23.0/3600, true)

	fmt.Println("=== SF exact times (UTC) ===")
	t1, _ := julian.JDToDateTime(sfJD1, "UTC")
	t2, _ := julian.JDToDateTime(sfJD2, "UTC")
	fmt.Printf("SF Exact 1: %s UTC (JD=%.6f)\n", t1, sfJD1)
	fmt.Printf("SF Exact 2: %s UTC (JD=%.6f)\n", t2, sfJD2)

	// Run with ASC/MC override
	events, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalJD:  2450800.900009,
		NatalLat: 30.9,
		NatalLon: 121.15,
		NatalASC: natalASC,
		NatalMC:  natalMC,
		NatalPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto, models.PlanetNorthNodeMean,
		},
		StartJD: sweph.JulDay(2026, 1, 1, 0, true),
		EndJD:   sweph.JulDay(2026, 12, 31, 0, true),
		TransitPlanets: []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto, models.PlanetNorthNodeMean,
		},
		SolarArcConfig: &models.SolarArcConfig{
			Enabled: true,
			Planets: []models.PlanetID{
				models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
				models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
				models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
				models.PlanetPluto, models.PlanetNorthNodeMean,
			},
		},
		SpecialPoints: &models.SpecialPointsConfig{
			NatalPoints:[]models.SpecialPointID{models.PointASC, models.PointMC},
			SolarArcPoints: []models.SpecialPointID{models.PointASC, models.PointMC},
		},
		EventConfig: models.EventConfig{
			IncludeTrNa: true,
			IncludeTrSa: true,
		},
		OrbConfigSolarArc: models.OrbConfig{
			Conjunction: 1, Opposition: 1, Square: 1, Trine: 1, Sextile: 1,
			SemiSextile: 1, Quincunx: 1,
		},
		HouseSystem: models.HousePlacidus,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Find Pluto Opposition ASC Tr-Sa exact events
	fmt.Println("\n=== Pluto Tr-Sa events ===")
	for _, e := range events {
		if e.Planet == "PLUTO" {
			dt, _ := julian.JDToDateTime(e.JD, "Australia/Perth")
			fmt.Printf("  %s %-10s %-12s %-10s [%d] %s AWST\n",
				e.EventType, e.Planet, e.AspectType, e.ChartType, e.ExactCount, dt)
		}
	}
	for _, e := range events {
		if e.Planet == "PLUTO" && e.AspectType == "Opposition" && e.ChartType == "Tr-Sa" && e.EventType == "ASPECT_EXACT" {
			dt, _ := julian.JDToDateTime(e.JD, "Australia/Perth")
			diff := e.JD - sfJD1
			if e.JD > sfJD1+0.5 {
				diff = e.JD - sfJD2
			}
			diffMin := diff * 24 * 60
			fmt.Printf("Our Exact: %s AWST (JD=%.6f)\n", dt, e.JD)
			fmt.Printf("  Diff from nearest SF: %.1f minutes\n", diffMin)
		}
	}

	// Compare Uranus Square Uranus Tr-Sa
	fmt.Println("\n=== Uranus Square Uranus Tr-Sa ===")
	// SF: Sep 5 06:15:27 AWST = Sep 4 22:15:27 UTC
	// SF: Sep 10 12:53:48 AWST = Sep 10 04:53:48 UTC
	sfUranusJD1 := sweph.JulDay(2026, 9, 4, 22+15.0/60+27.0/3600, true)
	sfUranusJD2 := sweph.JulDay(2026, 9, 10, 4+53.0/60+48.0/3600, true)

	for _, e := range events {
		if e.Planet == "URANUS" && e.AspectType == "Square" && e.ChartType == "Tr-Sa" && e.EventType == "ASPECT_EXACT" {
			dt, _ := julian.JDToDateTime(e.JD, "Australia/Perth")
			var diffMin float64
			if e.JD < sfUranusJD1+2 {
				diffMin = (e.JD - sfUranusJD1) * 24 * 60
			} else {
				diffMin = (e.JD - sfUranusJD2) * 24 * 60
			}
			fmt.Printf("Our Exact: %s AWST (JD=%.6f) diff=%.1f min\n", dt, e.JD, diffMin)
		}
	}
}
