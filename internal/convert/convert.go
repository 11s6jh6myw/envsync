// Package convert provides utilities for converting .env files to and from
// other common configuration formats such as JSON, YAML-style key=value exports,
// and shell export statements.
package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourorg/envsync/internal/parser"
)

// Format represents a supported output format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatExport Format = "export"
	FormatDotenv Format = "dotenv"
)

// ToJSON converts a slice of env entries to a JSON object string.
func ToJSON(entries []parser.Entry) (string, error) {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("convert: marshal json: %w", err)
	}
	return string(b), nil
}

// ToExport converts a slice of env entries to shell export statements.
// Values containing spaces or special characters are quoted.
func ToExport(entries []parser.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		value := e.Value
		if needsQuoting(value) {
			value = `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
		}
		fmt.Fprintf(&sb, "export %s=%s\n", e.Key, value)
	}
	return sb.String()
}

// FromJSON parses a flat JSON object into a slice of env entries.
func FromJSON(data string) ([]parser.Entry, error) {
	var m map[string]string
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, fmt.Errorf("convert: unmarshal json: %w", err)
	}
	entries := make([]parser.Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, parser.Entry{Key: k, Value: v})
	}
	return entries, nil
}

// needsQuoting reports whether a value should be wrapped in double quotes
// when rendered as a shell export.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '$' || c == '&' || c == '|' || c == ';' {
			return true
		}
	}
	return false
}
