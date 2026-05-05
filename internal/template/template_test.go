package template_test

import (
	"testing"

	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/template"
)

func entries(pairs ...string) []parser.Entry {
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestGenerate_ReplacesValues(t *testing.T) {
	in := entries("DB_HOST", "localhost", "DB_PASS", "secret")
	opts := template.DefaultOptions()
	out := template.Generate(in, opts)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Value != "CHANGE_ME" {
			t.Errorf("key %s: expected CHANGE_ME, got %q", e.Key, e.Value)
		}
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	in := entries("API_KEY", "abc123")
	opts := template.DefaultOptions()
	opts.Placeholder = "YOUR_VALUE_HERE"
	out := template.Generate(in, opts)

	if out[0].Value != "YOUR_VALUE_HERE" {
		t.Errorf("expected YOUR_VALUE_HERE, got %q", out[0].Value)
	}
}

func TestGenerate_KeepValues(t *testing.T) {
	in := entries("PORT", "8080")
	opts := template.DefaultOptions()
	opts.KeepValues = true
	out := template.Generate(in, opts)

	if out[0].Value != "8080" {
		t.Errorf("expected original value 8080, got %q", out[0].Value)
	}
}

func TestGenerate_KeepComments(t *testing.T) {
	in := []parser.Entry{
		{Comment: "# database config"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	opts := template.DefaultOptions()
	opts.KeepComments = true
	out := template.Generate(in, opts)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries (comment + key), got %d", len(out))
	}
	if out[0].Comment != "# database config" {
		t.Errorf("expected comment to be preserved")
	}
}

func TestGenerate_StripComments(t *testing.T) {
	in := []parser.Entry{
		{Comment: "# section header"},
		{Key: "FOO", Value: "bar"},
	}
	opts := template.DefaultOptions()
	opts.KeepComments = false
	out := template.Generate(in, opts)

	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "FOO" {
		t.Errorf("expected FOO entry, got %q", out[0].Key)
	}
}

func TestGenerate_AddsInlineComment(t *testing.T) {
	in := entries("DATABASE_URL", "postgres://localhost/db")
	opts := template.DefaultOptions()
	out := template.Generate(in, opts)

	if out[0].Comment == "" {
		t.Error("expected inline comment hint to be set")
	}
}
