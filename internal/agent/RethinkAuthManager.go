package agent

import (
	"encoding/base64"
	"fmt"
	cofig "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type rethinkAuthManager struct {
	sync.Mutex
	log slog.Instance
}

// MakeRethinkAuthManager creates an instance of Auth Manager that uses RethinkDB as storage
func MakeRethinkAuthManager(logger slog.Instance) interfaces.AuthManager {
	if logger == nil {
		logger = slog.Scope("RQL-AM")
	} else {
		logger = logger.SubScope("RQL-AM")
	}

	logger.Info("Creating RethinkDB Auth Manager")
	ram := &rethinkAuthManager{
		log: logger,
	}

	if !ram.UserExists("admin") {
		ram.log.Warn("User admin does not exists. Creating default")
		ram.addDefaultAdmin()
	}

	return ram
}

func (ram *rethinkAuthManager) addDefaultAdmin() {
	err := ram.LoginAdd("admin", "admin", "Administrator", cofig.AgentKeyFingerPrint)

	if err != nil {
		ram.log.Fatal("Error adding default admin: %v", err)
	}
}

// UserExists checks if a user with specified username exists in AuthManager
func (ram *rethinkAuthManager) UserExists(username string) bool {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, err := models.GetUser(conn, username)

	if err != nil || um == nil {
		return false
	}

	return true
}

// LoginAuth performs a login with the specified username and password
func (ram *rethinkAuthManager) LoginAuth(username, password string) (fingerPrint, fullname string, err error) {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, err := models.GetUser(conn, username)

	if err != nil || um == nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	hash, err := base64.StdEncoding.DecodeString(um.Password)
	if err != nil {
		ram.log.Error("Error decoding hash: %v", err)
		return "", "", fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	return um.FingerPrint, um.Fullname, nil
}

// LoginAdd creates a new user in AuthManager
func (ram *rethinkAuthManager) LoginAdd(username, password, fullname, fingerprint string) error {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, _ := models.GetUser(conn, username)

	if um != nil {
		return fmt.Errorf("already exists")
	}

	fp := fingerprint
	if fp == "" {
		fp = cofig.AgentKeyFingerPrint
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	encodedPassword := base64.StdEncoding.EncodeToString(pass)

	_, err = models.AddUser(conn, &models.User{
		Fingerprint: fp,
		Username:    username,
		Password:    encodedPassword,
		FullName:    fullname,
		CreatedAt:   time.Now(),
	})

	return err
}

// ChangePassword changes the password of the specified user
func (ram *rethinkAuthManager) ChangePassword(username, password string) error {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, err := models.GetUser(conn, username)

	if err != nil || um == nil {
		return fmt.Errorf("user does not exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	um.Password = base64.StdEncoding.EncodeToString(pass)

	return models.UpdateUser(conn, um)
}
