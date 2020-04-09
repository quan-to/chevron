package keymagic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/testdata"
	"testing"
)

func TestPutKeyPassword(t *testing.T) {
	ctx := context.Background()
	sm.PutKeyPassword(ctx, testdata.TestKeyFingerprint, testdata.TestKeyFingerprint)
	if len(sm.encryptedPasswords[testdata.TestKeyFingerprint]) == 0 {
		t.Errorf("Expected encrypted password in keystore, got empty")
		t.FailNow()
	}

	if sm.encryptedPasswords[testdata.TestKeyFingerprint] == testdata.TestKeyFingerprint {
		t.Errorf("BIG FUCKING MISTAKE ERROR: Passwords are in plaintext!!!!!!!")
		t.FailNow()
	}

	dec, err := pgpMan.Decrypt(ctx, sm.encryptedPasswords[testdata.TestKeyFingerprint], config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Got error decrypting password: %s", err)
		t.FailNow()
	}

	bytePass, err := base64.StdEncoding.DecodeString(dec.Base64Data)

	if err != nil {
		t.Errorf("Got error unbase64 decrypted password: %s", err)
		t.FailNow()
	}

	if string(bytePass) != testdata.TestKeyFingerprint {
		t.Errorf("Expected stored password to be %s but got %s", testdata.TestKeyFingerprint, string(bytePass))
	}
}

func TestPutEncryptedPassword(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", testdata.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(ctx, filename, sm.masterKeyFingerPrint, []byte(testdata.TestKeyFingerprint), config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(ctx, testdata.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[testdata.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}
}

func TestGetPasswords(t *testing.T) {
	ctx := context.Background()
	sm.PutKeyPassword(ctx, testdata.TestKeyFingerprint, testdata.TestKeyFingerprint)

	passwords := sm.GetPasswords(ctx)

	if passwords[testdata.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", testdata.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", testdata.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(ctx, filename, sm.masterKeyFingerPrint, []byte(testdata.TestKeyFingerprint), config.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(ctx, testdata.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[testdata.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}

	sm.UnlockLocalKeys(ctx, pgpMan)
}
