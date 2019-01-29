package keymagic

import (
	"crypto"
	"encoding/base64"
	"github.com/quan-to/remote-signer"
	"testing"
)

// region Tests
func TestVerifySign(t *testing.T) {
	valid, err := pgpMan.VerifySignature(testData, remote_signer.TestSignatureSignature)
	if err != nil || !valid {
		t.Errorf("Signature not valid or error found: %s", err)
	}

	valid, err = pgpMan.VerifySignatureStringData(remote_signer.TestSignatureData, remote_signer.TestSignatureSignature)
	if err != nil || !valid {
		t.Errorf("Signature not valid or error found: %s", err)
	}

	invalidTestData := []byte("huebr for the win!" + "makemeinvalid")

	valid, err = pgpMan.VerifySignature(invalidTestData, remote_signer.TestSignatureSignature)

	if valid || err == nil {
		t.Error("A invalid test data passed to verify has been validated!")
	}
}

func TestSign(t *testing.T) {
	_, err := pgpMan.SignData(remote_signer.TestKeyFingerprint, testData, crypto.SHA512)
	if err != nil {
		t.Error(err)
	}
}

func TestDecrypt(t *testing.T) {
	g, err := pgpMan.Decrypt(remote_signer.TestDecryptDataAscii, false)
	if err != nil {
		t.Error(err)
	}

	gd, err := base64.StdEncoding.DecodeString(g.Base64Data)
	if err != nil {
		t.Error(err)
	}

	if string(gd) != remote_signer.TestSignatureData {
		t.Errorf("Decrypted data does no match. Expected \"%s\" got \"%s\"", string(gd), remote_signer.TestSignatureData)
	}

	g, err = pgpMan.Decrypt(remote_signer.TestDecryptDataOnly, true)
	if err != nil {
		t.Error(err)
	}

	gd, err = base64.StdEncoding.DecodeString(g.Base64Data)
	if err != nil {
		t.Error(err)
	}

	if string(gd) != remote_signer.TestSignatureData {
		t.Errorf("Decrypted data does no match. Expected \"%s\" got \"%s\"", string(gd), remote_signer.TestSignatureData)
	}
}

func TestEncrypt(t *testing.T) {
	d, err := pgpMan.Encrypt("testing", remote_signer.TestKeyFingerprint, testData, false)

	if err != nil {
		t.Error(err)
	}

	// region Test Decrypt
	g, err := pgpMan.Decrypt(d, false)
	if err != nil {
		t.Error(err)
	}

	gd, err := base64.StdEncoding.DecodeString(g.Base64Data)
	if err != nil {
		t.Error(err)
	}

	if string(gd) != remote_signer.TestSignatureData {
		t.Errorf("Decrypted data does no match. Expected \"%s\" got \"%s\"", string(gd), remote_signer.TestSignatureData)
	}
	// endregion
	d, err = pgpMan.Encrypt("testing", remote_signer.TestKeyFingerprint, testData, true)

	if err != nil {
		t.Error(err)
	}

	// region Test Decrypt
	g, err = pgpMan.Decrypt(d, true)
	if err != nil {
		t.Error(err)
	}

	gd, err = base64.StdEncoding.DecodeString(g.Base64Data)
	if err != nil {
		t.Error(err)
	}

	if string(gd) != remote_signer.TestSignatureData {
		t.Errorf("Decrypted data does no match. Expected \"%s\" got \"%s\"", string(gd), remote_signer.TestSignatureData)
	}
	// endregion
}

func TestGenerateKey(t *testing.T) {
	key, err := pgpMan.GeneratePGPKey("HUE", remote_signer.TestKeyPassword, pgpMan.MinKeyBits())

	if err != nil {
		t.Error(err)
	}

	// Load key
	err, _ = pgpMan.LoadKey(key)
	if err != nil {
		t.Error(err)
	}

	fp, _ := remote_signer.GetFingerPrintFromKey(key)

	// Unlock Key
	err = pgpMan.UnlockKey(fp, remote_signer.TestKeyPassword)
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
		_, err := pgpMan.SignData(remote_signer.TestKeyFingerprint, testData, crypto.SHA512)
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkVerifySignature(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.VerifySignature(testData, remote_signer.TestSignatureSignature)
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkVerifySignatureStringData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.VerifySignatureStringData(remote_signer.TestSignatureData, remote_signer.TestSignatureSignature)
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkEncryptASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.Encrypt("", remote_signer.TestKeyFingerprint, testData, false)
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkEncryptDataOnly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pgpMan.Encrypt("", remote_signer.TestKeyFingerprint, testData, true)
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
