package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// FormatOptions controls how a diff result is rendered.
type FormatOptions struct {
	Color   bool
	Redact  bool
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{Color: true, Redact: false}
}

// Format writes a human-readable diff to w.
func Format(w io.Writer, result *Result, opts FormatOptions) {
	if !result.HasChanges() {
		fmt.Fprintln(w, "No differences found.")
		return
	}
	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			line := fmt.Sprintf("+ %s=%s", c.Key, maybeRedact(c.NewValue, opts.Redact))
			fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
		case Removed:
			line := fmt.Sprintf("- %s=%s", c.Key, maybeRedact(c.OldValue, opts.Redact))
			fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
		case Modified:
			old := maybeRedact(c.OldValue, opts.Redact)
			new := maybeRedact(c.NewValue, opts.Redact)
			line := fmt.Sprintf("~ %s: %s -> %s", c.Key, old, new)
			fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))
		}
	}
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func maybeRedact(val string, redact bool) string {
	if !redact {
		return val
	}
	if len(val) == 0 {
		return ""
	}
	return strings.Repeat("*", len(val))
}
