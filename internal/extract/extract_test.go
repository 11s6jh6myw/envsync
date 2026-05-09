package extract_test

import (
	"testing"

	"github.com/user/envsync/internal/extract"
	"github.com/user/envsync/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func findKey(t *testing.T, entries []parser.Entry, key string) parser.Entry {
	t.Helper()
	for _, e := range entries {
		if e.Key == key {
			return e
		}
	}
	t.Fatalf("key %q not found in entries", key)
	return parser.Entry{}
}

func TestExtract_ReturnsRequestedKeys(t *testing.T) {
	src := entries("A", "1", "B", "2", "C", "3")
	r, err := extract.Extract(src, []string{"A", "C"}, extract.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	findKey(t, r.Entries, "A")
	findKey(t, r.Entries, "C")
}

func TestExtract_MissingKeyReturnsError(t *testing.T) {
	src := entries("A", "1")
	_, err := extract.Extract(src, []string{"A", "MISSING"}, extract.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestExtract_IgnoreMissing_NoError(t *testing.T) {
	src := entries("A", "1")
	opts := extract.DefaultOptions()
	opts.IgnoreMissing = true
	r, err := extract.Extract(src, []string{"A", "GHOST"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "GHOST" {
		t.Errorf("expected GHOST in Missing, got %v", r.Missing)
	}
}

func TestExtract_PreserveOrder(t *testing.T) {
	src := entries("C", "3", "A", "1", "B", "2")
	opts := extract.DefaultOptions()
	opts.PreserveOrder = true
	r, err := extract.Extract(src, []string{"A", "B", "C"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Entries[0].Key != "C" || r.Entries[1].Key != "A" || r.Entries[2].Key != "B" {
		t.Errorf("order not preserved: %v", r.Entries)
	}
}

func TestExtract_KeyOrder_FollowsKeysSlice(t *testing.T) {
	src := entries("C", "3", "A", "1", "B", "2")
	opts := extract.DefaultOptions()
	opts.PreserveOrder = false
	r, err := extract.Extract(src, []string{"A", "B", "C"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Entries[0].Key != "A" || r.Entries[1].Key != "B" || r.Entries[2].Key != "C" {
		t.Errorf("expected keys-slice order, got %v", r.Entries)
	}
}

func TestExtract_SkipsCommentEntries(t *testing.T) {
	src := []parser.Entry{
		{Comment: "# header"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	r, err := extract.Extract(src, []string{"A"}, extract.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range r.Entries {
		if e.Key == "" {
			t.Error("comment/blank entry leaked into result")
		}
	}
}
