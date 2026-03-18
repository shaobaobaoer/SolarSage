// Package export formats transit events as CSV and JSON output.
//
// EventsToCSV renders a slice of TransitEvent as a complete CSV document
// with headers. EventsToJSON produces the equivalent JSON array.
// EventToCSVRow and CSVRowToString handle individual row conversion.
package export
