package convert

import (
	"bytes"
	"strings"
	"testing"
)

func TestReport_NoSkipped(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFormatOptions()
	opts.Color = false

	Report(&buf, "json", 5, 0, opts)

	out := buf.String()
	if !strings.Contains(out, "JSON") {
		t.Errorf("expected format name in output, got: %q", out)
	}
	if !strings.Contains(out, "entries : 5") {
		t.Errorf("expected entry count in output, got: %q", out)
	}
	if strings.Contains(out, "skipped") {
		t.Errorf("expected no skipped line when skipped=0, got: %q", out)
	}
}

func TestReport_WithSkipped(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFormatOptions()
	opts.Color = false

	Report(&buf, "export", 3, 2, opts)

	out := buf.String()
	if !strings.Contains(out, "EXPORT") {
		t.Errorf("expected EXPORT in output, got: %q", out)
	}
	if !strings.Contains(out, "entries : 3") {
		t.Errorf("expected entry count, got: %q", out)
	}
	if !strings.Contains(out, "skipped : 2") {
		t.Errorf("expected skipped count, got: %q", out)
	}
}

func TestReport_HideCount(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFormatOptions()
	opts.Color = false
	opts.ShowCount = false

	Report(&buf, "json", 10, 1, opts)

	out := buf.String()
	if strings.Contains(out, "entries") {
		t.Errorf("expected no entries line when ShowCount=false, got: %q", out)
	}
	if strings.Contains(out, "skipped") {
		t.Errorf("expected no skipped line when ShowCount=false, got: %q", out)
	}
}

func TestReport_ColorOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFormatOptions()
	opts.Color = true

	// Should not panic with color enabled
	Report(&buf, "json", 4, 1, opts)

	out := buf.String()
	if !strings.Contains(out, "JSON") {
		t.Errorf("expected JSON in color output, got: %q", out)
	}
}
