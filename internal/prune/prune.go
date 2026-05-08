// Package prune removes stale or unused keys from a .env file based on a
// reference set of allowed keys.
package prune

import "github.com/yourusername/envsync/internal/parser"

// Options controls Prune behaviour.
type Options struct {
	// DryRun reports what would be removed without modifying the slice.
	DryRun bool
	// KeepComments preserves comment/blank entries even when adjacent keys are
	// pruned.
	KeepComments bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:       false,
		KeepComments: true,
	}
}

// Summary describes the result of a Prune operation.
type Summary struct {
	Removed []string
	Kept    []string
}

// Prune filters entries, keeping only those whose keys appear in allowed.
// Comment and blank entries are kept or dropped according to opts.KeepComments.
// When opts.DryRun is true the original entries slice is returned unchanged but
// Summary is still populated.
func Prune(entries []parser.Entry, allowed []string, opts Options) ([]parser.Entry, Summary) {
	allowSet := make(map[string]struct{}, len(allowed))
	for _, k := range allowed {
		allowSet[k] = struct{}{}
	}

	var result []parser.Entry
	var summary Summary

	for _, e := range entries {
		if e.Key == "" {
			// blank line or comment
			if opts.KeepComments {
				result = append(result, e)
			}
			continue
		}
		if _, ok := allowSet[e.Key]; ok {
			result = append(result, e)
			summary.Kept = append(summary.Kept, e.Key)
		} else {
			summary.Removed = append(summary.Removed, e.Key)
		}
	}

	if opts.DryRun {
		return entries, summary
	}
	return result, summary
}
