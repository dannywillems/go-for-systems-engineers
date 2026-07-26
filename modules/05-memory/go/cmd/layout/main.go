// Command layout prints the in-memory SIZE of Go's core types, in machine words
// (8 bytes on 64-bit). The sizes are the falsifiable evidence for the header
// diagrams in the README: a slice is 3 words (ptr,len,cap), a string and an
// interface are 2, a map/chan/pointer are 1. unsafe.Sizeof is a compile-time
// constant, so this is deterministic.
package main

import (
	"fmt"
	"unsafe"
)

type oneMethod interface{ M() }

func main() {
	var (
		sl  []int
		str string
		emp any
		err error
		im  oneMethod
		mp  map[int]int
		ch  chan int
		ptr *int
		b   bool
	)
	const word = 8
	row := func(name string, size uintptr) {
		fmt.Printf("%-22s %2d bytes  (%d word%s)\n",
			name, size, size/word, map[bool]string{true: "s", false: ""}[size/word != 1])
	}
	row("[]int (slice)", unsafe.Sizeof(sl))
	row("string", unsafe.Sizeof(str))
	row("any / interface{}", unsafe.Sizeof(emp))
	row("error (interface)", unsafe.Sizeof(err))
	row("interface{M()}", unsafe.Sizeof(im))
	row("map[int]int", unsafe.Sizeof(mp))
	row("chan int", unsafe.Sizeof(ch))
	row("*int (pointer)", unsafe.Sizeof(ptr))
	fmt.Printf("%-22s %2d byte\n", "bool", unsafe.Sizeof(b))
}
