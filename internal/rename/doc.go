// Package rename provides key-rename operations for .env entry slices.
//
// It allows callers to rename one or more keys in a parsed .env file while
// preserving the original order, blank lines, inline comments, and values.
//
// Basic usage:
//
//	result, err := rename.Rename(entries, map[string]string{
//		"DB_HOST": "DATABASE_HOST",
//		"DB_PORT": "DATABASE_PORT",
//	}, rename.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("renamed %d key(s)\n", len(result.Renamed))
//
// By default, keys that appear in the mapping but are absent from the entries
// are silently collected in Result.Skipped. Set Options.FailOnMissing = true
// to receive a *MissingKeyError instead.
package rename
