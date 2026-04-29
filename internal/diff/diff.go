package diff

import "github.com/user/envsync/internal/parser"

// ChangeType represents the type of difference between two env files.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single difference between two env files.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between a source and target env file.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any differences.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Diff computes the difference between a source and target set of env entries.
// Keys present in source but not target are Added.
// Keys present in target but not source are Removed.
// Keys present in both but with different values are Modified.
func Diff(source, target []parser.Entry) *Result {
	sourceMap := toMap(source)
	targetMap := toMap(target)

	var changes []Change

	for _, entry := range source {
		if !entry.IsComment && entry.Key != "" {
			if targetVal, ok := targetMap[entry.Key]; !ok {
				changes = append(changes, Change{
					Key:      entry.Key,
					Type:     Added,
					NewValue: entry.Value,
				})
			} else if targetVal != entry.Value {
				changes = append(changes, Change{
					Key:      entry.Key,
					Type:     Modified,
					OldValue: targetVal,
					NewValue: entry.Value,
				})
			}
		}
	}

	for _, entry := range target {
		if !entry.IsComment && entry.Key != "" {
			if _, ok := sourceMap[entry.Key]; !ok {
				changes = append(changes, Change{
					Key:      entry.Key,
					Type:     Removed,
					OldValue: entry.Value,
				})
			}
		}
	}

	return &Result{Changes: changes}
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
