package memory

import (
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	"sync"
)

// MemoryDBDriver is a database driver for in-memory database for testing
type MemoryDBDriver struct {
	log    slog.Instance
	users  []models.User
	tokens []models.UserToken
	keys   []models.GPGKey
	lock   sync.RWMutex
}

// MakeMemoryDBDriver creates a new database driver for rethinkdb
func MakeMemoryDBDriver(log slog.Instance) *MemoryDBDriver {
	if log == nil {
		log = slog.Scope("MemoryDB")
	} else {
		log = log.SubScope("MemoryDB")
	}

	return &MemoryDBDriver{
		log:  log,
		lock: sync.RWMutex{},
	}
}

func (h *MemoryDBDriver) HealthCheck() error {
	return nil
}
