// Package compare provides multi-environment .env comparison.
//
// It accepts parsed entries from two or more environments and produces a
// unified matrix showing each key's value per environment, along with a
// conflict map highlighting keys whose values differ (or are absent) across
// any two environments.
//
// Basic usage:
//
//	result := compare.Compare(map[string][]parser.Entry{
//		"dev":  devEntries,
//		"stg":  stgEntries,
//		"prod": prodEntries,
//	})
//
//	compare.Report(os.Stdout, result, compare.DefaultFormatOptions())
//
// The Result.Conflicts map contains every key that is either missing from at
// least one environment or has differing values across environments, making it
// easy to identify drift before promoting a release.
package compare
