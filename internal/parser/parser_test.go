package parser_test

import (
	"strings"
	"testing"

	"github.com/envsync/envsync/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_NAME=envsync
APP_ENV=production
`
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)
	assert.Len(t, ef.Entries, 2)

	e, ok := ef.Get("APP_NAME")
	require.True(t, ok)
	assert.Equal(t, "envsync", e.Value)
}

func TestParse_QuotedValues(t *testing.T) {
	input := `DB_PASS="secret password"
API_KEY='abc123'
`
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)

	dbPass, ok := ef.Get("DB_PASS")
	require.True(t, ok)
	assert.Equal(t, "secret password", dbPass.Value)

	apiKey, ok := ef.Get("API_KEY")
	require.True(t, ok)
	assert.Equal(t, "abc123", apiKey.Value)
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	input := `# Database config
DB_HOST=localhost

DB_PORT=5432
`
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)
	assert.Len(t, ef.Entries, 2)

	host, ok := ef.Get("DB_HOST")
	require.True(t, ok)
	assert.Equal(t, "# Database config", host.Comment)

	port, ok := ef.Get("DB_PORT")
	require.True(t, ok)
	assert.Empty(t, port.Comment)
}

func TestParse_InlineComment(t *testing.T) {
	input := `TIMEOUT=30 # seconds
`
	ef, err := parser.Parse(strings.NewReader(input))
	require.NoError(t, err)

	e, ok := ef.Get("TIMEOUT")
	require.True(t, ok)
	assert.Equal(t, "30", e.Value)
}

func TestParse_InvalidLine(t *testing.T) {
	input := `INVALID_LINE_NO_EQUALS
`
	_, err := parser.Parse(strings.NewReader(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}

func TestParse_EmptyFile(t *testing.T) {
	ef, err := parser.Parse(strings.NewReader(""))
	require.NoError(t, err)
	assert.Empty(t, ef.Entries)
}

func TestGet_MissingKey(t *testing.T) {
	ef, err := parser.Parse(strings.NewReader("FOO=bar\n"))
	require.NoError(t, err)

	_, ok := ef.Get("MISSING")
	assert.False(t, ok)
}
