package cache

import (
	"github.com/quan-to/chevron/pkg/models"
)

func (h *Driver) InitCursor() error {
	return h.proxy.InitCursor()
}

func (h *Driver) FinishCursor() error {
	return h.proxy.FinishCursor()
}

func (h *Driver) NextGPGKey(key *models.GPGKey) bool {
	return h.proxy.NextGPGKey(key)
}

func (h *Driver) NextUser(user *models.User) bool {
	return h.proxy.NextUser(user)
}

func (h *Driver) NumGPGKeys() (int, error) {
	return h.proxy.NumGPGKeys()
}
