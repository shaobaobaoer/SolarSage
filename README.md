# solarsage-mcp

The most comprehensive open-source astrology calculation engine. **40 MCP tools, 40 REST endpoints, 38 packages, 824+ tests, 11 house systems, 5 ayanamsas, 50+ fixed stars, 27 Nakshatras, 15+ Arabic lots, 7 aspect patterns** - all with sub-arcsecond accuracy. Usable as a **Go library**, an **MCP server** for AI assistants, or a **RESTful HTTP API** for web/mobile clients.

Built on the [Swiss Ephemeris](https://www.astro.com/swisseph/). Faster, more modern, and more comprehensive than traditional desktop astrology software. Independently validated at 100% accuracy (247/247 transit events).

## Why solarsage?

| | solarsage | flatlib (Python) | Kerykeion (Python) | Swiss Ephemeris (C) |
|---|---|---|---|---|
| Language | Go | Python | Python | C |
| Transit detection | 7 types, 1s precision | Basic | None | Manual |
| Solar/Lunar returns | Series support | Single | Single | Manual |
| Composite charts | Midpoint method | None | None | Manual |
| Synastry scoring | Category breakdown | None | Basic | None |
| Eclipse detection | Solar + Lunar | None | None | Low-level |
| Profections | Annual + monthly | None | None | None |
| Arabic lots | 15+ with day/night | None | None | None |
| Essential dignities | Full + mutual reception | Basic | Basic | None |
| Aspect patterns | 7 types | None | None | None |
| Fixed stars | 50+ catalog | None | None | Low-level |
| Midpoints | 90deg dial + activations | None | None | None |
| Harmonic charts | 1-180 | None | None | None |
| Planetary hours | Chaldean | None | None | None |
| House systems | 11 | 7 | 3 | All |
| Sidereal/Vedic | Nakshatras + Dasha | None | None | Manual |
| Dispositors | Full chains | None | None | None |
| One-call report | Everything combined | None | None | None |
| Chart visualization | Wheel coordinates | None | None | None |
| MCP server | 40 tools | None | None | None |
| REST API | 40 endpoints | None | None | None |
| Accuracy validated | 247/247 (100%) | No | No | N/A |
| Thread-safe | Yes (mutex) | No | No | No |

## Features

### Chart Calculations
- **Natal Charts** - Positions, houses (11 systems), angles, aspects (9 types)
- **Double Charts** - Synastry/transit with cross-aspects
- **Composite Charts** - Midpoint method for relationship analysis
- **Harmonic Charts** - Divisional charts (5th quintile, 7th septile, 9th novile, etc.)

### Predictive Techniques
- **Transit Detection** - Tr-Na, Tr-Tr, Tr-Sp, Tr-Sa, Sp-Na, Sp-Sp, Sa-Na with 1-second precision
- **Secondary Progressions** - Day-for-a-year progressed positions and events
- **Solar Arc Directions** - Solar arc directed positions and events
- **Solar & Lunar Returns** - Exact return charts with series support
- **Annual Profections** - Time-lord technique with monthly sub-profections
- **Sign/House Ingress** - Planet sign and house change detection
- **Stations** - Retrograde and direct station detection

### Traditional Astrology
- **Essential Dignities** - Rulership, exaltation, detriment, fall with scoring
- **Mutual Receptions** - Rulership and exaltation mutual receptions
- **Sect** - Diurnal/nocturnal planet alignment analysis
- **Arabic Lots** - 15+ lots (Fortune, Spirit, Eros, Victory, etc.) with day/night reversal
- **Decans & Terms** - Chaldean decans and Egyptian/Ptolemaic term boundaries
- **Planetary Hours** - Chaldean hours with computed sunrise/sunset
- **Antiscia** - Solstice and equinox mirror points with pair detection
- **Dispositors** - Dispositorship chains, final dispositor, mutual dispositors

### Pattern Detection
- **Aspect Patterns** - Grand Trine, T-Square, Grand Cross, Yod, Kite, Mystic Rectangle, Stellium
- **Fixed Stars** - 50+ major star catalog with precession-corrected conjunctions
- **Midpoint Analysis** - Full midpoint tree, 90-degree Cosmobiology dial, activations

### Advanced Predictive
- **Primary Directions** - Ptolemy semi-arc method with Naibod key
- **Symbolic Directions** - 1-degree/year, Naibod, Profection, custom rate
- **Firdaria** - Planetary period system (day/night sequences) with timeline

### Traditional (continued)
- **Bonification & Maltreatment** - Aspect-based planetary condition analysis
- **Heliacal Risings/Settings** - Swiss Ephemeris visibility algorithms

### Vedic / Sidereal
- **Sidereal Charts** - 5 ayanamsa systems (Lahiri, Raman, Krishnamurti, Fagan-Bradley, Yukteshwar)
- **Nakshatras** - All 27 lunar mansions with padas and Vimshottari lords
- **Vimshottari Dasha** - Full Maha Dasha period sequence from Moon's Nakshatra
- **Divisional Charts** - 16 Varga charts (D1-D60, Navamsa, Dasamsa, etc.)
- **Ashtakavarga** - Bindu tables and Sarvashtakavarga
- **Yoga Detection** - Mahapurusha, Raja, Dhana, Gajakesari, and more

### Astronomical
- **Lunar Phases** - New/full moon finder, phase angle, illumination percentage
- **Eclipse Finder** - Solar and lunar eclipse detection with type classification
- **Void of Course Moon** - Automatic VOC detection with aspect context

### Relationship
- **Synastry Scoring** - Compatibility analysis with category breakdown (love, passion, communication, commitment)
- **Composite Charts** - Midpoint method with aspects
- **Davison Chart** - Midpoint in time and space relationship chart

### Visualization
- **Chart Wheel Coordinates** - Planet x/y positions, house cusp lines, aspect lines, sign segments for SVG/Canvas rendering

### Supported Bodies

Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, Chiron, North Node (True/Mean), South Node, Lilith (Mean/True)

**Special Points:** ASC, MC, DSC, IC, Vertex, East Point, Lot of Fortune, Lot of Spirit

**House Systems:** Placidus, Koch, Equal, Whole Sign, Campanus, Regiomontanus, Porphyry, Morinus, Topocentric, Alcabitius, Meridian

**Output:** JSON and CSV for all chart types. Unicode astrology glyphs (♈♉♊♋♌♍♎♏♐♑♒♓, ☉☽☿♀♂♃♄, ☌☍△□✱).

## Quick Start

### Prerequisites

- Go 1.25+
- GCC (for CGO / Swiss Ephemeris compilation)

### Build

```bash
git clone https://github.com/shaobaobaoer/solarsage-mcp.git
cd solarsage-mcp
make build          # MCP server
make build-api      # REST API server
```

### Run as MCP Server

```bash
./bin/solarsage-mcp

# Or with a custom ephemeris path
SWISSEPH_EPHE_PATH=/path/to/ephe ./bin/solarsage-mcp
```

### Claude Desktop Integration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "astrology": {
      "command": "/path/to/solarsage-mcp",
      "env": {
        "SWISSEPH_EPHE_PATH": "/path/to/ephe"
      }
    }
  }
}
```

### Run as REST API Server

```bash
./bin/solarsage-api --port 8080

