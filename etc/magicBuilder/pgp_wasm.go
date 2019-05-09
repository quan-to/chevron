package magicBuilder

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/keymagic"
)

func MakePGP() etc.PGPInterface {
	kb := keyBackend.MakeSaveToDiskBackend(remote_signer.PrivateKeyFolder, remote_signer.KeyPrefix)

	return keymagic.MakePGPManagerWithKRM(kb, keymagic.MakeKeyRingManager())
}

func MakeVoidPGP() etc.PGPInterface {
	return keymagic.MakePGPManagerWithKRM(keyBackend.MakeVoidBackend(), keymagic.MakeKeyRingManager())
}
