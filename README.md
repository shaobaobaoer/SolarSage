# swisseph-mcp

High-precision astrology calculation engine exposed as a [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server. Built on the Swiss Ephemeris library with sub-arcsecond accuracy.

## Features

- **Natal Charts** - Planet positions, house cusps (7 systems), angles, and aspects
- **Transit Detection** - All major transit types: Tr-Na, Tr-Tr, Tr-Sp, Tr-Sa, Sp-Na, Sp-Sp, Sa-Na
- **Secondary Progressions** - Day-for-a-year progressed positions and events
- **Solar Arc Directions** - Solar arc directed positions and events
- **Sign & House Ingress** - Detect when planets enter new signs or houses
- **Stations** - Retrograde and direct station detection
- **Void of Course Moon** - Automatic VOC detection with aspect context
- **Geocoding** - Location name to coordinates via OpenStreetMap Nominatim
- **CSV Export** - Solar Fire compatible output format
- **1-second precision** - Bisection algorithm for exact event timing

## Supported Bodies

Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, Chiron, North Node (True/Mean), South Node, Lilith (Mean/True)

**Special Points:** ASC, MC, DSC, IC, Vertex, East Point, Lot of Fortune, Lot of Spirit

## House Systems

Placidus, Koch, Equal, Whole Sign, Campanus, Regiomontanus, Porphyry

## Quick Start

### Prerequisites

- Go 1.21+
- GCC (for CGO/Swiss Ephemeris compilation)

### Build

```bash
make build
```

### Run as MCP Server

```bash
# Uses ephemeris files from third_party/swisseph/ephe/
./bin/swisseph-mcp

# Or specify a custom ephemeris path
SWISSEPH_EPHE_PATH=/path/to/ephe ./bin/swisseph-mcp
```

### Claude Desktop Integration

Add to your Claude Desktop config (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "astrology": {
      "command": "/path/to/swisseph-mcp",
      "env": {
        "SWISSEPH_EPHE_PATH": "/path/to/ephe"
      }
    }
  }
}
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `geocode` | Location name to coordinates and timezone |
| `datetime_to_jd` | ISO 8601 datetime to Julian Day (UT/TT) |
| `jd_to_datetime` | Julian Day to ISO 8601 datetime |
| `calc_planet_position` | Single planet position at a given time |
| `calc_single_chart` | Full natal/event chart calculation |
| `calc_double_chart` | Synastry/transit double chart with cross-aspects |
| `calc_progressions` | Secondary progressed planet positions |
| `calc_solar_arc` | Solar arc directed planet positions |
| `calc_transit` | Full transit event search over a time range |

## Example: Calculate a Natal Chart

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "calc_single_chart",
    "arguments": {
      "latitude": 51.5074,
      "longitude": -0.1278,
      "jd_ut": 2451545.0,
      "house_system": "PLACIDUS"
    }
  }
}
```

## Example: Search Transit Events

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "calc_transit",
    "arguments": {
      "natal_latitude": 51.5074,
      "natal_longitude": -0.1278,
      "natal_jd_ut": 2451545.0,
      "transit_latitude": 51.5074,
      "transit_longitude": -0.1278,
      "start_jd_ut": 2460676.5,
      "end_jd_ut": 2460706.5,
      "format": "json"
    }
  }
}
```

## Architecture

```
cmd/server/     MCP server entry point
pkg/mcp/        MCP protocol (JSON-RPC over stdio)
pkg/chart/      Chart calculations (positions, houses, aspects)
pkg/transit/    Transit event detection engine
pkg/progressions/ Secondary progressions & solar arc
pkg/models/     Core data types and constants
pkg/julian/     Julian Day conversions
pkg/geo/        Geocoding and timezone lookup
pkg/export/     CSV/JSON export
pkg/sweph/      Swiss Ephemeris C bindings (CGO)
internal/aspect/ Aspect calculation engine
```

## Accuracy

Validated against Solar Fire 9 with **100% exact event match** (247/247 events) over a 1-year transit period including all 7 chart type combinations.

## License

MIT
