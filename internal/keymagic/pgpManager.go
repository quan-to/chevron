package keymagic

import (
	"bufio"
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/interfaces"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/quan-to/chevron/pkg/openpgp"
	"github.com/quan-to/chevron/pkg/openpgp/armor"
	"github.com/quan-to/chevron/pkg/openpgp/packet"
	"github.com/quan-to/slog"

	// Include MD5 hashing algorithm by default
	// skipcq: GSC-G501
	_ "crypto/md5"
	// Include SHA1 hashing algorithm by default
	// skipcq: GSC-G505
	_ "crypto/sha1"
	// Include SHA256 hashing algorithm by default
	_ "crypto/sha256"
	// Include SHA512 hashing algorithm by default
	_ "crypto/sha512"
	// Include RIPEMD160 hashing algorithm by default
	// skipcq: SCC-SA1019
	_ "golang.org/x/crypto/ripemd160"
)

const MinKeyBits = 2048 // Should be safe until we have decent Quantum Computers

type pgpManager struct {
	sync.Mutex
	KeysBase64Encoded    bool
	keyIdentity          map[string][]*openpgp.Identity
	decryptedPrivateKeys map[string]*packet.PrivateKey
	entities             map[string]*openpgp.Entity
	fp8to16              map[string]string
	subKeyToKey          map[string]string
	krm                  interfaces.KeyRingManager
	kbkend               interfaces.StorageBackend
	log                  slog.Instance
}

// MakePGPManager creates a new PGPManager with the specified keyBackend, log and KeyRingManager
func MakePGPManager(log slog.Instance, keyBackend interfaces.StorageBackend, krm interfaces.KeyRingManager) interfaces.PGPManager {
	if log == nil {
		log = slog.Scope("PGPMan")
	} else {
		log = log.SubScope("PGPMan")
	}

	if keyBackend == nil {
		log.Fatal("No keyBackend specified")
	}

	return &pgpManager{
		kbkend:               keyBackend,
		KeysBase64Encoded:    config.KeysBase64Encoded,
		keyIdentity:          make(map[string][]*openpgp.Identity),
		decryptedPrivateKeys: make(map[string]*packet.PrivateKey),
		entities:             make(map[string]*openpgp.Entity),
		fp8to16:              make(map[string]string),
		subKeyToKey:          make(map[string]string),
		krm:                  krm,
		log:                  log,
	}
}

// MinKeyBits returns the minimum key bits allowed for generating PGP Keys
func (pm *pgpManager) MinKeyBits() int {
	return MinKeyBits
}

// LoadKeys loads the keys stored on the PGP Manager key backend
func (pm *pgpManager) LoadKeys(ctx context.Context) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	pm.log.DebugNote("LoadKeys()")
	pm.Lock()
	defer pm.Unlock()

	if config.OnDemandKeyLoad {
		log.Warn("On Demand Key load enabled. Skipping loading keys.")
	} else {
		log.Info("Loading keys from %s -> %s", pm.kbkend.Name(), pm.kbkend.Path())

		files, err := pm.kbkend.List()
		if err != nil {
			log.Error("Error listing keys: %s", err)
			return
		}

		keysLoaded := 0

		for _, file := range files {
			log.Info("Loading key %s", file)
			keyData, m, err := pm.kbkend.Read(file)
			if err != nil {
				log.Error("Error loading key %s: %s", file, err)
				continue
			}

			if pm.KeysBase64Encoded {
				b, err := base64.StdEncoding.DecodeString(keyData)
				if err != nil {
					log.Error("Error base64 decoding %s: %s", file, err)
					continue
				}
				keyData = string(b)
			}

			kl, err := pm.LoadKeyWithMetadata(ctx, keyData, m)
			if err != nil {
				log.Error("Error decoding key %s: %s", file, err)
				continue
			}

			keysLoaded += kl
		}

		log.Info("Loaded %d private keys.", keysLoaded)
	}
}