# With API key authentication
./bin/solarsage-api --port 8080 --api-key your-secret-key

# Example request
curl -X POST http://localhost:8080/api/v1/chart/natal \
  -H "Content-Type: application/json" \
  -d '{"latitude": 51.5074, "longitude": -0.1278, "jd_ut": 2451545.0}'
```

All 40 endpoints are available under `/api/v1/`. CORS enabled. Optional API key auth via `X-API-Key` header.

## Use as a Go Library

### Quick API (recommended)

The `solarsage` package provides a high-level API with sensible defaults. Pass datetime strings instead of Julian Day numbers:

```go
package main

import (
    "fmt"
    "github.com/shaobaobaoer/solarsage-mcp/pkg/solarsage"
)

func main() {
    solarsage.Init("/path/to/ephe")
    defer solarsage.Close()

    // Natal chart
    chart, _ := solarsage.NatalChart(51.5074, -0.1278, "1990-06-15T14:30:00Z")
    for _, p := range chart.Planets {
        fmt.Printf("%s in %s (house %d)\n", p.PlanetID, p.Sign, p.House)
    }

    // Solar return for 2025
    sr, _ := solarsage.SolarReturn(51.5074, -0.1278, "1990-06-15T14:30:00Z", 2025)
    fmt.Printf("Solar return: age %.1f\n", sr.Age)

    // Moon phase right now
    phase, _ := solarsage.MoonPhase("2025-03-18T12:00:00Z")
    fmt.Printf("Moon: %s (%.0f%% illuminated)\n", phase.PhaseName, phase.Illumination*100)

    // Eclipses in 2025
    eclipses, _ := solarsage.Eclipses("2025-01-01", "2026-01-01")
    for _, e := range eclipses {
        fmt.Printf("Eclipse: %s in %s\n", e.Type, e.MoonSign)
    }

    // Relationship compatibility
    score, _ := solarsage.Compatibility(
        51.5074, -0.1278, "1990-06-15T14:30:00Z",
        40.7128, -74.006, "1992-03-22T08:00:00Z",
    )
    fmt.Printf("Compatibility: %.0f%%\n", score.Compatibility)

    // Single planet position
    pos, _ := solarsage.PlanetPosition("Venus", "2025-03-18T12:00:00Z")
    fmt.Printf("Venus: %s at %.2f\n", pos.Sign, pos.SignDegree)

    // Vedic sidereal chart with Nakshatras
    vedic, _ := solarsage.SiderealChart(51.5074, -0.1278, "1990-06-15T14:30:00Z")
    for _, p := range vedic.Planets {
        fmt.Printf("%s: %s (Nakshatra: %s, Pada %d)\n",
            p.PlanetID, p.SiderealSign, p.Nakshatra, p.NakshatraPada)
    }

    // Vimshottari Dasha periods
    periods, _ := solarsage.Dasha(51.5074, -0.1278, "1990-06-15T14:30:00Z")
    for _, d := range periods {
        fmt.Printf("Age %.0f-%.0f: %s Dasha\n", d.StartAge, d.StartAge+d.Years, d.Lord)
    }

    // Chart wheel coordinates for SVG/Canvas rendering
    wheel, _ := solarsage.ChartWheel(51.5074, -0.1278, "1990-06-15T14:30:00Z")
    for _, p := range wheel.Planets {
        fmt.Printf("%s at (%.2f, %.2f)\n", p.PlanetID, p.Position.X, p.Position.Y)
    }

    // Comprehensive report (everything in one call)
    report, _ := solarsage.FullReport(51.5074, -0.1278, "1990-06-15T14:30:00Z")
    fmt.Printf("Elements: Fire=%d Earth=%d Air=%d Water=%d\n",
        report.ElementBalance["Fire"], report.ElementBalance["Earth"],
        report.ElementBalance["Air"], report.ElementBalance["Water"])
}
```

### Low-level API

Every calculation package can also be imported directly for full control:

```go
import (
    "github.com/shaobaobaoer/solarsage-mcp/pkg/chart"
    "github.com/shaobaobaoer/solarsage-mcp/pkg/models"
    "github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
    "github.com/shaobaobaoer/solarsage-mcp/pkg/transit"
)

