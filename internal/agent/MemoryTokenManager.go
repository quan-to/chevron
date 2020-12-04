package agent

import (
	"fmt"
	"github.com/google/uuid"
	remote_signer "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
	"sync"
	"time"
)

type memoryUser struct {
	username    string
	fullname    string
	token       string
	createdAt   time.Time
	fingerPrint string
	expiration  time.Time
}

func (mu *memoryUser) GetId() string {
	return mu.username
}

func (mu *memoryUser) GetUsername() string {
	return mu.username
}

func (mu *memoryUser) GetToken() string {
	return mu.token
}

func (mu *memoryUser) GetFullName() string {
	return mu.fullname
}

func (mu *memoryUser) GetUserdata() interface{} {
	return nil
}

func (mu *memoryUser) GetCreatedAt() time.Time {
	return mu.createdAt
}

func (mu *memoryUser) GetFingerPrint() string {
	return mu.fingerPrint
}

func (mu *memoryUser) GetExpiration() time.Time {
	return mu.expiration
}

type MemoryTokenManager struct {
	storedTokens map[string]*memoryUser
	lock         sync.Mutex
	log          slog.Instance
}

// MakeMemoryTokenManager creates an instance of TokenManager managed in memory
func MakeMemoryTokenManager(logger slog.Instance) *MemoryTokenManager {
	if logger == nil {
		logger = slog.Scope("Memory-TM")
	} else {
		logger = logger.SubScope("MTM")
	}
	logger.Info("Creating Memory Token Manager")
	return &MemoryTokenManager{
		lock:         sync.Mutex{},
		storedTokens: map[string]*memoryUser{},
		log:          logger,
	}
}

func (mtm *MemoryTokenManager) AddUserWithExpiration(user interfaces.UserData, expiration int) string {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	tokenUuid, _ := uuid.NewRandom()
	token := tokenUuid.String()

	mtm.storedTokens[token] = &memoryUser{
		username:    user.GetUsername(),
		token:       token,
		createdAt:   user.GetCreatedAt(),
		fingerPrint: user.GetFingerPrint(),
		fullname:    user.GetFullName(),
		expiration:  user.GetCreatedAt().Add(time.Duration(expiration) * time.Second),
	}

	return token
}

func (mtm *MemoryTokenManager) AddUser(user interfaces.UserData) string {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	tokenUuid, _ := uuid.NewUUID()
	token := tokenUuid.String()

	mtm.storedTokens[token] = &memoryUser{
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

// InvalidateToken removes a token from the internal memory making it unusable in the future
func (mtm *MemoryTokenManager) InvalidateToken(token string) error {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	u := mtm.storedTokens[token]

	if u == nil {
		return fmt.Errorf("not exists")
	}

	delete(mtm.storedTokens, token)

	return nil
}

func (mtm *MemoryTokenManager) GetUserData(token string) interfaces.UserData {
	mtm.lock.Lock()
	defer mtm.lock.Unlock()

	return mtm.storedTokens[token]
}
