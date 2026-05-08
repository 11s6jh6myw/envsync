// Package prune provides utilities for removing stale or disallowed keys from
// a parsed .env file.
//
// Given a list of allowed keys, Prune returns only the entries whose keys
// appear in that list. Keys not present in the allowed set are considered
// stale and are dropped from the output.
//
// # Dry-run mode
//
// When Options.DryRun is true the original entry slice is returned unchanged
// so callers can inspect the Summary before committing any writes.
//
// # Comments and blank lines
//
// By default comment and blank entries are preserved in the output regardless
// of which keys are pruned. Set Options.KeepComments to false to drop them.
//
// # Usage
//
//	allowed := []string{"APP_ENV", "PORT", "DATABASE_URL"}
//	cleaned, summary := prune.Prune(entries, allowed, prune.DefaultOptions())
//	fmt.Printf("removed %d stale keys\n", len(summary.Removed))
package prune
