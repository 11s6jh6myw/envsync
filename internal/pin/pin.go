// Package pin provides functionality to lock (pin) specific env entries
// to fixed values, preventing them from being overwritten during sync or merge.
package pin

import (
	"fmt"

	"github.com/yourusername/envsync/internal/parser"
)

// PinnedKey represents a key that is locked to a specific value.
type PinnedKey struct {
	Key   string
	Value string
}

// Result holds the outcome of a pin operation.
type Result struct {
	Pinned   []PinnedKey
	Skipped  []string // keys not found in entries
	Entries  []parser.Entry
}

// DefaultOptions returns sensible defaults for Pin.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
	}
}

// Options controls Pin behaviour.
type Options struct {
	// FailOnMissing causes Pin to return an error if a requested key is absent.
	FailOnMissing bool
}

// Pin locks the given keys to their current values in entries, returning
// a new slice where those keys are marked with a "# pinned" inline comment.
// Keys not present in entries are reported in Result.Skipped.
func Pin(entries []parser.Entry, keys []string, opts Options) (Result, error) {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	res := Result{}
	pinned := make(map[string]bool)
	out := make([]parser.Entry, 0, len(entries))

	for _, e := range entries {
		if keySet[e.Key] {
			e.Comment = "pinned"
			res.Pinned = append(res.Pinned, PinnedKey{Key: e.Key, Value: e.Value})
			pinned[e.Key] = true
		}
		out = append(out, e)
	}

	for _, k := range keys {
		if !pinned[k] {
			if opts.FailOnMissing {
				return Result{}, fmt.Errorf("pin: key %q not found in entries", k)
			}
			res.Skipped = append(res.Skipped, k)
		}
	}

	res.Entries = out
	return res, nil
}

// IsPinned reports whether the given entry carries a pin marker.
func IsPinned(e parser.Entry) bool {
	return e.Comment == "pinned"
}
