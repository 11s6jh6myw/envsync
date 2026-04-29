package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline or preceding comment
	Line    int
}

// EnvFile holds all parsed entries from an .env file.
type EnvFile struct {
	Entries []Entry
	Index   map[string]int // key -> index in Entries
}

// Parse reads an .env file from the given reader and returns an EnvFile.
func Parse(r io.Reader) (*EnvFile, error) {
	ef := &EnvFile{
		Index: make(map[string]int),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0
	var pendingComment string

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			pendingComment = trimmed
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, trimmed)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripInlineComment(value)
		value = unquote(value)

		entry := Entry{
			Key:     key,
			Value:   value,
			Comment: pendingComment,
			Line:    lineNum,
		}
		ef.Index[key] = len(ef.Entries)
		ef.Entries = append(ef.Entries, entry)
		pendingComment = ""
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning error: %w", err)
	}

	return ef, nil
}

// Get returns the Entry for the given key, and whether it was found.
func (ef *EnvFile) Get(key string) (Entry, bool) {
	if idx, ok := ef.Index[key]; ok {
		return ef.Entries[idx], true
	}
	return Entry{}, false
}

func stripInlineComment(s string) string {
	if idx := strings.Index(s, " #"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
