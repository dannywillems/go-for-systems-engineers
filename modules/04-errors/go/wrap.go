package errs

import (
	"errors"
	"fmt"
)

// ErrNotFound is a sentinel matched with errors.Is up the wrap chain.
var ErrNotFound = errors.New("not found")

// MyError is a concrete error type used both for errors.As and for the typed-nil
// trap below.
type MyError struct{ Key string }

func (e *MyError) Error() string { return "no such key: " + e.Key }

// Lookup wraps the sentinel with %w so callers can errors.Is it, and also
// carries structured data via a typed error for errors.As.
func Lookup(key string) error {
	if key == "" {
		return fmt.Errorf("lookup %q: %w", key, ErrNotFound)
	}
	return fmt.Errorf("lookup %q: %w", key, &MyError{Key: key})
}

// Combined joins two errors into a DAG (errors.Join); errors.Is walks all
// branches.
func Combined() error {
	return errors.Join(fmt.Errorf("disk: %w", ErrNotFound), errors.New("net down"))
}

// region:typednil:start

// BuggyTypedNil returns a *MyError typed nil. Assigning a nil *MyError to the
// error interface yields a NON-nil error (non-nil itab, nil data), so a caller's
// `if err != nil` fires even though "there is no error". The value is returned
// through the interface at the boundary, which is exactly where the trap hides.
func BuggyTypedNil() error {
	var e *MyError // nil
	return e       // becomes a non-nil error
}

// CorrectNil returns the interface nil explicitly.
func CorrectNil() error { return nil }

// region:typednil:end
