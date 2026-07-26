// Package store is under an internal/ directory, so the compiler only lets code
// rooted at its parent (rejectinternal/lib) import it. It is a normal package
// with an exported function; the boundary is purely the internal/ rule.
package store

// Secret is exported, but reachable only from within rejectinternal/lib.
func Secret() string { return "hunter2" }
