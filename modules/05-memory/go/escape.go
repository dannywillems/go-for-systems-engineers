package mem

// region:escape:start

// Sink is an escape target for the boxing case below. It is exported so
// staticcheck does not flag it as write-only/unused.
var Sink any

// NoEscape: the array is used only locally, so escape analysis keeps it on the
// STACK (no allocation), like Rust's default.
func NoEscape() int {
	v := [4]int{1, 2, 3, 4}
	return v[0] + v[3]
}

// EscapesReturn: returning a pointer to a local forces the local to the HEAP —
// the value must outlive the frame.
func EscapesReturn() *int {
	x := 42
	return &x
}

// EscapesInterface: putting a concrete value into an interface boxes it on the
// HEAP, even though the value itself is tiny. A frequent, surprising source of
// allocation.
func EscapesInterface(n int) {
	Sink = n
}

// region:escape:end

// AliasParent shows a subslice sharing the parent's backing array: appending
// into a subslice with spare capacity OVERWRITES the parent.
func AliasParent() (before, after []int) {
	orig := []int{1, 2, 3, 4, 5}
	before = append([]int(nil), orig...) // a copy for display
	sub := orig[:2]                      // len 2, cap 5 (shares orig)
	_ = append(sub, 99)                  // writes orig[2]!
	after = orig
	return before, after
}
