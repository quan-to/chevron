package agent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mewkiz/pkg/osutil"
	"github.com/quan-to/chevron"
	"github.com/quan-to/slog"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"sync"
)

const jamFileName = "users.json"
const jamFilePerm = 0600

type jsonUser struct {
	Username    string
	Password    string
	FullName    string
	FingerPrint string
}

type JSONAuthManager struct {
	sync.Mutex
	users map[string]jsonUser
	log   slog.Instance
}

// MakeJSONAuthManager creates an instance of AuthManager that uses JSON Storage
func MakeJSONAuthManager(logger slog.Instance) *JSONAuthManager {
	if logger == nil {
		logger = slog.Scope("JSON-AM")
	} else {
		logger = logger.SubScope("JSON-AM")
	}

	logger.Info("Creating JSON Auth Manager")
	jam := JSONAuthManager{log: logger}
	jam.loadFile()
	return &jam
}

func (jam *JSONAuthManager) loadFile() {
	// Load From Self folder
	if !osutil.Exists(jamFileName) {
		jam.log.Warn("File %s does not exists. Creating one...", jamFileName)
		err := ioutil.WriteFile(jamFileName, []byte("{}"), jamFilePerm)
		if err != nil {
			jam.log.Fatal("Error writing file %s: %s", jamFileName, err)
		}
	}

	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		jam.log.Fatal("Error writing file %s: %s", jamFileName, err)
	}

	err = json.Unmarshal(data, &jam.users)

	if err != nil {
		jam.log.Fatal("Corrupted or invalid JSON at %s: %s", jamFileName, err)
	}

	jam.log.Info("Loaded %d users from %s", len(jam.users), jamFileName)

	if len(jam.users) == 0 {
		jam.addDefaultAdmin()
	}
}

func (jam *JSONAuthManager) addDefaultAdmin() {
	err := jam.LoginAdd("admin", "admin", "Administrator", remote_signer.AgentKeyFingerPrint)

	if err != nil {
		jam.log.Fatal("Error adding default admin: %v", err)
	}
}

func (jam *JSONAuthManager) flushFile() {
	jam.log.Warn("Saving credentials to %s", jamFileName)
	data, _ := json.Marshal(jam.users)
	err := ioutil.WriteFile(jamFileName, data, jamFilePerm)
	if err != nil {
		jam.log.Error("Error saving credentials: %s", err)
	}
}

func (jam *JSONAuthManager) UserExists(username string) bool {
	jam.Lock()
	defer jam.Unlock()

	_, exists := jam.users[username]

	return exists
}

func (jam *JSONAuthManager) LoginAuth(username, password string) (fingerPrint, fullname string, err error) {
	jam.Lock()
	defer jam.Unlock()

	user, exists := jam.users[username]

	if !exists {
		return "", "", fmt.Errorf("invalid username or password")
	}

	hash, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		jam.log.Error("Error decoding hash: %v", err)
		return "", "", fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	return user.FingerPrint, user.FullName, nil
}

func (jam *JSONAuthManager) LoginAdd(username, password, fullname, fingerprint string) error {
	jam.Lock()
	defer jam.Unlock()
	_, exists := jam.users[username]

	if exists {
		return fmt.Errorf("user already exists")
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

	jam.users[username] = jsonUser{
		Username:    username,
		FullName:    fullname,
		FingerPrint: fp,
		Password:    encodedPassword,
	}

	jam.flushFile()

	return nil
}

func (jam *JSONAuthManager) ChangePassword(username, password string) error {
	jam.Lock()
	defer jam.Unlock()

	user, exists := jam.users[username]

	if !exists {
		return fmt.Errorf("user does not exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error generating hash: %v", err)
	}

	user.Password = base64.StdEncoding.EncodeToString(pass)

	jam.users[username] = user

	jam.flushFile()

	return nil
}
