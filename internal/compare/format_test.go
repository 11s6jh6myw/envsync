package compare

import (
	"strings"
	"testing"

	"github.com/yourusername/envsync/internal/parser"
)

func makeResult(envs map[string][]parser.Entry) Result {
	return Compare(envs)
}

func TestReport_NoKeys(t *testing.T) {
	var sb strings.Builder
	Report(&sb, Result{}, DefaultFormatOptions())
	if !strings.Contains(sb.String(), "No keys") {
		t.Error("expected 'No keys' message")
	}
}

func TestReport_ShowsEnvHeaders(t *testing.T) {
	r := makeResult(map[string][]parser.Entry{
		"dev":  {{Key: "A", Value: "1"}},
		"prod": {{Key: "A", Value: "1"}},
	})
	var sb strings.Builder
	Report(&sb, r, DefaultFormatOptions())
	out := sb.String()
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Error("expected env labels in header")
	}
}

func TestReport_MarksConflicts(t *testing.T) {
	r := makeResult(map[string][]parser.Entry{
		"dev":  {{Key: "DB", Value: "localhost"}},
		"prod": {{Key: "DB", Value: "prod.db"}},
	})
	var sb strings.Builder
	opts := DefaultFormatOptions()
	Report(&sb, r, opts)
	out := sb.String()
	if !strings.Contains(out, opts.ConflictMark) {
		t.Error("expected conflict marker in output")
	}
}

func TestReport_HidesEqualWhenDisabled(t *testing.T) {
	r := makeResult(map[string][]parser.Entry{
		"dev":  {{Key: "SAME", Value: "x"}, {Key: "DIFF", Value: "1"}},
		"prod": {{Key: "SAME", Value: "x"}, {Key: "DIFF", Value: "2"}},
	})
	var sb strings.Builder
	opts := DefaultFormatOptions()
	opts.ShowEqual = false
	Report(&sb, r, opts)
	out := sb.String()
	if strings.Contains(out, "SAME") {
		t.Error("expected SAME key to be hidden")
	}
	if !strings.Contains(out, "DIFF") {
		t.Error("expected DIFF key to be shown")
	}
}

func TestReport_SummaryLine(t *testing.T) {
	r := makeResult(map[string][]parser.Entry{
		"dev":  {{Key: "A", Value: "1"}, {Key: "B", Value: "x"}},
		"prod": {{Key: "A", Value: "2"}, {Key: "B", Value: "x"}},
	})
	var sb strings.Builder
	Report(&sb, r, DefaultFormatOptions())
	out := sb.String()
	if !strings.Contains(out, "2 key(s) total") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}
