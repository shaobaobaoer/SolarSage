// Package lunar provides lunar phase calculation and eclipse detection.
//
// CalcLunarPhase returns the current Moon phase, illumination, and Sun-Moon
// angle for a given Julian Day. FindLunarPhases scans a date range for
// New Moons, Full Moons, and quarter phases. FindEclipses detects solar
// and lunar eclipses with type classification (partial, total, annular).
// NextNewMoon and NextFullMoon find the next occurrence of each phase.
package lunar
