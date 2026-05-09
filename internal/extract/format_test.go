package extract_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/extract"
	"github.com/user/envsync/internal/parser"
)

func makeResult(extracted []string, missing []string) extract.Result {
	var ents []parser.Entry
	for _, k := range extracted {
		ents = append(ents, parser.Entry{Key: k, Value: "val"})
	}
	return extract.Result{Entries: ents, Missing: missing}
}

func TestReport_NoKeys(t *testing.T) {
	var buf strings.Builder
	extract.Report(&buf, extract.Result{}, extract.DefaultFormatOptions())
	if !strings.Contains(buf.String(), "No keys") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestReport_ShowsExtractedCount(t *testing.T) {
	var buf strings.Builder
	r := makeResult([]string{"A", "B"}, nil)
	extract.Report(&buf, r, extract.DefaultFormatOptions())
	if !strings.Contains(buf.String(), "2 key(s)") {
		t.Errorf("expected count in output, got: %s", buf.String())
	}
}

func TestReport_ShowsMissingKeys(t *testing.T) {
	var buf strings.Builder
	r := makeResult([]string{"A"}, []string{"GHOST"})
	opts := extract.DefaultFormatOptions()
	opts.ShowMissing = true
	extract.Report(&buf, r, opts)
	if !strings.Contains(buf.String(), "GHOST") {
		t.Errorf("expected GHOST in output, got: %s", buf.String())
	}
}

func TestReport_HidesMissingWhenDisabled(t *testing.T) {
	var buf strings.Builder
	r := makeResult([]string{"A"}, []string{"GHOST"})
	opts := extract.DefaultFormatOptions()
	opts.ShowMissing = false
	extract.Report(&buf, r, opts)
	if strings.Contains(buf.String(), "GHOST") {
		t.Errorf("GHOST should be hidden, got: %s", buf.String())
	}
}

func TestReport_SummaryAlwaysPresent(t *testing.T) {
	var buf strings.Builder
	r := makeResult([]string{"X"}, nil)
	extract.Report(&buf, r, extract.DefaultFormatOptions())
	if !strings.Contains(buf.String(), "Summary:") {
		t.Errorf("expected Summary line, got: %s", buf.String())
	}
}
