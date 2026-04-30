// Package env provides utilities for loading .env files from disk
// and resolving environment-specific file paths.
package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envsync/internal/parser"
)

// LoadResult holds the parsed entries and the resolved file path.
type LoadResult struct {
	Path    string
	Entries []parser.Entry
}

// Load reads and parses a .env file at the given path.
func Load(path string) (*LoadResult, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("env: resolving path %q: %w", path, err)
	}

	f, err := os.Open(abs)
	if err != nil {
		return nil, fmt.Errorf("env: opening %q: %w", abs, err)
	}
	defer f.Close()

	entries, err := parser.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("env: parsing %q: %w", abs, err)
	}

	return &LoadResult{Path: abs, Entries: entries}, nil
}

// LoadPair loads two .env files and returns both results.
// Useful for diff/sync workflows where a source and target are required.
func LoadPair(sourcePath, targetPath string) (*LoadResult, *LoadResult, error) {
	src, err := Load(sourcePath)
	if err != nil {
		return nil, nil, fmt.Errorf("env: loading source: %w", err)
	}

	tgt, err := Load(targetPath)
	if err != nil {
		return nil, nil, fmt.Errorf("env: loading target: %w", err)
	}

	return src, tgt, nil
}

// Exists reports whether a file exists at the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
