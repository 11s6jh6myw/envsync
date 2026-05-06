// Package history tracks changes to .env files over time,
// allowing users to review past diffs and restore previous states.
package history

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/user/envsync/internal/parser"
)

// Entry represents a single point-in-time record of an env file.
type Entry struct {
	ID        string
	Timestamp time.Time
	Label     string
	Entries   []parser.Entry
	checksum  string
}

// History holds an ordered list of recorded env snapshots.
type History struct {
	records []Entry
}

// New returns an empty History.
func New() *History {
	return &History{}
}

// Record appends a new entry to the history.
func (h *History) Record(label string, entries []parser.Entry) Entry {
	cs := checksum(entries)
	e := Entry{
		ID:        cs[:8],
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   copyEntries(entries),
		checksum:  cs,
	}
	h.records = append(h.records, e)
	return e
}

// List returns all recorded entries in chronological order.
func (h *History) List() []Entry {
	out := make([]Entry, len(h.records))
	copy(out, h.records)
	return out
}

// Get returns the entry with the given ID, or false if not found.
func (h *History) Get(id string) (Entry, bool) {
	for _, r := range h.records {
		if r.ID == id {
			return r, true
		}
	}
	return Entry{}, false
}

// Len returns the number of recorded entries.
func (h *History) Len() int { return len(h.records) }

func checksum(entries []parser.Entry) string {
	h := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(h, "%s=%s\n", e.Key, e.Value)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func copyEntries(src []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(src))
	copy(out, src)
	return out
}
