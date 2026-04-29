package sync

import (
	"fmt"

	"github.com/yourusername/envsync/internal/diff"
	"github.com/yourusername/envsync/internal/parser"
)

// Strategy defines how missing or extra keys are handled during sync.
type Strategy int

const (
	// StrategyFill adds missing keys from source into target, leaving extras alone.
	StrategyFill Strategy = iota
	// StrategyOverwrite updates existing keys and adds missing ones.
	StrategyOverwrite
	// StrategyExact makes target exactly match source (adds missing, removes extra).
	StrategyExact
)

// Options controls sync behaviour.
type Options struct {
	Strategy Strategy
	// Placeholder is written as the value for newly added keys.
	Placeholder string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:    StrategyFill,
		Placeholder: "",
	}
}

// Apply takes a source and target entry slice, applies the sync strategy and
// returns the updated target entries ready to be written.
func Apply(source, target []parser.Entry, opts Options) ([]parser.Entry, error) {
	srcMap := toMap(source)
	tgtMap := toMap(target)

	// Work on a copy of target so we preserve order and comments.
	result := make([]parser.Entry, len(target))
	copy(result, target)

	// Update or fill existing keys.
	for i, e := range result {
		if e.Key == "" {
			continue
		}
		srcVal, exists := srcMap[e.Key]
		if !exists {
			if opts.Strategy == StrategyExact {
				// Mark for removal by zeroing the key; we filter below.
				result[i].Key = ""
				result[i].Value = ""
			}
			continue
		}
		if opts.Strategy == StrategyOverwrite || opts.Strategy == StrategyExact {
			result[i].Value = srcVal
		}
	}

	// Add keys present in source but missing in target.
	for _, e := range source {
		if e.Key == "" {
			continue
		}
		if _, exists := tgtMap[e.Key]; !exists {
			newEntry := parser.Entry{
				Key:   e.Key,
				Value: opts.Placeholder,
			}
			result = append(result, newEntry)
		}
	}

	// Filter removed entries for StrategyExact.
	if opts.Strategy == StrategyExact {
		filtered := result[:0]
		for _, e := range result {
			if e.Key != "" || e.Raw != "" {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	return result, nil
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}

// Summary returns a human-readable summary of what Apply would change.
func Summary(source, target []parser.Entry, opts Options) string {
	changes := diff.Diff(source, target)
	var added, removed, modified int
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			added++
		case diff.Removed:
			removed++
		case diff.Modified:
			modified++
		}
	}
	return fmt.Sprintf("sync summary: +%d added, -%d removed, ~%d modified", added, removed, modified)
}
