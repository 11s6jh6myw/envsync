// Package merge provides utilities for merging two .env file entry sets.
// It supports multiple strategies: prefer source, prefer target, or interactive.
package merge

import (
	"github.com/user/envsync/internal/parser"
)

// Strategy controls how conflicts are resolved during a merge.
type Strategy int

const (
	// PreferSource uses the source value when a key exists in both files.
	PreferSource Strategy = iota
	// PreferTarget keeps the target value when a key exists in both files.
	PreferTarget
	// Union includes all keys from both files; source wins on conflict.
	Union
)

// Result holds the output of a merge operation.
type Result struct {
	Entries  []parser.Entry
	Conflicts []Conflict
}

// Conflict describes a key that existed in both source and target with different values.
type Conflict struct {
	Key         string
	SourceValue string
	TargetValue string
	Resolved    string
}

// Merge combines source and target entry slices using the given strategy.
func Merge(source, target []parser.Entry, strategy Strategy) Result {
	srcMap := toMap(source)
	tgtMap := toMap(target)

	var result Result
	seen := make(map[string]bool)

	// Walk source entries to preserve ordering.
	for _, e := range source {
		if !e.IsComment && e.Key != "" {
			seen[e.Key] = true
			tgtVal, inTarget := tgtMap[e.Key]
			if inTarget && tgtVal != e.Value {
				conflict := Conflict{
					Key:         e.Key,
					SourceValue: e.Value,
					TargetValue: tgtVal,
				}
				resolved := e.Value
				if strategy == PreferTarget {
					resolved = tgtVal
				}
				conflict.Resolved = resolved
				result.Conflicts = append(result.Conflicts, conflict)
				e.Value = resolved
			}
		}
		result.Entries = append(result.Entries, e)
	}

	// For Union strategy, append keys only in target.
	if strategy == Union {
		for _, e := range target {
			if !e.IsComment && e.Key != "" && !seen[e.Key] {
				result.Entries = append(result.Entries, e)
			}
		}
		// Also include target-only comments/blanks? Skip for simplicity.
	}

	// For PreferTarget with keys only in source, they are already added above.
	// Keys only in target and not Union are dropped (source is authoritative).
	if strategy == PreferTarget {
		for _, e := range target {
			if !e.IsComment && e.Key != "" && !seen[e.Key] {
				result.Entries = append(result.Entries, e)
			}
		}
	}

	_ = srcMap
	return result
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.IsComment && e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}
