package kbBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/keyBackend"
)

// BuildKeyBackend returns a new instance of SaveToDisk KeyBackend defined by environment variables KeyPrefix, PrivateKeyFolder
func BuildKeyBackend(log slog.Instance) keyBackend.Backend {
	return keyBackend.MakeSaveToDiskBackend(log, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
}
