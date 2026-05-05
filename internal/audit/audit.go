// Package audit records and reports on env file operations.
package audit

import (
	"fmt"
	"time"
)

// EventType classifies the kind of operation that was performed.
type EventType string

const (
	EventDiff    EventType = "diff"
	EventSync    EventType = "sync"
	EventMerge   EventType = "merge"
	EventValidate EventType = "validate"
	EventConvert EventType = "convert"
)

// Event represents a single auditable operation.
type Event struct {
	Timestamp time.Time
	Type      EventType
	Source    string
	Target    string
	Details   map[string]string
}

// Log is an ordered collection of audit events.
type Log struct {
	events []Event
}

// Record appends a new event to the log.
func (l *Log) Record(typ EventType, source, target string, details map[string]string) {
	l.events = append(l.events, Event{
		Timestamp: time.Now().UTC(),
		Type:      typ,
		Source:    source,
		Target:    target,
		Details:   details,
	})
}

// Events returns a copy of all recorded events.
func (l *Log) Events() []Event {
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Len returns the number of recorded events.
func (l *Log) Len() int { return len(l.events) }

// FormatEvent returns a human-readable single-line summary of an event.
func FormatEvent(e Event) string {
	ts := e.Timestamp.Format("2006-01-02T15:04:05Z")
	base := fmt.Sprintf("[%s] %s  source=%q target=%q", ts, e.Type, e.Source, e.Target)
	for k, v := range e.Details {
		base += fmt.Sprintf(" %s=%s", k, v)
	}
	return base
}
