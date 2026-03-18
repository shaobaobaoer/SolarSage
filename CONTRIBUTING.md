# Contributing

Contributions are welcome! Here's how to get started.

## Development Setup

```bash
git clone https://github.com/shaobaobaoer/solarsage-mcp.git
cd solarsage-mcp
make build
make check   # runs vet + test
```

### Prerequisites

- Go 1.25+
- GCC (for CGO / Swiss Ephemeris)

## Making Changes

1. Fork the repo and create a feature branch
2. Make your changes
3. Run `make check` to ensure quality
4. Run `make test-race` to verify no race conditions
5. Submit a pull request

## Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build binary |
| `make test` | Run all tests |
| `make test-race` | Run tests with race detector |
| `make test-cover` | Show coverage summary |
| `make cover-html` | Generate HTML coverage report |
| `make bench` | Run benchmarks |
| `make vet` | Run go vet |
| `make check` | Run vet + test |

## Guidelines

- Keep changes focused and minimal
- Add tests for new functionality
- Maintain backward compatibility for the MCP tool API
- Changes to `pkg/transit/transit.go` must pass the accuracy validation test (247/247 events)
- All packages must pass race detection (`make test-race`)
- Target 80%+ test coverage for new packages

## Architecture

See [CLAUDE.md](CLAUDE.md) for architecture details and key constraints.
