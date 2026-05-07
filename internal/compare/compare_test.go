package compare

import (
	"testing"

	"github.com/yourusername/envsync/internal/parser"
)

func envEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_UnionOfKeys(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev":  envEntries("A", "1", "B", "2"),
		"prod": envEntries("A", "1", "C", "3"),
	})
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
}

func TestCompare_NoConflictWhenEqual(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev":  envEntries("A", "same"),
		"prod": envEntries("A", "same"),
	})
	if r.Conflicts["A"] {
		t.Error("expected no conflict for identical values")
	}
}

func TestCompare_ConflictWhenDifferent(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev":  envEntries("DB_URL", "localhost"),
		"prod": envEntries("DB_URL", "prod.example.com"),
	})
	if !r.Conflicts["DB_URL"] {
		t.Error("expected conflict for differing values")
	}
}

func TestCompare_MissingKeyIsConflict(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev":  envEntries("ONLY_DEV", "yes"),
		"prod": envEntries("OTHER", "no"),
	})
	if !r.Conflicts["ONLY_DEV"] {
		t.Error("expected conflict for key missing in prod")
	}
}

func TestCompare_MatrixValues(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev":  envEntries("PORT", "3000"),
		"prod": envEntries("PORT", "8080"),
	})
	if r.Matrix["PORT"]["dev"] != "3000" {
		t.Errorf("expected 3000, got %s", r.Matrix["PORT"]["dev"])
	}
	if r.Matrix["PORT"]["prod"] != "8080" {
		t.Errorf("expected 8080, got %s", r.Matrix["PORT"]["prod"])
	}
}

func TestCompare_KeysAreSorted(t *testing.T) {
	r := Compare(map[string][]parser.Entry{
		"dev": envEntries("Z_KEY", "z", "A_KEY", "a", "M_KEY", "m"),
	})
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range r.Keys {
		if k != expected[i] {
			t.Errorf("pos %d: expected %s got %s", i, expected[i], k)
		}
	}
}
