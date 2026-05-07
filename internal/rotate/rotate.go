// Package rotate provides utilities for rotating secret values in .env files.
// It replaces the values of sensitive keys with newly generated or provided
// values, preserving structure and non-sensitive keys.
package rotate

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/redact"
)

// Options controls rotation behaviour.
type Options struct {
	// Replacements maps key names to their new values.
	// Keys not present in this map are left unchanged.
	Replacements map[string]string

	// OnlySensitive restricts rotation to keys flagged as sensitive
	// by the default redact patterns. Keys in Replacements are always rotated.
	OnlySensitive bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Replacements:  make(map[string]string),
		OnlySensitive: false,
	}
}

// Result holds the outcome of a rotation pass.
type Result struct {
	Entries  []parser.Entry
	Rotated  []string // keys that were rotated
	Skipped  []string // keys that were candidates but had no replacement
}

// Rotate applies replacements to entries according to opts.
// It returns a Result containing the updated entries and metadata.
func Rotate(entries []parser.Entry, opts Options) (Result, error) {
	if opts.Replacements == nil {
		opts.Replacements = make(map[string]string)
	}

	redactor := redact.New(nil)

	result := Result{
		Entries: make([]parser.Entry, len(entries)),
	}
	copy(result.Entries, entries)

	for i, e := range result.Entries {
		if e.Key == "" {
			continue
		}

		newVal, explicit := opts.Replacements[e.Key]
		if !explicit && opts.OnlySensitive && !redactor.IsSensitive(e.Key) {
			continue
		}

		if !explicit {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		result.Entries[i].Value = newVal
		result.Rotated = append(result.Rotated, e.Key)
	}

	return result, nil
}

// Summary returns a human-readable one-line summary of the result.
func Summary(r Result) string {
	parts := []string{
		fmt.Sprintf("%d rotated", len(r.Rotated)),
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped (no replacement provided)", len(r.Skipped)))
	}
	return strings.Join(parts, ", ")
}
