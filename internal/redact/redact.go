// Package redact provides utilities for redacting sensitive values
// in .env file entries based on configurable key patterns.
package redact

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns is the default list of patterns that indicate
// a key contains sensitive data.
var DefaultSensitivePatterns = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
}

// Redactor determines whether a key is sensitive and masks its value.
type Redactor struct {
	patterns []*regexp.Regexp
	mask     string
}

// New creates a Redactor using the provided key patterns (case-insensitive).
// If patterns is empty, DefaultSensitivePatterns is used.
func New(patterns []string, mask string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}
	if mask == "" {
		mask = "***"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(p))
		if err == nil {
			compiled = append(compiled, re)
		}
	}
	return &Redactor{patterns: compiled, mask: mask}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, re := range r.patterns {
		if re.MatchString(upper) {
			return true
		}
	}
	return false
}

// Redact returns the masked value if the key is sensitive, otherwise the original value.
func (r *Redactor) Redact(key, value string) string {
	if r.IsSensitive(key) {
		return r.mask
	}
	return value
}

// Mask returns the configured mask string.
func (r *Redactor) Mask() string {
	return r.mask
}
