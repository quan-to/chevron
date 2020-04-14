package keymagic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/test"
	"testing"
)

func TestPutKeyPassword(t *testing.T) {
	ctx := context.Background()
	sm.PutKeyPassword(ctx, test.TestKeyFingerprint, test.TestKeyFingerprint)
	if len(sm.encryptedPasswords[test.TestKeyFingerprint]) == 0 {
		t.Errorf("Expected encrypted password in keystore, got empty")
		t.FailNow()
	}

	if sm.encryptedPasswords[test.TestKeyFingerprint] == test.TestKeyFingerprint {
		t.Errorf("BIG FUCKING MISTAKE ERROR: Passwords are in plaintext!!!!!!!")
		t.FailNow()
	}

	dec, err := pgpMan.Decrypt(ctx, sm.encryptedPasswords[test.TestKeyFingerprint], config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Got error decrypting password: %s", err)
		t.FailNow()
	}

	bytePass, err := base64.StdEncoding.DecodeString(dec.Base64Data)

	if err != nil {
		t.Errorf("Got error unbase64 decrypted password: %s", err)
		t.FailNow()
	}

	if string(bytePass) != test.TestKeyFingerprint {
		t.Errorf("Expected stored password to be %s but got %s", test.TestKeyFingerprint, string(bytePass))
	}
}

func TestPutEncryptedPassword(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", test.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(ctx, filename, sm.masterKeyFingerPrint, []byte(test.TestKeyFingerprint), config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(ctx, test.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[test.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}
}

func TestGetPasswords(t *testing.T) {
	ctx := context.Background()
	sm.PutKeyPassword(ctx, test.TestKeyFingerprint, test.TestKeyFingerprint)

	passwords := sm.GetPasswords(ctx)

	if passwords[test.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", test.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", test.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(ctx, filename, sm.masterKeyFingerPrint, []byte(test.TestKeyFingerprint), config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(ctx, test.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[test.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}

	sm.UnlockLocalKeys(ctx, pgpMan)
}
