package agent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	"time"
)

type DatabaseTokenManager struct {
	log     slog.Instance
	dbToken DBToken
}

type DBToken interface {
	GetUser(username string) (um *models.User, err error)
	AddUserToken(ut models.UserToken) (string, error)
	RemoveUserToken(token string) (err error)
	GetUserToken(token string) (ut *models.UserToken, err error)
	InvalidateUserTokens() (int, error)
}

// MakeDatabaseTokenManager creates an instance of TokenManager that stores data in RethinkDB
func MakeDatabaseTokenManager(logger slog.Instance, dbToken DBToken) *DatabaseTokenManager {
	if logger == nil {
		logger = slog.Scope("DB-TM")
	} else {
		logger = logger.SubScope("DB-TM")
	}
	logger.Info("Creating Database Token Manager")
	return &DatabaseTokenManager{
		log:     logger,
		dbToken: dbToken,
	}
}

// AddUserWithExpiration adds an user to Token Manager that will expires in `expiration` seconds.
func (rtm *DatabaseTokenManager) AddUserWithExpiration(user interfaces.UserData, expiration int) string {
	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	_, _ = rtm.dbToken.AddUserToken(models.UserToken{
		Fingerprint: user.GetFingerPrint(),
		Username:    user.GetUsername(),
		CreatedAt:   user.GetCreatedAt(),
		Fullname:    user.GetFullName(),
		Expiration:  user.GetCreatedAt().Add(time.Duration(expiration) * time.Second),
		Token:       token,
	})

	return token
}

// AddUser adds a user to Token Manager and returns a login token
func (rtm *DatabaseTokenManager) AddUser(user interfaces.UserData) string {
	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	_, _ = rtm.dbToken.AddUserToken(models.UserToken{
		Fingerprint: user.GetFingerPrint(),
		Username:    user.GetUsername(),
		CreatedAt:   user.GetCreatedAt(),
		Fullname:    user.GetFullName(),
		Expiration:  user.GetCreatedAt().Add(time.Duration(config.AgentTokenExpiration) * time.Second),
		Token:       token,
	})

	return token
}

func (rtm *DatabaseTokenManager) invalidateTokens() {
	rtm.log.Await("Checking for invalid tokens")
	n, err := rtm.dbToken.InvalidateUserTokens()
	if err != nil {
		rtm.log.Error(err)
		return
	}

	rtm.log.Done("Cleaned %d invalid / expired tokens", n)
}

// Verify verifies if the specified token is valid
func (rtm *DatabaseTokenManager) Verify(token string) error {
	ut, err := rtm.dbToken.GetUserToken(token)

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
func (rtm *DatabaseTokenManager) GetUserData(token string) interfaces.UserData {
	udata, _ := rtm.dbToken.GetUserToken(token)
	return udata
}

// InvalidateToken removes a token from the database making it unusable in the future
func (rtm *DatabaseTokenManager) InvalidateToken(token string) error {
	u, _ := rtm.dbToken.GetUserToken(token)

	if u == nil {
		return fmt.Errorf("not exists")
	}

	return rtm.dbToken.RemoveUserToken(token)
}
