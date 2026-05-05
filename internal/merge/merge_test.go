package merge_test

import (
	"testing"

	"github.com/user/envsync/internal/merge"
	"github.com/user/envsync/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func findKey(t *testing.T, result []parser.Entry, key string) string {
	t.Helper()
	for _, e := range result {
		if e.Key == key {
			return e.Value
		}
	}
	t.Fatalf("key %q not found in result", key)
	return ""
}

func TestMerge_PreferSource_NoConflict(t *testing.T) {
	src := entries("APP", "prod", "PORT", "8080")
	tgt := entries("APP", "prod", "DEBUG", "false")
	r := merge.Merge(src, tgt, merge.PreferSource)
	if findKey(t, r.Entries, "APP") != "prod" {
		t.Error("expected APP=prod")
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
}

func TestMerge_PreferSource_WinsConflict(t *testing.T) {
	src := entries("DB", "postgres://src")
	tgt := entries("DB", "postgres://tgt")
	r := merge.Merge(src, tgt, merge.PreferSource)
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].Resolved != "postgres://src" {
		t.Errorf("expected source value, got %q", r.Conflicts[0].Resolved)
	}
}

func TestMerge_PreferTarget_WinsConflict(t *testing.T) {
	src := entries("SECRET", "src-secret")
	tgt := entries("SECRET", "tgt-secret")
	r := merge.Merge(src, tgt, merge.PreferTarget)
	if findKey(t, r.Entries, "SECRET") != "tgt-secret" {
		t.Error("expected target value to win")
	}
}

func TestMerge_Union_IncludesAllKeys(t *testing.T) {
	src := entries("A", "1", "B", "2")
	tgt := entries("B", "99", "C", "3")
	r := merge.Merge(src, tgt, merge.Union)
	keys := map[string]bool{}
	for _, e := range r.Entries {
		if e.Key != "" {
			keys[e.Key] = true
		}
	}
	for _, k := range []string{"A", "B", "C"} {
		if !keys[k] {
			t.Errorf("expected key %q in union result", k)
		}
	}
}

func TestMerge_ConflictFields(t *testing.T) {
	src := entries("KEY", "from-src")
	tgt := entries("KEY", "from-tgt")
	r := merge.Merge(src, tgt, merge.PreferSource)
	if len(r.Conflicts) != 1 {
		t.Fatal("expected one conflict")
	}
	c := r.Conflicts[0]
	if c.Key != "KEY" || c.SourceValue != "from-src" || c.TargetValue != "from-tgt" {
		t.Errorf("unexpected conflict fields: %+v", c)
	}
}
