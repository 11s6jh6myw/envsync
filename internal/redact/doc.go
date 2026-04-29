// Package redact provides key-based redaction for .env file values.
//
// It is used to mask sensitive environment variable values when displaying
// diffs or syncing files, ensuring secrets such as passwords, tokens, and
// API keys are never printed in plain text.
//
// Usage:
//
//	// Create a redactor with default sensitive patterns and mask
//	r := redact.New(nil, "")
//
//	// Check whether a key is considered sensitive
//	if r.IsSensitive("DB_PASSWORD") {
//	    fmt.Println("sensitive!")
//	}
//
//	// Redact a value (returns mask if sensitive, original otherwise)
//	display := r.Redact("DB_PASSWORD", "s3cr3t") // returns "***"
//
// Custom patterns and masks can be provided:
//
//	r := redact.New([]string{"INTERNAL", "PRIVATE"}, "[hidden]")
package redact
