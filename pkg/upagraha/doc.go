// Package upagraha computes the eight classical Vedic sub-planets (Upagrahas).
//
// Upagrahas are sensitive points derived from the planetary hour sequence and
// the birth time — they require no separate ephemeris lookup. They are widely
// used in Jyotish natal and horary analysis.
//
// The eight Upagrahas implemented are:
//   - Gulika / Mandi   — child of Saturn; the most important Upagraha
//   - Dhuma            — derived from Sun
//   - Vyatipaata       — derived from Dhuma
//   - Parivesha        — derived from Vyatipaata
//   - Indrachaapa      — derived from Parivesha
//   - Upaketu          — derived from Indrachaapa
//   - Kaala            — child of Saturn (alternate)
//   - Yamaghantaka     — child of Jupiter
//
// Usage:
//
//	result := upagraha.Calc(jdUT, lat, lon, isDayChart)
//	fmt.Println(result.Gulika.Longitude, result.Gulika.Sign)
package upagraha
