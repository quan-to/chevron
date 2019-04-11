package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
func TestGetPasswords(t *testing.T) {
	sm.PutKeyPassword(remote_signer.TestKeyFingerprint, remote_signer.TestKeyPassword)

	passwords := sm.GetPasswords()

	if passwords[remote_signer.TestKeyFingerprint] == "" {
		t.Errorf("Expected %s key password to be in password list.", remote_signer.TestKeyFingerprint)
	}
}

func TestUnlockLocalKeys(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", remote_signer.TestKeyFingerprint)

	encPass, err := sm.gpg.Encrypt(filename, sm.masterKeyFingerPrint, []byte(remote_signer.TestKeyPassword), smEncryptedDataOnly)

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
*/

func TestGetUnlockPasswords(t *testing.T) {
	filename := fmt.Sprintf("key-password-utf8-%s.txt", remote_signer.TestKeyFingerprint)

	encPass, err := gpg.Encrypt(filename, sm.GetMasterKeyFingerPrint(), []byte(remote_signer.TestKeyPassword), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	sm.PutEncryptedPassword(remote_signer.TestKeyFingerprint, encPass)

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
		if k == remote_signer.TestKeyFingerprint {
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
	filename := fmt.Sprintf("key-password-utf8-%s.txt", remote_signer.TestKeyFingerprint)

	encPass, err := gpg.Encrypt(filename, sm.GetMasterKeyFingerPrint(), []byte(remote_signer.TestKeyPassword), remote_signer.SMEncryptedDataOnly)

	if err != nil {
		t.Errorf("Error saving password: %s", err)
		t.FailNow()
	}

	payload := map[string]string{
		remote_signer.TestKeyFingerprint: encPass,
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

	passwords := sm.GetPasswords()
	if passwords[remote_signer.TestKeyFingerprint] != encPass {
		t.Errorf("Expected key %s to be in the password list.", remote_signer.TestKeyFingerprint)
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
