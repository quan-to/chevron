package keymagic

import (
	"encoding/base64"
	"fmt"
	remote_signer "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc"
	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/slog"
	"io/ioutil"
	"path"
	"sync"
)

var smLog = slog.Scope("secretsManager")

type secretsManager struct {
	sync.Mutex
	encryptedPasswords   map[string]string
	gpg                  etc.PGPInterface
	masterKeyFingerPrint string
	amIUseless           bool
}

func MakeSecretsManager(log slog.Instance) interfaces.SMInterface {
	var kb interfaces.Backend

	kb = keybackend.MakeSaveToDiskBackend(path.Dir(remote_signer.MasterGPGKeyPath), "__master__")

	var sm = &secretsManager{
		amIUseless:         false,
		encryptedPasswords: map[string]string{},
	}
	masterKeyBytes, err := ioutil.ReadFile(remote_signer.MasterGPGKeyPath)

	originalKeyBytes := masterKeyBytes

	if err != nil {
		sm.log.Warn("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
		sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
		sm.amIUseless = true
		return sm
	}

	if remote_signer.MasterGPGKeyBase64Encoded {
		masterKeyBytes, err = base64.StdEncoding.DecodeString(string(masterKeyBytes))
		if err != nil {
			smLog.Warn("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
			sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
			sm.amIUseless = true
			return sm
		}
	}

	masterKeyFp, err := remote_signer.GetFingerPrintFromKey(string(masterKeyBytes))

	if err != nil {
		smLog.Warn("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
		sm.log.Warn("Secrets Manager cannot load a master key. Cluster mode will not work.")
		sm.amIUseless = true
		return sm
	}

	smLog.Info("Master Key FingerPrint: %s", masterKeyFp)

	sm.masterKeyFingerPrint = masterKeyFp

	sm.gpg = MakePGPManagerWithKRM(kb, MakeKeyRingManager())
	sm.gpg.SetKeysBase64Encoded(remote_signer.MasterGPGKeyBase64Encoded)

	err, n := sm.gpg.LoadKey(string(masterKeyBytes))

	if err != nil {
		smLog.Fatal("Error loading private master key: %s", err)
	}

	if n == 0 {
		smLog.Fatal("The specified key doesnt have any private keys inside.")
	}

	sm.gpg.LoadKeys()

	masterKeyPassBytes, err := ioutil.ReadFile(remote_signer.MasterGPGKeyPasswordPath)

	if err != nil {
		smLog.Fatal("Error loading key password from %s: %s", remote_signer.MasterGPGKeyPasswordPath, err)
	}

	err = sm.gpg.UnlockKey(masterKeyFp, string(masterKeyPassBytes))

	if err != nil {
		smLog.Fatal("Error unlocking master key: %s", err)
	}

	err = sm.gpg.SaveKey(masterKeyFp, string(originalKeyBytes), string(masterKeyPassBytes))

	if err != nil {
		smLog.Fatal("Error saving master key to default backend: %s", err)
	}

	return sm
}

// PutKeyPassword stores the password for the specified key fingerprint in the key backend encrypted with the master key
func (sm *secretsManager) PutKeyPassword(fingerPrint, password string) {
	if sm.amIUseless {
		smLog.Warn("Not saving password. Master Key not loaded")
		return
	}

	sm.Lock()
	defer sm.Unlock()

	smLog.Info("Saving password for key %s", fingerPrint)

	filename := fmt.Sprintf("key-password-utf8-%s.txt", fingerPrint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(password), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		smLog.Error("Error saving key %s password: %s", fingerPrint, err)
		return
	}

	sm.encryptedPasswords[fingerPrint] = encPass
}

// PutEncryptedPassword stores in memory a master key encrypted password for the specified fingerprint
func (sm *secretsManager) PutEncryptedPassword(fingerPrint, encryptedPassword string) {
	if sm.amIUseless {
		smLog.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	defer sm.Unlock()

	sm.encryptedPasswords[fingerPrint] = encryptedPassword
}

// GetPasswords returns a list of master key encrypted passwords stored in memory
func (sm *secretsManager) GetPasswords() map[string]string {
	pss := make(map[string]string) // Force copy

	for fp, pass := range sm.encryptedPasswords {
		pss[fp] = pass
	}

	return pss
}

// UnlockLocalKeys unlocks the local private keys using memory stored master key encrypted passwords
func (sm *secretsManager) UnlockLocalKeys(gpg etc.PGPInterface) {
	if sm.amIUseless {
		smLog.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	passwords := sm.GetPasswords()
	sm.Unlock()

	for fp, pass := range passwords {
		if gpg.IsKeyLocked(fp) {
			continue
		}

		smLog.Info("Unlocking key %s", fp)
		g, err := sm.gpg.Decrypt(pass, remote_signer.SMEncryptedDataOnly)

		if err != nil {
			smLog.Error("Error decrypting password for key %s: %s", fp, err)
			continue
		}

		pass, err := base64.StdEncoding.DecodeString(g.Base64Data)

		if err != nil {
			// Shouldn't happen
			smLog.Error("Error decoding decrypted data: %s", err)
		}

		err = gpg.UnlockKey(fp, string(pass))
		if err != nil {
			smLog.Error("Error unlocking key %s: %s", fp, err)
		}
	}
}

// GetMasterKeyFingerPrint returns the fingerprint of the master key
func (sm *secretsManager) GetMasterKeyFingerPrint() string {
	return sm.masterKeyFingerPrint
}
