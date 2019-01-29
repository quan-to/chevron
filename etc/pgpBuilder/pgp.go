package pgpBuilder

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/etc/krmBuilder"
	"github.com/quan-to/remote-signer/keyBackend"
	"github.com/quan-to/remote-signer/pgp"
	"github.com/quan-to/remote-signer/vaultManager"
)

func MakePGP() etc.PGPInterface {
	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(remote_signer.KeyPrefix)
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)
	}

	return pgp.MakePGPManagerWithKRM(kb, krmBuilder.MakeKRM())
}
