package validate_test

import (
	"testing"

	"github.com/yourusername/envsync/internal/parser"
	"github.com/yourusername/envsync/internal/validate"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i < len(pairs)-1; i += 2 {
		entries = append(entries, parser.Entry{
			Key:   pairs[i],
			Value: pairs[i+1],
			Line:  i/2 + 1,
		})
	}
	return entries
}

func TestValidate_NoIssues(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "8080")
	issues := validate.Validate(entries, validate.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	entries := []parser.Entry{
		{Key: "HOST", Value: "localhost", Line: 1},
		{Key: "PORT", Value: "8080", Line: 2},
		{Key: "HOST", Value: "remotehost", Line: 3},
	}
	opts := validate.DefaultOptions()
	issues := validate.Validate(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	entries := makeEntries("SECRET", "", "PORT", "3000")
	opts := validate.DefaultOptions()
	issues := validate.Validate(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %s", issues[0].Severity)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	entries := makeEntries("PORT", "8080")
	opts := validate.DefaultOptions()
	opts.RequiredKeys = []string{"DATABASE_URL", "PORT"}
	issues := validate.Validate(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "" {
		t.Errorf("missing key issue should have empty Key field")
	}
}

func TestValidate_SkipsCommentAndBlankLines(t *testing.T) {
	entries := []parser.Entry{
		{Comment: true, Raw: "# comment", Line: 1},
		{Blank: true, Line: 2},
		{Key: "HOST", Value: "localhost", Line: 3},
	}
	issues := validate.Validate(entries, validate.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}
