// Package prefix provides utilities for adding or removing a common prefix
// from the keys of a set of .env entries.
//
// # Adding a prefix
//
// Use Add to prepend a string to every key:
//
//	out, summary := prefix.Add(entries, "APP_", prefix.DefaultOptions())
//
// # Stripping a prefix
//
// Use Strip to remove a known prefix from every matching key:
//
//	out, summary := prefix.Strip(entries, "APP_", prefix.DefaultOptions())
//
// Keys that do not carry the prefix are skipped when Options.SkipNoMatch is
// true (the default). Set it to false to retain unmatched keys in the output
// unchanged.
//
// # Summary
//
// Both operations return a Summary that reports how many keys were modified
// and how many were skipped, which can be used for reporting or audit
// purposes.
package prefix
