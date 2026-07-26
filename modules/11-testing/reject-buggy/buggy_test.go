package buggy

import (
	"strings"
	"testing"
)

// The same invariants as the real module. Run against the seed corpus, they
// deterministically catch the planted bug on the first seed with a double space.
func FuzzNormalize(f *testing.F) {
	f.Add("  Hello   World  ")
	f.Fuzz(func(t *testing.T, s string) {
		got := Normalize(s)
		if strings.Contains(got, "  ") {
			t.Errorf("invariant violated: double space in %q", got)
		}
		if strings.TrimSpace(got) != got {
			t.Errorf("invariant violated: leading/trailing space in %q", got)
		}
	})
}
