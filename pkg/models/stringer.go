package models

import "fmt"

// String returns a human-readable planet position with glyphs like "☉ 10°15'♑ (H10)"
func (p PlanetPosition) String() string {
	return FormatPlanetGlyph(p)
}

// String returns a human-readable aspect with glyphs like "☉ □ ☽ 2.5° (applying)"
func (a AspectInfo) String() string {
	applying := "separating"
	if a.IsApplying {
		applying = "applying"
	}
	return fmt.Sprintf("%s (%s)",
		FormatAspectGlyph(a),
		applying,
	)
}

// String returns a human-readable cross-aspect with glyphs
func (ca CrossAspectInfo) String() string {
	applying := "separating"
	if ca.IsApplying {
		applying = "applying"
	}
	pA := PlanetGlyph(PlanetID(ca.InnerBody))
	pB := PlanetGlyph(PlanetID(ca.OuterBody))
	asp := AspectGlyph(ca.AspectType)
	return fmt.Sprintf("%s %s %s %.1f° (%s)", pA, asp, pB, ca.Orb, applying)
}

// String returns a human-readable transit event with glyphs
func (te TransitEvent) String() string {
	pg := PlanetGlyph(te.Planet)
	switch te.EventType {
	case EventAspectExact, EventAspectEnter, EventAspectLeave, EventAspectBegin:
		ag := AspectGlyph(te.AspectType)
		tg := PlanetGlyph(PlanetID(te.Target))
		return fmt.Sprintf("%s %s %s %s %s (%s-%s)",
			EventTypeCSV(te.EventType, te.StationType),
			pg, ag, tg,
			BodyDisplayName(te.Target),
			ChartTypeShort(te.ChartType),
			ChartTypeShort(te.TargetChartType),
		)
	case EventSignIngress:
		sg := SignGlyph(te.ToSign)
		return fmt.Sprintf("%s -> %s %s", pg, sg, te.ToSign)
	case EventHouseIngress:
		return fmt.Sprintf("%s -> H%d", pg, te.ToHouse)
	case EventStation:
		retro := "\u211E" // ℞
		if te.StationType == StationDirect {
			retro = "D"
		}
		return fmt.Sprintf("%s Station %s", pg, retro)
	case EventVoidOfCourse:
		return fmt.Sprintf("VOC %s (%s to %s %s)",
			PlanetGlyph(PlanetMoon),
			te.LastAspectType,
			SignGlyph(te.NextSign), te.NextSign,
		)
	default:
		return fmt.Sprintf("%s %s", te.EventType, pg)
	}
}

// String returns angles with glyphs
func (a AnglesInfo) String() string {
	return fmt.Sprintf("AC %s, MC %s",
		FormatLonGlyph(a.ASC),
		FormatLonGlyph(a.MC),
	)
}
