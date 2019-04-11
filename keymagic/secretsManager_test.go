package keymagic

import (
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron"
	"testing"
)

func TestPutKeyPassword(t *testing.T) {
	sm.PutKeyPassword(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)
	if len(sm.encryptedPasswords[remote_signer.TestKeyFingerprint]) == 0 {
		t.Errorf("Expected encrypted password in keystore, got empty")
		t.FailNow()
	}

	if sm.encryptedPasswords[remote_signer.TestKeyFingerprint] == remote_signer.TestKeyPassword {
		t.Errorf("BIG FUCKING MISTAKE ERROR: Passwords are in plaintext!!!!!!!")
		t.FailNow()
	}

	dec, err := pgpMan.Decrypt(sm.encryptedPasswords[remote_signer.TestKeyFingerprint], remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Got error decrypting password: %s", err)
		t.FailNow()
	}

	bytePass, err := base64.StdEncoding.DecodeString(dec.Base64Data)

	if err != nil {
		t.Errorf("Got error unbase64 decrypted password: %s", err)
		t.FailNow()
	}

	if string(bytePass) != remote_signer.TestKeyPassword {
		t.Errorf("Expected stored password to be %s but got %s", remote_signer.TestKeyPassword, string(bytePass))
	}
}

func TestPutEncryptedPassword(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", remote_signer.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(remote_signer.TestKeyPassword), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(remote_signer.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[remote_signer.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}
}

func TestGetPasswords(t *testing.T) {
	sm.PutKeyPassword(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	passwords := sm.GetPasswords()

	if passwords[remote_signer.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", remote_signer.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", remote_signer.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(remote_signer.TestKeyPassword), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(remote_signer.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[remote_signer.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}

	sm.UnlockLocalKeys(pgpMan)
}
