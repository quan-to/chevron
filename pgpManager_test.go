package remote_signer

import (
	"crypto"
	"github.com/quan-to/remote-signer/SLog"
	"testing"
)

const testKeyFingerprint = "0016A9CA870AFA59"
const testKeyPassword = "I think you will never guess"

const testSignatureData = "huebr for the win!"
const testSignatureSignature = `-----BEGIN PGP SIGNATURE-----

wsFcBAABCgAQBQJcJMPWCRAAFqnKhwr6WQAA3kwQAB6pxQkN+5yMt0LSkpIcjeOS
UPqcMabEQlkD2HQrzisXlUZgllqP4jYAjFLCeErt0uu598LXO6pNTw7MnFQSfgcJ
dJF2S05GwI4k00mMNzCTn7PbJe3d96QwjbTeanoMAjHhypZKi/StbtkFpIa+t9WI
zm+EE5trFdZoE1SMOr5j85afDecl0DsGHEkKdmJ2mLK4ja3uaxsijtLd8d7mdI+Y
LbI8UnpGyWMLkK8FpjBm+BaVeNicUvqkt/LO3LwslbKAViKpdL6Gu5x7x6Q+tAyO
PZ6P6DQKjuGJl8aSv0eoKQ1TQz6vasBZNsYlasU0fM6dXny9XIucUD5sTsUpbMhw
uO/xap6i3mBtFpzSfQCo/23KHeQajXS23Al56iUr85jlSQ9+JvJhZFrU9NQa+ypq
Xi/IxrqTTvttVurXAVME1m06JirpiuD8fDdQTTboekaqLg8rXQ5eKqW0pAMIqHvf
aq97YCqxH4F3T2EE77v6D9iLnbx/+7EGHoCehTMUYiAIAhlo93Xf/hnj40Hl/N18
gYr2Yd/IYVsAoGH6AHrIyUykXgsK6RXiBy0Sa7LN14TMCnQYzG2AUvXCDf184YAQ
1obsUVANy+qxH4lwMbEoznEsAU0ppqLchX1Ixdru5/SEgSV13Qv34rMEHCdVy4Oe
1Jcr1AyB3KmDhw76PaBh
=D//n
-----END PGP SIGNATURE-----`

var testData = []byte(testSignatureData)

var pgpMan *PGPManager

func init() {
	SLog.SetTestMode()

	PrivateKeyFolder = "."
	KeyPrefix = "testkey_"
	KeysBase64Encoded = false

	pgpMan = MakePGPManager()
	pgpMan.LoadKeys()

	err := pgpMan.UnlockKey(testKeyFingerprint, testKeyPassword)
	if err != nil {
		panic(err)
	}
}

// region Tests
func TestVerifySign(t *testing.T) {
	valid, err := pgpMan.VerifySignature(testData, testSignatureSignature)
	if err != nil || !valid {
		t.Errorf("Signature not valid or error found: %s", err)
	}

	valid, err = pgpMan.VerifySignatureStringData(testSignatureData, testSignatureSignature)
	if err != nil || !valid {
		t.Errorf("Signature not valid or error found: %s", err)
	}

	invalidTestData := []byte("huebr for the win!" + "makemeinvalid")

	valid, err = pgpMan.VerifySignature(invalidTestData, testSignatureSignature)

	if valid || err == nil {
		t.Error("A invalid test data passed to verify has been validated!")
	}
}

func TestSign(t *testing.T) {
	_, err := pgpMan.SignData(testKeyFingerprint, testData, crypto.SHA512)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateKey(t *testing.T) {
	key, err := pgpMan.GeneratePGPKey("HUE", testKeyPassword, minKeyBits)

	if err != nil {
		t.Error(err)
	}

	// Load key
	err, _ = pgpMan.LoadKey(key)
	if err != nil {
		t.Error(err)
	}

	fp := GetFingerPrintFromKey(key)

	t.Logf("Key Fingerprint: %s", fp)

	// Unlock Key
	err = pgpMan.UnlockKey(fp, testKeyPassword)
	if err != nil {
		t.Error(err)
	}

	// Try sign
	signature, err := pgpMan.SignData(fp, testData, crypto.SHA512)
	if err != nil {
		t.Error(err)
	}
	// Try verify
	valid, err := pgpMan.VerifySignature(testData, signature)
	if err != nil {
		t.Error(err)
	}
	if !valid {
		t.Error("Generated signature is not valid!")
	}
}

// endregion
// region Benchmarks
func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.SignData(testKeyFingerprint, testData, crypto.SHA512)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.VerifySignature(testData, testSignatureSignature)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkVerifySignatureStringData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.VerifySignatureStringData(testSignatureData, testSignatureSignature)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkKeyGenerate2048(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.GeneratePGPKey("", "123456789", 2048)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkKeyGenerate3072(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.GeneratePGPKey("", "123456789", 3072)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkKeyGenerate4096(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.GeneratePGPKey("", "123456789", 4096)
		if err != nil {
			b.Error(err)
		}
	}
}

// endregion
