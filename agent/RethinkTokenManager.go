package agent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/slog"
	"time"
)

var rtmLog = slog.Scope("RQL-TM")

type RethinkTokenManager struct {
}

func MakeRethinkTokenManager() *RethinkTokenManager {
	rtmLog.Info("Creating RethinkDB Token Manager")
	return &RethinkTokenManager{}
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
	rtmLog.Info("Checking for invalid tokens")
	conn := etc.GetConnection()
	n, err := models.InvalidateUserTokens(conn)
	if err != nil {
		rtmLog.Error(err)
		return
	}

	rtmLog.Warn("Cleaned %d invalid / expired tokens", n)
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
