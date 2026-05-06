// Package lint provides style and convention checks for .env files.
package lint

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envsync/internal/parser"
)

// Severity represents the level of a lint issue.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a single lint finding.
type Issue struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

// Options controls which lint rules are enabled.
type Options struct {
	RequireUppercase bool // keys should be UPPER_SNAKE_CASE
	ForbidSpaces     bool // values must not contain unquoted leading/trailing spaces
	MaxKeyLength     int  // 0 means no limit
	ForbidEmptyValue bool // keys must have non-empty values
}

// DefaultOptions returns sensible lint defaults.
func DefaultOptions() Options {
	return Options{
		RequireUppercase: true,
		ForbidSpaces:     true,
		MaxKeyLength:     64,
		ForbidEmptyValue: false,
	}
}

var upperSnake = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint runs style checks against a slice of parsed entries and returns all findings.
func Lint(entries []parser.Entry, opts Options) []Issue {
	var issues []Issue

	for i, e := range entries {
		if e.Key == "" {
			continue // blank line or comment
		}
		lineNum := i + 1

		if opts.RequireUppercase && !upperSnake.MatchString(e.Key) {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q should be UPPER_SNAKE_CASE", e.Key),
				Severity: SeverityWarning,
			})
		}

		if opts.MaxKeyLength > 0 && len(e.Key) > opts.MaxKeyLength {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q exceeds max length of %d", e.Key, opts.MaxKeyLength),
				Severity: SeverityError,
			})
		}

		if opts.ForbidSpaces && (strings.HasPrefix(e.Value, " ") || strings.HasSuffix(e.Value, " ")) {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("value for %q has leading or trailing spaces", e.Key),
				Severity: SeverityWarning,
			})
		}

		if opts.ForbidEmptyValue && e.Value == "" {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q has an empty value", e.Key),
				Severity: SeverityError,
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue has severity Error.
func HasErrors(issues []Issue) bool {
	for _, iss := range issues {
		if iss.Severity == SeverityError {
			return true
		}
	}
	return false
}
