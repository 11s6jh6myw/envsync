// Package pin provides key-pinning support for envsync.
//
// Pinning a key marks it with an inline "pinned" comment so that downstream
// operations (sync, merge, rotate) can detect and skip it, preserving the
// current value across environment boundaries.
//
// Basic usage:
//
//	res, err := pin.Pin(entries, []string{"DB_HOST", "SECRET_KEY"}, pin.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// res.Entries now has DB_HOST and SECRET_KEY marked as pinned.
//	// res.Pinned lists each successfully pinned key/value pair.
//	// res.Skipped lists keys that were not found (when FailOnMissing is false).
//
// Detecting pins in other packages:
//
//	for _, e := range entries {
//	    if pin.IsPinned(e) {
//	        // skip this entry during sync
//	    }
//	}
package pin
