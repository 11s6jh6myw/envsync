package extract

import (
	"fmt"
	"io"
	"strings"
)

// DefaultFormatOptions returns format options with colour disabled.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color:       false,
		ShowMissing: true,
	}
}

// FormatOptions controls how a Result is rendered.
type FormatOptions struct {
	Color       bool
	ShowMissing bool
}

// Report writes a human-readable summary of an Extract result to w.
func Report(w io.Writer, r Result, opts FormatOptions) {
	if len(r.Entries) == 0 && len(r.Missing) == 0 {
		fmt.Fprintln(w, "No keys extracted.")
		return
	}

	if len(r.Entries) > 0 {
		fmt.Fprintf(w, "Extracted %d key(s):\n", len(r.Entries))
		for _, e := range r.Entries {
			line := fmt.Sprintf("  + %s", e.Key)
			if opts.Color {
				line = colorize("32", line)
			}
			fmt.Fprintln(w, line)
		}
	}

	if opts.ShowMissing && len(r.Missing) > 0 {
		fmt.Fprintf(w, "Missing %d key(s):\n", len(r.Missing))
		for _, k := range r.Missing {
			line := fmt.Sprintf("  ! %s", k)
			if opts.Color {
				line = colorize("33", line)
			}
			fmt.Fprintln(w, line)
		}
	}

	parts := []string{}
	parts = append(parts, fmt.Sprintf("%d extracted", len(r.Entries)))
	if len(r.Missing) > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", len(r.Missing)))
	}
	fmt.Fprintf(w, "Summary: %s\n", strings.Join(parts, ", "))
}

func colorize(code, s string) string {
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}
