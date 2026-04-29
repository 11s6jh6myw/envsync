// Package diff provides functionality for comparing two sets of environment
// file entries and producing a structured list of changes.
//
// # Overview
//
// The core Diff function accepts two slices of parser.Entry values (source and
// destination) and returns a slice of Change records, each tagged with one of
// four kinds:
//
//   - KindUnchanged – key exists in both files with the same value
//   - KindAdded     – key is present only in the destination
//   - KindRemoved   – key is present only in the source
//   - KindModified  – key exists in both files but with different values
//
// # Formatting
//
// The Format function renders a slice of Change records as a human-readable
// diff string. FormatOptions controls:
//
//   - Color        – enable ANSI colour codes (green/red/yellow)
//   - Redact       – mask sensitive values using the redact package
//   - ShowUnchanged – include unchanged keys in the output
//
// # Usage
//
//	src, _ := parser.Parse(srcBytes)
//	dst, _ := parser.Parse(dstBytes)
//	changes := diff.Diff(src, dst)
//	fmt.Print(diff.Format(changes, diff.DefaultFormatOptions()))
package diff
