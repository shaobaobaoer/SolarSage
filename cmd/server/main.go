package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anthropic/swisseph-mcp/pkg/mcp"
)

func main() {
	// Determine ephemeris path
	ephePath := os.Getenv("SWISSEPH_EPHE_PATH")
	if ephePath == "" {
		// Default: look for ephe directory relative to executable
		exe, err := os.Executable()
		if err == nil {
			ephePath = filepath.Join(filepath.Dir(exe), "..", "..", "third_party", "swisseph", "ephe")
		}
		// Fallback to current directory
		if _, err := os.Stat(ephePath); err != nil {
			ephePath = filepath.Join(".", "third_party", "swisseph", "ephe")
		}
	}

	if _, err := os.Stat(ephePath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: ephemeris path not found: %s\n", ephePath)
		fmt.Fprintf(os.Stderr, "Set SWISSEPH_EPHE_PATH environment variable or place ephe files correctly\n")
	}

	server := mcp.NewServer(ephePath)
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
