package mask_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/mask"
	"github.com/yourorg/envsync/internal/parser"
)

func TestValue_ShortValue_FullyMasked(t *testing.T) {
	opts := mask.DefaultOptions()
	result := mask.Value("abc", opts)
	if result != "***" {
		t.Errorf("expected '***', got %q", result)
	}
}

func TestValue_LongValue_PartiallyMasked(t *testing.T) {
	opts := mask.DefaultOptions()
	result := mask.Value("supersecretkey1234", opts)
	if result[len(result)-4:] != "1234" {
		t.Errorf("expected last 4 chars to be visible, got %q", result)
	}
	for _, ch := range result[:len(result)-4] {
		if ch != '*' {
			t.Errorf("expected masked chars to be '*', got %q", string(ch))
		}
	}
}

func TestValue_EmptyString_ReturnsEmpty(t *testing.T) {
	opts := mask.DefaultOptions()
	result := mask.Value("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestValue_FullMask_MasksAll(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.FullMask = true
	result := mask.Value("supersecretkey1234", opts)
	for _, ch := range result {
		if ch != '*' {
			t.Errorf("expected all chars masked, got %q", result)
		}
	}
}

func TestValue_CustomMaskChar(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.MaskChar = '#'
	result := mask.Value("supersecretkey1234", opts)
	for _, ch := range result[:len(result)-4] {
		if ch != '#' {
			t.Errorf("expected '#' mask char, got %q", string(ch))
		}
	}
}

func TestEntries_MasksAllKeys_WhenNilKeySet(t *testing.T) {
	entries := []parser.Entry{
		{Key: "API_KEY", Value: "supersecretvalue"},
		{Key: "HOST", Value: "localhost"},
	}
	opts := mask.DefaultOptions()
	result := mask.Entries(entries, nil, opts)
	for _, e := range result {
		if e.Key != "" && e.Value == entries[0].Value {
			t.Errorf("expected value to be masked for key %q", e.Key)
		}
	}
}

func TestEntries_MasksOnlyMatchingKeys(t *testing.T) {
	entries := []parser.Entry{
		{Key: "API_KEY", Value: "supersecretvalue"},
		{Key: "HOST", Value: "localhost1234"},
	}
	keys := map[string]struct{}{"API_KEY": {}}
	opts := mask.DefaultOptions()
	result := mask.Entries(entries, keys, opts)
	if result[0].Value == "supersecretvalue" {
		t.Error("expected API_KEY value to be masked")
	}
	if result[1].Value != "localhost1234" {
		t.Errorf("expected HOST value to be unchanged, got %q", result[1].Value)
	}
}

func TestEntries_PreservesCommentEntries(t *testing.T) {
	entries := []parser.Entry{
		{Key: "", Value: "", Comment: "# section header"},
		{Key: "TOKEN", Value: "mytoken1234"},
	}
	opts := mask.DefaultOptions()
	result := mask.Entries(entries, nil, opts)
	if result[0].Comment != "# section header" {
		t.Errorf("expected comment entry to be preserved, got %+v", result[0])
	}
}
