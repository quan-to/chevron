// +build !js,!wasm

package kbBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/vaultManager"
)

func BuildKeyBackend() keyBackend.Backend {
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	return kb
}
