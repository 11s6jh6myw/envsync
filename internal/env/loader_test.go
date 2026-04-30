package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envsync/internal/env"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_BasicFile(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	res, err := env.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Path == "" {
		t.Error("expected non-empty resolved path")
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := env.Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadPair_BothFiles(t *testing.T) {
	srcPath := writeTempEnv(t, "A=1\nB=2\n")
	tgtPath := writeTempEnv(t, "A=1\nC=3\n")

	src, tgt, err := env.LoadPair(srcPath, tgtPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(src.Entries) != 2 {
		t.Errorf("src: expected 2 entries, got %d", len(src.Entries))
	}
	if len(tgt.Entries) != 2 {
		t.Errorf("tgt: expected 2 entries, got %d", len(tgt.Entries))
	}
}

func TestLoadPair_MissingTarget(t *testing.T) {
	srcPath := writeTempEnv(t, "A=1\n")
	_, _, err := env.LoadPair(srcPath, "/no/such/file")
	if err == nil {
		t.Fatal("expected error when target is missing")
	}
}

func TestExists(t *testing.T) {
	p := writeTempEnv(t, "X=1\n")
	if !env.Exists(p) {
		t.Errorf("Exists(%q) = false, want true", p)
	}
	if env.Exists("/no/such/file.env") {
		t.Error("Exists returned true for non-existent file")
	}
}
