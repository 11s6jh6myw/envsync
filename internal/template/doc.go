// Package template generates .env template files from existing environment
// configurations. It is useful for onboarding new developers or creating
// example .env.example files that can be safely committed to source control.
//
// # Overview
//
// Given a set of parsed [parser.Entry] values, Generate produces a new slice
// where every key's value is replaced with a configurable placeholder string
// (default: "CHANGE_ME"). Inline comment hints are automatically derived from
// the key name so developers know what each variable represents.
//
// # Usage
//
//	entries, _ := env.Load(".env.production")
//	opts := template.DefaultOptions()
//	opts.Placeholder = "FILL_IN"
//	tmpl := template.Generate(entries, opts)
//	_ = parser.Write(os.Stdout, tmpl, parser.DefaultWriteOptions())
//
// # Options
//
//   - Placeholder – value written for every key (default "CHANGE_ME")
//   - KeepComments – whether to copy comment-only lines to the output
//   - KeepValues – when true the original values are preserved (useful for
//     generating a .env.example that shows realistic defaults)
package template
