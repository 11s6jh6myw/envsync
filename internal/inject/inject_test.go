package inject

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestIntoMap_BasicInjection(t *testing.T) {
	p := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	dst := map[string]string{}
	res, err := IntoMap(p, dst, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Fatalf("expected 2 injected, got %d", len(res.Injected))
	}
	if dst["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST: got %q, want %q", dst["APP_HOST"], "localhost")
	}
	if dst["APP_PORT"] != "8080" {
		t.Errorf("APP_PORT: got %q, want %q", dst["APP_PORT"], "8080")
	}
}

func TestIntoMap_SkipsExistingWithoutOverwrite(t *testing.T) {
	p := writeTempEnv(t, "KEY=new\n")
	dst := map[string]string{"KEY": "original"}
	res, err := IntoMap(p, dst, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if dst["KEY"] != "original" {
		t.Errorf("KEY should remain %q, got %q", "original", dst["KEY"])
	}
}

func TestIntoMap_OverwriteExisting(t *testing.T) {
	p := writeTempEnv(t, "KEY=new\n")
	dst := map[string]string{"KEY": "original"}
	opts := DefaultOptions()
	opts.Overwrite = true
	res, err := IntoMap(p, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Injected) != 1 {
		t.Fatalf("expected 1 injected, got %d", len(res.Injected))
	}
	if dst["KEY"] != "new" {
		t.Errorf("KEY: got %q, want %q", dst["KEY"], "new")
	}
}

func TestIntoMap_WithPrefix(t *testing.T) {
	p := writeTempEnv(t, "NAME=world\n")
	dst := map[string]string{}
	opts := DefaultOptions()
	opts.Prefix = "app_"
	_, err := IntoMap(p, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["APP_NAME"] != "world" {
		t.Errorf("APP_NAME: got %q, want %q", dst["APP_NAME"], "world")
	}
}

func TestIntoMap_MissingFile(t *testing.T) {
	_, err := IntoMap("/nonexistent/.env", map[string]string{}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestIntoMap_SkipsBlankKeys(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nVALID=yes\n")
	dst := map[string]string{}
	res, err := IntoMap(p, dst, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Injected) != 1 {
		t.Errorf("expected 1 injected, got %d", len(res.Injected))
	}
}
