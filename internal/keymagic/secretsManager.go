// +build !js,!wasm

package keymagic

import (
	"context"
	"encoding/base64"
	"fmt"
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/internal/vaultManager"
	"github.com/quan-to/chevron/pkg/interfaces"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	"github.com/quan-to/slog"
)

type secretsManager struct {
	sync.Mutex
	encryptedPasswords   map[string]string
	gpg                  interfaces.PGPInterface
	masterKeyFingerPrint string
	amIUseless           bool
	log                  slog.Instance
}

// MakeSecretsManager creates an instance of the backend secrets manager
func MakeSecretsManager(log slog.Instance) interfaces.SMInterface {
	if log == nil {
		log = slog.Scope("SM")
	} else {
		log = log.SubScope("SM")
	}

	ctx := context.Background()

	var kb interfaces.Backend

	if config.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, "__master__")
	} else {
		kb = keybackend.MakeSaveToDiskBackend(log, path.Dir(config.MasterGPGKeyPath), "__master__")
	}

	var sm = &secretsManager{
		amIUseless:         false,
		encryptedPasswords: map[string]string{},
		log:                log,
	}
	masterKeyBytes, err := ioutil.ReadFile(config.MasterGPGKeyPath)

	if err != nil {
		sm.log.Warn("Error loading master key from %s: %s", config.MasterGPGKeyPath, err)
		sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
		sm.amIUseless = true
		return sm
	}

	if config.MasterGPGKeyBase64Encoded {
		masterKeyBytes, err = base64.StdEncoding.DecodeString(string(masterKeyBytes))
		if err != nil {
			sm.log.Warn("Error loading master key from %s: %s", config.MasterGPGKeyPath, err)
			sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
			sm.amIUseless = true
			return sm
		}
	}

	masterKeyFp, err := tools.GetFingerPrintFromKey(string(masterKeyBytes))

	if err != nil {
		sm.log.Warn("Error loading master key from %s: %s", config.MasterGPGKeyPath, err)
		sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
		sm.amIUseless = true
		return sm
	}

	sm.log.Info("Master Key FingerPrint: %s", masterKeyFp)

	sm.masterKeyFingerPrint = masterKeyFp

	sm.gpg = MakePGPManagerWithKRM(log, kb, MakeKeyRingManager(log))
	sm.gpg.SetKeysBase64Encoded(config.MasterGPGKeyBase64Encoded)

	err, n := sm.gpg.LoadKey(ctx, string(masterKeyBytes))

	if err != nil {
		sm.log.Fatal("Error loading private master key: %s", err)
	}

	if n == 0 {
		sm.log.Fatal("The specified key doesnt have any private keys inside.")
	}

	sm.gpg.LoadKeys(ctx)

	masterKeyPassBytes, err := ioutil.ReadFile(config.MasterGPGKeyPasswordPath)

	if err != nil {
		sm.log.Fatal("Error loading key password from %s: %s", config.MasterGPGKeyPasswordPath, err)
	}

	if config.MasterGPGKeyBase64Encoded { // If key is encoded, the password should be to
		masterKeyPassBytes, err = base64.StdEncoding.DecodeString(string(masterKeyPassBytes))
		if err != nil {
			sm.log.Fatal("Error decoding key password from %s: %s", config.MasterGPGKeyPasswordPath, err)
		}
	}

	err = sm.gpg.UnlockKey(ctx, masterKeyFp, strings.Trim(string(masterKeyPassBytes), "\n\r"))

	if err != nil {
		sm.log.Fatal("Error unlocking master key: %s", err)
	}

	err = sm.gpg.SaveKey(masterKeyFp, string(masterKeyBytes), string(masterKeyPassBytes))

	if err != nil {
		sm.log.Fatal("Error saving master key to default backend: %s", err)
	}

	return sm
}

// PutKeyPassword stores the password for the specified key fingerprint in the key backend encrypted with the master key
func (sm *secretsManager) PutKeyPassword(ctx context.Context, fingerprint, password string) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PutKeyPassword(%s, ---)", fingerprint)
	if sm.amIUseless {
		sm.log.Warn("Not saving password. Master Key not loaded")
		return
	}

	sm.Lock()
	defer sm.Unlock()

	sm.log.Info("Saving password for key %s", fingerprint)

	filename := fmt.Sprintf("key-password-utf8-%s.txt", fingerprint)

	encPass, err := sm.gpg.Encrypt(ctx, filename, sm.masterKeyFingerPrint, []byte(password), config.SMEncryptedDataOnly)

	if err != nil {
		sm.log.Error("Error saving key %s password: %s", fingerprint, err)
		return
	}

	sm.encryptedPasswords[fingerprint] = encPass
}

// PutEncryptedPassword stores in memory a master key encrypted password for the specified fingerprint
func (sm *secretsManager) PutEncryptedPassword(ctx context.Context, fingerprint, encryptedPassword string) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PutEncryptedPassword(%s, ---)", fingerprint)
	if sm.amIUseless {
		log.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	defer sm.Unlock()

	sm.encryptedPasswords[fingerprint] = encryptedPassword
}

// GetPasswords returns a list of master key encrypted passwords stored in memory
func (sm *secretsManager) GetPasswords(ctx context.Context) map[string]string {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("GetPasswords()")
	pss := make(map[string]string) // Force copy

	for fp, pass := range sm.encryptedPasswords {
		pss[fp] = pass
	}

	return pss
}

// UnlockLocalKeys unlocks the local private keys using memory stored master key encrypted passwords
func (sm *secretsManager) UnlockLocalKeys(ctx context.Context, gpg interfaces.PGPInterface) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("UnlockLocalKeys(---)")
	if sm.amIUseless {
		log.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	passwords := sm.GetPasswords(ctx)
	sm.Unlock()

	for fp, pass := range passwords {
		if !gpg.IsKeyLocked(fp) {
			continue
		}

		log.Info("Unlocking key %s", fp)
		g, err := sm.gpg.Decrypt(ctx, pass, config.SMEncryptedDataOnly)

		if err != nil {
			log.Error("Error decrypting password for key %s: %s", fp, err)
			continue
		}

		pass, err := base64.StdEncoding.DecodeString(g.Base64Data)

		if err != nil {
			// Shouldn't happen
			log.Error("Error decoding decrypted data: %s", err)
		}

		err = gpg.UnlockKey(ctx, fp, string(pass))
		if err != nil {
			log.Error("Error unlocking key %s: %s", fp, err)
		}
	}
}

// GetMasterKeyFingerPrint returns the fingerprint of the master key
func (sm *secretsManager) GetMasterKeyFingerPrint(ctx context.Context) string {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("GetMasterKeyFingerPrint()")
	return sm.masterKeyFingerPrint
}
