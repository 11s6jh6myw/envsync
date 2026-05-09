package filter_test

import (
	"testing"

	"github.com/yourusername/envsync/internal/filter"
	"github.com/yourusername/envsync/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	in := entries("APP_HOST", "localhost", "DB_URL", "postgres://")
	r := filter.Filter(in, filter.DefaultOptions())
	if r.Matched != 2 {
		t.Fatalf("expected 2 matched, got %d", r.Matched)
	}
	if r.Skipped != 0 {
		t.Fatalf("expected 0 skipped, got %d", r.Skipped)
	}
}

func TestFilter_Prefix_FiltersCorrectly(t *testing.T) {
	in := entries("APP_HOST", "localhost", "DB_URL", "postgres://", "APP_PORT", "8080")
	opts := filter.DefaultOptions()
	opts.Prefix = "APP_"
	r := filter.Filter(in, opts)
	if r.Matched != 2 {
		t.Fatalf("expected 2 matched, got %d", r.Matched)
	}
	if r.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", r.Skipped)
	}
}

func TestFilter_Suffix_FiltersCorrectly(t *testing.T) {
	in := entries("DB_URL", "postgres://", "API_URL", "https://", "APP_HOST", "localhost")
	opts := filter.DefaultOptions()
	opts.Suffix = "_URL"
	r := filter.Filter(in, opts)
	if r.Matched != 2 {
		t.Fatalf("expected 2 matched, got %d", r.Matched)
	}
}

func TestFilter_Contains_FiltersCorrectly(t *testing.T) {
	in := entries("AWS_SECRET_KEY", "abc", "AWS_ACCESS_KEY", "def", "DB_HOST", "localhost")
	opts := filter.DefaultOptions()
	opts.Contains = "SECRET"
	r := filter.Filter(in, opts)
	if r.Matched != 1 {
		t.Fatalf("expected 1 matched, got %d", r.Matched)
	}
}

func TestFilter_Invert_NegatesMatch(t *testing.T) {
	in := entries("APP_HOST", "localhost", "DB_URL", "postgres://", "APP_PORT", "8080")
	opts := filter.DefaultOptions()
	opts.Prefix = "APP_"
	opts.Invert = true
	r := filter.Filter(in, opts)
	if r.Matched != 1 {
		t.Fatalf("expected 1 matched (inverted), got %d", r.Matched)
	}
	if r.Entries[0].Key != "DB_URL" {
		t.Fatalf("expected DB_URL, got %s", r.Entries[0].Key)
	}
}

func TestFilter_KeepComments_PreservesBlankAndCommentLines(t *testing.T) {
	in := []parser.Entry{
		{IsComment: true, Raw: "# section"},
		{Key: "APP_HOST", Value: "localhost"},
		{IsBlank: true},
		{Key: "DB_URL", Value: "postgres://"},
	}
	opts := filter.DefaultOptions()
	opts.Prefix = "APP_"
	r := filter.Filter(in, opts)
	// expect: comment + APP_HOST + blank = 3 entries
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries (with comments), got %d", len(r.Entries))
	}
}

func TestFilter_CaseInsensitivePrefix(t *testing.T) {
	in := entries("app_host", "localhost", "DB_URL", "postgres://")
	opts := filter.DefaultOptions()
	opts.Prefix = "APP_"
	r := filter.Filter(in, opts)
	if r.Matched != 1 {
		t.Fatalf("expected case-insensitive match, got %d", r.Matched)
	}
}
