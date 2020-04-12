package magicbuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
)

// MakePGP creates a new PGPManager using environment variables KeyPrefix, PrivateKeyFolder
func MakePGP() interfaces.PGPManager {
	kb := keybackend.MakeSaveToDiskBackend(nil, config.PrivateKeyFolder, config.KeyPrefix)

	return keymagic.MakePGPManagerWithKRM(nil, kb, keymagic.MakeKeyRingManager())
}

// MakeVoidPGP creates a PGPManager that does not store anything anywhere
func MakeVoidPGP() interfaces.PGPManager {
	return keymagic.MakePGPManagerWithKRM(nil, keybackend.MakeVoidBackend(), keymagic.MakeKeyRingManager())
}
