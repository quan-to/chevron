package kbBuilder

import (
	remote_signer "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// BuildKeyBackend returns a new instance of SaveToDisk KeyBackend defined by environment variables KeyPrefix, PrivateKeyFolder
func BuildKeyBackend(log slog.Instance) interfaces.StorageBackend {
	return keybackend.MakeSaveToDiskBackend(log, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
}
