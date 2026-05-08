// Package rename provides utilities for renaming keys in .env files
// while preserving comments, ordering, and value formatting.
package rename

import "github.com/user/envsync/internal/parser"

// Options controls the behaviour of the Rename operation.
type Options struct {
	// FailOnMissing causes Rename to return an error when a source key is not
	// found in the entries instead of silently skipping it.
	FailOnMissing bool
}

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
	}
}

// Result holds the outcome of a Rename call.
type Result struct {
	// Entries is the updated slice (same order, same length as input).
	Entries []parser.Entry
	// Renamed contains every key that was successfully renamed, mapping old
	// name to new name.
	Renamed map[string]string
	// Skipped contains source keys that were not found in the entries.
	Skipped []string
}

// Rename applies the given mapping (oldKey -> newKey) to entries.
// Comment and blank entries are passed through unchanged.
// If a key appears in the mapping but is not present in entries it is added to
// Result.Skipped (or an error is returned when opts.FailOnMissing is true).
func Rename(entries []parser.Entry, mapping map[string]string, opts Options) (Result, error) {
	result := Result{
		Renamed: make(map[string]string),
	}

	// Build a quick lookup of which old keys exist.
	existing := make(map[string]bool, len(entries))
	for _, e := range entries {
		if !e.IsComment && !e.IsBlank {
			existing[e.Key] = true
		}
	}

	for oldKey, newKey := range mapping {
		if !existing[oldKey] {
			if opts.FailOnMissing {
				return Result{}, &MissingKeyError{Key: oldKey}
			}
			result.Skipped = append(result.Skipped, oldKey)
		}
	}

	updated := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if !e.IsComment && !e.IsBlank {
			if newKey, ok := mapping[e.Key]; ok {
				result.Renamed[e.Key] = newKey
				e.Key = newKey
			}
		}
		updated[i] = e
	}

	result.Entries = updated
	return result, nil
}

// MissingKeyError is returned when FailOnMissing is true and a source key
// cannot be found.
type MissingKeyError struct {
	Key string
}

func (e *MissingKeyError) Error() string {
	return "rename: key not found: " + e.Key
}
