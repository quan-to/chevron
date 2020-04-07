// +build !js,!wasm

package magicBuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keyBackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// MakePGP creates a new PGPInterface using environment variables VaultStorage, KeyPrefix, PrivateKeyFolder
func MakePGP(log slog.Instance) interfaces.PGPInterface {
	var kb interfaces.Backend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, config.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(log, config.PrivateKeyFolder, config.KeyPrefix)
	}

	return keymagic.MakePGPManagerWithKRM(log, kb, keymagic.MakeKeyRingManager(log))
}

// MakeVoidPGP creates a PGPInterface that does not store anything anywhere
func MakeVoidPGP(log slog.Instance) interfaces.PGPInterface {
	return keymagic.MakePGPManagerWithKRM(log, keyBackend.MakeVoidBackend(), keymagic.MakeKeyRingManager(log))
}
