package rotate

import (
	"fmt"
	"io"
	"strings"
)

// FormatOptions controls how a rotation report is rendered.
type FormatOptions struct {
	Color   bool
	Verbose bool // show skipped keys
}

// DefaultFormatOptions returns options suitable for terminal output.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{Color: true, Verbose: false}
}

// Report writes a human-readable rotation report to w.
func Report(w io.Writer, r Result, opts FormatOptions) {
	if len(r.Rotated) == 0 && len(r.Skipped) == 0 {
		fmt.Fprintln(w, "No keys rotated.")
		return
	}

	if len(r.Rotated) > 0 {
		header := fmt.Sprintf("Rotated (%d):", len(r.Rotated))
		if opts.Color {
			header = colorize("\033[32m", header)
		}
		fmt.Fprintln(w, header)
		for _, k := range r.Rotated {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	if opts.Verbose && len(r.Skipped) > 0 {
		header := fmt.Sprintf("Skipped (%d):", len(r.Skipped))
		if opts.Color {
			header = colorize("\033[33m", header)
		}
		fmt.Fprintln(w, header)
		for _, k := range r.Skipped {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 30))
	fmt.Fprintln(w, Summary(r))
}

func colorize(code, s string) string {
	return code + s + "\033[0m"
}
