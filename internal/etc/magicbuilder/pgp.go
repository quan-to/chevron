// +build !js,!wasm

package magicbuilder

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
)

type DatabaseHandler interface {
	GetUser(username string) (um *models.User, err error)
	AddUserToken(ut models.UserToken) (string, error)
	RemoveUserToken(token string) (err error)
	GetUserToken(token string) (ut *models.UserToken, err error)
	InvalidateUserTokens() (int, error)
	AddUser(um models.User) (string, error)
	UpdateUser(um models.User) error
	AddGPGKey(key models.GPGKey) (string, bool, error)
	FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error)
	HealthCheck() error
}

// MakePGP creates a new PGPManager using environment variables VaultStorage, KeyPrefix, PrivateKeyFolder
func MakePGP(log slog.Instance, dbHandler DatabaseHandler) interfaces.PGPManager {
	var kb interfaces.StorageBackend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, config.KeyPrefix)
	} else {
		kb = keybackend.MakeSaveToDiskBackend(log, config.PrivateKeyFolder, config.KeyPrefix)
	}

	return keymagic.MakePGPManager(log, kb, keymagic.MakeKeyRingManager(log, dbHandler))
}

// MakeVoidPGP creates a PGPManager that does not store anything anywhere
func MakeVoidPGP(log slog.Instance, dbHandler DatabaseHandler) interfaces.PGPManager {
	return keymagic.MakePGPManager(log, keybackend.MakeVoidBackend(), keymagic.MakeKeyRingManager(log, dbHandler))
}