// LoadKeyWithMetadata loads a armored ascii key with the specified json metadata
func (pm *pgpManager) LoadKeyWithMetadata(ctx context.Context, armoredKey, metadata string) (int, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("LoadKeyWithMetadata(---, ---)")
	n, err := pm.LoadKey(ctx, armoredKey)

	if err != nil {
		return n, err
	}

	fp, err := tools.GetFingerPrintFromKey(armoredKey)
	if err != nil {
		log.Error("Cannot get fingerprint from key: %s", err)
		return n, nil
	}

	if metadata != "" {
		var meta map[string]string
		err = json.Unmarshal([]byte(metadata), &meta)
		if err != nil {
			log.Warn("Error decoding metadata: %s", err)
			return n, nil
		}

		if meta["password"] != "" {
			err = pm.unlockKey(ctx, fp, meta["password"])
			if err != nil {
				log.Error("Cannot unlock key %s using metadata: %s", fp, err)
				return n, nil
			}
			log.Debug("Key %s unlocked using metadata.", fp)
			return n, nil
		}
	}

	log.Debug("No metadata for key %s. Skipping unlock...", fp)

	return n, nil
}

// SetKeysBase64Encoded sets if keys should be stored in Base64 Encoded format
func (pm *pgpManager) SetKeysBase64Encoded(k bool) {
	pm.log.DebugNote("SetKeysBase64Encoded(%t)", k)
	pm.KeysBase64Encoded = k
}

// LoadKey loads a armored ascii key
func (pm *pgpManager) LoadKey(ctx context.Context, armoredKey string) (int, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("LoadKey(---)")
	keysLoaded := 0
	kr := strings.NewReader(armoredKey)
	keys, err := openpgp.ReadArmoredKeyRing(kr)
	if err != nil {
		return 0, err
	}

	for _, key := range keys {
		if key.PrimaryKey != nil { // Cache Public Key
			fp := tools.ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
			log.Info("Loaded public key %s", fp)
			pm.krm.AddKey(ctx, key, true) // Add sticky public keys
			ids := make([]*openpgp.Identity, 0)
			for _, v := range key.Identities {
				// Get only first
				c := *v // copy
				ids = append(ids, &c)
			}
			pm.keyIdentity[fp] = ids
			pm.fp8to16[fp[8:]] = fp
			pm.entities[fp] = key
		}
		if key.PrivateKey != nil {
			fp := tools.ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
			log.Info("Loaded private key %s", fp)

			for _, sub := range key.Subkeys {
				subKeyFp := tools.IssuerKeyIdToFP16(sub.PublicKey.KeyId)
				log.Info("	Loaded subkey %s for %s", subKeyFp, fp)
				pm.subKeyToKey[subKeyFp] = fp
			}

			pm.krm.AddKey(ctx, key, true) // Add sticky public keys

			keysLoaded++
		}
	}

	return keysLoaded, err
}

func (pm *pgpManager) sanitizeFingerprint(fp string) string {
	if len(fp) > 16 {
		fp = fp[len(fp)-16:]
	}
	if len(fp) == 8 {
		fp = pm.fp8to16[fp]
	}
	if len(fp) != 16 {
		//pm.log.Fatal("Cannot find key or invalid fingerprint: %s", fp)
		return ""
	}

	return fp
}

// FixFingerPrint fixes and trims the fingerprint to 16 Char Hex
func (pm *pgpManager) FixFingerPrint(fp string) string {
	pm.Lock()
	defer pm.Unlock()

	return pm.sanitizeFingerprint(fp)
}

// IsKeyLocked returns if the specified private key is currently locked inside the PGP Manager
func (pm *pgpManager) IsKeyLocked(fp string) bool {
	pm.Lock()
	defer pm.Unlock()

	fp = pm.sanitizeFingerprint(fp)
	return pm.decryptedPrivateKeys[fp] == nil
}

