package agent

import (
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

var ramLog = slog.Scope("RQL-AM")

type RethinkAuthManager struct {
	sync.Mutex
}

func MakeRethinkAuthManager() *RethinkAuthManager {
	ramLog.Info("Creating RethinkDB Auth Manager")
	ram := &RethinkAuthManager{}

	if !ram.UserExists("admin") {
		ramLog.Warn("User admin does not exists. Creating default")
		ram.addDefaultAdmin()
	}

	return ram
}

func (ram *RethinkAuthManager) addDefaultAdmin() {
	err := ram.LoginAdd("admin", "admin", "Administrator", remote_signer.AgentKeyFingerPrint)

	if err != nil {
		ramLog.Fatal("Error adding default admin: %v", err)
	}
}

func (ram *RethinkAuthManager) UserExists(username string) bool {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, err := models.GetUser(conn, username)

	if err != nil || um == nil {
		return false
	}

	return true
}

func (ram *RethinkAuthManager) LoginAuth(username, password string) (fingerPrint, fullname string, err error) {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, err := models.GetUser(conn, username)

	if err != nil || um == nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	hash, err := base64.StdEncoding.DecodeString(um.Password)
	if err != nil {
		ramLog.Error("Error decoding hash: %v", err)
		return "", "", fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	return um.FingerPrint, um.Fullname, nil
}

func (ram *RethinkAuthManager) LoginAdd(username, password, fullname, fingerprint string) error {
	ram.Lock()
	defer ram.Unlock()

	conn := etc.GetConnection()

	um, _ := models.GetUser(conn, username)

	if um != nil {
		return fmt.Errorf("already exists")
	}

	fp := fingerprint
	if fp == "" {
		fp = remote_signer.AgentKeyFingerPrint
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	encodedPassword := base64.StdEncoding.EncodeToString(pass)

	_, err = models.AddUser(conn, &models.UserModel{
		FingerPrint: fp,
		Username:    username,
		Password:    encodedPassword,
		Fullname:    fullname,
		CreatedAt:   time.Now(),
	})

	return err
}

func (ram *RethinkAuthManager) ChangePassword(username, password string) error {
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
