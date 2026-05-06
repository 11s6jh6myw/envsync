package snapshot

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// DefaultFormatOptions returns sensible defaults for snapshot report formatting.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color:       false,
		ShowKeys:    true,
		TimeFormat:  time.RFC3339,
	}
}

// FormatOptions controls how a snapshot report is rendered.
type FormatOptions struct {
	Color      bool
	ShowKeys   bool
	TimeFormat string
}

// Report writes a human-readable summary of a Snapshot to w.
func Report(w io.Writer, s Snapshot, opts FormatOptions) {
	label := func(c, text string) string {
		if opts.Color {
			return c + text + "\033[0m"
		}
		return text
	}

	header := fmt.Sprintf("Snapshot: %s", s.Label)
	if s.Label == "" {
		header = "Snapshot"
	}
	fmt.Fprintln(w, label("\033[1m", header))

	taken := s.TakenAt.Format(opts.TimeFormat)
	fmt.Fprintf(w, "  Taken : %s\n", taken)
	fmt.Fprintf(w, "  Keys  : %d\n", len(s.Entries))

	if opts.ShowKeys && len(s.Entries) > 0 {
		fmt.Fprintln(w, "  Keys:")
		keys := keyEntries(s.Entries)
		for _, k := range keys {
			line := fmt.Sprintf("    - %s", k)
			fmt.Fprintln(w, strings.TrimRight(line, " "))
		}
	}
}