func (pm *pgpManager) unlockKey(ctx context.Context, fp, password string) error {
	fp = pm.sanitizeFingerprint(fp)
	_ = pm.LoadKeyFromKB(ctx, fp)

	ent := pm.entities[fp]

	if ent == nil {
		pm.log.Error("No such key with fingerprint %s", fp)
		return fmt.Errorf("no such key %s", fp)
	}

	pk := ent.PrivateKey

	if pk == nil {
		return errors.New(fmt.Sprintf("private key %s not found", fp))
	}

	vpk := *pk // Copy data, for safety (aka: not unlock key at encrypted keys list)

	err := vpk.Decrypt([]byte(password))

	if err != nil {
		return err
	}

	if config.AgentKeyFingerPrint == "" { // set default fingerprint
		pm.log.Warn("No Agent Key FingerPrint specified. Using %s", fp)
		config.AgentKeyFingerPrint = fp
	}

	if pm.decryptedPrivateKeys[fp] != nil {
		pm.log.Info("Key %s already unlocked.", fp)
		return nil
	}

	z := pm.entities[fp]

	for _, kz := range z.Subkeys {
		subkeyfp := tools.IssuerKeyIdToFP16(kz.PublicKey.KeyId)
		pm.log.Info("		Decrypting subkey %s from %s", subkeyfp, fp)
		err := kz.PrivateKey.Decrypt([]byte(password))
		if err != nil {
			return err
		}
		pm.decryptedPrivateKeys[subkeyfp] = kz.PrivateKey
		pm.log.Debug("		Creating virtual entity for subkey %s from %s", subkeyfp, fp)
		pm.entities[subkeyfp] = tools.CreateEntityFromKeys(fmt.Sprintf("Subkey for %s", fp), "", "", 0, kz.PublicKey, kz.PrivateKey)
	}

	pm.decryptedPrivateKeys[fp] = &vpk

	return nil
}

// UnlockKey unlocks the specified key with the specified password
func (pm *pgpManager) UnlockKey(ctx context.Context, fp, password string) error {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("UnlockKey(%s, ---)", fp)
	pm.Lock()
	defer pm.Unlock()

	return pm.unlockKey(ctx, fp, password)
}

func (pm *pgpManager) LoadKeyFromKB(ctx context.Context, fingerPrint string) error {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.Info("Loading key %s", fingerPrint)

	if pm.decryptedPrivateKeys[fingerPrint] != nil || pm.entities[fingerPrint] != nil {
		log.Warn("Public Key %s is already loaded", fingerPrint)
		return nil
	}

	keyData, m, err := pm.kbkend.Read(fingerPrint)
	if err != nil {
		return err
	}

	if pm.KeysBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(keyData)
		if err != nil {
			return err
		}
		keyData = string(b)
	}

	_, err = pm.LoadKeyWithMetadata(ctx, keyData, m)
	if err != nil {
		return err
	}

	return nil
}

// GetPrivateKeyInfo returns the information of the specified private key
func (pm *pgpManager) GetPrivateKeyInfo(ctx context.Context, fingerPrint string) *models.KeyInfo {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetPrivateKeyInfo(%s)", fingerPrint)
	for k, e := range pm.entities {
		v := e.PrivateKey
		if v == nil {
			continue
		}

		if tools.CompareFingerPrint(k, fingerPrint) {
			z, _ := v.BitLength()
			return &models.KeyInfo{
				FingerPrint:           k,
				Identifier:            tools.SimpleIdentitiesToString(pm.keyIdentity[k]),
				Bits:                  int(z),
				ContainsPrivateKey:    true,
				PrivateKeyIsDecrypted: pm.decryptedPrivateKeys[k] != nil,
			}
		}
	}

	return nil
}

