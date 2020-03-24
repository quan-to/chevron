package keymagic

import (
	"context"
	"fmt"
	"sync"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/chevron/openpgp"
	"github.com/quan-to/slog"
)

type KeyRingManager struct {
	sync.Mutex
	fingerPrints []string
	entities     map[string]*openpgp.Entity
	keyInfo      map[string]models.KeyInfo
	log          slog.Instance
}

// MakeKeyRingManager creates a new instance of KeyRingManager
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

func (krm *KeyRingManager) AddKey(ctx context.Context, key *openpgp.Entity, nonErasable bool) {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("AddKey(---, %v)", nonErasable)
	krm.Lock()
	fp := remote_signer.ByteFingerPrint2FP16(key.PrimaryKey.Fingerprint[:])
	if krm.containsFp(fp) {
		log.Debug("Key %s already in keyring", fp)
		krm.Unlock()
		return
	}
	if !nonErasable {
		if len(krm.fingerPrints)+1 > remote_signer.MaxKeyRingCache {
			lastFp := krm.fingerPrints[0]
			log.Debug("	There are more cached keys than allowed. Removing first key %s", lastFp)
			krm.removeFp(lastFp)
		}
		krm.addFp(fp)
	}

	log.Info("Adding Public Key %s to the cache", fp)

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
		log.Debug("	Adding also subkey %s", subfp)
		krm.AddKey(ctx, subE, nonErasable)
	}
}

func (krm *KeyRingManager) DeleteKey(ctx context.Context, fp string) error {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("DeleteKey(%s)", fp)
	krm.Lock()
	if _, ok := krm.entities[fp]; ok {
		log.Info("Deleting key %s from memory", fp)
		delete(krm.entities, fp)
	}
	krm.Unlock()

	return fmt.Errorf("key %s not found", fp)
}

func (krm *KeyRingManager) GetCachedKeys(ctx context.Context) []models.KeyInfo {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("GetCachedKeys()")
	krm.Lock()
	defer krm.Unlock()
	arr := make([]models.KeyInfo, 0)

	for _, v := range krm.keyInfo {
		arr = append(arr, v)
	}

	return arr
}

func (krm *KeyRingManager) ContainsKey(ctx context.Context, fp string) bool {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("ContainsKey(%s)", fp)
	krm.Lock()
	defer krm.Unlock()

	return krm.entities[fp] != nil
}

func (krm *KeyRingManager) GetKey(ctx context.Context, fp string) *openpgp.Entity {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("GetKey(%s)", fp)
	krm.Lock()
	ent := krm.entities[fp]
	krm.Unlock()

	if ent != nil {
		return ent
	}

	// Try fetch SKS
	log.Await("Key %s not found in local cache. Trying fetch KeyStore", fp)

	asciiArmored, err := PKSGetKey(ctx, fp)

	if err != nil {
		log.Error("Error fetching from KeyStore: %s", err)
		log.Error(err)
		return nil
	}

	if len(asciiArmored) > 0 {
		k, err := remote_signer.ReadKeyToEntity(asciiArmored)
		if err != nil {
			log.Error("Invalid key received from PKS! Error: %s", err)
			return nil
		}
		log.Info("Key %s found in PKS. Adding to local cache", fp)
		ent = k
		krm.AddKey(ctx, k, false)
	}

	return ent
}

func (krm *KeyRingManager) GetFingerPrints(ctx context.Context) []string {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := krm.log.Tag(requestID)
	log.DebugNote("GetFingerPrints()")
	return krm.fingerPrints
}
