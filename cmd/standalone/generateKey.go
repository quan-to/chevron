package main

import (
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc/magicBuilder"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"syscall"
)

func GenerateFlow(password, output, identifier string, bits int) {
	pgpMan := magicBuilder.MakePGP()
	if password == "" {
		_, _ = fmt.Fprint(os.Stderr, "Please enter the password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(fmt.Sprintf("Error reading password: %s", err))
		}
		password = string(bytePassword)
		fmt.Println("")
	}

	_, _ = fmt.Fprintln(os.Stderr, "Generating key. This might take a while...")

	key, err := pgpMan.GeneratePGPKey(identifier, password, bits)

	if err != nil {
		panic(fmt.Sprintf("Error creating key: %s\n", err))
	}

	fingerPrint, _ := remote_signer.GetFingerPrintFromKey(key)

	_, _ = fmt.Fprintf(os.Stderr, "Generated key fingerprint: %s\n", fingerPrint)

	if output == "-" {
		fmt.Println(key)
	} else if output == "+" {
		err := pgpMan.SaveKey(fingerPrint, key, nil)
		if err != nil {
			panic(fmt.Sprintf("Error saving key to default backend: %s", err))
		}
		_, _ = fmt.Fprintf(os.Stderr, "Key %s saved to default backend", fingerPrint)
	} else {
		err = ioutil.WriteFile(output, []byte(key), 0770)
		if err != nil {
			panic(fmt.Sprintf("Error saving file %s: %s\n", output, err))
		}
		_, _ = fmt.Fprintf(os.Stderr, "Key saved to %s", output)
	}
}
