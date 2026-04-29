package parser

import (
	"fmt"
	"io"
	"strings"
)

// WriteOptions controls serialization behaviour.
type WriteOptions struct {
	QuoteValues bool // wrap all values in double quotes
	KeepComments bool
}

// DefaultWriteOptions returns sensible defaults.
func DefaultWriteOptions() WriteOptions {
	return WriteOptions{
		QuoteValues:  false,
		KeepComments: true,
	}
}

// Write serializes an EnvFile to the given writer.
func Write(w io.Writer, ef *EnvFile, opts WriteOptions) error {
	for i, entry := range ef.Entries {
		if opts.KeepComments && entry.Comment != "" {
			if i > 0 {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w, entry.Comment); err != nil {
				return err
			}
		}

		value := entry.Value
		if opts.QuoteValues || needsQuoting(value) {
			value = fmt.Sprintf("%q", value)
		}

		if _, err := fmt.Fprintf(w, "%s=%s\n", entry.Key, value); err != nil {
			return err
		}
	}
	return nil
}

// needsQuoting returns true when a value contains spaces or special characters.
func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t#")
}