// GetLoadedPrivateKeys returns the information of each loaded private key
func (pm *pgpManager) GetLoadedPrivateKeys(ctx context.Context) []models.KeyInfo {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetLoadedPrivateKeys()")
	keyInfos := make([]models.KeyInfo, 0)

	for k, e := range pm.entities {
		v := e.PrivateKey
		if v == nil {
			continue
		}

		z, _ := v.BitLength()
		keyInfo := models.KeyInfo{
			FingerPrint:           k,
			Identifier:            tools.SimpleIdentitiesToString(pm.keyIdentity[k]),
			Bits:                  int(z),
			ContainsPrivateKey:    true,
			PrivateKeyIsDecrypted: pm.decryptedPrivateKeys[k] != nil,
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	return keyInfos
}

// GetLoadedKeys returns the information for all keys in PGP Manager
func (pm *pgpManager) GetLoadedKeys() []models.KeyInfo {
	pm.log.DebugNote("GetLoadedKeys()")
	keyInfos := make([]models.KeyInfo, 0)

	for k, e := range pm.entities {
		z, _ := e.PrimaryKey.BitLength()
		keyInfo := models.KeyInfo{
			FingerPrint:           k,
			Identifier:            tools.SimpleIdentitiesToString(pm.keyIdentity[k]),
			Bits:                  int(z),
			ContainsPrivateKey:    e.PrivateKey != nil,
			PrivateKeyIsDecrypted: pm.decryptedPrivateKeys[k] != nil,
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	return keyInfos
}

// SaveKey saves the specified key in PGP Manager Key Backend
func (pm *pgpManager) SaveKey(fingerPrint, armoredData string, password interface{}) error {
	pm.log.DebugNote("SaveKey(%s, %s, ---)", fingerPrint, tools.TruncateFieldForDisplay(armoredData))
	filename := fmt.Sprintf("%s.key", fingerPrint)
	if pm.KeysBase64Encoded {
		filename = fmt.Sprintf("%s.b64", fingerPrint)
	}

	filePath := path.Join(config.PrivateKeyFolder, filename)

	pm.log.Info("Saving key at %s", filePath)

	data := []byte(armoredData)

	if pm.KeysBase64Encoded {
		pm.log.Debug("Base64 Encoding enabled. Encoding key.")
		data = []byte(base64.StdEncoding.EncodeToString(data))
	}
	metadataJson := ""
	if password != nil {
		metadata := map[string]string{}
		metadata["password"] = password.(string)
		mj, _ := json.Marshal(metadata)
		metadataJson = string(mj)
	}

	rd, rm, err := pm.kbkend.Read(fingerPrint)

	if rd == "" || rm == "" || rm != metadataJson || string(data) != rd || err != nil {
		return pm.kbkend.SaveWithMetadata(fingerPrint, string(data), metadataJson)
	}

	pm.log.Warn("Key %s already in KeyBackend. Skipping add.", fingerPrint)

	return nil
}

// DeleteKey removes the specified key from the memory and key backend
func (pm *pgpManager) DeleteKey(ctx context.Context, fingerPrint string) error {
	pm.log.DebugAwait("Deleting key %s from KeyBackend", fingerPrint)
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)

	pm.Lock()
	if _, ok := pm.decryptedPrivateKeys[fingerPrint]; ok {
		pm.log.Info("Erasing private key %s from memory", fingerPrint)
		delete(pm.decryptedPrivateKeys, fingerPrint)
	}
	pm.Unlock()

	_ = pm.krm.DeleteKey(ctx, fingerPrint)

	_, _, err := pm.kbkend.Read(fingerPrint)
	if err != nil {
		pm.log.ErrorDone("Error reading key %s from KeyBackend, key not exist", fingerPrint)
		return nil
	}

	err = pm.kbkend.Delete(fingerPrint)
	if err != nil {
		pm.log.ErrorDone("Error deleting key %s from KeyBackend", fingerPrint)
	}

	return nil
}

// SignData signs the specified data with a unlocked private key
func (pm *pgpManager) SignData(ctx context.Context, fingerPrint string, data []byte, hashAlgorithm crypto.Hash) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("SignData(%s, ---, %v)", fingerPrint, hashAlgorithm)
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)
	pm.Lock()
	pk := pm.decryptedPrivateKeys[fingerPrint]

	if pk == nil {
		pm.Unlock()
		log.Warn("Private key %s not loaded or decrypted. Trying to load from keybackend", fingerPrint)
		err := pm.LoadKeyFromKB(ctx, fingerPrint)
		if err != nil {
			return "", errors.New(fmt.Sprintf("key %s is not decrypt or not loaded", fingerPrint))
		}
		pm.Lock()
		pk = pm.decryptedPrivateKeys[fingerPrint]
	}

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

	c := &packet.Config{
		DefaultHash: hashAlgorithm,
	}

	err := openpgp.ArmoredDetachSign(bw, &ent, d, c)
	if err != nil {
		return "", err
	}
	err = bw.Flush()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// GetPublicKeyEntity returns the public key entity
func (pm *pgpManager) GetPublicKeyEntity(ctx context.Context, fingerPrint string) *openpgp.Entity {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetPublicKeyEntity(%s)", fingerPrint)
	pm.Lock()
	defer pm.Unlock()
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)

	ent := pm.entities[fingerPrint]

	if ent == nil {
		// Try fetch subkey
		subMaster := pm.subKeyToKey[fingerPrint]
		if len(subMaster) > 0 {
			ent = pm.entities[subMaster]
		} else {
			// Try PKS
			ent = pm.krm.GetKey(ctx, fingerPrint)
		}
	}

	if ent != nil {
		pm.entities[fingerPrint] = ent
	}

	return ent
}

// GetPublicKey returns the public key
func (pm *pgpManager) GetPublicKey(ctx context.Context, fingerPrint string) *packet.PublicKey {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetPublicKey(%s)", fingerPrint)
	var pubKey *packet.PublicKey
	pm.Lock()
	defer pm.Unlock()
	log.Debug("Sanitizing fingerprint %s", fingerPrint)
	fingerPrint = pm.sanitizeFingerprint(fingerPrint)
	log.Debug("Sanitized %s", fingerPrint)

	ent := pm.entities[fingerPrint]

	if ent == nil {
		log.Debug("Not found in local cache as direct fingerprint. Trying by subkey")
		// Try fetch subkey
		subMaster := pm.subKeyToKey[fingerPrint]
		if len(subMaster) > 0 {
			ent = pm.entities[subMaster]
			pubKey = ent.PrimaryKey
			log.Note("Found as master key %s", fingerPrint)
		} else {
			// Try PKS
			log.Await("Not found as subkey. Checking in KeyRingManager")
			ent = pm.krm.GetKey(ctx, fingerPrint)
			if ent == nil {
				log.WarnDone("Not found in KeyRingManager")
			} else {
				log.Success("Found in Key Ring Manager")
			}
		}
	}

	if ent != nil {
		pubKey = ent.PrimaryKey
		pm.entities[fingerPrint] = ent
	}

	return pubKey
}

func (pm *pgpManager) GetSubKeys(fingerPrint string, decrypted bool) openpgp.EntityList {
	pm.log.DebugNote("GetSubKeys(%s, %v)", fingerPrint, decrypted)
	list := make([]*openpgp.Entity, 0)
	for k, v := range pm.subKeyToKey {
		if v == fingerPrint {
			ent := *pm.entities[k]
			if decrypted && pm.decryptedPrivateKeys[k] != nil {
				ent.PrivateKey = pm.decryptedPrivateKeys[k]
			}
			list = append(list, &ent)
		}
	}
	return list
}

func (pm *pgpManager) GetKey(ctx context.Context, fingerPrint string) *openpgp.Entity {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetKey(%s)", fingerPrint)
	fingerPrint = pm.FixFingerPrint(fingerPrint)

	// Try directly
	_ = pm.LoadKeyFromKB(ctx, fingerPrint)
	decv := pm.entities[fingerPrint]
	if decv != nil {
		return decv
	}

	// Try subkeys
	subKeyMaster := pm.subKeyToKey[fingerPrint]
	if subKeyMaster != fingerPrint {
		return pm.GetKey(ctx, subKeyMaster)
	}

	return nil
}

// GetPrivate returns the private key entity list for a specified private key
func (pm *pgpManager) GetPrivate(ctx context.Context, fingerPrint string) openpgp.EntityList {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetPrivate(%s)", fingerPrint)
	var ent openpgp.Entity
	fingerPrint = pm.FixFingerPrint(fingerPrint)

	// Try directly
	_ = pm.LoadKeyFromKB(ctx, fingerPrint)
	decv := pm.decryptedPrivateKeys[fingerPrint]
	if decv != nil {
		ent = *pm.entities[fingerPrint]
		ent.PrivateKey = pm.decryptedPrivateKeys[fingerPrint]
		keys := pm.GetSubKeys(fingerPrint, true)
		keys = append(keys, &ent)
		return keys
	}

	// Try subkeys
	subKeyMaster := pm.subKeyToKey[fingerPrint]
	if subKeyMaster != fingerPrint {
		return pm.GetPrivate(ctx, subKeyMaster)
	}

	return nil
}

// GetPublicKeyASCII returns the public key in ASCII Armored format
func (pm *pgpManager) GetPublicKeyASCII(ctx context.Context, fingerPrint string) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.Note("GetPublicKeyASCII(%s)", fingerPrint)
	key := ""
	pubKey := pm.GetPublicKey(ctx, fingerPrint)

	if pubKey == nil {
		return "", fmt.Errorf("not found")
	}

	ent := pm.GetPublicKeyEntity(ctx, fingerPrint)

	if ent != nil { // Try get full entity first
		serializedEntity := bytes.NewBuffer(nil)
		err := ent.Serialize(serializedEntity)

		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)
		headers := map[string]string{
			"Version": "GnuPG v2",
			"Comment": "Generated by Chevron",
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

		key = buf.String()
	} else { // If not, get just the public key
		serializedEntity := bytes.NewBuffer(nil)
		err := pubKey.Serialize(serializedEntity)

		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)
		headers := map[string]string{
			"Version": "GnuPG v2",
			"Comment": "Generated by Chevron",
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

		key = buf.String()
	}

	return key, nil
}

