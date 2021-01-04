package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
)

const userTokenPrefix = "userToken-"

// AddUserToken adds a new user token to be valid and returns its token ID
func (h *Driver) AddUserToken(ut models.UserToken) (string, error) {
	h.log.Debug("AddUserToken(%s, %s)", ut.Username, ut.Fingerprint)
	ut.ID = uuid.EnsureUUID(h.log)
	exp := ut.Expiration.Sub(time.Now())

	if err := h.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userTokenPrefix + ut.Token,
		Value: &ut,
		TTL:   exp,
	}); err != nil {
		return "", err
	}

	return ut.ID, nil
}

// RemoveUserToken removes a user token from the database
func (h *Driver) RemoveUserToken(token string) (err error) {
	h.log.Debug("RemoveUserToken(%s)", token)
	return h.cache.Delete(context.TODO(), userTokenPrefix+token)
}

// GetUserToken fetch a UserToken object by the specified token
func (h *Driver) GetUserToken(token string) (ut *models.UserToken, err error) {
	h.log.Debug("GetUserToken(%s)", token)
	err = h.cache.Get(context.TODO(), userTokenPrefix+token, &ut)
	return ut, err
}

// InvalidateUserTokens removes all user tokens that had been already expired
// Does nothing on REDIS due automatic expiration using TTL
func (h *Driver) InvalidateUserTokens() (int, error) {
	// Not needed for redis, automatic expiration due TTL
	return 0, nil
}
