package agent

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/database/rql"
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

func MakeDatabaseHandler(logger slog.Instance) (DatabaseHandler, error) {
	if config.RethinkAuthManager {
		logger.Info("RethinkDB Database Enabled. Creating handler")
		rdb := rql.MakeRethinkDBDriver(logger)
		logger.Info("Connecting to RethinkDB")
		err := rdb.Connect(config.RethinkDBHost, config.RethinkDBUsername, config.RethinkDBPassword, config.DatabaseName, config.RethinkDBPort, config.RethinkDBPoolSize)
		if err != nil {
			return nil, err
		}
		logger.Info("Initializing database")
		err = rdb.InitDatabase()
		if err != nil {
			return nil, err
		}
		logger.Info("RethinkDB Handler done!")
		return rdb, nil
	}
	logger.Warn("No database handler specified. Using memory database")

	mdb := memory.MakeMemoryDBDriver(logger)

	return mdb, nil
}

// MakeTokenManager creates an instance of token manager. If Rethink is enabled returns an RethinkTokenManager, if not a MemoryTokenManager
func MakeTokenManager(logger slog.Instance, dbHandler DatabaseHandler) interfaces.TokenManager {
	if dbHandler != nil {
		return MakeDatabaseTokenManager(logger, dbHandler)
	}

	return MakeMemoryTokenManager(logger)
}

// MakeAuthManager creates an instance of auth manager. If Rethink is enabled returns an RethinkAuthManager, if not a JSONAuthManager
func MakeAuthManager(logger slog.Instance, dbHandler DatabaseHandler) interfaces.AuthManager {
	if dbHandler != nil {
		return NewDatabaseAuthManager(logger, dbHandler)
	}

	return MakeJSONAuthManager(logger)
}