// GetPublicKeyASCII returns the encrypted private key in ASCII Armored format changing it's password
func (pm *pgpManager) GetPrivateKeyASCIIReencrypt(ctx context.Context, fingerPrint, currentPassword, newPassword string) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetPrivateKeyASCII(%s, ---)", fingerPrint)
	key := ""
	ent := pm.GetKey(ctx, fingerPrint)

	if ent != nil && ent.PrivateKey != nil { // Try get full entity first
		// Decrypt / Encrypt to initialize Signer
		err := ent.PrivateKey.Decrypt([]byte(currentPassword))

		if err != nil {
			return "", err
		}

		_ = ent.PrivateKey.Encrypt([]byte(newPassword))

		serializedEntity := bytes.NewBuffer(nil)
		err = ent.SerializePrivate(serializedEntity, &packet.Config{
			DefaultHash: crypto.SHA512,
		})
		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)
		headers := map[string]string{
			"Version": "GnuPG v2",
			"Comment": "Generated by Chevron",
		}

		w, err := armor.Encode(buf, openpgp.PrivateKeyType, headers)
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

		key = buf.String()
	} else {
		return "", fmt.Errorf("cannot find private key for %s", fingerPrint)
	}

	return key, nil
}

// GetPublicKeyASCII returns the encrypted private key in ASCII Armored format
func (pm *pgpManager) GetPrivateKeyASCII(ctx context.Context, fingerPrint, password string) (string, error) {
	return pm.GetPrivateKeyASCIIReencrypt(ctx, fingerPrint, password, password)
}

