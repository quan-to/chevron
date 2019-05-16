package keymagic

import (
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/rstest"
	"testing"
)

func TestPutKeyPassword(t *testing.T) {
	sm.PutKeyPassword(rstest.TestKeyFingerprint, rstest.TestKeyFingerprint)
	if len(sm.encryptedPasswords[rstest.TestKeyFingerprint]) == 0 {
		t.Errorf("Expected encrypted password in keystore, got empty")
		t.FailNow()
	}

	if sm.encryptedPasswords[rstest.TestKeyFingerprint] == rstest.TestKeyFingerprint {
		t.Errorf("BIG FUCKING MISTAKE ERROR: Passwords are in plaintext!!!!!!!")
		t.FailNow()
	}

	dec, err := pgpMan.Decrypt(sm.encryptedPasswords[rstest.TestKeyFingerprint], remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Got error decrypting password: %s", err)
		t.FailNow()
	}

	bytePass, err := base64.StdEncoding.DecodeString(dec.Base64Data)

	if err != nil {
		t.Errorf("Got error unbase64 decrypted password: %s", err)
		t.FailNow()
	}

	if string(bytePass) != rstest.TestKeyFingerprint {
		t.Errorf("Expected stored password to be %s but got %s", rstest.TestKeyFingerprint, string(bytePass))
	}
}

func TestPutEncryptedPassword(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", rstest.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(rstest.TestKeyFingerprint), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(rstest.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[rstest.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}
}

func TestGetPasswords(t *testing.T) {
	sm.PutKeyPassword(rstest.TestKeyFingerprint, rstest.TestKeyFingerprint)

	passwords := sm.GetPasswords()

	if passwords[rstest.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", rstest.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", rstest.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(rstest.TestKeyFingerprint), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(rstest.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[rstest.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}

	sm.UnlockLocalKeys(pgpMan)
}
