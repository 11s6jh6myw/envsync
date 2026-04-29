package diff_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/parser"
)

func makeEntries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestFormat_NoChanges(t *testing.T) {
	src := makeEntries("HOST", "localhost", "PORT", "5432")
	changes := diff.Diff(src, src)
	opts := diff.DefaultFormatOptions()
	opts.Color = false
	out := diff.Format(changes, opts)
	if strings.Contains(out, "+") || strings.Contains(out, "-") {
		t.Errorf("expected no diff markers, got:\n%s", out)
	}
}

func TestFormat_ShowsAdded(t *testing.T) {
	src := makeEntries("HOST", "localhost")
	dst := makeEntries("HOST", "localhost", "PORT", "5432")
	changes := diff.Diff(src, dst)
	opts := diff.DefaultFormatOptions()
	opts.Color = false
	out := diff.Format(changes, opts)
	if !strings.Contains(out, "+ PORT") {
		t.Errorf("expected '+ PORT' in output, got:\n%s", out)
	}
}

func TestFormat_ShowsRemoved(t *testing.T) {
	src := makeEntries("HOST", "localhost", "PORT", "5432")
	dst := makeEntries("HOST", "localhost")
	changes := diff.Diff(src, dst)
	opts := diff.DefaultFormatOptions()
	opts.Color = false
	out := diff.Format(changes, opts)
	if !strings.Contains(out, "- PORT") {
		t.Errorf("expected '- PORT' in output, got:\n%s", out)
	}
}

func TestFormat_RedactsSecrets(t *testing.T) {
	src := makeEntries("DB_PASSWORD", "old-secret")
	dst := makeEntries("DB_PASSWORD", "new-secret")
	changes := diff.Diff(src, dst)
	opts := diff.DefaultFormatOptions()
	opts.Color = false
	opts.Redact = true
	out := diff.Format(changes, opts)
	if strings.Contains(out, "old-secret") || strings.Contains(out, "new-secret") {
		t.Errorf("expected secrets to be redacted, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key name to appear in output, got:\n%s", out)
	}
}

func TestFormat_UnchangedHidden(t *testing.T) {
	src := makeEntries("HOST", "localhost", "PORT", "5432")
	dst := makeEntries("HOST", "remotehost", "PORT", "5432")
	changes := diff.Diff(src, dst)
	opts := diff.DefaultFormatOptions()
	opts.Color = false
	opts.ShowUnchanged = false
	out := diff.Format(changes, opts)
	if strings.Contains(out, "PORT") {
		t.Errorf("expected unchanged key PORT to be hidden, got:\n%s", out)
	}
}
