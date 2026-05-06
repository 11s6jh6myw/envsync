package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/snapshot"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTake_SetsMetadata(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take("prod.env", entries("FOO", "bar"))
	if s.Source != "prod.env" {
		t.Errorf("expected source prod.env, got %s", s.Source)
	}
	if s.KeyCount != 1 {
		t.Errorf("expected key count 1, got %d", s.KeyCount)
	}
	if s.Timestamp.Before(before) {
		t.Error("timestamp should be set to current time")
	}
	if s.Checksum == "" {
		t.Error("checksum should not be empty")
	}
}

func TestEqual_IdenticalContent(t *testing.T) {
	a := snapshot.Take("a.env", entries("FOO", "bar", "BAZ", "qux"))
	b := snapshot.Take("b.env", entries("FOO", "bar", "BAZ", "qux"))
	if !snapshot.Equal(a, b) {
		t.Error("expected snapshots with identical content to be equal")
	}
}

func TestEqual_DifferentValues(t *testing.T) {
	a := snapshot.Take("a.env", entries("FOO", "bar"))
	b := snapshot.Take("b.env", entries("FOO", "changed"))
	if snapshot.Equal(a, b) {
		t.Error("expected snapshots with different values to be unequal")
	}
}

func TestEqual_OrderIndependent(t *testing.T) {
	a := snapshot.Take("a.env", entries("A", "1", "B", "2"))
	b := snapshot.Take("b.env", entries("B", "2", "A", "1"))
	if !snapshot.Equal(a, b) {
		t.Error("checksum should be order-independent")
	}
}

func TestSerializeDeserialize_RoundTrip(t *testing.T) {
	orig := snapshot.Take("staging.env", entries("KEY", "value", "OTHER", "data"))
	data, err := snapshot.Serialize(orig)
	if err != nil {
		t.Fatalf("serialize error: %v", err)
	}
	got, err := snapshot.Deserialize(data)
	if err != nil {
		t.Fatalf("deserialize error: %v", err)
	}
	if got.Source != orig.Source {
		t.Errorf("source mismatch: got %s, want %s", got.Source, orig.Source)
	}
	if got.Checksum != orig.Checksum {
		t.Errorf("checksum mismatch after round-trip")
	}
	if got.KeyCount != orig.KeyCount {
		t.Errorf("key count mismatch: got %d, want %d", got.KeyCount, orig.KeyCount)
	}
}

func TestDeserialize_InvalidJSON(t *testing.T) {
	_, err := snapshot.Deserialize([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestTake_IgnoresCommentEntries(t *testing.T) {
	commentEntry := parser.Entry{Key: "", Value: "", Raw: "# comment"}
	s := snapshot.Take("test.env", append(entries("FOO", "bar"), commentEntry))
	if s.KeyCount != 1 {
		t.Errorf("expected key count 1 (ignoring comment), got %d", s.KeyCount)
	}
}
