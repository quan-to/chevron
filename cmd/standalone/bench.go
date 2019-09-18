package main

import (
	"fmt"
	"github.com/quan-to/chevron/etc/magicBuilder"
	"runtime"
	"time"
)

// BenchmarkGeneration benchmarks the key generation
func BenchmarkGeneration(runs, bits int) {
	pgpMan := magicBuilder.MakePGP(nil)

	fmt.Printf("Benchmarking GPG Key Generation with %d bits and %d runs.\n", bits, runs)
	fmt.Printf("Running on %s-%s\n", runtime.GOOS, runtime.GOARCH)

	startTime := time.Now()
	for i := 0; i < runs; i++ {
		_, _ = pgpMan.GeneratePGPKey("", "", bits)
	}
	delta := time.Since(startTime)
	keyTime := delta.Seconds() / float64(runs)

	fmt.Printf("Took average of %f seconds to generate a %d bits key.\n", keyTime, bits)
}
