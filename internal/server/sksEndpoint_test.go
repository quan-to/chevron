package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/test"
)

func TestSKSGetKey(t *testing.T) {
	// region Test Get Key
	req, err := http.NewRequest("GET", "/sks/getKey", nil)
	q := req.URL.Query()
	q.Add("fingerPrint", test.TestKeyFingerprint)
	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)
	errorDie(err, t)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)
	// endregion
	// region Test Unknown Key
	req, err = http.NewRequest("GET", "/sks/getKey", nil)
	q = req.URL.Query()
	q.Add("fingerPrint", "ABCDDEFGH")
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

func BaseTestSearch(name, value, endpoint string, t *testing.T) {
	config.PushVariables()
	defer config.PopVariables()

	config.EnableRethinkSKS = true

	// region Test Get Key
	req, err := http.NewRequest("GET", endpoint, nil)
	q := req.URL.Query()
	q.Add(name, value)

	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)
	errorDie(err, t)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	var gpgKeys []models.GPGKey

	err = json.Unmarshal(d, &gpgKeys)

	errorDie(err, t)

	if len(gpgKeys) == 0 {
		errorDie(fmt.Errorf("expected to find at least one key, got 0"), t)
	}
	// endregion
	// region Fetch Non Existent key
	req, err = http.NewRequest("GET", endpoint, nil)
	q = req.URL.Query()
	q.Add(name, "WOLOLO937091273092")
	q.Add("pageStart", "0")
	q.Add("pageEnd", "10")

	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res = executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)
	errorDie(err, t)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		fmt.Println(errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	err = json.Unmarshal(d, &gpgKeys)

	errorDie(err, t)

	if len(gpgKeys) != 0 {
		errorDie(fmt.Errorf("expected to find no keys. Got %d", len(gpgKeys)), t)
	}
	// endregion
	// region Fetch Empty Name
	req, err = http.NewRequest("GET", endpoint, nil)
	q = req.URL.Query()
	q.Add(name, "")
	q.Add("pageStart", "0")
	q.Add("pageEnd", "10")

	req.URL.RawQuery = q.Encode()

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}

func TestSKSSearchByName(t *testing.T) {
	BaseTestSearch("name", test.TestKeyName, "/sks/searchByName", t)
}

func TestSKSSearchByFingerPrint(t *testing.T) {
	BaseTestSearch("fingerPrint", test.TestKeyFingerprint, "/sks/searchByFingerPrint", t)
}

func TestSKSSearchByEmail(t *testing.T) {
	BaseTestSearch("email", test.TestKeyEmail, "/sks/searchByEmail", t)
}

func TestSKSSearch(t *testing.T) {
	BaseTestSearch("valueData", test.TestKeyEmail, "/sks/search", t)
	BaseTestSearch("valueData", test.TestKeyFingerprint, "/sks/search", t)
	BaseTestSearch("valueData", test.TestKeyName, "/sks/search", t)
}

func TestAddKey(t *testing.T) {
	ctx := context.Background()
	config.PushVariables()
	defer config.PopVariables()

	config.EnableRethinkSKS = true
	// region Test Add Key
	pubKey, _ := gpg.GetPublicKeyASCII(ctx, test.TestKeyFingerprint)

	payload := models.SKSAddKey{
		PublicKey: pubKey,
	}

	body, _ := json.Marshal(payload)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/sks/addKey", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	if string(d) != "OK" {
		errorDie(fmt.Errorf("expected OK got %s", string(d)), t)
	}
	// endregion
	// region Test Add Invalid Key
	payload.PublicKey = "huebrbrbrbrbr"
	body, _ = json.Marshal(payload)

	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/sks/addKey", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)
	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected errorCode %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}
