package parser_test

import (
	"strings"
	"testing"

	"github.com/envsync/envsync/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite_RoundTrip(t *testing.T) {
	input := "APP=hello\nDEBUG=true\n"
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)

	var sb strings.Builder
	err = parser.Write(&sb, ef, parser.DefaultWriteOptions())
	require.NoError(t, err)

	out := sb.String()
	assert.Contains(t, out, "APP=hello")
	assert.Contains(t, out, "DEBUG=true")
}

func TestWrite_QuoteValues(t *testing.T) {
	ef := &parser.EnvFile{
		Entries: []parser.Entry{
			{Key: "MSG", Value: "hello world"},
		},
		Index: map[string]int{"MSG": 0},
	}

	var sb strings.Builder
	opts := parser.WriteOptions{QuoteValues: true, KeepComments: true}
	err := parser.Write(&sb, ef, opts)
	require.NoError(t, err)
	assert.Contains(t, sb.String(), `MSG="hello world"`)
}

func TestWrite_KeepComments(t *testing.T) {
	input := "# section\nKEY=val\n"
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)

	var sb strings.Builder
	err = parser.Write(&sb, ef, parser.DefaultWriteOptions())
	require.NoError(t, err)
	assert.Contains(t, sb.String(), "# section")
}

func TestWrite_NeedsQuotingSpace(t *testing.T) {
	input := "GREETING=hello world\n"
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)

	var sb strings.Builder
	err = parser.Write(&sb, ef, parser.DefaultWriteOptions())
	require.NoError(t, err)
	// value with space should be auto-quoted
	assert.Contains(t, sb.String(), `GREETING="hello world"`)
}
