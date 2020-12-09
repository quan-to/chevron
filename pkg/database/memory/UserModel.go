package memory

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/quan-to/chevron/pkg/models"
)

func (h *MemoryDBDriver) AddUser(um models.User) (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	um.ID = uuid.New().String()

	h.users = append(h.users, um)

	return um.ID, nil
}

func (h *MemoryDBDriver) GetUser(username string) (um *models.User, err error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, v := range h.users {
		if strings.EqualFold(username, v.Username) {
			n := v // Copy
			return &n, nil
		}
	}

	return um, fmt.Errorf("not found")
}

func (h *MemoryDBDriver) UpdateUser(um models.User) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if um.ID == "" {
		return fmt.Errorf("not found")
	}

	for i, v := range h.users {
		if v.ID == um.ID {
			h.users[i] = um
			return nil
		}
	}

	return fmt.Errorf("not found")
}
