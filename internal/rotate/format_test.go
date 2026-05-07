package rotate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envsync/internal/rotate"
)

func TestReport_NoRotations(t *testing.T) {
	var buf bytes.Buffer
	rotate.Report(&buf, rotate.Result{}, rotate.DefaultFormatOptions())
	if !strings.Contains(buf.String(), "No keys rotated") {
		t.Errorf("expected 'No keys rotated' message, got: %s", buf.String())
	}
}

func TestReport_ShowsRotatedKeys(t *testing.T) {
	r := rotate.Result{
		Rotated: []string{"API_KEY", "DB_PASSWORD"},
	}
	var buf bytes.Buffer
	rotate.Report(&buf, r, rotate.DefaultFormatOptions())
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output")
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, "Rotated (2)") {
		t.Errorf("expected rotated header with count")
	}
}

func TestReport_VerboseShowsSkipped(t *testing.T) {
	r := rotate.Result{
		Rotated: []string{"API_KEY"},
		Skipped: []string{"SECRET_TOKEN"},
	}
	opts := rotate.DefaultFormatOptions()
	opts.Verbose = true

	var buf bytes.Buffer
	rotate.Report(&buf, r, opts)
	out := buf.String()
	if !strings.Contains(out, "SECRET_TOKEN") {
		t.Errorf("expected SECRET_TOKEN in verbose output")
	}
	if !strings.Contains(out, "Skipped") {
		t.Errorf("expected Skipped section in verbose output")
	}
}

func TestReport_HidesSkippedByDefault(t *testing.T) {
	r := rotate.Result{
		Skipped: []string{"SOME_SECRET"},
	}
	opts := rotate.DefaultFormatOptions()
	opts.Verbose = false

	var buf bytes.Buffer
	rotate.Report(&buf, r, opts)
	if strings.Contains(buf.String(), "SOME_SECRET") {
		t.Errorf("skipped keys should be hidden when Verbose=false")
	}
}

func TestReport_SummaryAlwaysPresent(t *testing.T) {
	r := rotate.Result{
		Rotated: []string{"X"},
	}
	var buf bytes.Buffer
	rotate.Report(&buf, r, rotate.DefaultFormatOptions())
	if !strings.Contains(buf.String(), "rotated") {
		t.Errorf("expected summary line containing 'rotated'")
	}
}
