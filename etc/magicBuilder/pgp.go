// +build !js,!wasm

package magicBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/keymagic"
	"github.com/quan-to/chevron/vaultManager"
)

func MakePGP() etc.PGPInterface {
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	return keymagic.MakePGPManagerWithKRM(kb, keymagic.MakeKeyRingManager())
}
