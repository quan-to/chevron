package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron/internal/etc/magicBuilder"
	"io"
	"io/ioutil"
	"os"
)

func Decrypt(input, output string) {
	var err error
	var data []byte

	pgpMan := magicBuilder.MakePGP(nil)
	pgpMan.LoadKeys(ctx)

	if input == "-" {
		// Read from stdin
		_, _ = fmt.Fprintf(os.Stderr, "Reading from stdin:\n")
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

	d, err := pgpMan.Decrypt(ctx, string(data), false)

	if err != nil {
		panic(err)
	}

	var out *bufio.Writer
	var f *os.File

	if output == "-" {
		// Write to Stdout
		out = bufio.NewWriter(os.Stdout)
		f = os.Stdout
	} else {
		f, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		out = bufio.NewWriter(f)
	}

	data, _ = base64.StdEncoding.DecodeString(d.Base64Data)

	_, err = out.Write(data)

	if err != nil {
		panic(err)
	}

	_ = out.Flush()
	_ = f.Close()
}
