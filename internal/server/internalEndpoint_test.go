package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	remote_signer "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/testdata"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
func TestGetPasswords(t *testing.T) {
	sm.PutKeyPassword(testdata.TestKeyFingerprint, testdata.TestKeyFingerprint)

	passwords := sm.GetPasswords()

	if passwords[testdata.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", testdata.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", testdata.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(testdata.TestKeyFingerprint), smEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(testdata.TestKeyFingerprint, encPass)

	if sm.encryptedPasswords[testdata.TestKeyFingerprint] != encPass {
		t.Errorf("Expected stored encrypted password to be set.")
	}

	sm.UnlockLocalKeys(pgpMan)
}
*/

func TestGetUnlockPasswords(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", testdata.TestKeyFingerprint)

	encPass, err := gpg.Encrypt(ctx, filename, sm.GetMasterKeyFingerPrint(ctx), []byte(testdata.TestKeyPassword), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(ctx, testdata.TestKeyFingerprint, encPass)

	req, err := http.NewRequest("GET", "/__internal/__getUnlockPasswords", nil)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	var data map[string]string

	err = json.Unmarshal(d, &data)

	errorDie(err, t)

	found := false

	for k, v := range data {
		if k == testdata.TestKeyFingerprint {
			found = true
			if v != encPass {
				t.Errorf("The encrypted password is not the expected one.")
			}
			break
		}
	}

	if !found {
		t.Errorf("The added password was not found at the password list")
	}
}

func TestPostUnlockPassword(t *testing.T) {
	ctx := context.Background()
	filename := fmt.Sprintf("key-password-utf8-%s.txt", testdata.TestKeyFingerprint)

	encPass, err := gpg.Encrypt(ctx, filename, sm.GetMasterKeyFingerPrint(ctx), []byte(testdata.TestKeyFingerprint), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	payload := map[string]string{
		testdata.TestKeyFingerprint: encPass,
	}

	d, _ := json.Marshal(payload)

	r := bytes.NewReader(d)

	req, err := http.NewRequest("POST", "/__internal/__postEncryptedPasswords", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}

	passwords := sm.GetPasswords(ctx)
	if passwords[testdata.TestKeyFingerprint] != encPass {
		t.Errorf("Expected key %s to be in the password list.", testdata.TestKeyFingerprint)
	}
}

func TestTriggerKeyUnlock(t *testing.T) {
	req, err := http.NewRequest("POST", "/__internal/__triggerKeyUnlock", nil)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}

	// TODO: Check if the key was really unlocked
}
