package redact_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/redact"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redact.New(nil, "")

	sensitive := []string{
		"DB_PASSWORD",
		"API_SECRET",
		"GITHUB_TOKEN",
		"STRIPE_API_KEY",
		"PRIVATE_KEY_PATH",
		"AWS_CREDENTIAL",
		"OAUTH_AUTH_CODE",
	}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NonSensitive(t *testing.T) {
	r := redact.New(nil, "")

	insensitive := []string{
		"APP_NAME",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
	}
	for _, key := range insensitive {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestRedact_MasksValue(t *testing.T) {
	r := redact.New(nil, "[REDACTED]")

	got := r.Redact("DB_PASSWORD", "supersecret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestRedact_PassthroughNonSensitive(t *testing.T) {
	r := redact.New(nil, "***")

	got := r.Redact("APP_NAME", "myapp")
	if got != "myapp" {
		t.Errorf("expected myapp, got %q", got)
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r := redact.New([]string{"INTERNAL"}, "<hidden>")

	if !r.IsSensitive("INTERNAL_KEY") {
		t.Error("expected INTERNAL_KEY to be sensitive with custom pattern")
	}
	if r.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to NOT be sensitive with custom patterns")
	}
}

func TestMask_Default(t *testing.T) {
	r := redact.New(nil, "")
	if r.Mask() != "***" {
		t.Errorf("expected default mask ***, got %q", r.Mask())
	}
}
