// Package snapshot provides functionality to capture and compare
// point-in-time snapshots of .env file state for change tracking.
package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/user/envsync/internal/parser"
)

// Snapshot represents a point-in-time capture of an env file's state.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Checksum  string            `json:"checksum"`
	Entries   []parser.Entry    `json:"entries"`
	KeyCount  int               `json:"key_count"`
}

// Take creates a new Snapshot from a slice of parsed entries.
func Take(source string, entries []parser.Entry) Snapshot {
	keys := keyEntries(entries)
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Checksum:  checksum(keys),
		Entries:   entries,
		KeyCount:  len(keys),
	}
}

// Equal reports whether two snapshots have identical key-value content.
func Equal(a, b Snapshot) bool {
	return a.Checksum == b.Checksum
}

// Serialize encodes a Snapshot to JSON bytes.
func Serialize(s Snapshot) ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Deserialize decodes a Snapshot from JSON bytes.
func Deserialize(data []byte) (Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode failed: %w", err)
	}
	return s, nil
}

// keyEntries filters entries to only those with non-empty keys.
func keyEntries(entries []parser.Entry) []parser.Entry {
	var out []parser.Entry
	for _, e := range entries {
		if e.Key != "" {
			out = append(out, e)
		}
	}
	return out
}

// checksum produces a deterministic SHA-256 digest of sorted key=value pairs.
func checksum(entries []parser.Entry) string {
	sorted := make([]parser.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	h := sha256.New()
	for _, e := range sorted {
		fmt.Fprintf(h, "%s=%s\n", e.Key, e.Value)
	}
	return hex.EncodeToString(h.Sum(nil))
}
