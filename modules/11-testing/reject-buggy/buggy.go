// Package buggy has a PLANTED BUG: this Normalize only lowercases, forgetting to
// collapse whitespace and trim. The fuzz test's invariants catch it -- the point
// of property/fuzz testing is that a machine finds the violation you did not
// think to write a unit test for. `go test` here FAILS (that is the demo); the
// captured output is the failure the tooling reports.
package buggy

import "strings"

func Normalize(s string) string { return strings.ToLower(s) } // BUG: no collapse/trim
