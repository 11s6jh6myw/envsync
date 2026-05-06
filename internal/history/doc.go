// Package history provides lightweight change tracking for .env files.
//
// It records snapshots of parsed env entries over time, each identified
// by a short content-based ID and an optional human-readable label.
//
// # Usage
//
//	h := history.New()
//
//	// Record a snapshot after loading
//	entries, _ := env.Load("production.env")
//	h.Record("pre-deploy", entries)
//
//	// ... apply changes ...
//
//	h.Record("post-deploy", updated)
//
//	// List all snapshots
//	for _, e := range h.List() {
//		fmt.Println(e.ID, e.Label, e.Timestamp)
//	}
//
//	// Render a report
//	history.Report(os.Stdout, h, history.DefaultFormatOptions())
package history
