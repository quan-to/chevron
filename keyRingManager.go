package remote_signer

import (
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	"github.com/quan-to/remote-signer/openpgp"
	"sync"
)

var krmLog = SLog.Scope("KeyRingManager")

type KeyRingManager struct {
	sync.Mutex
	fingerPrints []string
	entities     map[string]*openpgp.Entity
	keyInfo      map[string]models.KeyInfo
}

func MakeKeyRingManager() *KeyRingManager {
	return &KeyRingManager{
		fingerPrints: make([]string, 0),
		entities:     make(map[string]*openpgp.Entity),
		keyInfo:      make(map[string]models.KeyInfo),
	}
}

func (krm *KeyRingManager) containsFp(fp string) bool {
	for _, v := range krm.fingerPrints {
		if v == fp {
			return true
		}
	}

	return false
}

func (krm *KeyRingManager) removeFp(fp string) {
	for i, v := range krm.fingerPrints {
		if v == fp {
			krm.fingerPrints = append(krm.fingerPrints[:i], krm.fingerPrints[i+1:]...)
			delete(krm.entities, fp)
			delete(krm.keyInfo, fp)
			return
		}
	}
}

func (krm *KeyRingManager) addFp(fp string) {
	krm.fingerPrints = append(krm.fingerPrints, fp)
}

func (krm *KeyRingManager) AddKey(key *openpgp.Entity, nonErasable bool) {
	krm.Lock()
	fp := ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
	if krm.containsFp(fp) {
		krmLog.Debug("Key %s already in keyring", fp)
		krm.Unlock()
		return
	}
	if !nonErasable {
		krm.addFp(fp)
	}

	krmLog.Info("Adding Public Key %s to the cache", fp)

	krm.entities[fp] = key

	keyBits, _ := key.PrimaryKey.BitLength()

	krm.keyInfo[fp] = models.KeyInfo{
		FingerPrint:           fp,
		Identifier:            SimpleIdentitiesToString(IdentityMapToArray(key.Identities)),
		Bits:                  int(keyBits),
		ContainsPrivateKey:    false,
		PrivateKeyIsDecrypted: false,
	}

	if len(krm.fingerPrints) > MaxKeyRingCache {
		lastFp := krm.fingerPrints[0]
		krmLog.Debug("	There are more cached keys than allowed. Removing first key %s", lastFp)
		krm.removeFp(lastFp)
	}

	krm.Unlock()

	for _, sub := range key.Subkeys {
		subfp := ByteFingerPrint2FP16(sub.PublicKey.Fingerprint[:])
		subE := CreateEntityForSubKey(fp, sub.PublicKey, sub.PrivateKey)
		krmLog.Debug("	Adding also subkey %s", subfp)
		krm.AddKey(subE, nonErasable)
	}
}

func (krm *KeyRingManager) GetCachedKeys() []models.KeyInfo {
	krm.Lock()
	defer krm.Unlock()
	arr := make([]models.KeyInfo, 0)

	for _, v := range krm.keyInfo {
		arr = append(arr, v)
	}

	return arr
}

func (krm *KeyRingManager) ContainsKey(fp string) bool {
	krm.Lock()
	defer krm.Unlock()

	return krm.entities[fp] != nil
}

func (krm *KeyRingManager) GetKey(fp string) *openpgp.Entity {
	krm.Lock()
	ent := krm.entities[fp]
	krm.Unlock()

	if ent != nil {
		return ent
	}

	// Try fetch SKS
	krmLog.Info("Key %s not found in local cache. Trying fetch KeyStore", fp)

	asciiArmored := PKSGetKey(fp)

	if len(asciiArmored) > 0 {
		k, err := ReadKeyToEntity(asciiArmored)
		if err != nil {
			krmLog.Error("Invalid key received from PKS! Error: %s", err)
			return nil
		}
		krmLog.Info("Key %s found in PKS. Adding to local cache", fp)
		ent = k
		krm.AddKey(k, false)
	}

	return ent
}
