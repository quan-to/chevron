package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/etc"
	"sync"
	"time"
)

type MemoryUser struct {
	username    string
	token       string
	createdAt   time.Time
	fingerPrint string
}

func (mu *MemoryUser) GetUsername() string {
	return mu.username
}

func (mu *MemoryUser) GetToken() string {
	return mu.token
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

type MemoryTokenManager struct {
	storedTokens map[string]etc.UserData
	lock         sync.Mutex
}

func MakeMemoryTokenManager() *MemoryTokenManager {
	return &MemoryTokenManager{
		lock:         sync.Mutex{},
		storedTokens: map[string]etc.UserData{},
	}
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

	if time.Since(user.GetCreatedAt()) > time.Duration(remote_signer.AgentTokenExpiration)*time.Second {
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
