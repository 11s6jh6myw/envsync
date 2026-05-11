package prefix

import "github.com/jasonuc/envsync/internal/parser"

// Options controls how prefix operations are applied.
type Options struct {
	// Strip removes the given prefix from matching keys instead of adding it.
	Strip bool
	// SkipNoMatch silently ignores keys that don't match when stripping.
	SkipNoMatch bool
}

// DefaultOptions returns sensible defaults for prefix operations.
func DefaultOptions() Options {
	return Options{
		Strip:       false,
		SkipNoMatch: true,
	}
}

// Summary holds the result of a prefix operation.
type Summary struct {
	Modified int
	Skipped  int
}

// Add prepends pfx to every key in entries and returns the updated slice
// along with a summary of changes.
func Add(entries []parser.Entry, pfx string, opts Options) ([]parser.Entry, Summary) {
	if pfx == "" {
		return entries, Summary{}
	}
	out := make([]parser.Entry, len(entries))
	var s Summary
	for i, e := range entries {
		if e.Key == "" {
			out[i] = e
			continue
		}
		out[i] = e
		out[i].Key = pfx + e.Key
		s.Modified++
	}
	return out, s
}

// Strip removes pfx from the start of every key in entries. Keys that do not
// carry the prefix are skipped when opts.SkipNoMatch is true, otherwise they
// are left unchanged and counted as skipped.
func Strip(entries []parser.Entry, pfx string, opts Options) ([]parser.Entry, Summary) {
	if pfx == "" {
		return entries, Summary{}
	}
	out := make([]parser.Entry, 0, len(entries))
	var s Summary
	for _, e := range entries {
		if e.Key == "" {
			out = append(out, e)
			continue
		}
		if len(e.Key) >= len(pfx) && e.Key[:len(pfx)] == pfx {
			e.Key = e.Key[len(pfx):]
			s.Modified++
		} else {
			s.Skipped++
			if opts.SkipNoMatch {
				out = append(out, e)
				continue
			}
		}
		out = append(out, e)
	}
	return out, s
}
