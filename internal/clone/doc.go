// Package clone provides the Clone function for duplicating a set of parsed
// .env entries into a new slice.
//
// # Overview
//
// Clone is useful when you need to produce a sanitised or template copy of an
// existing environment file — for example, to commit a redacted version to
// source control or to hand off a starter file to a new team member.
//
// # Redaction
//
// When Options.Redact is true, any entry whose key matches a sensitive pattern
// (e.g. "secret", "password", "token") has its value replaced by the
// configured Placeholder string (default: "REDACTED").  You can supply your
// own SensitiveKeys list to override the built-in defaults.
//
// # Comment Stripping
//
// Setting Options.StripComments to true removes blank lines and comment-only
// entries from the output, producing a compact file that contains only
// key=value pairs.
//
// # Usage
//
//	entries, _ := parser.Parse(src)
//	opts := clone.DefaultOptions()
//	opts.Redact = true
//	out, summary, err := clone.Clone(entries, opts)
package clone
