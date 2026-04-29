package redact_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/redact"
)

// TestRedact_DiffIntegration simulates how the redactor would be used
// when formatting a diff — sensitive values should be masked, others shown.
func TestRedact_DiffIntegration(t *testing.T) {
	r := redact.New(nil, "***")

	type kv struct {
		key      string
		value    string
		wantMask bool
	}

	cases := []kv{
		{"DB_PASSWORD", "hunter2", true},
		{"STRIPE_SECRET", "sk_live_abc", true},
		{"APP_ENV", "production", false},
		{"PORT", "8080", false},
		{"GITHUB_TOKEN", "ghp_xyz", true},
		{"LOG_LEVEL", "debug", false},
	}

	for _, tc := range cases {
		result := r.Redact(tc.key, tc.value)
		if tc.wantMask {
			if result != "***" {
				t.Errorf("key %q: expected masked value, got %q", tc.key, result)
			}
		} else {
			if result != tc.value {
				t.Errorf("key %q: expected %q, got %q", tc.key, tc.value, result)
			}
		}
	}
}

// TestRedact_CaseInsensitiveKeys ensures keys are matched regardless of case.
func TestRedact_CaseInsensitiveKeys(t *testing.T) {
	r := redact.New(nil, "***")

	variants := []string{
		"db_password",
		"Db_Password",
		"DB_PASSWORD",
		"dB_pAsSwOrD",
	}
	for _, key := range variants {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive (case-insensitive match)", key)
		}
	}
}
