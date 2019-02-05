package server

import (
	"encoding/json"
	"fmt"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/models"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestKREGetKey(t *testing.T) {
	// region Test Get Key
	req, err := http.NewRequest("GET", "/keyRing/getKey", nil)
	q := req.URL.Query()
	q.Add("fingerPrint", testKeyFingerprint)
	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)
	// endregion
	// region Test Inexistent Get Key
	req, err = http.NewRequest("GET", "/keyRing/getKey", nil)
	q = req.URL.Query()
	q.Add("fingerPrint", "WOLOLO")
	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.NotFound {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.NotFound, errObj.ErrorCode), t)
	}
	// endregion
}

func TestKREGetCachedKeys(t *testing.T) {
	// region Test Get Cached Keys
	req, err := http.NewRequest("GET", "/keyRing/cachedKeys", nil)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	var keyInfo []models.KeyInfo

	err = json.Unmarshal(d, &keyInfo)
	errorDie(err, t)

	found := false
	for _, v := range keyInfo {
		if v.FingerPrint == testKeyFingerprint {
			found = true
			break
		}
	}

	if !found {
		errorDie(fmt.Errorf("expected key %s to be in cached keys", testKeyFingerprint), t)
	}
	// endregion
}

func TestKREGetLoadedPrivateKeys(t *testing.T) {
	// region Test Get Loaded Private Keys
	req, err := http.NewRequest("GET", "/keyRing/privateKeys", nil)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	var keyInfo []models.KeyInfo

	err = json.Unmarshal(d, &keyInfo)
	errorDie(err, t)

	found := false
	for _, v := range keyInfo {
		if v.FingerPrint == testKeyFingerprint {
			found = true
			break
		}
	}

	if !found {
		errorDie(fmt.Errorf("expected key %s to be in private keys", testKeyFingerprint), t)
	}
	// endregion
}

func TestKREAddPrivateKey(t *testing.T) {
	key, err := gpg.GenerateTestKey()
	errorDie(err, t)

}
