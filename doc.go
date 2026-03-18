// Package swisseph is the solarsage-mcp module root.
//
// For the high-level convenience API, import pkg/solarsage:
//
//	import "github.com/shaobaobaoer/solarsage-mcp/pkg/solarsage"
//
//	solarsage.Init("/path/to/ephe")
//	defer solarsage.Close()
//
//	chart, _ := solarsage.NatalChart(51.5, -0.1, "2000-01-01T12:00:00Z")
//	phase, _ := solarsage.MoonPhase("2025-03-18T12:00:00Z")
//	score, _ := solarsage.Compatibility(lat1, lon1, dt1, lat2, lon2, dt2)
//
// For lower-level control, import individual packages from pkg/.
// For MCP server usage, run the binary directly or via Docker.
package swisseph
