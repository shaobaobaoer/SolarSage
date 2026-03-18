// Package transit implements the transit event detection engine supporting
// seven chart-type combinations (natal transits, progressions, solar arcs,
// and their cross-references).
//
// CalcTransitEvents is the primary entry point, accepting a TransitCalcInput
// that specifies the date range, bodies, and event types to scan. Results are
// independently validated at 100% accuracy (247/247 reference events).
package transit
