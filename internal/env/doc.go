// Package env handles loading and resolving .env files from the filesystem.
//
// It wraps the parser package to provide convenient file-based entry points
// used by higher-level commands such as diff and sync.
//
// Typical usage:
//
//	// Load a single file
//	res, err := env.Load(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Entries)
//
//	// Load a source/target pair for diffing
//	src, tgt, err := env.LoadPair(".env", ".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = src
//	_ = tgt
package env
