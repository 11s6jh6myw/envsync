package validate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envsync/internal/validate"
)

func TestReport_NoIssues(t *testing.T) {
	var buf strings.Builder
	opts := validate.FormatOptions{Color: false, File: "prod.env"}
	n := validate.Report(&buf, nil, opts)
	if n != 0 {
		t.Errorf("expected 0 errors, got %d", n)
	}
	if !strings.Contains(buf.String(), "no issues found") {
		t.Errorf("expected 'no issues found' in output, got: %s", buf.String())
	}
}

func TestReport_ShowsErrors(t *testing.T) {
	issues := []validate.Issue{
		{Line: 3, Key: "HOST", Message: "duplicate key \"HOST\" (first seen at line 1)", Severity: "error"},
	}
	var buf strings.Builder
	opts := validate.FormatOptions{Color: false}
	n := validate.Report(&buf, issues, opts)
	if n != 1 {
		t.Errorf("expected 1 error count, got %d", n)
	}
	if !strings.Contains(buf.String(), "[error]") {
		t.Errorf("expected [error] tag in output")
	}
	if !strings.Contains(buf.String(), "duplicate") {
		t.Errorf("expected duplicate message in output")
	}
}

func TestReport_ShowsWarnings(t *testing.T) {
	issues := []validate.Issue{
		{Line: 2, Key: "SECRET", Message: "key \"SECRET\" has an empty value", Severity: "warning"},
	}
	var buf strings.Builder
	opts := validate.FormatOptions{Color: false, File: "dev.env"}
	n := validate.Report(&buf, issues, opts)
	if n != 0 {
		t.Errorf("warnings should not increment error count, got %d", n)
	}
	if !strings.Contains(buf.String(), "[warn]") {
		t.Errorf("expected [warn] tag in output")
	}
}

func TestReport_SummaryLine(t *testing.T) {
	issues := []validate.Issue{
		{Severity: "error", Message: "some error"},
		{Severity: "warning", Message: "some warning"},
		{Severity: "warning", Message: "another warning"},
	}
	var buf strings.Builder
	opts := validate.FormatOptions{Color: false}
	validate.Report(&buf, issues, opts)
	out := buf.String()
	if !strings.Contains(out, "1 error(s), 2 warning(s)") {
		t.Errorf("unexpected summary line in: %s", out)
	}
}
