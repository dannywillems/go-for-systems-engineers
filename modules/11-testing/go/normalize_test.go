package testkit

import (
	"strings"
	"testing"
	"testing/quick"
)

// Table-driven tests with subtests: the Go default.
func TestNormalizeTable(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"trim", "  hi  ", "hi"},
		{"collapse", "a\t\n  b", "a b"},
		{"lower", "MiXeD", "mixed"},
		{"empty", "   ", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Normalize(c.in); got != c.want {
				t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// region:property:start

// Property-based test with the stdlib testing/quick: for ARBITRARY strings,
// Normalize is idempotent. quick.Check generates random inputs and fails with a
// counterexample if the property is ever violated.
func TestNormalizeIdempotent(t *testing.T) {
	idempotent := func(s string) bool {
		return Normalize(Normalize(s)) == Normalize(s)
	}
	if err := quick.Check(idempotent, nil); err != nil {
		t.Fatal(err)
	}
}

// region:property:end

// region:fuzz:start

// Native fuzz test: `go test -fuzz=FuzzNormalize` mutates the seed corpus to
// hunt for an input that breaks an invariant; plain `go test` runs it against
// the seeds. The invariants: idempotent, trimmed, and no double spaces -- for
// ANY input, including invalid UTF-8.
func FuzzNormalize(f *testing.F) {
	f.Add("  Hello   World  ")
	f.Add("")
	f.Add("\t\nMiXeD\r\n")
	f.Fuzz(func(t *testing.T, s string) {
		got := Normalize(s)
		if Normalize(got) != got {
			t.Errorf("not idempotent on %q", s)
		}
		if strings.TrimSpace(got) != got {
			t.Errorf("leading/trailing space in %q", got)
		}
		if strings.Contains(got, "  ") {
			t.Errorf("double space in %q", got)
		}
	})
}

// region:fuzz:end
