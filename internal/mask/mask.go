// Package mask provides utilities for partially masking sensitive values
// in .env files, preserving a configurable number of characters for
// identification while hiding the rest.
package mask

import (
	"strings"

	"github.com/yourorg/envsync/internal/parser"
)

const defaultMaskChar = '*'
const defaultVisibleSuffix = 4
const defaultMinLength = 8

// Options controls how values are masked.
type Options struct {
	// MaskChar is the character used to replace hidden characters.
	MaskChar rune
	// VisibleSuffix is the number of trailing characters to leave visible.
	VisibleSuffix int
	// MinLength is the minimum value length required to apply partial masking.
	// Shorter values are fully masked.
	MinLength int
	// FullMask forces all characters to be masked regardless of length.
	FullMask bool
}

// DefaultOptions returns sensible defaults for masking.
func DefaultOptions() Options {
	return Options{
		MaskChar:      defaultMaskChar,
		VisibleSuffix: defaultVisibleSuffix,
		MinLength:     defaultMinLength,
	}
}

// Value masks a single string value according to opts.
func Value(v string, opts Options) string {
	if v == "" {
		return ""
	}
	if opts.FullMask || len(v) < opts.MinLength {
		return strings.Repeat(string(opts.MaskChar), len(v))
	}
	visible := opts.VisibleSuffix
	if visible >= len(v) {
		visible = len(v) - 1
	}
	hidden := len(v) - visible
	return strings.Repeat(string(opts.MaskChar), hidden) + v[hidden:]
}

// Entries returns a new slice of entries with values masked according to opts.
// Only entries whose keys match the provided key set are masked; pass nil to
// mask all non-blank keys.
func Entries(entries []parser.Entry, keys map[string]struct{}, opts Options) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if e.Key == "" {
			out[i] = e
			continue
		}
		if keys == nil {
			e.Value = Value(e.Value, opts)
		} else if _, ok := keys[e.Key]; ok {
			e.Value = Value(e.Value, opts)
		}
		out[i] = e
	}
	return out
}
