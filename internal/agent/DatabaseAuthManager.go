package agent

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	"golang.org/x/crypto/bcrypt"
)

type DatabaseAuthManager struct {
	sync.Mutex
	log    slog.Instance
	dbAuth DBAuth
}

type DBAuth interface {
	GetUser(username string) (um *models.User, err error)
	AddUser(um models.User) (string, error)
	UpdateUser(um models.User) error
}

// NewDatabaseAuthManager creates an instance of Auth Manager that uses RethinkDB as storage
func NewDatabaseAuthManager(logger slog.Instance, dbAuth DBAuth) *DatabaseAuthManager {
	if logger == nil {
		logger = slog.Scope("DB-AM")
	} else {
		logger = logger.SubScope("DB-AM")
	}

	logger.Info("Creating Database Auth Manager")
	ram := &DatabaseAuthManager{
		log:    logger,
		dbAuth: dbAuth,
	}

	if !ram.UserExists("admin") {
		ram.log.Warn("User admin does not exists. Creating default")
		ram.addDefaultAdmin()
	}

	return ram
}

func (ram *DatabaseAuthManager) addDefaultAdmin() {
	err := ram.LoginAdd("admin", "admin", "Administrator", config.AgentKeyFingerPrint)

	if err != nil {
		ram.log.Fatal("Error adding default admin: %v", err)
	}
}

// UserExists checks if a user with specified username exists in AuthManager
func (ram *DatabaseAuthManager) UserExists(username string) bool {
	ram.Lock()
	defer ram.Unlock()

	um, err := ram.dbAuth.GetUser(username)

	if err != nil || um == nil {
		return false
	}

	return true
}

// LoginAuth performs a login with the specified username and password
func (ram *DatabaseAuthManager) LoginAuth(username, password string) (fingerPrint, fullname string, err error) {
	ram.Lock()
	defer ram.Unlock()

	um, err := ram.dbAuth.GetUser(username)

	if err != nil || um == nil {
		// To avoid timing attack pretend that we have a valid user
		um = &models.User{ID: invalidUserId}
	}

	hash, err := base64.StdEncoding.DecodeString(um.Password)
	if err != nil {
		ram.log.Error("Error decoding hash: %v", err)
		// Same here, let's pretend we have a hash
		hash = []byte(bcryptFallbackHash)
		// Unset password to avoid knowing the correct password from the hash above
		password = ""
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password)) // Execute even if we know it's invalid
	if um.ID == invalidUserId || err != nil {                   // Now we check for a invalid user
		return "", "", fmt.Errorf("invalid username or password")
	}

	return um.Fingerprint, um.FullName, nil
}

// LoginAdd creates a new user in AuthManager
func (ram *DatabaseAuthManager) LoginAdd(username, password, fullname, fingerprint string) error {
	ram.Lock()
	defer ram.Unlock()

	um, _ := ram.dbAuth.GetUser(username)

	if um != nil {
		return fmt.Errorf("already exists")
	}

	fp := fingerprint
	if fp == "" {
		fp = config.AgentKeyFingerPrint
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	encodedPassword := base64.StdEncoding.EncodeToString(pass)

	_, err = ram.dbAuth.AddUser(models.User{
		Fingerprint: fp,
		Username:    username,
		Password:    encodedPassword,
		FullName:    fullname,
		CreatedAt:   time.Now(),
	})

	return err
}

// ChangePassword changes the password of the specified user
func (ram *DatabaseAuthManager) ChangePassword(username, password string) error {
	ram.Lock()
	defer ram.Unlock()

	um, err := ram.dbAuth.GetUser(username)

	if err != nil || um == nil {
		return fmt.Errorf("user does not exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	um.Password = base64.StdEncoding.EncodeToString(pass)

	return ram.dbAuth.UpdateUser(*um)
}
