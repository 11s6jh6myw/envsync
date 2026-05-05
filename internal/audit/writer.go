package audit

import (
	"fmt"
	"io"
)

// WriteOptions controls how the audit log is rendered.
type WriteOptions struct {
	// ShowAll includes events of every type; when false only errors are shown.
	ShowAll bool
	// Prefix is prepended to every output line.
	Prefix string
}

// DefaultWriteOptions returns sensible defaults.
func DefaultWriteOptions() WriteOptions {
	return WriteOptions{ShowAll: true}
}

// Write renders all events in the log to w.
func Write(w io.Writer, l *Log, opts WriteOptions) error {
	events := l.Events()
	if len(events) == 0 {
		_, err := fmt.Fprintln(w, opts.Prefix+"audit: no events recorded")
		return err
	}
	for _, e := range events {
		line := FormatEvent(e)
		if _, err := fmt.Fprintln(w, opts.Prefix+line); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "%saudit: %d event(s) total\n", opts.Prefix, len(events))
	return err
}
