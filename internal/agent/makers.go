package agent

import (
	"crypto/tls"

	"github.com/go-redis/redis/v8"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/database/cache"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/database/pg"
	"github.com/quan-to/chevron/pkg/database/rql"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
)

type MigrationHandler interface {
	InitCursor() error
	FinishCursor() error
	NextGPGKey(key *models.GPGKey) bool
	NextUser(user *models.User) bool
	NumGPGKeys() (int, error)
}

type GPGRepository interface {
	AddGPGKey(key models.GPGKey) (string, bool, error)
	// AddGPGKey adds a list GPG Key to the database or update an existing one by fingerprint
	// Same as AddGPGKey but in a single transaction
	AddGPGKeys(keys []models.GPGKey) (ids []string, addeds []bool, err error)
	FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error)
	FetchGPGKeysWithoutSubKeys() (res []models.GPGKey, err error)
	DeleteGPGKey(key models.GPGKey) error
	UpdateGPGKey(key models.GPGKey) (err error)
}

type UserRepository interface {
	GetUser(username string) (um *models.User, err error)
	AddUserToken(ut models.UserToken) (string, error)
	RemoveUserToken(token string) (err error)
	GetUserToken(token string) (ut *models.UserToken, err error)
	InvalidateUserTokens() (int, error)
	AddUser(um models.User) (string, error)
	UpdateUser(um models.User) error
}

type HealthChecker interface {
	HealthCheck() error
}

type DatabaseHandler interface {
	MigrationHandler
	GPGRepository
	UserRepository
	HealthChecker
}

func makeRethinkDBHandler(logger slog.Instance) (*rql.RethinkDBDriver, error) {
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

func makePostgresDBHandler(logger slog.Instance) (*pg.PostgreSQLDBDriver, error) {
	logger.Info("PostgreSQL Database Enabled. Creating handler")
	rdb := pg.MakePostgreSQLDBDriver(logger)
	logger.Info("Initializing database")
	err := rdb.Connect(config.ConnectionString)
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

// MakeDatabaseHandler initializes a Database Access Handler based on the current configuration
func MakeDatabaseHandler(logger slog.Instance) (dbh DatabaseHandler, err error) {
	if config.EnableDatabase {
		switch config.DatabaseDialect {
		case "rethinkdb":
			dbh, err = makeRethinkDBHandler(logger)
			if err != nil {
				return nil, err
			}
		case "postgres":
			dbh, err = makePostgresDBHandler(logger)
			if err != nil {
				return nil, err
			}
		case "memory":
			dbh = memory.MakeMemoryDBDriver(logger)
		default:
			logger.Fatal("Unknown database dialect %q", config.DatabaseDialect)
		}
	}
	if dbh == nil {
		logger.Warn("No database handler enabled. Using memory database")
		dbh = memory.MakeMemoryDBDriver(logger)
	}

	if config.EnableRedis {
		logger.Info("Redis enabled. Wrapping cache layer")
		redisDriver := cache.MakeRedisDriver(dbh, logger)
		var tlsConfig *tls.Config
		if config.RedisTLSEnabled {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		if config.RedisClusterMode {
			cluster := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:     []string{config.RedisHost},
				Username:  config.RedisUser,
				Password:  config.RedisPass,
				TLSConfig: tlsConfig,
			})
			err = redisDriver.Setup(cluster, config.RedisMaxLocalObjects, config.RedisLocalObjectTTL)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:      config.RedisHost,
				Username:  config.RedisUser,
				Password:  config.RedisPass,
				TLSConfig: tlsConfig,
			})
			err = redisDriver.Setup(client, config.RedisMaxLocalObjects, config.RedisLocalObjectTTL)
		}

		if err != nil {
			return nil, err
		}
		dbh = redisDriver
	}

	return dbh, nil
}

// MakeTokenManager creates an instance of token manager. If Rethink is enabled returns an DatabaseTokenManager, if not a MemoryTokenManager
func MakeTokenManager(logger slog.Instance, dbHandler DatabaseHandler) interfaces.TokenManager {
	if dbHandler != nil {
		return MakeDatabaseTokenManager(logger, dbHandler)
	}

	return MakeMemoryTokenManager(logger)
}

// MakeAuthManager creates an instance of auth manager. If Rethink is enabled returns an DatabaseAuthManager, if not a JSONAuthManager
func MakeAuthManager(logger slog.Instance, dbHandler DatabaseHandler) interfaces.AuthManager {
	if dbHandler != nil {
		return NewDatabaseAuthManager(logger, dbHandler)
	}

	return MakeJSONAuthManager(logger)
}
