// Package julian converts between ISO 8601 datetime strings and Julian Day
// numbers used by the Swiss Ephemeris.
//
// DateTimeToJD parses an ISO 8601 datetime and returns a Julian Day with
// calendar metadata. JDToDateTime converts a Julian Day back to an ISO 8601
// string in the specified timezone.
package julian
