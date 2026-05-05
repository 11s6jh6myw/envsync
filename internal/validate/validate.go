// Package validate checks .env files for common issues such as
// missing required keys, duplicate keys, and malformed entries.
package validate

import (
	"fmt"
	"strings"

	"github.com/yourusername/envsync/internal/parser"
)

// Issue represents a single validation problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
	Severity string // "error" or "warning"
}

func (i Issue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("%s (line %d): %s", i.Severity, i.Line, i.Message)
	}
	return fmt.Sprintf("%s: %s", i.Severity, i.Message)
}

// IsError reports whether the issue has error severity.
func (i Issue) IsError() bool {
	return i.Severity == "error"
}

// HasErrors returns true if any of the provided issues have error severity.
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.IsError() {
			return true
		}
	}
	return false
}

// Options controls which checks are performed.
type Options struct {
	RequiredKeys  []string
	ForbidEmpty   bool
	ForbidDupes   bool
}

// DefaultOptions returns sensible validation defaults.
func DefaultOptions() Options {
	return Options{
		ForbidEmpty: true,
		ForbidDupes: true,
	}
}

// Validate runs all enabled checks against the provided entries and returns
// any issues found. A nil or empty slice means the file is valid.
func Validate(entries []parser.Entry, opts Options) []Issue {
	var issues []Issue

	seen := make(map[string]int) // key -> first line number
	present := make(map[string]bool)

	for _, e := range entries {
		if e.Comment || e.Blank {
			continue
		}

		norm := strings.ToUpper(e.Key)
		present[norm] = true

		if opts.ForbidDupes {
			if first, dup := seen[norm]; dup {
				issues = append(issues, Issue{
					Line:     e.Line,
					Key:      e.Key,
					Message:  fmt.Sprintf("duplicate key %q (first seen at line %d)", e.Key, first),
					Severity: "error",
				})
			} else {
				seen[norm] = e.Line
			}
		}

		if opts.ForbidEmpty && e.Value == "" {
			issues = append(issues, Issue{
				Line:     e.Line,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q has an empty value", e.Key),
				Severity: "warning",
			})
		}
	}

	for _, req := range opts.RequiredKeys {
		if !present[strings.ToUpper(req)] {
			issues = append(issues, Issue{
				Message:  fmt.Sprintf("required key %q is missing", req),
				Severity: "error",
			})
		}
	}

	return issues
}
