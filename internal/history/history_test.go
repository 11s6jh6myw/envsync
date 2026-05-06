package history_test

import (
	"testing"

	"github.com/user/envsync/internal/history"
	"github.com/user/envsync/internal/parser"
)

func envEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRecord_AppendsEntry(t *testing.T) {
	h := history.New()
	e := h.Record("initial", envEntries("FOO", "bar"))
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
	if e.Label != "initial" {
		t.Errorf("unexpected label: %s", e.Label)
	}
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	h := history.New()
	h.Record("v1", envEntries("A", "1"))
	h.Record("v2", envEntries("A", "2"))
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
}

func TestGet_FindsById(t *testing.T) {
	h := history.New()
	recorded := h.Record("find-me", envEntries("X", "y"))
	found, ok := h.Get(recorded.ID)
	if !ok {
		t.Fatal("expected to find entry by ID")
	}
	if found.Label != "find-me" {
		t.Errorf("unexpected label: %s", found.Label)
	}
}

func TestGet_MissingId(t *testing.T) {
	h := history.New()
	_, ok := h.Get("nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestList_ReturnsCopy(t *testing.T) {
	h := history.New()
	h.Record("a", envEntries("K", "v"))
	list := h.List()
	list[0].Label = "mutated"
	original := h.List()
	if original[0].Label == "mutated" {
		t.Error("List should return a copy, not a reference")
	}
}

func TestRecord_ChecksumDiffersForDifferentValues(t *testing.T) {
	h := history.New()
	e1 := h.Record("v1", envEntries("FOO", "aaa"))
	e2 := h.Record("v2", envEntries("FOO", "bbb"))
	if e1.ID == e2.ID {
		t.Error("expected different IDs for different values")
	}
}
