package compare

import (
	"fmt"
	"io"
	"strings"
)

// DefaultFormatOptions returns sensible defaults for Report.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color:        false,
		ShowEqual:    true,
		MissingMark:  "(missing)",
		ConflictMark: "!!",
	}
}

// FormatOptions controls Report output.
type FormatOptions struct {
	Color        bool
	ShowEqual    bool
	MissingMark  string
	ConflictMark string
}

// Report writes a human-readable matrix table to w.
func Report(w io.Writer, r Result, opts FormatOptions) {
	if len(r.Keys) == 0 {
		fmt.Fprintln(w, "No keys found.")
		return
	}

	// Header
	header := fmt.Sprintf("%-30s", "KEY")
	for _, env := range r.Envs {
		header += fmt.Sprintf("  %-20s", env)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	conflictCount := 0
	for _, key := range r.Keys {
		isConflict := r.Conflicts[key]
		if !opts.ShowEqual && !isConflict {
			continue
		}
		if isConflict {
			conflictCount++
		}
		prefix := "  "
		if isConflict {
			prefix = opts.ConflictMark
			if opts.Color {
				prefix = colorize("\033[33m", prefix)
			}
		}
		line := fmt.Sprintf("%s %-29s", prefix, key)
		for _, env := range r.Envs {
			v, ok := r.Matrix[key][env]
			if !ok {
				v = opts.MissingMark
			}
			if len(v) > 18 {
				v = v[:15] + "..."
			}
			line += fmt.Sprintf("  %-20s", v)
		}
		fmt.Fprintln(w, line)
	}

	fmt.Fprintf(w, "\n%d key(s) total, %d conflict(s)\n", len(r.Keys), conflictCount)
}

func colorize(code, s string) string {
	return code + s + "\033[0m"
}
