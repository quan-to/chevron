package main

import (
	"bufio"
	"fmt"
	"github.com/quan-to/chevron/internal/etc/magicbuilder"
	"github.com/quan-to/chevron/internal/tools"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ImportKey(filename, keyPassword string, keyPasswordFd int) {
	var data []byte
	var err error
	pgpMan := magicbuilder.MakePGP(nil, mem)
	pgpMan.LoadKeys(ctx)

	if filename == "-" {
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
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("Error loading file %s: %s\n", filename, err))
		}
	}

	n, err := pgpMan.LoadKey(ctx, string(data))

	if err != nil {
		panic(fmt.Sprintf("Error loading file %s: %s\n", filename, err))
	}

	if keyPasswordFd != -1 {
		// Load from FD
		_, _ = fmt.Fprintf(os.Stderr, "Reading key password from FD %d\n", keyPasswordFd)

		f := os.NewFile(uintptr(keyPasswordFd), "kp")
		d, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		keyPassword = strings.Trim(string(d), "\n\r")
	}

	fps, _ := tools.GetFingerPrintsFromKey(string(data))

	for _, v := range fps {
		if keyPassword != "" {
			// Try get a private key
			private, err := pgpMan.GetPrivateKeyASCII(ctx, v, keyPassword)

			if err == nil {
				err = pgpMan.SaveKey(v, private, keyPassword)
				if err == nil {
					_, _ = fmt.Fprintf(os.Stderr, "Imported private key %s\n", v)
					continue
				}
			}

			_, _ = fmt.Fprintf(os.Stderr, "Cannot import private key %s with supplied password: %s\n", v, err)
		} else if n > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "File contains private key, but no password supplied. Skipping saving private key.\n")
		}
		// Try public if no private
		public, err := pgpMan.GetPublicKeyASCII(ctx, v)
		if err == nil {
			err = pgpMan.SaveKey(v, public, nil)
			if err == nil {
				_, _ = fmt.Fprintf(os.Stderr, "Imported public key %s\n", v)
				continue
			}
		}

		_, _ = fmt.Fprintf(os.Stderr, "Cannot import public key %s: %s\n", v, err)
	}
}
