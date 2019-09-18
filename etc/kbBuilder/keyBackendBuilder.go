// +build !js,!wasm

package kbBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
)

// BuildKeyBackend returns a new instance of KeyBackend defined by environment variables VaultStorage, KeyPrefix, PrivateKeyFolder
func BuildKeyBackend(log slog.Instance) keyBackend.Backend {
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(log, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	return kb
}
