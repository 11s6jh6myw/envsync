// Package extract provides utilities for extracting a subset of key-value
// entries from a parsed .env file based on a list of key names.
package extract

import (
	"fmt"

	"github.com/user/envsync/internal/parser"
)

// Options controls the behaviour of Extract.
type Options struct {
	// IgnoreMissing suppresses errors when a requested key is absent.
	IgnoreMissing bool
	// PreserveOrder returns entries in the order they appear in src rather
	// than the order of the keys slice.
	PreserveOrder bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		IgnoreMissing: false,
		PreserveOrder: true,
	}
}

// Result holds the output of an Extract call.
type Result struct {
	Entries []parser.Entry
	// Missing contains keys that were requested but not found in src.
	Missing []string
}

// Extract pulls the entries whose keys appear in keys from src.
// Comment and blank entries are never included in the output.
func Extract(src []parser.Entry, keys []string, opts Options) (Result, error) {
	want := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		want[k] = struct{}{}
	}

	var result Result
	found := make(map[string]bool, len(keys))

	if opts.PreserveOrder {
		for _, e := range src {
			if e.Key == "" {
				continue
			}
			if _, ok := want[e.Key]; ok {
				result.Entries = append(result.Entries, e)
				found[e.Key] = true
			}
		}
	} else {
		index := make(map[string]parser.Entry, len(src))
		for _, e := range src {
			if e.Key != "" {
				index[e.Key] = e
			}
		}
		for _, k := range keys {
			if e, ok := index[k]; ok {
				result.Entries = append(result.Entries, e)
				found[k] = true
			}
		}
	}

	for _, k := range keys {
		if !found[k] {
			result.Missing = append(result.Missing, k)
		}
	}

	if !opts.IgnoreMissing && len(result.Missing) > 0 {
		return result, fmt.Errorf("extract: keys not found: %v", result.Missing)
	}

	return result, nil
}
