package agent

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/etc"
	"sync"
	"time"
)

var mtmLog = SLog.Scope("Memory-TM")

type MemoryUser struct {
	username    string
	fullname    string
	token       string
	createdAt   time.Time
	fingerPrint string
	expiration  time.Time
}

func (mu *MemoryUser) GetUsername() string {
	return mu.username
}

func (mu *MemoryUser) GetToken() string {
	return mu.token
}

func (mu *MemoryUser) GetFullName() string {
	return mu.fullname
}

func (mu *MemoryUser) GetUserdata() interface{} {
	return nil
}

func (mu *MemoryUser) GetCreatedAt() time.Time {
	return mu.createdAt
}

func (mu *MemoryUser) GetFingerPrint() string {
	return mu.fingerPrint
}

func (mu *MemoryUser) GetExpiration() time.Time {
	return mu.expiration
}

type MemoryTokenManager struct {
	storedTokens map[string]*MemoryUser
	lock         sync.Mutex
}

func MakeMemoryTokenManager() *MemoryTokenManager {
	mtmLog.Info("Creating Memory Token Manager")
	return &MemoryTokenManager{
		lock:         sync.Mutex{},
		storedTokens: map[string]*MemoryUser{},
	}
}

func (mtm *MemoryTokenManager) AddUserWithExpiration(user etc.UserData, expiration int) string {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	mtm.storedTokens[token] = &MemoryUser{
		username:    user.GetUsername(),
		token:       token,
		createdAt:   user.GetCreatedAt(),
		fingerPrint: user.GetFingerPrint(),
		fullname:    user.GetFullName(),
		expiration:  user.GetCreatedAt().Add(time.Duration(expiration) * time.Second),
	}

	return token
}

func (mtm *MemoryTokenManager) AddUser(user etc.UserData) string {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	tokenUuid, _ := uuid.NewUUID()
	token := tokenUuid.String()

	mtm.storedTokens[token] = &MemoryUser{
		username:    user.GetUsername(),
		token:       token,
		createdAt:   user.GetCreatedAt(),
		fingerPrint: user.GetFingerPrint(),
		fullname:    user.GetFullName(),
		expiration:  user.GetCreatedAt().Add(time.Duration(remote_signer.AgentTokenExpiration) * time.Second),
	}

	return token
}

func (mtm *MemoryTokenManager) Verify(token string) error {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	user := mtm.storedTokens[token]
	if user == nil {
		return fmt.Errorf("not found")
	}

	if time.Now().After(user.GetExpiration()) {
		delete(mtm.storedTokens, token)
		return fmt.Errorf("expired")
	}

	return nil
}

func (mtm *MemoryTokenManager) GetUserData(token string) etc.UserData {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	return mtm.storedTokens[token]
}