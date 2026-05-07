package rotate_test

import (
	"testing"

	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/rotate"
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

func TestRotate_ExplicitReplacement(t *testing.T) {
	opts := rotate.DefaultOptions()
	opts.Replacements["API_KEY"] = "new-secret"

	res, err := rotate.Rotate(entries("API_KEY", "old-secret", "APP_NAME", "myapp"), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findKey(t, res.Entries, "API_KEY").Value != "new-secret" {
		t.Errorf("expected API_KEY to be rotated")
	}
	if findKey(t, res.Entries, "APP_NAME").Value != "myapp" {
		t.Errorf("expected APP_NAME to remain unchanged")
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "API_KEY" {
		t.Errorf("expected 1 rotated key, got %v", res.Rotated)
	}
}

func TestRotate_SkipsWhenNoReplacement(t *testing.T) {
	opts := rotate.DefaultOptions()
	opts.OnlySensitive = true

	res, err := rotate.Rotate(entries("SECRET_TOKEN", "abc123"), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 0 {
		t.Errorf("expected 0 rotated, got %v", res.Rotated)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", res.Skipped)
	}
}

func TestRotate_PreservesCommentEntries(t *testing.T) {
	input := []parser.Entry{
		{Raw: "# header comment"},
		{Key: "DB_PASSWORD", Value: "old"},
	}
	opts := rotate.DefaultOptions()
	opts.Replacements["DB_PASSWORD"] = "new"

	res, err := rotate.Rotate(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Raw != "# header comment" {
		t.Errorf("comment entry should be preserved")
	}
	if res.Entries[1].Value != "new" {
		t.Errorf("DB_PASSWORD should be rotated")
	}
}

func TestSummary_Format(t *testing.T) {
	r := rotate.Result{
		Rotated: []string{"A", "B"},
		Skipped: []string{"C"},
	}
	s := rotate.Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
