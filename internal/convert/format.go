package convert

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// DefaultFormatOptions returns sensible defaults for format output.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Color: true,
		ShowCount: true,
	}
}

// FormatOptions controls how conversion reports are rendered.
type FormatOptions struct {
	Color     bool
	ShowCount bool
}

// Report writes a human-readable summary of a conversion operation to w.
// format is the target format (e.g. "json", "export"), count is the number
// of entries converted, and skipped is the number of blank/invalid entries
// that were omitted.
func Report(w io.Writer, format string, count, skipped int, opts FormatOptions) {
	header := fmt.Sprintf("Converted to %s", strings.ToUpper(format))
	if opts.Color {
		header = color.New(color.FgCyan, color.Bold).Sprint(header)
	}
	fmt.Fprintln(w, header)

	if opts.ShowCount {
		countLine := fmt.Sprintf("  entries : %d", count)
		if opts.Color {
			countLine = color.GreenString(countLine)
		}
		fmt.Fprintln(w, countLine)

		if skipped > 0 {
			skipLine := fmt.Sprintf("  skipped : %d (blank keys)", skipped)
			if opts.Color {
				skipLine = color.YellowString(skipLine)
			}
			fmt.Fprintln(w, skipLine)
		}
	}
}
