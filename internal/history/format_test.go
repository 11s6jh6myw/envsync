package history_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/history"
)

func TestReport_NoEntries(t *testing.T) {
	h := history.New()
	var sb strings.Builder
	history.Report(&sb, h, history.DefaultFormatOptions())
	if !strings.Contains(sb.String(), "No history") {
		t.Errorf("expected empty message, got: %s", sb.String())
	}
}

func TestReport_ShowsHeader(t *testing.T) {
	h := history.New()
	h.Record("deploy", envEntries("APP", "prod"))
	var sb strings.Builder
	history.Report(&sb, h, history.DefaultFormatOptions())
	out := sb.String()
	if !strings.Contains(out, "History") {
		t.Errorf("expected header, got: %s", out)
	}
}

func TestReport_ShowsLabel(t *testing.T) {
	h := history.New()
	h.Record("staging-deploy", envEntries("X", "1"))
	var sb strings.Builder
	history.Report(&sb, h, history.DefaultFormatOptions())
	if !strings.Contains(sb.String(), "staging-deploy") {
		t.Error("expected label in output")
	}
}

func TestReport_ShowsKeys(t *testing.T) {
	h := history.New()
	h.Record("v1", envEntries("FOO", "1", "BAR", "2"))
	opts := history.DefaultFormatOptions()
	opts.ShowKeys = true
	var sb strings.Builder
	history.Report(&sb, h, opts)
	out := sb.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected keys in output, got: %s", out)
	}
}

func TestReport_HidesKeys(t *testing.T) {
	h := history.New()
	h.Record("v1", envEntries("SECRET", "val"))
	opts := history.DefaultFormatOptions()
	opts.ShowKeys = false
	var sb strings.Builder
	history.Report(&sb, h, opts)
	if strings.Contains(sb.String(), "SECRET") {
		t.Error("keys should be hidden")
	}
}

func TestReport_MaxEntries(t *testing.T) {
	h := history.New()
	for i := 0; i < 5; i++ {
		h.Record("entry", envEntries("K", "v"))
	}
	opts := history.DefaultFormatOptions()
	opts.MaxEntries = 2
	opts.ShowKeys = false
	var sb strings.Builder
	history.Report(&sb, h, opts)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	// header + 2 entries
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %v", len(lines), lines)
	}
}
