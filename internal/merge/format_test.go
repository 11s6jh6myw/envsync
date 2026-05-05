package merge_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/merge"
)

func TestReport_NoConflicts(t *testing.T) {
	result := merge.Result{}
	var buf strings.Builder
	merge.Report(&buf, result, merge.FormatOptions{UseColor: false})
	if !strings.Contains(buf.String(), "no conflicts") {
		t.Errorf("expected no-conflict message, got: %q", buf.String())
	}
}

func TestReport_ShowsConflictCount(t *testing.T) {
	result := merge.Result{
		Conflicts: []merge.Conflict{
			{Key: "DB_URL", SourceValue: "src", TargetValue: "tgt", Resolved: "src"},
			{Key: "SECRET", SourceValue: "a", TargetValue: "b", Resolved: "a"},
		},
	}
	var buf strings.Builder
	merge.Report(&buf, result, merge.FormatOptions{UseColor: false})
	if !strings.Contains(buf.String(), "2 conflict") {
		t.Errorf("expected conflict count in output, got: %q", buf.String())
	}
}

func TestReport_ShowsKeyAndValues(t *testing.T) {
	result := merge.Result{
		Conflicts: []merge.Conflict{
			{Key: "PORT", SourceValue: "8080", TargetValue: "9090", Resolved: "8080"},
		},
	}
	var buf strings.Builder
	merge.Report(&buf, result, merge.FormatOptions{UseColor: false})
	out := buf.String()
	for _, want := range []string{"PORT", "8080", "9090"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestReport_ColorOutput(t *testing.T) {
	result := merge.Result{
		Conflicts: []merge.Conflict{
			{Key: "X", SourceValue: "1", TargetValue: "2", Resolved: "1"},
		},
	}
	var buf strings.Builder
	merge.Report(&buf, result, merge.FormatOptions{UseColor: true})
	if !strings.Contains(buf.String(), "\033[") {
		t.Error("expected ANSI escape codes in color output")
	}
}
