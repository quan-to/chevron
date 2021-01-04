package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mewkiz/pkg/osutil"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
)

var cli struct {
	FromConfigFile string `arg:"" name:"path" help:"JSON Config file for the source database" type:"path"`
	ToConfigFile   string `arg:"" name:"path" help:"JSON Config file for the destination database" type:"path"`
}

type dbSource interface {
	InitCursor() error
	FinishCursor() error
	NextGPGKey(key *models.GPGKey) bool
	NextUser(user *models.User) bool
	NumGPGKeys() (int, error)
}

type dbDestination interface {
	AddGPGKey(key models.GPGKey) (string, bool, error)
	AddGPGKeys(keys []models.GPGKey) ([]string, []bool, error)
	AddUser(um models.User) (string, error)
	GetUser(username string) (um *models.User, err error)
	UpdateUser(um models.User) error
}

type migrateHandler interface {
	dbSource
	dbDestination
}

func main() {
	//var err error
	ctx := kong.Parse(&cli,
		kong.Name("dbmigrate"),
		kong.Description("Database migration tool for Chevron.\nThe configuration JSON files has the same fields as environment variables"))

	if !osutil.Exists(cli.FromConfigFile) {
		ctx.Fatalf("No such file: %s", cli.FromConfigFile)
	}
	if !osutil.Exists(cli.ToConfigFile) {
		ctx.Fatalf("No such file: %s", cli.ToConfigFile)
	}

	slog.SetDebug(false)

	logger := slog.Scope("Migrate").WithCustomWriter(ctx.Stdout)

	src, err := makeHandler(cli.FromConfigFile, logger)
	if err != nil {
		logger.Fatal("Error initializing source handler: %s", err)
	}
	dst, err := makeHandler(cli.ToConfigFile, logger)
	if err != nil {
		logger.Fatal("Error initializing destination handler: %s", err)
	}
	logger.Info("Couting keys...")
	err = src.InitCursor()
	if err != nil {
		logger.Fatal("Error initializing cursor: %s", err)
	}
	totalKeys, err := src.NumGPGKeys()
	if err != nil {
		logger.Fatal("Error getting number of keys: %s", err)
	}
	logger.Info("Starting migration of %d keys", totalKeys)
	migratedKeys := 0

	gpgKey := models.GPGKey{}
	var keys []models.GPGKey

	logger.Info("Fetching <= 500 keys from source")
	for src.NextGPGKey(&gpgKey) {
		if len(keys) >= 500 {
			logger.Info("Saving %d keys to destination", len(keys))
			_, _, err = dst.AddGPGKeys(keys)
			if err != nil {
				logger.Fatal("error migrating keys: %s", err)
			}
			migratedKeys += len(keys)
			logger.Info("Migrated %6d from %6d keys... [%d]", migratedKeys, totalKeys, len(keys))
			keys = nil
			logger.Info("Fetching <= 500 keys from source")
		}
		gpgKey.ID = "" // Re-generate the ID
		keys = append(keys, gpgKey)
	}

	if len(keys) > 0 {
		logger.Info("Saving %d keys to destination", len(keys))
		_, _, err = dst.AddGPGKeys(keys)
		if err != nil {
			logger.Fatal("error migrating keys: %s", err)
		}
		migratedKeys += len(keys)
		keys = nil
	}

	logger.Info("Migrated %d keys...", migratedKeys)
	logger.Info("Migrating users...")

	user := models.User{}

	migratedUsers := 0

	for src.NextUser(&user) {
		user.ID = "" // Re-generate the ID
		_, err := dst.AddUser(user)
		if err != nil && strings.EqualFold("already exists", err.Error()) {
			logger.Info("User %s already exists. Updating it...", user.Username)
			oldUser, _ := dst.GetUser(user.Username)
			user.ID = oldUser.ID
			err := dst.UpdateUser(user)
			if err != nil {
				logger.Fatal("error migrating user %s: %s", user.Username, err)
			}
			// Update user TODO
		} else if err != nil {
			logger.Fatal("error migrating user %s: %s", user.Username, err)
		}

		migratedUsers++
		if migratedUsers%10 == 0 {
			logger.Info("Migrated %d users", migratedUsers)
		}
	}
	logger.Info("Migrated %d users...", migratedUsers)

	err = src.FinishCursor()
	if err != nil {
		logger.Fatal("error closing cursor")
	}
}
