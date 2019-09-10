// +build !js,!wasm

package keymagic

import (
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/vaultManager"
	"github.com/quan-to/slog"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

type SecretsManager struct {
	sync.Mutex
	encryptedPasswords   map[string]string
	gpg                  etc.PGPInterface
	masterKeyFingerPrint string
	amIUseless           bool
	log                  slog.Instance
}

func MakeSecretsManager(log slog.Instance) *SecretsManager {
	if log == nil {
		log = slog.Scope("SM")
	} else {
		log = log.SubScope("SM")
	}

	var kb keyBackend.Backend

	if remote_signer.VaultStorage {
		kb = vaultManager.MakeVaultManager(log, "__master__")
	} else {
		kb = keyBackend.MakeSaveToDiskBackend(log, path.Dir(remote_signer.MasterGPGKeyPath), "__master__")
	}

	var sm = &SecretsManager{
		amIUseless:         false,
		encryptedPasswords: map[string]string{},
		log:                log,
	}
	masterKeyBytes, err := ioutil.ReadFile(remote_signer.MasterGPGKeyPath)

	if err != nil {
		sm.log.Error("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
		sm.log.Error("I'm useless :(")
		sm.amIUseless = true
		return sm
	}

	if remote_signer.MasterGPGKeyBase64Encoded {
		masterKeyBytes, err = base64.StdEncoding.DecodeString(string(masterKeyBytes))
		if err != nil {
			sm.log.Error("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
			sm.log.Error("I'm useless :(")
			sm.amIUseless = true
			return sm
		}
	}

	masterKeyFp, err := remote_signer.GetFingerPrintFromKey(string(masterKeyBytes))

	if err != nil {
		sm.log.Error("Error loading master key from %s: %s", remote_signer.MasterGPGKeyPath, err)
		sm.log.Error("I'm useless :(")
		sm.amIUseless = true
		return sm
	}

	sm.log.Info("Master Key FingerPrint: %s", masterKeyFp)

	sm.masterKeyFingerPrint = masterKeyFp

	sm.gpg = MakePGPManagerWithKRM(log, kb, MakeKeyRingManager(log))
	sm.gpg.SetKeysBase64Encoded(remote_signer.MasterGPGKeyBase64Encoded)

	err, n := sm.gpg.LoadKey(string(masterKeyBytes))

	if err != nil {
		sm.log.Fatal("Error loading private master key: %s", err)
	}

	if n == 0 {
		sm.log.Fatal("The specified key doesnt have any private keys inside.")
	}

	sm.gpg.LoadKeys()

	masterKeyPassBytes, err := ioutil.ReadFile(remote_signer.MasterGPGKeyPasswordPath)

	if err != nil {
		sm.log.Fatal("Error loading key password from %s: %s", remote_signer.MasterGPGKeyPasswordPath, err)
	}

	if remote_signer.MasterGPGKeyBase64Encoded { // If key is encoded, the password should be to
		masterKeyPassBytes, err = base64.StdEncoding.DecodeString(string(masterKeyPassBytes))
		if err != nil {
			sm.log.Fatal("Error decoding key password from %s: %s", remote_signer.MasterGPGKeyPasswordPath, err)
		}
	}

	err = sm.gpg.UnlockKey(masterKeyFp, strings.Trim(string(masterKeyPassBytes), "\n\r"))

	if err != nil {
		sm.log.Fatal("Error unlocking master key: %s", err)
	}

	err = sm.gpg.SaveKey(masterKeyFp, string(masterKeyBytes), string(masterKeyPassBytes))

	if err != nil {
		sm.log.Fatal("Error saving master key to default backend: %s", err)
	}

	return sm
}

func (sm *SecretsManager) PutKeyPassword(fingerPrint, password string) {
	pksLog.DebugNote("PutKeyPassword(%s, ---)", fingerPrint)
	if sm.amIUseless {
		sm.log.Warn("Not saving password. Master Key not loaded")
		return
	}

	sm.Lock()
	defer sm.Unlock()

	sm.log.Info("Saving password for key %s", fingerPrint)

	filename := fmt.Sprintf("key-password-utf8-%s.txt", fingerPrint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(password), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		sm.log.Error("Error saving key %s password: %s", fingerPrint, err)
		return
	}

	sm.encryptedPasswords[fingerPrint] = encPass
}

func (sm *SecretsManager) PutEncryptedPassword(fingerPrint, encryptedPassword string) {
	pksLog.DebugNote("PutEncryptedPassword(%s, ---)", fingerPrint)
	if sm.amIUseless {
		sm.log.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	defer sm.Unlock()

	sm.encryptedPasswords[fingerPrint] = encryptedPassword
}

func (sm *SecretsManager) GetPasswords() map[string]string {
	pksLog.DebugNote("GetPasswords()")
	pss := make(map[string]string) // Force copy

	for fp, pass := range sm.encryptedPasswords {
		pss[fp] = pass
	}

	return pss
}

func (sm *SecretsManager) UnlockLocalKeys(gpg etc.PGPInterface) {
	pksLog.DebugNote("UnlockLocalKeys(---)")
	if sm.amIUseless {
		sm.log.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	passwords := sm.GetPasswords()
	sm.Unlock()

	for fp, pass := range passwords {
		if gpg.IsKeyLocked(fp) {
			continue
		}

		sm.log.Info("Unlocking key %s", fp)
		g, err := sm.gpg.Decrypt(pass, remote_signer.SMEncryptedDataOnly)

		if err != nil {
			sm.log.Error("Error decrypting password for key %s: %s", fp, err)
			continue
		}

		pass, err := base64.StdEncoding.DecodeString(g.Base64Data)

		if err != nil {
			// Shouldn't happen
			sm.log.Error("Error decoding decrypted data: %s", err)
		}

		err = gpg.UnlockKey(fp, string(pass))
		if err != nil {
			sm.log.Error("Error unlocking key %s: %s", fp, err)
		}
	}
}

func (sm *SecretsManager) GetMasterKeyFingerPrint() string {
	pksLog.DebugNote("GetMasterKeyFingerPrint()")
	return sm.masterKeyFingerPrint
}
