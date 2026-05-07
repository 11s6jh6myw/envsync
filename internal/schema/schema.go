// Package schema provides functionality for defining and validating
// expected keys in a .env file against a schema definition.
package schema

import (
	"fmt"
	"strings"

	"github.com/your-org/envsync/internal/parser"
)

// Rule describes expectations for a single key.
type Rule struct {
	Key         string
	Required    bool
	Description string
	Allowed     []string // if non-empty, value must be one of these
}

// Issue represents a schema violation found during Check.
type Issue struct {
	Key     string
	Message string
	Error   bool // true = error, false = warning
}

// Schema holds a collection of rules keyed by variable name.
type Schema struct {
	rules map[string]Rule
}

// New creates a Schema from a slice of Rules.
func New(rules []Rule) *Schema {
	m := make(map[string]Rule, len(rules))
	for _, r := range rules {
		m[strings.ToUpper(r.Key)] = r
	}
	return &Schema{rules: m}
}

// Check validates a set of parsed entries against the schema.
// It returns issues for missing required keys and disallowed values.
func (s *Schema) Check(entries []parser.Entry) []Issue {
	var issues []Issue

	present := make(map[string]string)
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		present[strings.ToUpper(e.Key)] = e.Value
	}

	// Check required keys and allowed values.
	for upperKey, rule := range s.rules {
		val, exists := present[upperKey]
		if rule.Required && !exists {
			issues = append(issues, Issue{
				Key:     rule.Key,
				Message: "required key is missing",
				Error:   true,
			})
			continue
		}
		if exists && len(rule.Allowed) > 0 {
			if !contains(rule.Allowed, val) {
				issues = append(issues, Issue{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q not in allowed set [%s]", val, strings.Join(rule.Allowed, ", ")),
					Error:   true,
				})
			}
		}
	}

	// Warn about keys not defined in schema.
	for upperKey := range present {
		if _, defined := s.rules[upperKey]; !defined {
			issues = append(issues, Issue{
				Key:     upperKey,
				Message: "key is not defined in schema",
				Error:   false,
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue is an error-level violation.
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Error {
			return true
		}
	}
	return false
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
