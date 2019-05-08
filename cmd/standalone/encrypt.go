package main

import (
	"bufio"
	"fmt"
	"github.com/quan-to/chevron/etc/magicBuilder"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func EncryptFile(input, output, recipient string) {
	var err error
	var data []byte
	pgpMan := magicBuilder.MakePGP()
	pgpMan.LoadKeys()

	ent := pgpMan.GetPublicKeyEntity(recipient)

	if ent == nil {
		panic(fmt.Sprintf("Cannot find key \"%s\"\n", recipient))
	}

	filename := input

	if input == "-" {
		// Read from stdin
		fmt.Fprintf(os.Stderr, "Reading from stdin:\n")
		filename = fmt.Sprintf("stdin-%s", time.Now())
		fio := bufio.NewReader(os.Stdin)
		chunk := make([]byte, 4096)
		data = make([]byte, 0)
		for {
			n, err := fio.Read(chunk)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if n > 0 {
				data = append(data, chunk[:n]...)
			}
		}
	} else {
		data, err = ioutil.ReadFile(input)
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(os.Stderr, "Encrypting to %s\n", recipient)

	var out *bufio.Writer
	var f *os.File

	if output == "-" {
		// Write to Stdout
		out = bufio.NewWriter(os.Stdout)
		f = os.Stdout
	} else {
		f, err := os.Create(output)
		if err != nil {

		}
		out = bufio.NewWriter(f)
	}

	var d string

	d, err = pgpMan.Encrypt(filename, recipient, data, false)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Done encrypting to %s\n", recipient)

	_, err = out.WriteString(d)

	if err != nil {
		panic(err)
	}

	out.Flush()
	f.Close()
}
