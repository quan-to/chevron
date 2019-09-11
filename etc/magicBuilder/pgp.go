// +build !js,!wasm

package magicBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/keymagic"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
)

// MakePGP creates a new PGPInterface using environment variables VaultStorage, KeyPrefix, PrivateKeyFolder
func MakePGP(log slog.Instance) etc.PGPInterface {
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(log, remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	return keymagic.MakePGPManagerWithKRM(log, kb, keymagic.MakeKeyRingManager(log))
}

// MakeVoidPGP creates a PGPInterface that does not store anything anywhere
func MakeVoidPGP(log slog.Instance) etc.PGPInterface {
	return keymagic.MakePGPManagerWithKRM(log, keyBackend.MakeVoidBackend(), keymagic.MakeKeyRingManager(log))
}
