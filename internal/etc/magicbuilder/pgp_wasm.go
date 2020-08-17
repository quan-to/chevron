package magicbuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// MakePGP creates a new PGPManager using environment variables KeyPrefix, PrivateKeyFolder
func MakePGP(log slog.Instance) interfaces.PGPManager {
	kb := keybackend.MakeSaveToDiskBackend(log, config.PrivateKeyFolder, config.KeyPrefix)

	return keymagic.MakePGPManager(log, kb, keymagic.MakeKeyRingManager(log))
}

// MakeVoidPGP creates a PGPManager that does not store anything anywhere
func MakeVoidPGP(log slog.Instance) interfaces.PGPManager {
	return keymagic.MakePGPManager(log, keybackend.MakeVoidBackend(), keymagic.MakeKeyRingManager(log))
}