func main() {
    sweph.Init("/path/to/ephe")
    defer sweph.Close()

    planets := []models.PlanetID{
        models.PlanetSun, models.PlanetMoon, models.PlanetVenus,
    }

    // Full control over orbs, house system, planet selection
    info, _ := chart.CalcSingleChart(
        51.5074, -0.1278, 2451545.0,
        planets, models.OrbConfig{Conjunction: 10, Trine: 8, Square: 8},
        models.HouseKoch,
    )

    // Transit search with all options
    events, _ := transit.CalcTransitEvents(transit.TransitCalcInput{
        NatalLat: 51.5074, NatalLon: -0.1278,
        NatalJD:  2451545.0, NatalPlanets: planets,
        TransitLat: 51.5074, TransitLon: -0.1278,
        StartJD: 2460676.5, EndJD: 2460706.5,
        TransitPlanets: planets,
        EventConfig:    models.DefaultEventConfig(),
        OrbConfigTransit: models.DefaultOrbConfig(),
        HouseSystem:    models.HousePlacidus,
    })
    _ = info
    _ = events
}
```

## MCP Tools (31)

| Tool | Description |
|------|-------------|
| **Utilities** | |
| `geocode` | Location name to coordinates and timezone |
| `datetime_to_jd` | ISO 8601 datetime to Julian Day (UT/TT) |
| `jd_to_datetime` | Julian Day to ISO 8601 datetime |
| **Chart Calculations** | |
| `calc_planet_position` | Single planet position at a given time |
| `calc_single_chart` | Full natal/event chart with positions, houses, and aspects |
| `calc_double_chart` | Synastry/transit double chart with cross-aspects |
| `calc_composite_chart` | Composite (midpoint) chart for relationships |
| `calc_harmonic_chart` | Nth harmonic (divisional) chart |
| **Predictive** | |
| `calc_transit` | Full transit event search over a time range (JSON or CSV) |
| `calc_progressions` | Secondary progressed planet positions |
| `calc_solar_arc` | Solar arc directed planet positions |
| `calc_solar_return` | Solar return chart (exact Sun return + full chart) |
| `calc_lunar_return` | Lunar return chart (exact Moon return + full chart) |
| `calc_profection` | Annual/monthly profections with time-lord |
| **Traditional** | |
| `calc_dignity` | Essential dignities, mutual receptions, and sect |
| `calc_lots` | Arabic lots (Fortune, Spirit, Eros, etc.) |
| `calc_bounds` | Chaldean decans and Egyptian terms |
| `calc_planetary_hours` | Chaldean planetary hours with sunrise/sunset |
| `calc_antiscia` | Antiscia and contra-antiscia mirror points |
| **Pattern Detection** | |
| `calc_aspect_patterns` | Grand Trine, T-Square, Yod, Grand Cross, etc. |
| `calc_fixed_stars` | Fixed star conjunctions (50+ star catalog) |
| `calc_midpoints` | Midpoint tree with 90deg dial sort and activations |
| **Astronomical** | |
| `calc_lunar_phase` | Lunar phase, illumination, and angle |
| `calc_lunar_phases` | Find new/full moons and quarters in date range |
| `calc_eclipses` | Solar and lunar eclipse finder |
| **Relationship** | |
| `calc_synastry` | Relationship compatibility scoring |
| **Analysis** | |
| `calc_dispositors` | Dispositorship chains and final dispositor |
| `calc_natal_report` | Comprehensive natal analysis (all techniques combined) |
| **Vedic / Sidereal** | |
| `calc_sidereal_chart` | Sidereal chart with Nakshatras and padas |
| `calc_vimshottari_dasha` | Vimshottari Maha Dasha periods |
| **Visualization** | |
| `calc_chart_wheel` | Chart wheel coordinates for SVG/Canvas rendering |

## Architecture

```
cmd/server/        MCP server entry point (JSON-RPC over stdio)
pkg/
  solarsage/       High-level convenience API (recommended entry point)
  mcp/             MCP protocol handler (31 tools)
  chart/           Chart calculations (positions, houses, aspects)
  transit/         Transit event detection engine
  progressions/    Secondary progressions & solar arc
  returns/         Solar & lunar return charts
  composite/       Composite (midpoint) charts
  synastry/        Synastry compatibility scoring
  dispositor/      Dispositorship chains & final dispositor
  report/          Comprehensive chart analysis report
  vedic/           Sidereal charts, Nakshatras, Vimshottari Dasha
  render/          Chart wheel visualization coordinates
  dignity/         Essential dignities, mutual receptions, sect
  fixedstars/      Fixed star catalog & conjunction detection
  midpoint/        Midpoint analysis & Cosmobiology dial
  harmonic/        Harmonic (divisional) charts
  planetary/       Planetary hours & day ruler
  profection/      Annual & monthly profections
  antiscia/        Antiscia & contra-antiscia
  lots/            Arabic lots/parts calculator
  bounds/          Decans & Egyptian terms
  lunar/           Lunar phases & eclipse detection
  models/          Core data types and constants
  julian/          Julian Day conversions
  geo/             Geocoding and timezone lookup
  export/          CSV/JSON export
  sweph/           Swiss Ephemeris C bindings (CGO, thread-safe)
internal/
  aspect/          Aspect calculation & pattern detection engine
```

## Performance

| Operation | Time | Throughput |
|-----------|------|------------|
| Planet position | 380ns | 2.6M/sec |
| Natal chart (10 planets) | 80us | 12,400/sec |
| Double chart + cross-aspects | 347us | 2,880/sec |
| 30-day transit scan (5 planets) | 764ms | - |
| 1-year transit scan (outer planets) | 2.1s | - |

Run `make bench` to reproduce.

## Accuracy

Independently validated with **100% exact event match** (247/247 transit events) over a 1-year period including all 7 chart type combinations, benchmarked against industry-standard desktop astrology software.

## Docker

```bash
docker build -t solarsage-mcp .
docker run -i solarsage-mcp
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
