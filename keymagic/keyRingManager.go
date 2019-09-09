package keymagic

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
	"github.com/quan-to/slog"
	"sync"
)

type KeyRingManager struct {
	sync.Mutex
	fingerPrints []string
	entities     map[string]*openpgp.Entity
	keyInfo      map[string]models.KeyInfo
	log          slog.Instance
}

func MakeKeyRingManager(log slog.Instance) *KeyRingManager {
	if log == nil {
		log = slog.Scope("KRM")
	} else {
		log = log.SubScope("KRM")
	}

	return &KeyRingManager{
		fingerPrints: make([]string, 0),
		entities:     make(map[string]*openpgp.Entity),
		keyInfo:      make(map[string]models.KeyInfo),
		log:          log,
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
	krm.log.Operation(slog.NOTE).Debug("AddKey(---, %v)", nonErasable)
	krm.Lock()
	fp := remote_signer.ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
	if krm.containsFp(fp) {
		krm.log.Debug("Key %s already in keyring", fp)
		krm.Unlock()
		return
	}
	if !nonErasable {
		if len(krm.fingerPrints)+1 > remote_signer.MaxKeyRingCache {
			lastFp := krm.fingerPrints[0]
			krm.log.Debug("	There are more cached keys than allowed. Removing first key %s", lastFp)
			krm.removeFp(lastFp)
		}
		krm.addFp(fp)
	}

	krm.log.Info("Adding Public Key %s to the cache", fp)

	krm.entities[fp] = key

	keyBits, _ := key.PrimaryKey.BitLength()

	krm.keyInfo[fp] = models.KeyInfo{
		FingerPrint:           fp,
		Identifier:            remote_signer.SimpleIdentitiesToString(remote_signer.IdentityMapToArray(key.Identities)),
		Bits:                  int(keyBits),
		ContainsPrivateKey:    false,
		PrivateKeyIsDecrypted: false,
	}

	krm.Unlock()

	for _, sub := range key.Subkeys {
		subfp := remote_signer.ByteFingerPrint2FP16(sub.PublicKey.Fingerprint[:])
		subE := remote_signer.CreateEntityForSubKey(fp, sub.PublicKey, sub.PrivateKey)
		krm.log.Debug("	Adding also subkey %s", subfp)
		krm.AddKey(subE, nonErasable)
	}
}

func (krm *KeyRingManager) GetCachedKeys() []models.KeyInfo {
	krm.log.Operation(slog.NOTE).Debug("GetCachedKeys()")
	krm.Lock()
	defer krm.Unlock()
	arr := make([]models.KeyInfo, 0)

	for _, v := range krm.keyInfo {
		arr = append(arr, v)
	}

	return arr
}

func (krm *KeyRingManager) ContainsKey(fp string) bool {
	krm.log.Operation(slog.NOTE).Debug("ContainsKey(%s)", fp)
	krm.Lock()
	defer krm.Unlock()

	return krm.entities[fp] != nil
}

func (krm *KeyRingManager) GetKey(fp string) *openpgp.Entity {
	krm.log.Operation(slog.NOTE).Debug("GetKey(%s)", fp)
	krm.Lock()
	ent := krm.entities[fp]
	krm.Unlock()

	if ent != nil {
		return ent
	}

	// Try fetch SKS
	krm.log.Await("Key %s not found in local cache. Trying fetch KeyStore", fp)

	asciiArmored, err := PKSGetKey(fp)

	if err != nil {
		krm.log.Error("Error fetching from KeyStore: %s", err)
		krm.log.Error(err)
		return nil
	}

	if len(asciiArmored) > 0 {
		k, err := remote_signer.ReadKeyToEntity(asciiArmored)
		if err != nil {
			krm.log.Error("Invalid key received from PKS! Error: %s", err)
			return nil
		}
		krm.log.Info("Key %s found in PKS. Adding to local cache", fp)
		ent = k
		krm.AddKey(k, false)
	}

	return ent
}

func (krm *KeyRingManager) GetFingerPrints() []string {
	krm.log.Operation(slog.NOTE).Debug("GetFingerPrints()")
	return krm.fingerPrints
}
