package agent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
	"time"
)

type rethinkTokenManager struct {
	log slog.Instance
}

// MakeRethinkTokenManager creates an instance of TokenManager that stores data in RethinkDB
func MakeRethinkTokenManager(logger slog.Instance) interfaces.TokenManager {
	if logger == nil {
		logger = slog.Scope("RQL-TM")
	} else {
		logger = logger.SubScope("RQL-TM")
	}
	logger.Info("Creating RethinkDB Token Manager")
	return &rethinkTokenManager{
		log: logger,
	}
}

// AddUserWithExpiration adds an user to Token Manager that will expires in `expiration` seconds.
func (rtm *rethinkTokenManager) AddUserWithExpiration(user interfaces.UserData, expiration int) string {
	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	conn := etc.GetConnection()

	_, _ = models.AddUserToken(conn, &models.UserToken{
		FingerPrint: user.GetFingerPrint(),
		Username:    user.GetUsername(),
		CreatedAt:   user.GetCreatedAt(),
		Fullname:    user.GetFullName(),
		Expiration:  user.GetCreatedAt().Add(time.Duration(expiration) * time.Second),
		Token:       token,
	})

	return token
}

// AddUser adds a user to Token Manager and returns a login token
func (rtm *rethinkTokenManager) AddUser(user interfaces.UserData) string {
	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	conn := etc.GetConnection()

	_, _ = models.AddUserToken(conn, &models.UserToken{
		FingerPrint: user.GetFingerPrint(),
		Username:    user.GetUsername(),
		CreatedAt:   user.GetCreatedAt(),
		Fullname:    user.GetFullName(),
		Expiration:  user.GetCreatedAt().Add(time.Duration(config.AgentTokenExpiration) * time.Second),
		Token:       token,
	})

	return token
}

func (rtm *rethinkTokenManager) invalidateTokens() {
	rtm.log.Await("Checking for invalid tokens")
	conn := etc.GetConnection()
	n, err := models.InvalidateUserTokens(conn)
	if err != nil {
		rtm.log.Error(err)
		return
	}

	rtm.log.Done("Cleaned %d invalid / expired tokens", n)
}

// Verify verifies if the specified token is valid
func (rtm *rethinkTokenManager) Verify(token string) error {
	conn := etc.GetConnection()

	ut, err := models.GetUserToken(conn, token)

	if err != nil {
		return err
	}

	if ut == nil {
		return fmt.Errorf("not found")
	}

	if time.Now().After(ut.Expiration) {
		go rtm.invalidateTokens()
		return fmt.Errorf("expired")
	}

	return nil
}

// GetUserData returns the user data for the specified token
func (rtm *rethinkTokenManager) GetUserData(token string) interfaces.UserData {
	conn := etc.GetConnection()
	udata, _ := models.GetUserToken(conn, token)
	return udata
}

// InvalidateToken removes a token from the database making it unusable in the future
func (rtm *rethinkTokenManager) InvalidateToken(token string) error {
	conn := etc.GetConnection()
	u, _ := models.GetUserToken(conn, token)

	if u == nil {
		return fmt.Errorf("not exists")
	}

	return models.RemoveUserToken(conn, token)
}
