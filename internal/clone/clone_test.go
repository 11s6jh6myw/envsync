package clone

import (
	"testing"

	"github.com/your-org/envsync/internal/parser"
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

func TestClone_CopiesAllEntries(t *testing.T) {
	src := entries("APP_NAME", "myapp", "PORT", "8080")
	out, sum, err := Clone(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if sum.Total != 2 {
		t.Errorf("expected Total=2, got %d", sum.Total)
	}
	if sum.Redacted != 0 {
		t.Errorf("expected Redacted=0, got %d", sum.Redacted)
	}
}

func TestClone_RedactsSensitiveKeys(t *testing.T) {
	src := entries("APP_NAME", "myapp", "API_KEY", "super-secret", "DB_PASSWORD", "hunter2")
	opts := DefaultOptions()
	opts.Redact = true

	out, sum, err := Clone(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	appEntry := findKey(t, out, "APP_NAME")
	if appEntry.Value != "myapp" {
		t.Errorf("APP_NAME should not be redacted, got %q", appEntry.Value)
	}

	apiEntry := findKey(t, out, "API_KEY")
	if apiEntry.Value != "REDACTED" {
		t.Errorf("API_KEY should be redacted, got %q", apiEntry.Value)
	}

	dbEntry := findKey(t, out, "DB_PASSWORD")
	if dbEntry.Value != "REDACTED" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", dbEntry.Value)
	}

	if sum.Redacted != 2 {
		t.Errorf("expected Redacted=2, got %d", sum.Redacted)
	}
}

func TestClone_CustomPlaceholder(t *testing.T) {
	src := entries("SECRET_TOKEN", "abc123")
	opts := DefaultOptions()
	opts.Redact = true
	opts.Placeholder = "***"

	out, _, err := Clone(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := findKey(t, out, "SECRET_TOKEN")
	if e.Value != "***" {
		t.Errorf("expected placeholder ***, got %q", e.Value)
	}
}

func TestClone_StripComments(t *testing.T) {
	src := []parser.Entry{
		{Key: "", Value: "", Comment: "# header comment"},
		{Key: "APP", Value: "1"},
		{Key: "", Value: ""},
	}
	opts := DefaultOptions()
	opts.StripComments = true

	out, sum, err := Clone(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry after stripping, got %d", len(out))
	}
	if sum.Skipped != 2 {
		t.Errorf("expected Skipped=2, got %d", sum.Skipped)
	}
}

func TestClone_CustomSensitiveKeys(t *testing.T) {
	src := entries("MY_INTERNAL", "value", "NORMAL", "ok")
	opts := DefaultOptions()
	opts.Redact = true
	opts.SensitiveKeys = []string{"internal"}

	out, sum, err := Clone(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := findKey(t, out, "MY_INTERNAL")
	if e.Value != "REDACTED" {
		t.Errorf("expected MY_INTERNAL to be redacted")
	}
	n := findKey(t, out, "NORMAL")
	if n.Value != "ok" {
		t.Errorf("expected NORMAL to be unchanged")
	}
	if sum.Redacted != 1 {
		t.Errorf("expected Redacted=1, got %d", sum.Redacted)
	}
}
