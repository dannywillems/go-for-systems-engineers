// Package solutions is the corrigé for Module 04. Run via `make solutions`.
package solutions

import (
	"errors"
	"fmt"
	"io"
)

var ErrDivByZero = errors.New("division by zero")

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide %d/%d: %w", a, b, ErrDivByZero)
	}
	return a / b, nil
}

func Drain(r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 8)
	for {
		n, err := r.Read(buf)
		out = append(out, buf[:n]...) // consume n before inspecting err
		if err != nil {
			if err == io.EOF {
				return out, nil
			}
			return out, err
		}
	}
}
