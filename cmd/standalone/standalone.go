package main

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"
	"time"
)

func GenerateFlow(password, output, identifier string, bits int) {
	pgpMan := remote_signer.MakePGPManager()
	if password == "" {
		_, _ = fmt.Fprint(os.Stderr, "Please enter the password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(fmt.Sprintf("Error reading password: %s", err))
		}
		password = string(bytePassword)
	}

	key, err := pgpMan.GeneratePGPKey(identifier, password, bits)

	if err != nil {
		panic(fmt.Sprintf("Error creating key: %s\n", err))
	}

	if output == "-" {
		fmt.Println(key)
	} else {
		err = ioutil.WriteFile(output, []byte(key), 0770)
		if err != nil {
			panic(fmt.Sprintf("Error saving file %s: %s\n", output, err))
		}
	}
}

func BenchmarkGeneration(runs, bits int) {
	pgpMan := remote_signer.MakePGPManager()

	fmt.Printf("Benchmarking GPG Key Generation with %d bits and %d runs.\n", bits, runs)
	fmt.Printf("Running on %s-%s\n", runtime.GOOS, runtime.GOARCH)

	startTime := time.Now()
	for i := 0; i < runs; i++ {
		_, _ = pgpMan.GeneratePGPKey("", "", bits)
	}
	delta := time.Since(startTime)
	keyTime := delta.Seconds()  / float64(runs)


	fmt.Printf("Took average of %f seconds to generate a %d bits key.\n", keyTime, bits)
}

func main() {

	// region Generate
	gen := kingpin.Command("gen", "Generate GPG Key")
	genBits := gen.Flag("bits", "Number of bits").Default("2048").Uint16()
	genIdentifier := gen.Flag("id", "Key Identifier").Default("").String()
	genOutput := gen.Flag("output", "Filename of the output ( use - for stdout )").Default("-").String()
	genPassword := gen.Flag("password", "Key Password (if not provided, it will be prompted)").Default("").String()
	// endregion
	// region Benchmark Generate

	// endregion
	benchGen := kingpin.Command("benchgen", "Benchmark Key Generation")
	benchGenBits := benchGen.Flag("bits", "Number of bits").Default("2048").Uint16()
	benchGenRuns := benchGen.Flag("runs", "Number of runs").Default("20").Int()

	switch kingpin.Parse() {
	case "gen": GenerateFlow(*genPassword, *genOutput, *genIdentifier, int(*genBits))
	case "benchgen": BenchmarkGeneration(*benchGenRuns, int(*benchGenBits))
	}
}