// VerifySignatureStringData verifies signature of specified data in string format
func (pm *pgpManager) VerifySignatureStringData(ctx context.Context, data string, signature string) (bool, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("VerifySignatureStringData(---, %s)", tools.TruncateFieldForDisplay(signature))
	return pm.VerifySignature(ctx, []byte(data), signature)
}

// VerifySignatureStringData verifies signature of specified data
func (pm *pgpManager) VerifySignature(ctx context.Context, data []byte, signature string) (bool, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("VerifySignature(---, %s)", tools.TruncateFieldForDisplay(signature))
	var issuerKeyId uint64
	var publicKey *packet.PublicKey
	var fingerprint string

	signature = tools.SignatureFix(signature)
	b := bytes.NewReader([]byte(signature))
	block, err := armor.Decode(b)
	if err != nil {
		return false, err
	}

	if block.Type != openpgp.SignatureType {
		return false, errors.New("openpgp packet is not signature")
	}

	reader := packet.NewReader(block.Body)
	publicKey = nil
	foundSignatureFingerprints := make([]string, 0)

	for {
		pkt, err := reader.Next()

		if err != nil {
			if len(foundSignatureFingerprints) > 0 {
				break // We found signatures just not public keys
			} else {
				return false, err
			}
		}

		switch sig := pkt.(type) {
		case *packet.Signature:
			if sig.IssuerKeyId == nil {
				return false, errors.New("signature doesn't have an issuer")
			}
			issuerKeyId = *sig.IssuerKeyId
			fingerprint = tools.IssuerKeyIdToFP16(issuerKeyId)
			foundSignatureFingerprints = append(foundSignatureFingerprints, fingerprint)
		case *packet.SignatureV3:
			issuerKeyId = sig.IssuerKeyId
			fingerprint = tools.IssuerKeyIdToFP16(issuerKeyId)
			foundSignatureFingerprints = append(foundSignatureFingerprints, fingerprint)
		}

		if len(fingerprint) == 16 {
			publicKey = pm.GetPublicKey(ctx, fingerprint)
			if publicKey != nil {
				break
			}
		}
	}

	if publicKey == nil {
		return false, fmt.Errorf("cannot find public key for any of these signatures: %s", strings.Join(foundSignatureFingerprints, ", "))
	}

	keyRing := make(openpgp.EntityList, 1)
	keyRing[0] = pm.entities[fingerprint]

	dr := bytes.NewReader(data)
	sr := strings.NewReader(signature)

	_, err = openpgp.CheckArmoredDetachedSignature(keyRing, dr, sr)

	if err != nil {
		return false, err
	}

	return true, nil
}

