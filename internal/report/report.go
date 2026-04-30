// Package report provides formatted summary reporting for env sync operations.
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envsync/internal/sync"
)

// Options controls report output behaviour.
type Options struct {
	// ShowUnchanged includes keys that were not modified in the output.
	ShowUnchanged bool
	// UseColor enables ANSI colour codes in the output.
	UseColor bool
}

// DefaultOptions returns sensible defaults for report generation.
func DefaultOptions() Options {
	return Options{
		ShowUnchanged: false,
		UseColor:      true,
	}
}

// Write renders a human-readable sync summary to w.
func Write(w io.Writer, s sync.Summary, opts Options) error {
	lines := []string{}

	for _, k := range s.Added {
		line := fmt.Sprintf("  + %s", k)
		if opts.UseColor {
			line = colorize("\033[32m", line)
		}
		lines = append(lines, line)
	}

	for _, k := range s.Removed {
		line := fmt.Sprintf("  - %s", k)
		if opts.UseColor {
			line = colorize("\033[31m", line)
		}
		lines = append(lines, line)
	}

	for _, k := range s.Updated {
		line := fmt.Sprintf("  ~ %s", k)
		if opts.UseColor {
			line = colorize("\033[33m", line)
		}
		lines = append(lines, line)
	}

	if opts.ShowUnchanged {
		for _, k := range s.Unchanged {
			lines = append(lines, fmt.Sprintf("    %s", k))
		}
	}

	if len(lines) == 0 {
		_, err := fmt.Fprintln(w, "No changes applied.")
		return err
	}

	header := fmt.Sprintf(
		"Sync summary: +%d added  -%d removed  ~%d updated  %d unchanged",
		len(s.Added), len(s.Removed), len(s.Updated), len(s.Unchanged),
	)
	_, err := fmt.Fprintf(w, "%s\n%s\n", header, strings.Join(lines, "\n"))
	return err
}

func colorize(code, text string) string {
	return fmt.Sprintf("%s%s\033[0m", code, text)
}
