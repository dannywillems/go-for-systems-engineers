package errs

import (
	"bytes"
	"testing"
)

func TestCopyProduct(t *testing.T) {
	data := []byte("hi!")
	if got := CopyCorrect(NewEOFReader(data)); !bytes.Equal(got, data) {
		t.Errorf("CopyCorrect = %q, want %q", got, data)
	}
	// The bug is real: reading err-first drops the bytes that came with EOF.
	if got := CopyBuggy(NewEOFReader(data)); bytes.Equal(got, data) {
		t.Errorf("CopyBuggy = %q, expected to have dropped the EOF chunk", got)
	}
}

func TestTypedNil(t *testing.T) {
	// Confirm the trap at the type level: a nil *MyError as error is non-nil.
	box := []error{BuggyTypedNil()}
	if box[0] == nil {
		t.Error("BuggyTypedNil should be a non-nil (typed nil) error")
	}
	if CorrectNil() != nil {
		t.Error("CorrectNil should be a nil error")
	}
}
