// Package api provides a RESTful HTTP API server that exposes all SolarSage
// astrology tools as JSON endpoints. It wraps the same underlying packages
// used by the MCP server, providing a thin adapter layer over net/http.
//
// All 40 endpoints are POST-only under /api/v1/, with optional API key
// authentication via the X-API-Key header and CORS support.
package api
