# Contributing

Contributions are welcome! Here's how to get started.

## Development Setup

```bash
git clone https://github.com/anthropic/swisseph-mcp.git
cd swisseph-mcp
make build
make test
```

### Prerequisites

- Go 1.21+
- GCC (for CGO / Swiss Ephemeris)

## Making Changes

1. Fork the repo and create a feature branch
2. Make your changes
3. Run `make test` and `make vet` to ensure quality
4. Submit a pull request

## Guidelines

- Keep changes focused and minimal
- Add tests for new functionality
- Maintain backward compatibility for the MCP tool API
- Changes to `pkg/transit/transit.go` must pass the Solar Fire validation test

## Architecture

See [CLAUDE.md](CLAUDE.md) for architecture details and key constraints.
