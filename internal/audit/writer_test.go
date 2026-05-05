package audit_test

import (
	"strings"
	"testing"

	"github.com/your-org/envsync/internal/audit"
)

func TestWrite_NoEvents(t *testing.T) {
	var l audit.Log
	var sb strings.Builder
	if err := audit.Write(&sb, &l, audit.DefaultWriteOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no events") {
		t.Errorf("expected 'no events' message, got: %s", sb.String())
	}
}

func TestWrite_ShowsAllEvents(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventDiff, "a.env", "b.env", nil)
	l.Record(audit.EventSync, "b.env", "c.env", nil)
	var sb strings.Builder
	if err := audit.Write(&sb, &l, audit.DefaultWriteOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "diff") {
		t.Errorf("expected 'diff' in output")
	}
	if !strings.Contains(out, "sync") {
		t.Errorf("expected 'sync' in output")
	}
}

func TestWrite_SummaryLine(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventMerge, "a.env", "b.env", nil)
	l.Record(audit.EventConvert, "b.env", "out.json", nil)
	var sb strings.Builder
	_ = audit.Write(&sb, &l, audit.DefaultWriteOptions())
	if !strings.Contains(sb.String(), "2 event(s) total") {
		t.Errorf("expected summary line, got: %s", sb.String())
	}
}

func TestWrite_Prefix(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventValidate, "a.env", "", nil)
	var sb strings.Builder
	opts := audit.DefaultWriteOptions()
	opts.Prefix = ">>> "
	_ = audit.Write(&sb, &l, opts)
	for _, line := range strings.Split(strings.TrimSpace(sb.String()), "\n") {
		if !strings.HasPrefix(line, ">>> ") {
			t.Errorf("line missing prefix: %q", line)
		}
	}
}
