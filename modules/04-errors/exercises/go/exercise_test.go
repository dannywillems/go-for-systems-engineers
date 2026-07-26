package exercises

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestDivide(t *testing.T) {
	if got, err := Divide(10, 2); err != nil || got != 5 {
		t.Fatalf("Divide(10,2) = %d, %v", got, err)
	}
	_, err := Divide(1, 0)
	if !errors.Is(err, ErrDivByZero) {
		t.Fatalf("Divide(1,0) err = %v, want errors.Is ErrDivByZero", err)
	}
}

// eofReader returns "hi" then "!" WITH io.EOF, to exercise the product trap.
type eofReader struct {
	chunks [][]byte
	i      int
}

func (r *eofReader) Read(p []byte) (int, error) {
	if r.i >= len(r.chunks) {
		return 0, io.EOF
	}
	n := copy(p, r.chunks[r.i])
	r.i++
	if r.i == len(r.chunks) {
		return n, io.EOF
	}
	return n, nil
}

func TestDrain(t *testing.T) {
	r := &eofReader{chunks: [][]byte{[]byte("hi"), []byte("!")}}
	got, err := Drain(r)
	if err != nil {
		t.Fatalf("Drain err = %v", err)
	}
	if !bytes.Equal(got, []byte("hi!")) {
		t.Fatalf("Drain = %q, want %q (did you drop the EOF chunk?)", got, "hi!")
	}
}
