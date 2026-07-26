// Command demo runs the error demonstrations with deterministic output.
package main

import (
	"errors"
	"fmt"

	errs "github.com/dannywillems/go-for-systems-engineers/modules/04-errors/go"
)

func main() {
	// (T, error) is a product: the last Read returns bytes AND io.EOF together.
	data := []byte("hi!")
	fmt.Printf("CopyBuggy   got %q  (dropped the EOF chunk)\n",
		errs.CopyBuggy(errs.NewEOFReader(data)))
	fmt.Printf("CopyCorrect got %q\n", errs.CopyCorrect(errs.NewEOFReader(data)))

	// Typed nil: a nil *MyError returned as error is a non-nil interface. The
	// value goes through a slice so the comparison is a real runtime check.
	box := []error{errs.BuggyTypedNil()}
	fmt.Printf("BuggyTypedNil() == nil: %v  (phantom error)\n", box[0] == nil)
	fmt.Printf("CorrectNil()    == nil: %v\n", errs.CorrectNil() == nil)

	// Wrapping and the DAG unwrap: %w, errors.Is/As/Join.
	fmt.Printf("errors.Is(Lookup(\"\"), ErrNotFound): %v\n",
		errors.Is(errs.Lookup(""), errs.ErrNotFound))
	var me *errs.MyError
	fmt.Printf("errors.As(Lookup(\"k\"), *MyError):    %v\n",
		errors.As(errs.Lookup("k"), &me))
	fmt.Printf("errors.Is(Join(...), ErrNotFound):   %v\n",
		errors.Is(errs.Combined(), errs.ErrNotFound))

	// Result monad: lawful, but written inside-out because Map/AndThen cannot be
	// methods. Compare Rust's `Ok(3).map(|x| x*2).and_then(|x| Ok(x+1))`.
	r := errs.AndThen(
		errs.Map(errs.Ok(3), func(x int) int { return x * 2 }),
		func(x int) errs.Result[int] { return errs.Ok(x + 1) },
	)
	v, _ := r.Unwrap()
	fmt.Printf("AndThen(Map(Ok(3), *2), +1) = %d  (no method chaining)\n", v)
}
