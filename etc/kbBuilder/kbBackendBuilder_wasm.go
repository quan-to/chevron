package kbBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/keyBackend"
)

func BuildKeyBackend() keyBackend.Backend {
	return keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
}
