// Package clone provides functionality to clone a .env file into a new
// target file, optionally redacting sensitive values or replacing them
// with placeholders.
package clone

import (
	"fmt"

	"github.com/your-org/envsync/internal/parser"
)

// Options controls how a clone operation is performed.
type Options struct {
	// Redact replaces sensitive values with a redaction marker.
	Redact bool

	// Placeholder is used when Redact is true. Defaults to "REDACTED".
	Placeholder string

	// SensitiveKeys is a list of key substrings considered sensitive.
	// If empty, a built-in default list is used.
	SensitiveKeys []string

	// StripComments removes comment and blank lines from the output.
	StripComments bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Redact:      false,
		Placeholder: "REDACTED",
	}
}

// Summary describes the result of a Clone operation.
type Summary struct {
	Total    int
	Redacted int
	Skipped  int // blank/comment lines stripped
}

// Clone duplicates src entries into a new slice, applying options.
func Clone(src []parser.Entry, opts Options) ([]parser.Entry, Summary, error) {
	if opts.Placeholder == "" {
		opts.Placeholder = "REDACTED"
	}

	sensitive := opts.SensitiveKeys
	if len(sensitive) == 0 {
		sensitive = defaultSensitiveKeys
	}

	var out []parser.Entry
	var sum Summary

	for _, e := range src {
		if e.Key == "" {
			if opts.StripComments {
				sum.Skipped++
				continue
			}
			out = append(out, e)
			continue
		}

		cloned := parser.Entry{
			Key:     e.Key,
			Value:   e.Value,
			Comment: e.Comment,
		}

		if opts.Redact && isSensitive(e.Key, sensitive) {
			cloned.Value = opts.Placeholder
			sum.Redacted++
		}

		out = append(out, cloned)
		sum.Total++
	}

	return out, sum, nil
}

func isSensitive(key string, patterns []string) bool {
	lower := toLower(key)
	for _, p := range patterns {
		if contains(lower, toLower(p)) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) &&
		(s == sub || len(s) > 0 && fmt.Sprintf("%s", s) != "" &&
			func() bool {
				for i := 0; i <= len(s)-len(sub); i++ {
					if s[i:i+len(sub)] == sub {
						return true
					}
				}
				return false
			}())
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

var defaultSensitiveKeys = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "credential", "auth", "key",
}
