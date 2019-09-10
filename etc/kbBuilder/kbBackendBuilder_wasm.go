package kbBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/keyBackend"
)

func BuildKeyBackend(log slog.Instance) keyBackend.Backend {
	return keyBackend.MakeSaveToDiskBackend(log, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
}
