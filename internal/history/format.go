package history

import (
	"fmt"
	"io"
	"strings"
)

// FormatOptions controls how history is rendered.
type FormatOptions struct {
	Color      bool
	ShowKeys   bool
	MaxEntries int
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color:      false,
		ShowKeys:   true,
		MaxEntries: 20,
	}
}

// Report writes a human-readable history report to w.
func Report(w io.Writer, h *History, opts FormatOptions) {
	records := h.List()
	if len(records) == 0 {
		fmt.Fprintln(w, "No history recorded.")
		return
	}

	limit := len(records)
	if opts.MaxEntries > 0 && opts.MaxEntries < limit {
		limit = opts.MaxEntries
	}

	header := fmt.Sprintf("History (%d entries):", len(records))
	if opts.Color {
		header = "\033[1m" + header + "\033[0m"
	}
	fmt.Fprintln(w, header)

	for i := len(records) - 1; i >= len(records)-limit; i-- {
		r := records[i]
		ts := r.Timestamp.Format("2006-01-02 15:04:05")
		line := fmt.Sprintf("  [%s] %s  %s  (%d keys)", r.ID, ts, r.Label, len(r.Entries))
		if opts.Color {
			line = "\033[36m" + line + "\033[0m"
		}
		fmt.Fprintln(w, line)
		if opts.ShowKeys && len(r.Entries) > 0 {
			var keys []string
			for _, e := range r.Entries {
				if e.Key != "" {
					keys = append(keys, e.Key)
				}
			}
			if len(keys) > 0 {
				fmt.Fprintf(w, "      keys: %s\n", strings.Join(keys, ", "))
			}
		}
	}
}
