package magicBuilder

import (
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/keyBackend"
	"github.com/quan-to/remote-signer/keymagic"
)

func MakePGP() etc.PGPInterface {
	kb := keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)

	return keymagic.MakePGPManagerWithKRM(kb, keymagic.MakeKeyRingManager())
}
