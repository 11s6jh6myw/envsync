package snapshot

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"envsync/internal/parser"
)

func makeSnapshot(label string, kvs [][2]string) Snapshot {
	var entries []parser.Entry
	for _, kv := range kvs {
		entries = append(entries, parser.Entry{Key: kv[0], Value: kv[1]})
	}
	return Snapshot{
		Label:   label,
		Entries: entries,
		TakenAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestReport_ShowsHeader(t *testing.T) {
	s := makeSnapshot("production", [][2]string{{"FOO", "bar"}})
	var buf bytes.Buffer
	Report(&buf, s, DefaultFormatOptions())
	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected label in output, got:\n%s", buf.String())
	}
}

func TestReport_ShowsKeyCount(t *testing.T) {
	s := makeSnapshot("staging", [][2]string{{"A", "1"}, {"B", "2"}, {"C", "3"}})
	var buf bytes.Buffer
	Report(&buf, s, DefaultFormatOptions())
	if !strings.Contains(buf.String(), "Keys  : 3") {
		t.Errorf("expected key count in output, got:\n%s", buf.String())
	}
}

func TestReport_ShowsKeyNames(t *testing.T) {
	s := makeSnapshot("", [][2]string{{"DB_HOST", "localhost"}, {"PORT", "5432"}})
	opts := DefaultFormatOptions()
	opts.ShowKeys = true
	var buf bytes.Buffer
	Report(&buf, s, opts)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "PORT") {
		t.Errorf("expected key names in output, got:\n%s", out)
	}
}

func TestReport_HidesKeysWhenDisabled(t *testing.T) {
	s := makeSnapshot("dev", [][2]string{{"SECRET", "abc"}})
	opts := DefaultFormatOptions()
	opts.ShowKeys = false
	var buf bytes.Buffer
	Report(&buf, s, opts)
	if strings.Contains(buf.String(), "SECRET") {
		t.Errorf("expected key names to be hidden, got:\n%s", buf.String())
	}
}

func TestReport_FallbackLabelWhenEmpty(t *testing.T) {
	s := makeSnapshot("", [][2]string{})
	var buf bytes.Buffer
	Report(&buf, s, DefaultFormatOptions())
	if !strings.Contains(buf.String(), "Snapshot") {
		t.Errorf("expected fallback header, got:\n%s", buf.String())
	}
}

func TestReport_CustomTimeFormat(t *testing.T) {
	s := makeSnapshot("ci", [][2]string{{"X", "1"}})
	opts := DefaultFormatOptions()
	opts.TimeFormat = "2006-01-02"
	var buf bytes.Buffer
	Report(&buf, s, opts)
	if !strings.Contains(buf.String(), "2024-06-01") {
		t.Errorf("expected formatted date in output, got:\n%s", buf.String())
	}
}
