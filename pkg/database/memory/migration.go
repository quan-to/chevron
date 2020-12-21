package memory

import (
	"fmt"

	"github.com/quan-to/chevron/pkg/models"
)

func (h *DbDriver) InitCursor() error {
	return fmt.Errorf("not implemented")
}

func (h *DbDriver) FinishCursor() error {
	return fmt.Errorf("not implemented")
}

func (h *DbDriver) NextGPGKey(key *models.GPGKey) bool {
	return false
}

func (h *DbDriver) NextUser(user *models.User) bool {
	return false
}

func (h *DbDriver) NumGPGKeys() (int, error) {
	return 0, fmt.Errorf("not implemented")
}
