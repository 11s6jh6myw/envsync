package convert_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/convert"
	"github.com/yourorg/envsync/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestToJSON_BasicEntries(t *testing.T) {
	e := entries("APP_ENV", "production", "PORT", "8080")
	out, err := convert.ToJSON(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", m["PORT"])
	}
}

func TestToJSON_SkipsBlankKeys(t *testing.T) {
	e := []parser.Entry{{Key: "", Value: ""}, {Key: "FOO", Value: "bar"}}
	out, err := convert.ToJSON(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	_ = json.Unmarshal([]byte(out), &m)
	if _, ok := m[""]; ok {
		t.Error("blank key should not appear in JSON output")
	}
}

func TestToExport_SimpleValues(t *testing.T) {
	e := entries("DB_HOST", "localhost", "DEBUG", "true")
	out := convert.ToExport(e)
	if !strings.Contains(out, "export DB_HOST=localhost") {
		t.Errorf("expected export DB_HOST=localhost in output:\n%s", out)
	}
	if !strings.Contains(out, "export DEBUG=true") {
		t.Errorf("expected export DEBUG=true in output:\n%s", out)
	}
}

func TestToExport_QuotesValuesWithSpaces(t *testing.T) {
	e := entries("GREETING", "hello world")
	out := convert.ToExport(e)
	if !strings.Contains(out, `export GREETING="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestFromJSON_RoundTrip(t *testing.T) {
	original := entries("KEY1", "value1", "KEY2", "value2")
	jsonStr, err := convert.ToJSON(original)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	result, err := convert.FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	m := make(map[string]string)
	for _, e := range result {
		m[e.Key] = e.Value
	}
	for _, e := range original {
		if m[e.Key] != e.Value {
			t.Errorf("key %q: expected %q got %q", e.Key, e.Value, m[e.Key])
		}
	}
}

func TestFromJSON_InvalidInput(t *testing.T) {
	_, err := convert.FromJSON(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON input")
	}
}
