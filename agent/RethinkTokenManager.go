package agent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	"time"
)

type RethinkTokenManager struct {
	log slog.Instance
}

// MakeRethinkTokenManager creates an instance of TokenManager that stores data in RethinkDB
func MakeRethinkTokenManager(logger slog.Instance) *RethinkTokenManager {
	if logger == nil {
		logger = slog.Scope("RQL-TM")
	} else {
		logger = logger.SubScope("RQL-TM")
	}
	logger.Info("Creating RethinkDB Token Manager")
	return &RethinkTokenManager{
		log: logger,
	}
}

func (rtm *RethinkTokenManager) AddUserWithExpiration(user etc.UserData, expiration int) string {
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

func (rtm *RethinkTokenManager) AddUser(user etc.UserData) string {
	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	conn := etc.GetConnection()

	_, _ = models.AddUserToken(conn, &models.UserToken{
		FingerPrint: user.GetFingerPrint(),
		Username:    user.GetUsername(),
		CreatedAt:   user.GetCreatedAt(),
		Fullname:    user.GetFullName(),
		Expiration:  user.GetCreatedAt().Add(time.Duration(remote_signer.AgentTokenExpiration) * time.Second),
		Token:       token,
	})

	return token
}

func (rtm *RethinkTokenManager) invalidateTokens() {
	rtm.log.Await("Checking for invalid tokens")
	conn := etc.GetConnection()
	n, err := models.InvalidateUserTokens(conn)
	if err != nil {
		rtm.log.Error(err)
		return
	}

	rtm.log.Done("Cleaned %d invalid / expired tokens", n)
}

func (rtm *RethinkTokenManager) Verify(token string) error {
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

func (rtm *RethinkTokenManager) GetUserData(token string) etc.UserData {
	conn := etc.GetConnection()
	udata, _ := models.GetUserToken(conn, token)
	return udata
}

// InvalidateToken removes a token from the database making it unusable in the future
func (rtm *RethinkTokenManager) InvalidateToken(token string) error {
	conn := etc.GetConnection()
	u, _ := models.GetUserToken(conn, token)

	if u == nil {
		return fmt.Errorf("not exists")
	}

	return models.RemoveUserToken(conn, token)
}
