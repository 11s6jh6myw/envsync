// Package template provides functionality to generate .env template files
// from existing env entries, replacing values with placeholder descriptions.
package template

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/parser"
)

// Options controls template generation behaviour.
type Options struct {
	// Placeholder is the value written for each key in the template.
	// Defaults to "CHANGE_ME" when empty.
	Placeholder string

	// KeepComments preserves comment lines from the source entries.
	KeepComments bool

	// KeepValues retains the original values instead of replacing them.
	KeepValues bool
}

// DefaultOptions returns sensible defaults for template generation.
func DefaultOptions() Options {
	return Options{
		Placeholder:  "CHANGE_ME",
		KeepComments: true,
		KeepValues:   false,
	}
}

// Generate takes a slice of parsed entries and returns a new slice where
// each key's value is replaced with the configured placeholder.
func Generate(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.Placeholder == "" {
		opts.Placeholder = "CHANGE_ME"
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if e.Blank || (e.Comment != "" && e.Key == "") {
			if opts.KeepComments {
				out = append(out, e)
			}
			continue
		}
		if e.Key == "" {
			continue
		}
		copy := e
		if !opts.KeepValues {
			copy.Value = opts.Placeholder
			copy.Comment = describeKey(e.Key)
		}
		out = append(out, copy)
	}
	return out
}

// describeKey produces a short inline comment hint from a key name.
func describeKey(key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", " ")
	return fmt.Sprintf("set your %s", key)
}
