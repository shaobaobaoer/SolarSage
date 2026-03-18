// Package sweph provides thread-safe Go bindings to the Swiss Ephemeris C
// library via CGO. All calls are serialized through a global mutex because
// the underlying C library is not thread-safe.
//
// Init must be called before any calculations to set the ephemeris file path.
// CalcUT computes planet positions, Houses returns house cusps and angles,
// and JulDay/RevJul convert between calendar dates and Julian Day numbers.
package sweph
