package rename_test

import (
	"testing"

	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/rename"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func findKey(t *testing.T, es []parser.Entry, key string) parser.Entry {
	t.Helper()
	for _, e := range es {
		if e.Key == key {
			return e
		}
	}
	t.Fatalf("key %q not found in entries", key)
	return parser.Entry{}
}

func TestRename_BasicRename(t *testing.T) {
	es := entries("DB_HOST", "localhost", "DB_PORT", "5432")
	res, err := rename.Rename(es, map[string]string{"DB_HOST": "DATABASE_HOST"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 1 {
		t.Fatalf("expected 1 renamed, got %d", len(res.Renamed))
	}
	e := findKey(t, res.Entries, "DATABASE_HOST")
	if e.Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", e.Value)
	}
}

func TestRename_PreservesOtherKeys(t *testing.T) {
	es := entries("APP_NAME", "myapp", "APP_ENV", "prod")
	res, err := rename.Rename(es, map[string]string{"APP_NAME": "APPLICATION_NAME"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = findKey(t, res.Entries, "APP_ENV")
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestRename_SkipsMissingKeyByDefault(t *testing.T) {
	es := entries("FOO", "bar")
	res, err := rename.Rename(es, map[string]string{"MISSING": "NEW_KEY"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestRename_FailOnMissingKey(t *testing.T) {
	es := entries("FOO", "bar")
	opts := rename.DefaultOptions()
	opts.FailOnMissing = true
	_, err := rename.Rename(es, map[string]string{"MISSING": "NEW_KEY"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRename_PreservesCommentEntries(t *testing.T) {
	es := []parser.Entry{
		{IsComment: true, Raw: "# database config"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
	res, err := rename.Rename(es, map[string]string{"DB_URL": "DATABASE_URL"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Entries[0].IsComment {
		t.Error("expected first entry to remain a comment")
	}
	_ = findKey(t, res.Entries, "DATABASE_URL")
}

func TestRename_MultipleKeys(t *testing.T) {
	es := entries("OLD_A", "1", "OLD_B", "2", "OLD_C", "3")
	mapping := map[string]string{"OLD_A": "NEW_A", "OLD_B": "NEW_B"}
	res, err := rename.Rename(es, mapping, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renamed entries, got %d", len(res.Renamed))
	}
	_ = findKey(t, res.Entries, "NEW_A")
	_ = findKey(t, res.Entries, "NEW_B")
	_ = findKey(t, res.Entries, "OLD_C")
}
