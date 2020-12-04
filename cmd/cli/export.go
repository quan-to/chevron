package main

import (
	"fmt"
	"github.com/quan-to/chevron/internal/etc/magicbuilder"
	"github.com/quan-to/chevron/pkg/models"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ExportKey exports the specified public / secret key
func ExportKey(name, password string, secret bool) {
	var err error
	pgpMan := magicbuilder.MakePGP(nil, mem)
	pgpMan.LoadKeys(ctx)

	// First Search the key
	keys := pgpMan.GetLoadedKeys()

	var kInfo *models.KeyInfo

	for _, v := range keys {
		if strings.Contains(v.FingerPrint, strings.ToUpper(name)) || strings.Contains(strings.ToLower(v.Identifier), strings.ToLower(name)) {
			// Thats our key!
			kInfo = &v
			break
		}
	}

	if kInfo == nil {
		panic(fmt.Sprintf("Cannot find key with \"%s\"\n", name))
	}

	if !kInfo.ContainsPrivateKey && secret {
		panic(fmt.Sprintf("The key identified with \"%s\" does not have a private key (found fingerPrint: %s)\n", name, kInfo.FingerPrint))
	}

	k := ""

	if secret {
		if password == "" {
			_, _ = fmt.Fprint(os.Stderr, "Please enter the password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic(fmt.Sprintf("Error reading password: %s", err))
			}
			password = string(bytePassword)
			fmt.Println("")
		}

		k, err = pgpMan.GetPrivateKeyASCII(ctx, kInfo.FingerPrint, password)
		if err != nil {
			if strings.Contains(err.Error(), "checksum failure") {
				panic("Invalid key password")
			}
			panic(err)
		}
	} else {
		k, err = pgpMan.GetPublicKeyASCII(ctx, kInfo.FingerPrint)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(strings.Trim(k, "\n"))
}
