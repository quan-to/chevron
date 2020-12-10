package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/quan-to/chevron/internal/etc/magicbuilder"
)

// BenchmarkGeneration benchmarks the key generation
func BenchmarkGeneration(runs, bits int) {
	pgpMan := magicbuilder.MakePGP(nil, mem)

	fmt.Printf("Benchmarking GPG Key Generation with %d bits and %d runs.\n", bits, runs)
	fmt.Printf("Running on %s-%s\n", runtime.GOOS, runtime.GOARCH)

	startTime := time.Now()
	for i := 0; i < runs; i++ {
		_, _ = pgpMan.GeneratePGPKey(ctx, "", "", bits)
	}
	delta := time.Since(startTime)
	keyTime := delta.Seconds() / float64(runs)

	fmt.Printf("Took average of %f seconds to generate a %d bits key.\n", keyTime, bits)
}
