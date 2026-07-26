package errs

import (
	"errors"
	"testing"
	"testing/quick"
)

// The three monad laws for Result[int], checked with testing/quick over random
// ints. That they PASS proves the type is a lawful monad; the README shows why
// it is nonetheless unusable (no generic methods -> no chaining).

func TestLeftIdentity(t *testing.T) {
	// AndThen(Ok(a), f) == f(a)
	f := func(a int) Result[int] { return Ok(a + 1) }
	prop := func(a int) bool {
		return equalIntResult(AndThen(Ok(a), f), f(a))
	}
	if err := quick.Check(prop, nil); err != nil {
		t.Error(err)
	}
}

func TestRightIdentity(t *testing.T) {
	// AndThen(m, Ok) == m
	prop := func(a int, fail bool) bool {
		m := Ok(a)
		if fail {
			m = Err[int](errors.New("e"))
		}
		return equalIntResult(AndThen(m, Ok), m)
	}
	if err := quick.Check(prop, nil); err != nil {
		t.Error(err)
	}
}

func TestAssociativity(t *testing.T) {
	// AndThen(AndThen(m, f), g) == AndThen(m, x -> AndThen(f(x), g))
	f := func(x int) Result[int] { return Ok(x * 2) }
	g := func(x int) Result[int] { return Ok(x + 3) }
	prop := func(a int, fail bool) bool {
		m := Ok(a)
		if fail {
			m = Err[int](errors.New("e"))
		}
		lhs := AndThen(AndThen(m, f), g)
		rhs := AndThen(m, func(x int) Result[int] { return AndThen(f(x), g) })
		return equalIntResult(lhs, rhs)
	}
	if err := quick.Check(prop, nil); err != nil {
		t.Error(err)
	}
}
