// Package sync provides strategies for synchronising .env files across
// environments.
//
// Three strategies are supported:
//
//   - StrategyFill      – adds keys that are present in source but absent in
//     target; existing target values are never overwritten. This is the safest
//     option for onboarding a new environment without clobbering secrets.
//
//   - StrategyOverwrite – like Fill but also updates the values of keys that
//     already exist in target to match the source values.
//
//   - StrategyExact     – makes the target an exact structural mirror of the
//     source: missing keys are added (with the configured placeholder value)
//     and keys present in target but absent in source are removed.
//
// Typical usage:
//
//	src, _ := parser.Parse(srcReader)
//	tgt, _ := parser.Parse(tgtReader)
//
//	updated, err := sync.Apply(src, tgt, sync.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	parser.Write(outWriter, updated, parser.DefaultWriteOptions())
package sync
