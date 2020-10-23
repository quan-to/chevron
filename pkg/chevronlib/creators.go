package chevronlib

import (
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
)

// MakeSaveToDiskBackend creates an instance of a StorageBackend that
// saves the keys in the specified folder with the specified prefix
// log instance can be nil
func MakeSaveToDiskBackend(log slog.Instance, keysFolder, prefix string) interfaces.StorageBackend {
	return keybackend.MakeSaveToDiskBackend(log, keysFolder, prefix)
}

// MakeKeyRingManager creates a new instance of Key Ring Manager
// log instance can be nil
func MakeKeyRingManager(log slog.Instance) interfaces.KeyRingManager {
	return keymagic.MakeKeyRingManager(log)
}

// MakePGPManager creates a new instance of PGP Operations Manager
// log instance can be nil
func MakePGPManager(log slog.Instance, storage interfaces.StorageBackend, keyRingManager interfaces.KeyRingManager) interfaces.PGPManager {
	return keymagic.MakePGPManager(log, storage, keyRingManager)
}
