package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/envsync/internal/audit"
)

func TestRecord_AppendsEvent(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventDiff, "a.env", "b.env", nil)
	if l.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", l.Len())
	}
}

func TestRecord_MultipleEvents(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventSync, "a.env", "b.env", nil)
	l.Record(audit.EventMerge, "b.env", "c.env", nil)
	if l.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", l.Len())
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventValidate, "x.env", "", nil)
	evs := l.Events()
	evs[0].Source = "mutated"
	if l.Events()[0].Source == "mutated" {
		t.Fatal("Events() should return a copy, not a reference")
	}
}

func TestRecord_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	var l audit.Log
	l.Record(audit.EventConvert, "a.env", "out.json", nil)
	after := time.Now().UTC()
	ev := l.Events()[0]
	if ev.Timestamp.Before(before) || ev.Timestamp.After(after) {
		t.Fatalf("timestamp %v out of range [%v, %v]", ev.Timestamp, before, after)
	}
}

func TestFormatEvent_ContainsFields(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventDiff, "src.env", "tgt.env", map[string]string{"added": "3"})
	line := audit.FormatEvent(l.Events()[0])
	for _, want := range []string{"diff", "src.env", "tgt.env", "added=3"} {
		if !strings.Contains(line, want) {
			t.Errorf("FormatEvent output missing %q: %s", want, line)
		}
	}
}

func TestFormatEvent_NoDetails(t *testing.T) {
	var l audit.Log
	l.Record(audit.EventSync, "a.env", "b.env", nil)
	line := audit.FormatEvent(l.Events()[0])
	if !strings.Contains(line, "sync") {
		t.Errorf("expected 'sync' in output: %s", line)
	}
}
