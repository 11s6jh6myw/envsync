package merge

import (
	"fmt"
	"io"
	"strings"
)

// FormatOptions controls how merge results are rendered.
type FormatOptions struct {
	// UseColor enables ANSI color codes in output.
	UseColor bool
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{UseColor: true}
}

// Report writes a human-readable summary of the merge result to w.
func Report(w io.Writer, result Result, opts FormatOptions) {
	if len(result.Conflicts) == 0 {
		fmt.Fprintln(w, "Merge completed with no conflicts.")
		return
	}

	header := fmt.Sprintf("Merge completed with %d conflict(s):", len(result.Conflicts))
	fmt.Fprintln(w, header)

	for _, c := range result.Conflicts {
		srcLine := fmt.Sprintf("  - source: %s=%s", c.Key, c.SourceValue)
		tgtLine := fmt.Sprintf("  - target: %s=%s", c.Key, c.TargetValue)
		resLine := fmt.Sprintf("  ✓ resolved: %s=%s", c.Key, c.Resolved)

		if opts.UseColor {
			srcLine = colorize(srcLine, "\033[33m") // yellow
			tgtLine = colorize(tgtLine, "\033[33m")
			resLine = colorize(resLine, "\033[32m") // green
		}

		fmt.Fprintln(w, strings.Repeat("-", 40))
		fmt.Fprintln(w, srcLine)
		fmt.Fprintln(w, tgtLine)
		fmt.Fprintln(w, resLine)
	}
}

func colorize(s, code string) string {
	return code + s + "\033[0m"
}
