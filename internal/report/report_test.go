package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envsync/internal/report"
	"github.com/user/envsync/internal/sync"
)

func makeSummary() sync.Summary {
	return sync.Summary{
		Added:     []string{"NEW_KEY"},
		Removed:   []string{"OLD_KEY"},
		Updated:   []string{"CHANGED_KEY"},
		Unchanged: []string{"STABLE_KEY"},
	}
}

func TestWrite_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, sync.Summary{}, report.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %q", buf.String())
	}
}

func TestWrite_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.UseColor = false

	err := report.Write(&buf, makeSummary(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected '+1 added' in header, got: %q", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected '-1 removed' in header, got: %q", out)
	}
	if !strings.Contains(out, "~1 updated") {
		t.Errorf("expected '~1 updated' in header, got: %q", out)
	}
}

func TestWrite_ShowsUnchanged(t *testing.T) {
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.UseColor = false
	opts.ShowUnchanged = true

	err := report.Write(&buf, makeSummary(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "STABLE_KEY") {
		t.Errorf("expected STABLE_KEY in output when ShowUnchanged=true")
	}
}

func TestWrite_HidesUnchangedByDefault(t *testing.T) {
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.UseColor = false

	err := report.Write(&buf, makeSummary(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "STABLE_KEY") {
		t.Errorf("expected STABLE_KEY to be hidden when ShowUnchanged=false")
	}
}

func TestWrite_KeysListed(t *testing.T) {
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.UseColor = false

	err := report.Write(&buf, makeSummary(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"NEW_KEY", "OLD_KEY", "CHANGED_KEY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}
