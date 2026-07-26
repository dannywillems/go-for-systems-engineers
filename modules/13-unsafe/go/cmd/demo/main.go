// Command demo prints deterministic layout and zero-copy conversion results.
package main

import (
	"fmt"

	u "github.com/dannywillems/go-for-systems-engineers/modules/13-unsafe/go"
)

func main() {
	size, off, align := u.Layout()
	fmt.Printf("Record: size=%d bytes, ID offset=%d, align=%d\n", size, off, align)

	b := []byte("hello")
	s := u.BytesToString(b)
	fmt.Printf("BytesToString(%q) = %q (no copy)\n", b, s)

	b2 := u.StringToBytes("world")
	fmt.Printf("StringToBytes(\"world\") -> %v (%q)\n", b2, string(b2))
}
