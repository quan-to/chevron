package chevronlib

import (
	"fmt"
	"os"
)

func GenerateKey(password, identifier string, bits int) (result string, err error) {
	if password == "" {
		err = fmt.Errorf("no password supplied")
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, "Generating key. This might take a while...")

	result, err = pgpBackend.GeneratePGPKey(identifier, password, bits)

	return
}
