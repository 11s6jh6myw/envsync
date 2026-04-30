package validate

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// FormatOptions controls how validation results are rendered.
type FormatOptions struct {
	Color bool
	File  string // optional filename shown in header
}

// DefaultFormatOptions returns format options with color enabled.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{Color: true}
}

// Report writes a human-readable validation report to w.
// Returns the number of errors found.
func Report(w io.Writer, issues []Issue, opts FormatOptions) int {
	if len(issues) == 0 {
		label := "validation"
		if opts.File != "" {
			label = opts.File
		}
		fmt.Fprintf(w, "%s: no issues found\n", label)
		return 0
	}

	if opts.File != "" {
		fmt.Fprintf(w, "Validation report for %s:\n", opts.File)
	}

	errCount := 0
	for _, issue := range issues {
		prefix := severityPrefix(issue.Severity, opts.Color)
		fmt.Fprintf(w, "  %s %s\n", prefix, issue.String())
		if issue.Severity == "error" {
			errCount++
		}
	}

	summary := buildSummary(issues)
	fmt.Fprintf(w, "\n%s\n", summary)
	return errCount
}

func severityPrefix(severity string, color bool) string {
	switch strings.ToLower(severity) {
	case "error":
		if color {
			return colorRed + "[error]" + colorReset
		}
		return "[error]"
	case "warning":
		if color {
			return colorYellow + "[warn] " + colorReset
		}
		return "[warn] "
	default:
		return "[info] "
	}
}

func buildSummary(issues []Issue) string {
	errs, warns := 0, 0
	for _, i := range issues {
		if i.Severity == "error" {
			errs++
		} else {
			warns++
		}
	}
	return fmt.Sprintf("%d error(s), %d warning(s)", errs, warns)
}
