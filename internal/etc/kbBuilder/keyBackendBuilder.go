// +build !js,!wasm

package kbBuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keyBackend"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// BuildKeyBackend returns a new instance of KeyBackend defined by environment variables VaultStorage, KeyPrefix, PrivateKeyFolder
func BuildKeyBackend(log slog.Instance) interfaces.Backend {
	var kb interfaces.Backend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, config.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(log, config.PrivateKeyFolder, config.KeyPrefix)
	}

	return kb
}
