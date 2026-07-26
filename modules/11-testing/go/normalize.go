// Package testkit is Module 11: the testing techniques the whole repo relies on,
// made the lesson. The subjects are a pure function (Normalize) exercised by
// table-driven, property-based, and fuzz tests, and a time-dependent Limiter
// tested deterministically with virtual time (testing/synctest).
package testkit

import "strings"

// region:normalize:start

// Normalize collapses all whitespace runs to a single space, trims the ends, and
// lowercases. It is idempotent: Normalize(Normalize(s)) == Normalize(s) for all
// s -- the invariant the property and fuzz tests check.
func Normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// region:normalize:end
