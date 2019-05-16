package main

import (
	"fmt"
	"github.com/quan-to/chevron/etc/magicBuilder"
)

// ListKeys list the Public / Private keys stored in the default backend
func ListKeys() {
	pgpMan := magicBuilder.MakePGP()
	pgpMan.LoadKeys()

	keys := pgpMan.GetLoadedKeys()
	fmt.Printf("There is %d private keys stored.\n", len(keys))
	if len(keys) > 0 {
		fmt.Printf("%-18s %4s %12s     %-50s\n", "FingerPrint", "Bits", "Private", "Identifier")
		for _, key := range keys {
			fmt.Printf("%-18s %4d %12v     %-50s\n", key.FingerPrint, key.Bits, key.ContainsPrivateKey, key.Identifier)
		}
	}
}
