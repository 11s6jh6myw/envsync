package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envsync/internal/watch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestNew_MissingFile_ReturnsError(t *testing.T) {
	w := watch.New("/nonexistent/.env", watch.DefaultOptions())
	_, err := w.Start()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")

	opts := watch.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	w := watch.New(path, opts)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	// Give the watcher one tick before modifying.
	time.Sleep(80 * time.Millisecond)

	if err := os.WriteFile(path, []byte("FOO=baz\nBAR=1\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("path = %q, want %q", ev.Path, path)
		}
		if ev.PrevHash == ev.CurrHash {
			t.Error("hashes should differ after modification")
		}
		if len(ev.Entries) == 0 {
			t.Error("expected parsed entries in event")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")

	opts := watch.DefaultOptions()
	opts.PollInterval = 40 * time.Millisecond

	w := watch.New(path, opts)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	select {
	case ev := <-ch:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	case <-time.After(200 * time.Millisecond):
		// expected: no event
	}
}

func TestDefaultOptions_PollInterval(t *testing.T) {
	opts := watch.DefaultOptions()
	if opts.PollInterval <= 0 {
		t.Errorf("PollInterval should be positive, got %v", opts.PollInterval)
	}
}
