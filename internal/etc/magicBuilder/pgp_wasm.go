package magicBuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keyBackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
)

func MakePGP() interfaces.PGPInterface {
	kb := keyBackend.MakeSaveToDiskBackend(nil, config.PrivateKeyFolder, config.KeyPrefix)

	return keymagic.MakePGPManagerWithKRM(nil, kb, keymagic.MakeKeyRingManager())
}

func MakeVoidPGP() interfaces.PGPInterface {
	return keymagic.MakePGPManagerWithKRM(nil, keyBackend.MakeVoidBackend(), keymagic.MakeKeyRingManager())
}
