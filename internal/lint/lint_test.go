package lint_test

import (
	"testing"

	"github.com/user/envsync/internal/lint"
	"github.com/user/envsync/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func assertIssueCount(t *testing.T, got []lint.Issue, want int) {
	t.Helper()
	if len(got) != want {
		t.Fatalf("expected %d issue(s), got %d: %+v", want, len(got), got)
	}
}

func TestLint_NoIssues(t *testing.T) {
	got := lint.Lint(entries("DB_HOST", "localhost", "APP_PORT", "8080"), lint.DefaultOptions())
	assertIssueCount(t, got, 0)
}

func TestLint_LowercaseKey(t *testing.T) {
	opts := lint.DefaultOptions()
	opts.RequireUppercase = true
	got := lint.Lint(entries("db_host", "localhost"), opts)
	assertIssueCount(t, got, 1)
	if got[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning, got %s", got[0].Severity)
	}
}

func TestLint_KeyTooLong(t *testing.T) {
	opts := lint.DefaultOptions()
	opts.MaxKeyLength = 5
	got := lint.Lint(entries("TOOLONGKEY", "val"), opts)
	if len(got) == 0 {
		t.Fatal("expected at least one issue for long key")
	}
	if !lint.HasErrors(got) {
		t.Error("expected error severity for key length violation")
	}
}

func TestLint_LeadingSpace(t *testing.T) {
	opts := lint.DefaultOptions()
	opts.ForbidSpaces = true
	got := lint.Lint(entries("APP_NAME", " myapp"), opts)
	assertIssueCount(t, got, 1)
	if got[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning, got %s", got[0].Severity)
	}
}

func TestLint_EmptyValueForbidden(t *testing.T) {
	opts := lint.DefaultOptions()
	opts.ForbidEmptyValue = true
	got := lint.Lint(entries("MISSING_VAL", ""), opts)
	assertIssueCount(t, got, 1)
	if got[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %s", got[0].Severity)
	}
}

func TestLint_SkipsBlankEntries(t *testing.T) {
	entries := []parser.Entry{{Key: "", Value: ""}, {Key: "VALID_KEY", Value: "ok"}}
	got := lint.Lint(entries, lint.DefaultOptions())
	assertIssueCount(t, got, 0)
}

func TestHasErrors_False(t *testing.T) {
	issues := []lint.Issue{{Severity: lint.SeverityWarning}}
	if lint.HasErrors(issues) {
		t.Error("expected HasErrors to return false for warnings only")
	}
}

func TestHasErrors_True(t *testing.T) {
	issues := []lint.Issue{{Severity: lint.SeverityError}}
	if !lint.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}
