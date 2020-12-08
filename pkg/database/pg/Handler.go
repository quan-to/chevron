//go:generate go-bindata -pkg migrations -prefix migrations/ -o migrations/bindata.go -ignore=bindata.go -ignore=auto.go migrations
package pg

import (
	"context"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/quan-to/chevron/pkg/database/pg/migrations"
	"github.com/quan-to/slog"
)

// PostgreSQLDBDriver is a database driver for PostgreSQL
type PostgreSQLDBDriver struct {
	log      slog.Instance
	database string
	conn     *sqlx.DB
}

// MakeRethinkDBDriver creates a new database driver for rethinkdb
func MakePostgreSQLDBDriver(log slog.Instance) *PostgreSQLDBDriver {
	if log == nil {
		log = slog.Scope("PostgreSQL")
	} else {
		log = log.SubScope("PostgreSQL")
	}

	return &PostgreSQLDBDriver{
		log: log,
	}
}

func (h *PostgreSQLDBDriver) Connect(connectionString string) error {
	h.log.Info("Connecting to PostgreSQL to run migrations")
	h.log.Debug("Connection string %s", connectionString)
	bindatadriver, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		}))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", bindatadriver, connectionString)
	if err != nil {
		h.log.Error("Error loading migrations: %s", err)
		return err
	}
	err = m.Up()
	if err != nil && !strings.EqualFold("no change", err.Error()) {
		h.log.Error("Error running migrations: %s", err)
		return err
	}
	h.log.Info("Migrations finished")
	h.log.Info("Connecting to postgresql for usage")
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	h.conn = db
	h.log.Info("Connected!")
	return nil
}

func (h *PostgreSQLDBDriver) HealthCheck() error {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5)) // 5 second timeout
	return h.conn.PingContext(ctx)
}

func (h *PostgreSQLDBDriver) rollbackIfErrorCommitIfNot(err error, tx *sqlx.Tx) {
	if err != nil && tx != nil {
		h.log.Debug("and error ocurred, rollback transaction: %s", err)
		if slog.DebugEnabled() {
			debug.PrintStack()
		}
		_ = tx.Rollback()
	} else if tx != nil {
		err = tx.Commit()
	}
}

func errorIsNotNilAndNotNotFound(err error) bool {
	return err != nil && !strings.EqualFold("not found", err.Error())
}
