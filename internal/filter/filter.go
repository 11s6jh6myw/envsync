// Package filter provides utilities for selecting a subset of env entries
// based on key prefix, suffix, or pattern matching.
package filter

import (
	"strings"

	"github.com/yourusername/envsync/internal/parser"
)

// Options controls how entries are filtered.
type Options struct {
	// Prefix retains only keys that start with the given string (case-insensitive).
	Prefix string

	// Suffix retains only keys that end with the given string (case-insensitive).
	Suffix string

	// Contains retains only keys that contain the given substring (case-insensitive).
	Contains string

	// Invert inverts the match — retains entries that do NOT match.
	Invert bool

	// KeepComments preserves comment and blank entries regardless of match.
	KeepComments bool
}

// DefaultOptions returns an Options with no filters applied.
func DefaultOptions() Options {
	return Options{
		KeepComments: true,
	}
}

// Result holds the output of a Filter call.
type Result struct {
	Entries []parser.Entry
	Matched int
	Skipped int
}

// Filter returns a subset of entries according to opts.
// Comment and blank entries are passed through when KeepComments is true.
func Filter(entries []parser.Entry, opts Options) Result {
	var out []parser.Entry
	matched, skipped := 0, 0

	for _, e := range entries {
		if e.IsComment || e.IsBlank {
			if opts.KeepComments {
				out = append(out, e)
			}
			continue
		}

		if matches(e.Key, opts) {
			out = append(out, e)
			matched++
		} else {
			skipped++
		}
	}

	return Result{Entries: out, Matched: matched, Skipped: skipped}
}

func matches(key string, opts Options) bool {
	k := strings.ToLower(key)

	ok := true
	if opts.Prefix != "" {
		ok = ok && strings.HasPrefix(k, strings.ToLower(opts.Prefix))
	}
	if opts.Suffix != "" {
		ok = ok && strings.HasSuffix(k, strings.ToLower(opts.Suffix))
	}
	if opts.Contains != "" {
		ok = ok && strings.Contains(k, strings.ToLower(opts.Contains))
	}

	if opts.Invert {
		return !ok
	}
	return ok
}
