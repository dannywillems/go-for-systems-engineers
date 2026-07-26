// Package other DOES NOT COMPILE. It sits at rejectinternal/other, OUTSIDE the
// rejectinternal/lib subtree, so importing rejectinternal/lib/internal/store is
// forbidden by the internal/ rule and the compiler rejects it:
//
//	use of internal package rejectinternal/lib/internal/store not allowed
package other

import "rejectinternal/lib/internal/store"

func Leak() string { return store.Secret() }
