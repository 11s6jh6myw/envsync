package sync

import (
	"testing"

	"github.com/yourusername/envsync/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func findKey(entries []parser.Entry, key string) (string, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

func TestApply_Fill_AddsMissingKeys(t *testing.T) {
	src := entries("A", "1", "B", "2", "C", "3")
	tgt := entries("A", "old")

	result, err := Apply(src, tgt, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// A should keep its original value (fill strategy)
	if v, ok := findKey(result, "A"); !ok || v != "old" {
		t.Errorf("expected A=old, got %q", v)
	}
	// B and C should be added with placeholder
	if _, ok := findKey(result, "B"); !ok {
		t.Error("expected B to be added")
	}
	if _, ok := findKey(result, "C"); !ok {
		t.Error("expected C to be added")
	}
}

func TestApply_Overwrite_UpdatesExisting(t *testing.T) {
	src := entries("A", "new", "B", "2")
	tgt := entries("A", "old", "EXTRA", "keep")

	opts := Options{Strategy: StrategyOverwrite, Placeholder: ""}
	result, err := Apply(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v, ok := findKey(result, "A"); !ok || v != "new" {
		t.Errorf("expected A=new, got %q", v)
	}
	if _, ok := findKey(result, "EXTRA"); !ok {
		t.Error("expected EXTRA to be preserved")
	}
	if _, ok := findKey(result, "B"); !ok {
		t.Error("expected B to be added")
	}
}

func TestApply_Exact_RemovesExtraKeys(t *testing.T) {
	src := entries("A", "1", "B", "2")
	tgt := entries("A", "old", "EXTRA", "gone")

	opts := Options{Strategy: StrategyExact, Placeholder: ""}
	result, err := Apply(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := findKey(result, "EXTRA"); ok {
		t.Error("expected EXTRA to be removed")
	}
	if v, ok := findKey(result, "A"); !ok || v != "1" {
		t.Errorf("expected A=1, got %q", v)
	}
	if _, ok := findKey(result, "B"); !ok {
		t.Error("expected B to be added")
	}
}

func TestApply_Placeholder(t *testing.T) {
	src := entries("NEW", "secret")
	tgt := entries("OTHER", "val")

	opts := Options{Strategy: StrategyFill, Placeholder: "CHANGE_ME"}
	result, err := Apply(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v, ok := findKey(result, "NEW"); !ok || v != "CHANGE_ME" {
		t.Errorf("expected NEW=CHANGE_ME, got %q", v)
	}
}
