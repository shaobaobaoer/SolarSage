// Package kp implements the Krishnamurti Paddhati (KP) system of Vedic astrology.
//
// KP astrology subdivides each of the 27 Nakshatras into sub-periods proportional
// to the Vimshottari Dasha sequence (120-year cycle). Every degree of the zodiac
// belongs to a specific Nakshatra lord (Star lord) and Sub lord.
//
// Core features:
//   - Sub-lord and Sub-sub-lord calculation for any sidereal longitude
//   - Pre-computed 249-row reference table (KP cuspal sub-divisions)
//   - KP house placement by cusp longitude (distinct from sign-based placement)
//   - ABCD significator analysis for planets and houses
//
// Usage:
//
//	// Get Star lord / Sub lord / Sub-sub lord for a planet longitude
//	info := kp.SubLords(planetSiderealLon)
//	fmt.Println(info.NakshatraLord, info.SubLord, info.SubSubLord)
//
//	// Get KP house number for a planet (by cusp longitude)
//	house := kp.HouseNumber(planetLon, cuspLongitudes)
package kp
