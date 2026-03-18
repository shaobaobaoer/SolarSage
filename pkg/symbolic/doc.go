// Package symbolic implements symbolic directions, a predictive technique
// that advances each natal planet and angle by a fixed rate per year of life.
//
// Unlike secondary progressions (which use actual ephemeris positions),
// symbolic directions apply a mathematical formula:
//
//	Directed position = natal position + age x rate
//
// Supported methods include One Degree per Year (1.0/year), Naibod Arc
// (0.98556/year, the Sun's mean daily motion), Profection Arc (30/year),
// and user-defined custom rates.
//
// CalcSymbolicDirections is the main entry point. It computes the natal
// chart, applies the arc to every planet and angle, and returns aspects
// between directed and natal positions.
package symbolic
