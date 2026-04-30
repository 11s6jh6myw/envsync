// Package validate provides linting and validation for .env files.
//
// It checks for common problems including:
//
//   - Duplicate keys (configurable, default: error)
//   - Empty values (configurable, default: warning)
//   - Missing required keys (opt-in via Options.RequiredKeys)
//
// Basic usage:
//
//	entries, _ := parser.Parse(r)
//	opts := validate.DefaultOptions()
//	opts.RequiredKeys = []string{"DATABASE_URL", "SECRET_KEY"}
//	issues := validate.Validate(entries, opts)
//	validate.Report(os.Stdout, issues, validate.DefaultFormatOptions())
//
// Validation is non-destructive — it never modifies the entries slice.
package validate
