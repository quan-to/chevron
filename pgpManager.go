package remote_signer

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

const minKeyBits = 2048 // Should be safe until we have decent Quantum Computers

var pgpLog = SLog.Scope("PGPManager")

type PGPManager struct {
	sync.Mutex
	keyFolder            string
	keyIdentity          map[string][]openpgp.Identity
	decryptedPrivateKeys map[string]*packet.PrivateKey
	entities             map[string]*openpgp.Entity
	fp8to16              map[string]string
}

func MakePGPManager() *PGPManager {
	return &PGPManager{
		keyFolder:            PrivateKeyFolder,
		keyIdentity:          make(map[string][]openpgp.Identity),
		decryptedPrivateKeys: make(map[string]*packet.PrivateKey),
		entities:             make(map[string]*openpgp.Entity),
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
		if !file.IsDir() && len(fileName) > len(KeyPrefix) && fileName[:len(KeyPrefix)] == KeyPrefix {
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

			err, kl := pm.LoadKey(keyData)
			if err != nil {
				pgpLog.Error("Error decoding key %s: %s", fileName, err)
				continue
			}

			keysLoaded += kl
		}
	}

	pgpLog.Info("Loaded %d private keys.", keysLoaded)
}

func (pm *PGPManager) LoadKey(armoredKey string) (error, int) {
	keysLoaded := 0
	kr := strings.NewReader(armoredKey)
	keys, err := openpgp.ReadArmoredKeyRing(kr)
	if err != nil {
		return err, 0
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
			pm.fp8to16[fp[8:]] = fp
			pm.entities[fp] = key
			pgpLog.Info("Loaded private key %s", fp)
			keysLoaded++
		}
	}

	return nil, keysLoaded
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

	pk := pm.entities[fp].PrivateKey

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

	for k, e := range pm.entities {
		v := e.PrivateKey
		z, _ := v.BitLength()
		identifier := ""
		if len(pm.keyIdentity[k]) > 0 {
			ki := pm.keyIdentity[k][0]
			identifier = ki.Name
		}
		keyInfo := models.KeyInfo{
			FingerPrint:           k,
			Identifier:            identifier,
			Bits:                  int(z),
			ContainsPrivateKey:    true,
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

func (pm *PGPManager) SignData(fingerPrint string, data []byte, hashAlgorithm crypto.Hash) (string, error) {
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)
	pm.Lock()
	pk := pm.decryptedPrivateKeys[fingerPrint]

	if pk == nil {
		pm.Unlock()
		return "", errors.New(fmt.Sprintf("key %s is not decrypt or not loaded", fingerPrint))
	}

	vpk := *pk
	ent := *pm.entities[fingerPrint]
	ent.PrivateKey = &vpk
	pm.Unlock()

	d := bytes.NewReader(data)

	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	config := &packet.Config{
		DefaultHash: hashAlgorithm,
	}

	err := openpgp.ArmoredDetachSign(bw, &ent, d, config)
	if err != nil {
		return "", err
	}
	err = bw.Flush()
	if err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}

func (pm *PGPManager) GetPublicKey(fingerPrint string) *packet.PublicKey {
	var pubKey *packet.PublicKey
	pm.Lock()
	defer pm.Unlock()
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)

	ent := pm.entities[fingerPrint]

	if ent == nil {
		// Try fetch SKS
		// TODO
	} else {
		pubKey = ent.PrimaryKey
	}

	return pubKey
}

func (pm *PGPManager) VerifySignatureStringData(data string, signature string) (bool, error) {
	return pm.VerifySignature([]byte(data), signature)
}

func (pm *PGPManager) VerifySignature(data []byte, signature string) (bool, error) {
	var issuerKeyId uint64
	var publicKey *packet.PublicKey
	var fingerPrint string

	signature = signatureFix(signature)
	b := bytes.NewReader([]byte(signature))
	block, err := armor.Decode(b)
	if err != nil {
		return false, err
	}

	if block.Type != openpgp.SignatureType {
		return false, errors.New("openpgp packet is not signature")
	}

	reader := packet.NewReader(block.Body)
	for {
		pkt, err := reader.Next()

		if err != nil {
			return false, err
		}

		switch sig := pkt.(type) {
		case *packet.Signature:
			if sig.IssuerKeyId == nil {
				return false, errors.New("signature doesn't have an issuer")
			}
			issuerKeyId = *sig.IssuerKeyId
			fingerPrint = IssuerKeyIdToFP16(issuerKeyId)
		case *packet.SignatureV3:
			issuerKeyId = sig.IssuerKeyId
			fingerPrint = IssuerKeyIdToFP16(issuerKeyId)
		default:
			return false, errors.New("non signature packet found")
		}

		if len(fingerPrint) == 16 {
			publicKey = pm.GetPublicKey(fingerPrint)
			if publicKey != nil {
				break
			}
		}
	}

	if publicKey == nil {
		return false, errors.New("cannot find public key to verify signature")
	}

	keyRing := make(openpgp.EntityList, 1)
	keyRing[0] = pm.entities[fingerPrint]

	dr := bytes.NewReader(data)
	sr := strings.NewReader(signature)

	_, err = openpgp.CheckArmoredDetachedSignature(keyRing, dr, sr)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (pm *PGPManager) GeneratePGPKey(identifier, password string, numBits int) (string, error) {
	if numBits < minKeyBits {
		return "", errors.New(fmt.Sprintf("dont generate RSA keys with less than %d, its not safe. try use 3072 or higher", minKeyBits))
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, numBits)

	if err != nil {
		return "", err
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	e := CreateEntityFromKeys(identifier, "", "", 0, pgpPubKey, pgpPrivKey)

	serializedEntity := bytes.NewBuffer(nil)
	err = e.SerializePrivate(serializedEntity, &packet.Config{
		DefaultHash: crypto.SHA512,
	})

	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	headers := map[string]string{
		"Version": "GnuPG v2",
		"Comment": "Generated by Quanto Remote Signer",
	}

	w, err := armor.Encode(buf, openpgp.PublicKeyType, headers)
	if err != nil {
		return "", err
	}
	_, err = w.Write(serializedEntity.Bytes())
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
