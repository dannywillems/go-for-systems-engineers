// Package errs demonstrates that Go's (T, error) is a PRODUCT, not a coproduct.
//
// A Rust engineer reads `func F() (T, error)` as `fn f() -> Result<T, E>`, i.e.
// a sum: either a value or an error, never both. In Go it is a pair: BOTH
// components are always present, and the contract for which parts are valid is
// per-function convention, not a type. The canonical trap is io.Reader, whose
// Read may return n > 0 bytes AND a non-nil error (including io.EOF) in the SAME
// call. Treating it as a sum and bailing on err != nil discards the last bytes.
package errs

import "io"

// region:reader:start

// eofReader yields its data in chunks and reports io.EOF ON THE SAME Read that
// returns the FINAL chunk. This is explicitly permitted by the io.Reader
// contract, and many real readers do it. It is the shape that breaks sum-type
// intuition: the last Read is (n > 0, io.EOF) together.
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
		return n, io.EOF // last chunk: n > 0 AND err != nil, together
	}
	return n, nil
}

// region:reader:end

// CopyBuggy is the mistake a sum-type reflex produces: it checks err first and
// returns before consuming the n bytes that came WITH the EOF, silently losing
// data. It compiles and looks correct.
func CopyBuggy(r io.Reader) []byte {
	var out []byte
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return out // BUG: the n bytes read on this call are dropped
		}
		out = append(out, buf[:n]...)
	}
}

// CopyCorrect honors the product: it consumes the n bytes FIRST, then inspects
// the error. This is what io.Copy and every correct Read loop do.
func CopyCorrect(r io.Reader) []byte {
	var out []byte
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		out = append(out, buf[:n]...) // use n regardless of err
		if err != nil {
			return out
		}
	}
}

// NewEOFReader returns a reader that emits data as (head, last-byte) so the last
// Read carries both the final byte and io.EOF.
func NewEOFReader(data []byte) io.Reader {
	if len(data) <= 1 {
		return &eofReader{chunks: [][]byte{data}}
	}
	return &eofReader{chunks: [][]byte{data[:len(data)-1], data[len(data)-1:]}}
}
