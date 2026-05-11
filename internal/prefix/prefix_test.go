package prefix

import (
	"testing"

	"github.com/jasonuc/envsync/internal/parser"
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

func TestAdd_BasicPrefix(t *testing.T) {
	in := entries("HOST", "localhost", "PORT", "5432")
	out, s := Add(in, "APP_", DefaultOptions())
	if s.Modified != 2 {
		t.Fatalf("expected 2 modified, got %d", s.Modified)
	}
	if _, ok := findKey(out, "APP_HOST"); !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := findKey(out, "APP_PORT"); !ok {
		t.Error("expected APP_PORT")
	}
}

func TestAdd_EmptyPrefix_NoChange(t *testing.T) {
	in := entries("KEY", "val")
	out, s := Add(in, "", DefaultOptions())
	if s.Modified != 0 {
		t.Fatalf("expected 0 modified, got %d", s.Modified)
	}
	if out[0].Key != "KEY" {
		t.Errorf("key should be unchanged, got %s", out[0].Key)
	}
}

func TestAdd_SkipsBlankKeys(t *testing.T) {
	in := []parser.Entry{{Key: "", Value: "", Comment: "# comment"}}
	out, s := Add(in, "X_", DefaultOptions())
	if s.Modified != 0 {
		t.Fatalf("expected 0 modified, got %d", s.Modified)
	}
	if out[0].Comment != "# comment" {
		t.Error("comment entry should be preserved")
	}
}

func TestStrip_BasicStrip(t *testing.T) {
	in := entries("APP_HOST", "localhost", "APP_PORT", "5432")
	out, s := Strip(in, "APP_", DefaultOptions())
	if s.Modified != 2 {
		t.Fatalf("expected 2 modified, got %d", s.Modified)
	}
	if _, ok := findKey(out, "HOST"); !ok {
		t.Error("expected HOST")
	}
}

func TestStrip_SkipsNonMatching(t *testing.T) {
	in := entries("APP_HOST", "localhost", "PORT", "5432")
	out, s := Strip(in, "APP_", DefaultOptions())
	if s.Modified != 1 {
		t.Fatalf("expected 1 modified, got %d", s.Modified)
	}
	if s.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", s.Skipped)
	}
	// PORT should still be present (SkipNoMatch=true keeps it)
	if _, ok := findKey(out, "PORT"); !ok {
		t.Error("expected PORT to be retained")
	}
}

func TestStrip_EmptyPrefix_NoChange(t *testing.T) {
	in := entries("KEY", "val")
	out, s := Strip(in, "", DefaultOptions())
	if s.Modified != 0 {
		t.Error("expected 0 modified")
	}
	if out[0].Key != "KEY" {
		t.Errorf("key should be unchanged, got %s", out[0].Key)
	}
}
