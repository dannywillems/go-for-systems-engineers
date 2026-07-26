package solutions

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

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

func TestDivide(t *testing.T) {
	if _, err := Divide(1, 0); !errors.Is(err, ErrDivByZero) {
		t.Fatal("want ErrDivByZero")
	}
	if got, _ := Divide(10, 2); got != 5 {
		t.Fatal("want 5")
	}
}

func TestDrain(t *testing.T) {
	r := &eofReader{chunks: [][]byte{[]byte("hi"), []byte("!")}}
	got, err := Drain(r)
	if err != nil || !bytes.Equal(got, []byte("hi!")) {
		t.Fatalf("Drain = %q, %v", got, err)
	}
}
