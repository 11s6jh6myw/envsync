package prune

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

func findKey(es []parser.Entry, key string) (parser.Entry, bool) {
	for _, e := range es {
		if e.Key == key {
			return e, true
		}
	}
	return parser.Entry{}, false
}

func TestPrune_RemovesDisallowedKeys(t *testing.T) {
	in := entries("APP_ENV", "prod", "SECRET_KEY", "abc", "PORT", "8080")
	out, sum := Prune(in, []string{"APP_ENV", "PORT"}, DefaultOptions())
	if _, ok := findKey(out, "SECRET_KEY"); ok {
		t.Error("expected SECRET_KEY to be pruned")
	}
	if len(sum.Removed) != 1 || sum.Removed[0] != "SECRET_KEY" {
		t.Errorf("unexpected removed list: %v", sum.Removed)
	}
}

func TestPrune_KeepsAllowedKeys(t *testing.T) {
	in := entries("APP_ENV", "prod", "PORT", "8080")
	out, sum := Prune(in, []string{"APP_ENV", "PORT"}, DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if len(sum.Kept) != 2 {
		t.Errorf("expected 2 kept, got %v", sum.Kept)
	}
}

func TestPrune_DryRun_DoesNotModify(t *testing.T) {
	in := entries("APP_ENV", "prod", "STALE", "old")
	opts := DefaultOptions()
	opts.DryRun = true
	out, sum := Prune(in, []string{"APP_ENV"}, opts)
	if len(out) != len(in) {
		t.Errorf("dry-run should return original slice, got %d entries", len(out))
	}
	if len(sum.Removed) != 1 {
		t.Errorf("summary should still report removed keys")
	}
}

func TestPrune_KeepsComments(t *testing.T) {
	in := []parser.Entry{
		{Key: "", Raw: "# section header"},
		{Key: "APP_ENV", Value: "prod"},
		{Key: "STALE", Value: "old"},
	}
	out, _ := Prune(in, []string{"APP_ENV"}, DefaultOptions())
	if len(out) != 2 {
		t.Errorf("expected comment + APP_ENV, got %d entries", len(out))
	}
}

func TestPrune_DropsCommentsWhenDisabled(t *testing.T) {
	in := []parser.Entry{
		{Key: "", Raw: "# section header"},
		{Key: "APP_ENV", Value: "prod"},
	}
	opts := DefaultOptions()
	opts.KeepComments = false
	out, _ := Prune(in, []string{"APP_ENV"}, opts)
	if len(out) != 1 {
		t.Errorf("expected only APP_ENV, got %d entries", len(out))
	}
}
