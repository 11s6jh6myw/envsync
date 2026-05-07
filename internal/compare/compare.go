// Package compare provides multi-file .env comparison across more than two
// environments, producing a unified view of keys and their values.
package compare

import (
	"sort"

	"github.com/yourusername/envsync/internal/parser"
)

// Result holds the outcome of comparing multiple env files.
type Result struct {
	// Keys is the union of all keys found across all files.
	Keys []string
	// Envs is the ordered list of environment labels.
	Envs []string
	// Matrix maps key → (env label → value). Missing entries mean the key is
	// absent in that environment.
	Matrix map[string]map[string]string
	// Conflicts is the set of keys whose values differ across at least two envs.
	Conflicts map[string]bool
}

// Compare takes a map of label → parsed entries and returns a unified Result.
func Compare(envs map[string][]parser.Entry) Result {
	keySet := map[string]bool{}
	matrix := map[string]map[string]string{}

	for label, entries := range envs {
		for _, e := range entries {
			if e.Key == "" {
				continue
			}
			keySet[e.Key] = true
			if matrix[e.Key] == nil {
				matrix[e.Key] = map[string]string{}
			}
			matrix[e.Key][label] = e.Value
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	labels := make([]string, 0, len(envs))
	for l := range envs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	conflicts := map[string]bool{}
	for key, vals := range matrix {
		if hasConflict(vals, labels) {
			conflicts[key] = true
		}
	}

	return Result{
		Keys:      keys,
		Envs:      labels,
		Matrix:    matrix,
		Conflicts: conflicts,
	}
}

func hasConflict(vals map[string]string, labels []string) bool {
	if len(vals) < 2 {
		return true // missing in at least one env
	}
	first := vals[labels[0]]
	for _, l := range labels[1:] {
		if v, ok := vals[l]; !ok || v != first {
			return true
		}
	}
	return false
}
