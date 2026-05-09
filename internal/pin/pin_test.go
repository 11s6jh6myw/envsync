package pin_test

import (
	"testing"

	"github.com/yourusername/envsync/internal/parser"
	"github.com/yourusername/envsync/internal/pin"
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
	t.Fatalf("key %q not found", key)
	return parser.Entry{}
}

func TestPin_MarksComment(t *testing.T) {
	in := entries("DB_HOST", "localhost", "API_KEY", "secret")
	res, err := pin.Pin(in, []string{"DB_HOST"}, pin.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := findKey(t, res.Entries, "DB_HOST")
	if e.Comment != "pinned" {
		t.Errorf("expected comment 'pinned', got %q", e.Comment)
	}
}

func TestPin_PreservesOtherEntries(t *testing.T) {
	in := entries("DB_HOST", "localhost", "API_KEY", "secret")
	res, err := pin.Pin(in, []string{"DB_HOST"}, pin.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	e := findKey(t, res.Entries, "API_KEY")
	if e.Comment == "pinned" {
		t.Error("API_KEY should not be pinned")
	}
}

func TestPin_SkipsMissingKey(t *testing.T) {
	in := entries("DB_HOST", "localhost")
	res, err := pin.Pin(in, []string{"MISSING"}, pin.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestPin_FailOnMissing(t *testing.T) {
	in := entries("DB_HOST", "localhost")
	opts := pin.DefaultOptions()
	opts.FailOnMissing = true
	_, err := pin.Pin(in, []string{"MISSING"}, opts)
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestPin_RecordsPinnedKeys(t *testing.T) {
	in := entries("A", "1", "B", "2", "C", "3")
	res, err := pin.Pin(in, []string{"A", "C"}, pin.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d", len(res.Pinned))
	}
}

func TestIsPinned(t *testing.T) {
	e := parser.Entry{Key: "X", Value: "y", Comment: "pinned"}
	if !pin.IsPinned(e) {
		t.Error("expected IsPinned to return true")
	}
	e.Comment = ""
	if pin.IsPinned(e) {
		t.Error("expected IsPinned to return false")
	}
}
