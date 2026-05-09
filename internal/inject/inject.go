// Package inject provides functionality for injecting .env file values
// into a process environment or a map of substitution targets.
package inject

import (
	"fmt"
	"os"
	"strings"

	"github.com/nicholasgasior/envsync/internal/parser"
)

// DefaultOptions returns a default Options value.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Prefix:    "",
	}
}

// Options controls the behaviour of Inject.
type Options struct {
	// Overwrite existing environment variables when true.
	Overwrite bool
	// Prefix is prepended to every key before injection.
	Prefix string
}

// Result holds the outcome of an Inject call.
type Result struct {
	Injected []string
	Skipped  []string
}

// IntoOS injects entries from the given .env file into the current process
// environment using os.Setenv. Keys that already exist are skipped unless
// Options.Overwrite is true.
func IntoOS(path string, opts Options) (Result, error) {
	entries, err := loadEntries(path)
	if err != nil {
		return Result{}, err
	}

	var res Result
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		key := applyPrefix(opts.Prefix, e.Key)
		if _, exists := os.LookupEnv(key); exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, key)
			continue
		}
		if err := os.Setenv(key, e.Value); err != nil {
			return res, fmt.Errorf("inject: setenv %q: %w", key, err)
		}
		res.Injected = append(res.Injected, key)
	}
	return res, nil
}

// IntoMap injects entries from the given .env file into dst. Keys that already
// exist in dst are skipped unless Options.Overwrite is true.
func IntoMap(path string, dst map[string]string, opts Options) (Result, error) {
	entries, err := loadEntries(path)
	if err != nil {
		return Result{}, err
	}

	var res Result
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		key := applyPrefix(opts.Prefix, e.Key)
		if _, exists := dst[key]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, key)
			continue
		}
		dst[key] = e.Value
		res.Injected = append(res.Injected, key)
	}
	return res, nil
}

func loadEntries(path string) ([]parser.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("inject: open %q: %w", path, err)
	}
	defer f.Close()
	return parser.Parse(f)
}

func applyPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return strings.ToUpper(prefix) + key
}
