package schema_test

import (
	"testing"

	"github.com/your-org/envsync/internal/parser"
	"github.com/your-org/envsync/internal/schema"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCheck_NoIssues(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "APP_ENV", Required: true, Allowed: []string{"development", "production"}},
	})
	issues := s.Check(entries("APP_ENV", "development"))
	// Only the undefined-key warning is absent; APP_ENV is defined.
	for _, issue := range issues {
		if issue.Error {
			t.Errorf("unexpected error issue: %+v", issue)
		}
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "DATABASE_URL", Required: true},
	})
	issues := s.Check(entries("APP_ENV", "development"))
	if !schema.HasErrors(issues) {
		t.Fatal("expected error for missing required key")
	}
	found := false
	for _, i := range issues {
		if i.Key == "DATABASE_URL" && i.Error {
			found = true
		}
	}
	if !found {
		t.Error("expected issue for DATABASE_URL")
	}
}

func TestCheck_DisallowedValue(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "LOG_LEVEL", Required: true, Allowed: []string{"info", "debug", "error"}},
	})
	issues := s.Check(entries("LOG_LEVEL", "verbose"))
	if !schema.HasErrors(issues) {
		t.Fatal("expected error for disallowed value")
	}
}

func TestCheck_AllowedValue(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "LOG_LEVEL", Required: true, Allowed: []string{"info", "debug"}},
	})
	issues := s.Check(entries("LOG_LEVEL", "info"))
	for _, i := range issues {
		if i.Error {
			t.Errorf("unexpected error: %+v", i)
		}
	}
}

func TestCheck_UndefinedKeyWarning(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "APP_ENV", Required: false},
	})
	issues := s.Check(entries("APP_ENV", "dev", "UNKNOWN_KEY", "foo"))
	warning := false
	for _, i := range issues {
		if i.Key == "UNKNOWN_KEY" && !i.Error {
			warning = true
		}
	}
	if !warning {
		t.Error("expected warning for undefined key UNKNOWN_KEY")
	}
}

func TestHasErrors_Empty(t *testing.T) {
	if schema.HasErrors(nil) {
		t.Error("expected no errors for nil slice")
	}
}
