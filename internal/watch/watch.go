// Package watch provides file-watching functionality for .env files,
// detecting changes and emitting diff events when a file is modified.
package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envsync/internal/env"
	"github.com/user/envsync/internal/parser"
)

// Event describes a change detected in a watched .env file.
type Event struct {
	Path      string
	ChangedAt time.Time
	PrevHash  string
	CurrHash  string
	Entries   []parser.Entry
}

// Options configures the Watch behaviour.
type Options struct {
	// PollInterval is how often the file is stat-checked.
	PollInterval time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 2 * time.Second,
	}
}

// Watcher watches a single .env file for changes.
type Watcher struct {
	path    string
	opts    Options
	lastSum string
	stop    chan struct{}
}

// New creates a new Watcher for the given path.
func New(path string, opts Options) *Watcher {
	return &Watcher{
		path: path,
		opts: opts,
		stop: make(chan struct{}),
	}
}

// Start begins polling and sends Events on the returned channel.
// The caller must call Stop to release resources.
func (w *Watcher) Start() (<-chan Event, error) {
	if !env.Exists(w.path) {
		return nil, fmt.Errorf("watch: file not found: %s", w.path)
	}

	sum, err := hashFile(w.path)
	if err != nil {
		return nil, err
	}
	w.lastSum = sum

	ch := make(chan Event, 4)
	go w.poll(ch)
	return ch, nil
}

// Stop halts the watcher.
func (w *Watcher) Stop() { close(w.stop) }

func (w *Watcher) poll(ch chan<- Event) {
	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()
	defer close(ch)

	for {
		select {
		case <-w.stop:
			return
		case t := <-ticker.C:
			sum, err := hashFile(w.path)
			if err != nil || sum == w.lastSum {
				continue
			}
			entries, err := env.Load(w.path)
			if err != nil {
				continue
			}
			ch <- Event{
				Path:      w.path,
				ChangedAt: t,
				PrevHash:  w.lastSum,
				CurrHash:  sum,
				Entries:   entries,
			}
			w.lastSum = sum
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
