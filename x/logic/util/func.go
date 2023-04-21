package util

import "github.com/samber/lo"

// WhitelistBlacklistMatches returns a function that matches the given item according to the given whitelist and
// blacklist returning true if the item matches the whitelist and does not match the blacklist, and false otherwise.
// Note that if the whitelist is empty, the item is considered to match the whitelist.
func WhitelistBlacklistMatches[T any](whitelist []T, blacklist []T, predicate func(item T) func(b T) bool) func(T) bool {
	return func(item T) bool {
		matches := predicate(item)
		return ((len(whitelist) == 0) || lo.ContainsBy(whitelist, matches)) && !lo.ContainsBy(blacklist, matches)
	}
}

// Indexed returns a function that applies the given function to the given item and returns the result.
// It's a convenience function to be used with lo.Map which transforms a predicate function into a mapper function
// by adding an index argument (which is ignored).
func Indexed[T any, U any](f func(t T) U) func(T, int) U {
	return func(t T, _ int) U {
		return f(t)
	}
}
