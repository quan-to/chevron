package chevronlib

import (
	"crypto"
	"fmt"
	"io"

	"github.com/quan-to/chevron/internal/keybackend"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/openpgp"
	"github.com/quan-to/chevron/pkg/openpgp/packet"
	"github.com/quan-to/slog"
)

// MakeSaveToDiskBackend creates an instance of a StorageBackend that
// saves the keys in the specified folder with the specified prefix
// log instance can be nil
func MakeSaveToDiskBackend(log slog.Instance, keysFolder, prefix string) interfaces.StorageBackend {
	return keybackend.MakeSaveToDiskBackend(log, keysFolder, prefix)
}

// MakeKeyRingManager creates a new instance of Key Ring Manager
// log instance can be nil
func MakeKeyRingManager(log slog.Instance) interfaces.KeyRingManager {
	return keymagic.MakeKeyRingManager(log, mem)
}

// MakePGPManager creates a new instance of PGP Operations Manager
// log instance can be nil
func MakePGPManager(log slog.Instance, storage interfaces.StorageBackend, keyRingManager interfaces.KeyRingManager) interfaces.PGPManager {
	return keymagic.MakePGPManager(log, storage, keyRingManager)
}

// GetStreamEncrypter returns a PGP Encrypter IO Writer
func GetStreamEncrypter(signerFingerprint string, encryptToFingerprints []string, output io.Writer, filehint *openpgp.FileHints, config *packet.Config) (enc io.WriteCloser, err error) {
	var signer *openpgp.Entity
	var pubkeys []*openpgp.Entity

	if len(signerFingerprint) != 0 {
		list := GetPrivateKeyEntity(signerFingerprint)
		for _, k := range list {
			if k.PrivateKey != nil {
				signer = k
				break
			}
		}
		if signer == nil {
			return nil, fmt.Errorf("cannot find private key with fingerprint %s", signerFingerprint)
		}
		if signer.PrivateKey.Encrypted {
			return nil, fmt.Errorf("found private key %s but it's encrypted", signerFingerprint)
		}
	}

	for _, pubKeyFp := range encryptToFingerprints {
		pub := GetPublicKeyEntity(pubKeyFp)
		if pub == nil {
			return nil, fmt.Errorf("cannot find public key with fingerprint %s", pubKeyFp)
		}
		pubkeys = append(pubkeys, pub)
	}

	if config == nil {
		config = &packet.Config{
			DefaultHash:            crypto.SHA3_512,
			DefaultCipher:          packet.CipherAES256,
			DefaultCompressionAlgo: packet.CompressionZIP,
			S2KCount:               65536 * 4, // Symmetric Key Generator expansion, should be > 65536
		}
	}
	return openpgp.Encrypt(output, pubkeys, signer, filehint, config)
}
