package remote_signer

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

var pgpLog = SLog.Scope("PGPManager")

type PGPManager struct {
	sync.Mutex
	keyFolder            string
	keyIdentity			 map[string][]openpgp.Identity
	publicKeys           map[string]*packet.PublicKey
	privateKeys          map[string]*packet.PrivateKey
	decryptedPrivateKeys map[string]*packet.PrivateKey
	fp8to16              map[string]string
}

func MakePGPManager() *PGPManager {
	return &PGPManager{
		keyFolder:            PrivateKeyFolder,
		keyIdentity:		  make(map[string][]openpgp.Identity),
		publicKeys:           make(map[string]*packet.PublicKey),
		privateKeys:          make(map[string]*packet.PrivateKey),
		decryptedPrivateKeys: make(map[string]*packet.PrivateKey),
		fp8to16:              make(map[string]string),
	}
}

func (pm *PGPManager) LoadKeys() {
	pm.Lock()
	defer pm.Unlock()
	pgpLog.Info("Loading keys from %s", pm.keyFolder)
	files, err := ioutil.ReadDir(pm.keyFolder)
	if err != nil {
		pgpLog.Fatal("Error listing keys: %s", err)
	}
	keysLoaded := 0
	for _, file := range files {
		fileName := file.Name()
		filePath := path.Join(pm.keyFolder, fileName)
		if !file.IsDir() && fileName[:len(KeyPrefix)] == KeyPrefix {
			pgpLog.Info("Loading key %s", fileName)
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				pgpLog.Error("Error loading key %s: %s", fileName, err)
				continue
			}

			keyData := string(data)

			if KeysBase64Encoded {
				b, err := base64.StdEncoding.DecodeString(keyData)
				if err != nil {
					pgpLog.Error("Error base64 decoding %s: %s", fileName, err)
					continue
				}
				keyData = string(b)
			}

			kr := strings.NewReader(keyData)

			keys, err := openpgp.ReadArmoredKeyRing(kr)
			if err != nil {
				pgpLog.Error("Error decoding key %s: %s", fileName, err)
				continue
			}

			for _, key := range keys {
				if key.PrivateKey != nil {
					fp := ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
					ids := make([]openpgp.Identity, 0)
					for _, v := range key.Identities {
						// Get only first
						ids = append(ids, *v)
					}
					pm.keyIdentity[fp] = ids
					pm.privateKeys[fp] = key.PrivateKey
					pm.publicKeys[fp] = key.PrimaryKey
					pm.fp8to16[fp[8:]] = fp
					pgpLog.Info("Loaded private key %s", fp)
					keysLoaded++
				}
			}
		}
	}

	pgpLog.Info("Loaded %d private keys.", keysLoaded)
}

func (pm *PGPManager) sanitizeFingerprint(fp string) string {
	if len(fp) > 16 {
		fp = fp[len(fp)-16:]
	}
	if len(fp) == 8 {
		fp = pm.fp8to16[fp]
	}
	if len(fp) != 16 {
		pgpLog.Fatal("Cannot find key or invalid fingerprint: %s", fp)
	}

	return fp
}

func (pm *PGPManager) IsKeyLocked(fp string) bool {
	pm.Lock()
	defer pm.Unlock()

	fp = pm.sanitizeFingerprint(fp)
	return pm.decryptedPrivateKeys[fp] != nil
}

func (pm *PGPManager) UnlockKey(fp, password string) error {
	pm.Lock()
	defer pm.Unlock()

	fp = pm.sanitizeFingerprint(fp)

	if pm.decryptedPrivateKeys[fp] != nil {
		pgpLog.Info("Key %s already unlocked.", fp)
		return nil
	}

	pk := pm.privateKeys[fp]

	if pk == nil {
		return errors.New(fmt.Sprintf("private key %s not found", fp))
	}

	vpk := *pk // Copy data, for safety

	err := vpk.Decrypt([]byte(password))

	if err != nil {
		return err
	}

	pm.decryptedPrivateKeys[fp] = &vpk

	return nil
}


func (pm *PGPManager) GetLoadedPrivateKeys() []models.KeyInfo {
	keyInfos := make([]models.KeyInfo, 0)

	for k, v := range pm.privateKeys {
		z, _ := v.BitLength()
		identifier := ""
		if len(pm.keyIdentity[k]) > 0 {
			ki := pm.keyIdentity[k][0]
			identifier = ki.Name
		}
		keyInfo := models.KeyInfo{
			FingerPrint: k,
			Identifier: identifier,
			Bits: int(z),
			ContainsPrivateKey: true,
			PrivateKeyIsDecrypted: pm.decryptedPrivateKeys[k] != nil,
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	return keyInfos
}

func (pm *PGPManager) SavePrivateKey(fingerPrint, armoredData string) error {
	filename := fmt.Sprintf("%s.key", fingerPrint)
	if KeysBase64Encoded {
		filename = fmt.Sprintf("%s.b64", fingerPrint)
	}

	filePath := path.Join(PrivateKeyFolder, filename)

	pgpLog.Info("Saving private key at %s", filePath)

	data := []byte(armoredData)

	if KeysBase64Encoded {
		pgpLog.Debug("Base64 Encoding enabled. Encoding private key.")
		data = []byte(base64.StdEncoding.EncodeToString(data))
	}

	return ioutil.WriteFile(filePath, data, 0660)
}