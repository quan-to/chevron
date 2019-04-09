package main

import (
	"github.com/quan-to/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func main() {
	debugMode := kingpin.Flag("debug", "Enable debug mode").Bool()

	// region Generate
	gen := kingpin.Command("gen", "Generate GPG Key")
	genBits := gen.Flag("bits", "Number of bits").Default("4096").Uint16()
	genIdentifier := gen.Flag("id", "Key Identifier").Default("").String()
	genOutput := gen.Flag("output", "Filename of the output ( use - for stdout, use + for default key backend )").Default("+").String()
	genPassword := gen.Flag("password", "Key Password (if not provided, it will be prompted)").Default("").String()
	// endregion
	// region Benchmark Generate

	// endregion
	// region Benchmark Generate
	benchGen := kingpin.Command("benchgen", "Benchmark Key Generation")
	benchGenBits := benchGen.Flag("bits", "Number of bits").Default("2048").Uint16()
	benchGenRuns := benchGen.Flag("runs", "Number of runs").Default("20").Int()
	// endregion

	// region List Keys
	_ = kingpin.Command("list-keys", "List Stored Keys")
	// endregion

	// region Export
	exp := kingpin.Command("export", "Export Key")
	exportSecret := exp.Flag("secret", "Export private key instead of public").Bool()
	exportName := exp.Arg("fingerPrint or email", "Finger Print or email for the key you want to export").String()
	exportPass := exp.Flag("password", "Pass password on command line instead of asking when exporting secret key").String()
	// endregion

	// region Encrypt
	//encrypt := kingpin.Command("encrypt", "Encrypt Data")
	//encryptRecipient := encrypt.Arg("recipient", "Fingerprint of who to encrypt for")
	//encryptInput := encrypt.Arg("input", "Filename of the input")
	//encryptOutput := encrypt.Arg("output", "Filename of the output (use - to stdout)")
	// endregion

	selectedCmd := kingpin.Parse()

	slog.SetDefaultOutput(os.Stderr)
	if !*debugMode {
		slog.SetInfo(false)
		slog.SetDebug(false)
		slog.SetWarning(false)
	} else {
		slog.SetInfo(true)
		slog.SetDebug(true)
		slog.SetWarning(true)
		slog.Info("Debug Mode Enabled!")
	}

	switch selectedCmd {
	case "gen":
		GenerateFlow(*genPassword, *genOutput, *genIdentifier, int(*genBits))
	case "benchgen":
		BenchmarkGeneration(*benchGenRuns, int(*benchGenBits))
	case "list-keys":
		ListKeys()
	case "export":
		ExportKey(*exportName, *exportPass, *exportSecret)
	}
}
