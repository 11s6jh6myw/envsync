package lint

import (
	"fmt"
	"io"
	"strings"
)

// FormatOptions controls how lint results are rendered.
type FormatOptions struct {
	Color   bool
	Verbose bool // include a summary line even when no issues found
}

// DefaultFormatOptions returns colour-enabled format options.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{Color: true, Verbose: false}
}

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorBold   = "\033[1m"
)

// Report writes a human-readable lint report to w.
func Report(w io.Writer, issues []Issue, opts FormatOptions) {
	if len(issues) == 0 {
		if opts.Verbose {
			msg := "✔ no lint issues found"
			if opts.Color {
				msg = colorGreen + msg + colorReset
			}
			fmt.Fprintln(w, msg)
		}
		return
	}

	warnings, errors := 0, 0
	for _, iss := range issues {
		prefix := severityLabel(iss.Severity, opts.Color)
		line := fmt.Sprintf("%s line %d [%s]: %s", prefix, iss.Line, iss.Key, iss.Message)
		fmt.Fprintln(w, line)
		if iss.Severity == SeverityError {
			errors++
		} else {
			warnings++
		}
	}

	summary := buildSummary(errors, warnings)
	if opts.Color {
		if errors > 0 {
			summary = colorBold + colorRed + summary + colorReset
		} else {
			summary = colorBold + colorYellow + summary + colorReset
		}
	}
	fmt.Fprintln(w, summary)
}

func severityLabel(s Severity, color bool) string {
	switch s {
	case SeverityError:
		if color {
			return colorRed + "[error]" + colorReset
		}
		return "[error]"
	default:
		if color {
			return colorYellow + "[warn]" + colorReset
		}
		return "[warn]"
	}
}

func buildSummary(errors, warnings int) string {
	parts := []string{}
	if errors > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", errors))
	}
	if warnings > 0 {
		parts = append(parts, fmt.Sprintf("%d warning(s)", warnings))
	}
	return "lint: " + strings.Join(parts, ", ")
}