// GenerateTestKey generates a private key for testing
// Bits: MinKeyBits
// Password: 1234
// Identity: *empty string*
func (pm *pgpManager) GenerateTestKey() (string, error) {
	pm.log.DebugNote("GenerateTestKey()")
	privateKey, err := rsa.GenerateKey(rand.Reader, MinKeyBits)

	if err != nil {
		return "", err
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	err = pgpPrivKey.Encrypt([]byte("1234"))

	if err != nil {
		return "", err
	}

	e := tools.CreateEntityFromKeys("", "", "", 0, pgpPubKey, pgpPrivKey)

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
		"Comment": "Generated by Chevron",
	}

	w, err := armor.Encode(buf, openpgp.PrivateKeyType, headers)
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

// GeneratePGPKey generates a new PGP Key with the specified information
func (pm *pgpManager) GeneratePGPKey(ctx context.Context, identifier, password string, numBits int) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GeneratePGPKey(%s, ---, %d)", identifier, numBits)
	if numBits < MinKeyBits {
		return "", errors.New(fmt.Sprintf("dont generate RSA keys with less than %d, its not safe. try use 3072 or higher", MinKeyBits))
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, numBits)

	if err != nil {
		return "", err
	}

	var cTimestamp = time.Now()

	pgpPubKey := packet.NewRSAPublicKey(cTimestamp, &privateKey.PublicKey)
	pgpPrivKey := packet.NewRSAPrivateKey(cTimestamp, privateKey)

	err = pgpPrivKey.Encrypt([]byte(password))

	if err != nil {
		return "", err
	}

	identifier, comment, email := tools.ExtractIdentifierFields(identifier)

	if packet.HasInvalidCharacters(identifier) || packet.HasInvalidCharacters(comment) || packet.HasInvalidCharacters(email) {
		return "", fmt.Errorf("the identifier has invalid characters '(', ')', '<', '>'. If you're trying to use the full identifier format please check if its in the right format Name <email>")
	}

	e := tools.CreateEntityFromKeys(identifier, comment, email, 0, pgpPubKey, pgpPrivKey)

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
		"Comment": "Generated by Chevron",
	}

	w, err := armor.Encode(buf, openpgp.PrivateKeyType, headers)
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

