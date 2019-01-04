package remote_signer

import (
	"encoding/base64"
	"fmt"
	"github.com/quan-to/remote-signer/SLog"
	"io/ioutil"
	"path"
	"sync"
)

const smEncryptedDataOnly = false

var smLog = SLog.Scope("SecretsManager")

type SecretsManager struct {
	sync.Mutex
	encryptedPasswords   map[string]string
	gpg                  *PGPManager
	masterKeyFingerPrint string
	amIUseless           bool
}

func MakeSecretsManager() *SecretsManager {
	var sm = &SecretsManager{
		amIUseless: false,
	}

	masterKeyBytes, err := ioutil.ReadFile(MasterGPGKeyPath)

	if err != nil {
		smLog.Error("Error loading master key from %s: %s", MasterGPGKeyPath, err)
		smLog.Error("I'm useless :(")
		sm.amIUseless = true
		return sm
	}

	masterKeyFp, err := GetFingerPrintFromKey(string(masterKeyBytes))

	if err != nil {
		smLog.Error("Error loading master key from %s: %s", MasterGPGKeyPath, err)
		smLog.Error("I'm useless :(")
		sm.amIUseless = true
		return sm
	}

	smLog.Info("Master Key FingerPrint: %s", masterKeyFp)

	sm.masterKeyFingerPrint = masterKeyFp

	sm.gpg = MakePGPManager()
	sm.gpg.keyFolder = path.Dir(MasterGPGKeyPath)
	sm.gpg.keysBase64Encoded = MasterGPGKeyBase64Encoded

	sm.gpg.LoadKeys()

	if MasterGPGKeyBase64Encoded {
		masterKeyBytes, err = base64.StdEncoding.DecodeString(string(masterKeyBytes))
		if err != nil {
			smLog.Fatal("Error decoding master key from base64: %s\nIs Master Key really base64 encoded?", err)
		}
	}

	masterKeyPassBytes, err := ioutil.ReadFile(MasterGPGKeyPasswordPath)

	if err != nil {
		smLog.Fatal("Error loading key password from %s: %s", MasterGPGKeyPasswordPath, err)
	}

	err = sm.gpg.UnlockKey(masterKeyFp, string(masterKeyPassBytes))

	if err != nil {
		smLog.Fatal("Error unlocking master key: %s", err)
	}

	return sm
}

func (sm *SecretsManager) PutKeyPassword(fingerPrint, password string) {
	if sm.amIUseless {
		smLog.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	defer sm.Unlock()

	smLog.Info("Saving password for key %s", fingerPrint)

	filename := fmt.Sprintf("key-password-utf8-%s.txt", fingerPrint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(password), smEncryptedDataOnly)

	if err != nil {
		smLog.Error("Error saving key %s password: %s", fingerPrint, err)
		return
	}

	sm.encryptedPasswords[fingerPrint] = encPass
}

func (sm *SecretsManager) PutEncryptedPassword(fingerPrint, encryptedPassword string) {
	if sm.amIUseless {
		smLog.Warn("Not saving password. Master Key not loaded")
	}

	sm.Lock()
	defer sm.Unlock()

	sm.encryptedPasswords[fingerPrint] = encryptedPassword
}

func (sm *SecretsManager) GetPasswords() map[string]string {
	pss := make(map[string]string) // Force copy

	for fp, pass := range sm.encryptedPasswords {
		pss[fp] = pass
	}

	return pss
}

func (sm *SecretsManager) UnlockLocalKeys(gpg *PGPManager) {
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
		g, err := sm.gpg.Decrypt(pass, smEncryptedDataOnly)

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
