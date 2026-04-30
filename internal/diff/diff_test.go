package diff_test

import (
	"testing"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestDiff_NoChanges(t *testing.T) {
	src := entries("FOO", "bar", "BAZ", "qux")
	res := diff.Diff(src, src)
	if res.HasChanges() {
		t.Errorf("expected no changes, got %+v", res.Changes)
	}
}

func TestDiff_Added(t *testing.T) {
	src := entries("FOO", "bar", "NEW_KEY", "newval")
	tgt := entries("FOO", "bar")
	res := diff.Diff(src, tgt)
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != diff.Added || c.Key != "NEW_KEY" || c.NewValue != "newval" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_Removed(t *testing.T) {
	src := entries("FOO", "bar")
	tgt := entries("FOO", "bar", "OLD_KEY", "oldval")
	res := diff.Diff(src, tgt)
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != diff.Removed || c.Key != "OLD_KEY" || c.OldValue != "oldval" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_Modified(t *testing.T) {
	src := entries("FOO", "newbar")
	tgt := entries("FOO", "oldbar")
	res := diff.Diff(src, tgt)
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != diff.Modified || c.OldValue != "oldbar" || c.NewValue != "newbar" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_IgnoresComments(t *testing.T) {
	src := []parser.Entry{{IsComment: true, Raw: "# comment"}, {Key: "FOO", Value: "bar"}}
	tgt := []parser.Entry{{Key: "FOO", Value: "bar"}}
	res := diff.Diff(src, tgt)
	if res.HasChanges() {
		t.Errorf("expected no changes, got %+v", res.Changes)
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	src := entries("FOO", "newbar", "NEW_KEY", "newval")
	tgt := entries("FOO", "oldbar", "OLD_KEY", "oldval")
	res := diff.Diff(src, tgt)
	if len(res.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(res.Changes))
	}
	types := map[diff.ChangeType]int{}
	for _, c := range res.Changes {
		types[c.Type]++
	}
	if types[diff.Modified] != 1 || types[diff.Added] != 1 || types[diff.Removed] != 1 {
		t.Errorf("unexpected change type counts: %+v", types)
	}
}
