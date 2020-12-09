package memory

import (
	"sync"

	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
)

// DbDriver is a database driver for in-memory database for testing
type DbDriver struct {
	log    slog.Instance
	users  []models.User
	tokens []models.UserToken
	keys   []models.GPGKey
	lock   sync.RWMutex
}

// MakeMemoryDBDriver creates a new database driver for rethinkdb
func MakeMemoryDBDriver(log slog.Instance) *DbDriver {
	if log == nil {
		log = slog.Scope("MemoryDB")
	} else {
		log = log.SubScope("MemoryDB")
	}

	return &DbDriver{
		log:  log,
		lock: sync.RWMutex{},
	}
}

func (h *DbDriver) HealthCheck() error {
	return nil
}