// Encrypt encrypts data using the specified public key.
// Filename is a metadata from GPG
// dataOnly field specifies that it will encrypt as binary content instead ASCII Armored
func (pm *pgpManager) Encrypt(ctx context.Context, filename, fingerPrint string, data []byte, dataOnly bool) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("Encrypt(%s, %s, ---, %v)", filename, fingerPrint, dataOnly)
	var pubKey = pm.GetPublicKey(ctx, fingerPrint)

	if pubKey == nil {
		return "", fmt.Errorf("no public key for %s", fingerPrint)
	}
	fingerPrint = tools.ByteFingerPrint2FP16(pubKey.Fingerprint[:])
	var entity *openpgp.Entity

	pm.Lock()
	entity = pm.entities[fingerPrint]
	pm.Unlock()

	buf := bytes.NewBuffer(nil)

	hints := &openpgp.FileHints{
		FileName: filename,
		IsBinary: true,
		ModTime:  time.Now(),
	}

	c := &packet.Config{
		DefaultHash:            crypto.SHA512,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
	}

	closer, err := openpgp.Encrypt(buf, []*openpgp.Entity{entity}, nil, hints, c)

	if err != nil {
		return "", err
	}

	_, err = closer.Write(data)

	if err != nil {
		return "", err
	}

	err = closer.Close()
	if err != nil {
		return "", err
	}

	encData := buf.Bytes()

	if dataOnly {
		return base64.StdEncoding.EncodeToString(encData), nil
	}

	buf = bytes.NewBuffer(nil)
	headers := map[string]string{
		"Version": "GnuPG v2",
		"Comment": "Generated by Chevron",
	}

	w, err := armor.Encode(buf, "PGP MESSAGE", headers)
	if err != nil {
		return "", err
	}
	_, err = w.Write(encData)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Decrypt decrypts data using any available unlocked private key
func (pm *pgpManager) Decrypt(ctx context.Context, data string, dataOnly bool) (*models.GPGDecryptedData, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("Decrypt(%s, %v)", tools.TruncateFieldForDisplay(data), dataOnly)
	var err error
	var fps []string
	ret := &models.GPGDecryptedData{}

	if dataOnly {
		fps, err = tools.GetFingerPrintsFromEncryptedMessageRaw(data)
	} else {
		fps, err = tools.GetFingerPrintsFromEncryptedMessage(data)
	}

	if err != nil {
		return nil, err
	}

	if len(fps) == 0 {
		return nil, fmt.Errorf("no encrypted payloads found")
	}

	var decv *packet.PrivateKey
	var ent openpgp.Entity
	var subent *openpgp.Entity

	pm.LoadKeys(ctx)

	pm.Lock()
	for _, v := range fps {
		// Try directly
		_ = pm.LoadKeyFromKB(ctx, v)
		decv = pm.decryptedPrivateKeys[v]
		if decv != nil {
			ent = *pm.entities[v]
			break
		}

		// Try subkeys
		subKeyMaster := pm.subKeyToKey[v]
		if len(subKeyMaster) > 0 {
			_ = pm.LoadKeyFromKB(ctx, subKeyMaster)
			// Check if it is decrypted
			decv = pm.decryptedPrivateKeys[subKeyMaster]
			if decv != nil {
				ent = *pm.entities[subKeyMaster]
				subent = pm.entities[v]
				break
			}
		}
	}
	pm.Unlock()

	if decv == nil {
		return nil, fmt.Errorf("no unlocked key for decrypting packet")
	}

	keyRing := make(openpgp.EntityList, 1)
	ent.PrivateKey = decv
	keyRing[0] = &ent

	if subent != nil {
		keyRing[1] = subent
	}

	var rd io.Reader

	if dataOnly {
		d, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
		rd = bytes.NewReader(d)
	} else {
		if tools.IsASCIIArmored(data) {
			srd := strings.NewReader(data)
			p, err := armor.Decode(srd)
			if err != nil {
				return nil, err
			}

			rd = p.Body
		} else {
			rd = strings.NewReader(data)
		}
	}

	md, err := openpgp.ReadMessage(rd, keyRing, nil, nil)

	if err != nil {
		return nil, err
	}

	rawData, err := ioutil.ReadAll(md.LiteralData.Body)

	if err != nil {
		return nil, err
	}

	ret.FingerPrint = tools.IssuerKeyIdToFP16(ent.PrimaryKey.KeyId)
	ret.Base64Data = base64.StdEncoding.EncodeToString(rawData)
	ret.Filename = md.LiteralData.FileName

	return ret, nil
}

// GetCachedKeys returns all cached public keys in memory
func (pm *pgpManager) GetCachedKeys(ctx context.Context) []models.KeyInfo {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pm.log.Tag(requestID)
	log.DebugNote("GetCachedKeys()")
	return pm.krm.GetCachedKeys(ctx)
}
