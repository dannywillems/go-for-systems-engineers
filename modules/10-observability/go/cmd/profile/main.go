// Command profile runs the CPU-bound hot loop under a CPU profile and writes
// cpu.prof, which `go tool pprof` analyzes offline. Used by
// scripts/obs-bench.sh to extract the hottest function into measured.txt.
package main

import (
	"log"
	"os"

	obs "observability"
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal(cerr)
		}
	}()

	var sink uint64
	if err := obs.CPUProfile(f, func() {
		for i := 0; i < 2000; i++ {
			sink += obs.HotLoop(1_000_000)
		}
	}); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote cpu.prof (sink=%d)", sink)
}
