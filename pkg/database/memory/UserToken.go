package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/quan-to/chevron/pkg/models"
)

func (h *MemoryDBDriver) AddUserToken(ut models.UserToken) (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	ut.ID = uuid.New().String()

	h.tokens = append(h.tokens, ut)

	return ut.ID, nil
}

// RemoveUserToken removes a user token from the database
func (h *MemoryDBDriver) RemoveUserToken(token string) (err error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for i, v := range h.tokens {
		if strings.EqualFold(v.Token, token) {
			h.tokens = append(h.tokens[:i], h.tokens[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("not found")
}

func (h *MemoryDBDriver) GetUserToken(token string) (ut *models.UserToken, err error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, v := range h.tokens {
		if strings.EqualFold(v.Token, token) {
			t := v
			return &t, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (h *MemoryDBDriver) InvalidateUserTokens() (int, error) {
	h.lock.RLock()
	var tokensToDelete []string
	for _, v := range h.tokens {
		if time.Since(v.Expiration) >= 0 {
			tokensToDelete = append(tokensToDelete, v.Token)
		}
	}
	h.lock.RUnlock()

	for _, v := range tokensToDelete {
		_ = h.RemoveUserToken(v)
	}

	return len(tokensToDelete), nil
}
