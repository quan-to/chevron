package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/models/HKP"
	"github.com/quan-to/chevron/test"
)

func TestHKPAdd(t *testing.T) {
	ctx := wrapContextWithDatabaseHandler(dbh, context.Background())
	config.PushVariables()
	defer config.PopVariables()

	req, err := http.NewRequest("POST", "/pks/add", nil)

	form := url.Values{}
	form.Add("keytext", test.TestPublicKey2)
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	if string(d) != "OK" {
		errorDie(fmt.Errorf("expected OK got %s", string(d)), t)
	}

	pubKey := gpg.GetPublicKey(ctx, test.TestPublicKey2FingerPrint)

	if pubKey == nil {
		errorDie(fmt.Errorf("expected to find key %s", test.TestPublicKey2FingerPrint), t)
	}

	if tools.IssuerKeyIdToFP16(pubKey.KeyId) != test.TestPublicKey2FingerPrint {
		errorDie(fmt.Errorf("expected key fingerprint to be %s got %s", test.TestPublicKey2FingerPrint, tools.IssuerKeyIdToFP16(pubKey.KeyId)), t)
	}

	// region Test Corrupted Form
	req, err = http.NewRequest("POST", "/pks/add", nil)
	errorDie(err, t)

	res = executeRequest(req)
	errObj, err := ReadErrorObject(res.Body)
	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}

func MakeHKPLookup(op, mr, nm, fingerprint, exact, search string) (output string, errObj *QuantoError.ErrorObject, err error) {
	req, errx := http.NewRequest("GET", "/pks/lookup", nil)
	err = errx
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("op", op)
	q.Add("mr", mr)
	q.Add("nm", nm)
	q.Add("fingerprint", fingerprint)
	q.Add("exact", exact)
	q.Add("search", search)
	req.URL.RawQuery = q.Encode()

	res := executeRequest(req)
	d, errx := ioutil.ReadAll(res.Body)
	err = errx
	if err != nil {
		return
	}

	// try decode error
	errObj = &QuantoError.ErrorObject{}
	err = json.Unmarshal(d, errObj)

	if err == nil {
		return // Decoder Error Object
	}
	errObj = nil
	err = nil

	output = string(d)

	return
}

func TestLookup(t *testing.T) {
	ctx := wrapContextWithDatabaseHandler(dbh, context.Background())

	// Ensure key is in SKS
	key, _ := gpg.GetPublicKeyASCII(ctx, test.TestKeyFingerprint)
	_ = keymagic.PKSAdd(ctx, key)

	// region Operation GET
	output, errObj, err := MakeHKPLookup(HKP.OperationGet, "true", "true", "on", "true", "0x"+test.TestKeyFingerprint)

	errorDie(err, t)
	if errObj != nil {
		errorDie(fmt.Errorf("expected error object to be nil got %v", errObj), t)
	}

	fp, err := tools.GetFingerPrintFromKey(output)

	errorDie(err, t)

	if fp != test.TestKeyFingerprint {
		errorDie(fmt.Errorf("expected public key fingerprint to be %s got %s", test.TestKeyFingerprint, fp), t)
	}

	_, errObj, err = MakeHKPLookup(HKP.OperationGet, "true", "true", "on", "true", test.TestKeyName)

	errorDie(err, t)
	if errObj != nil {
		errorDie(fmt.Errorf("expected error object to be nil got %v", errObj), t)
	}

	// TODO: Extended tests when full implementation of lookup is made
	// endregion
	// region Operation VIndex
	output, errObj, err = MakeHKPLookup(HKP.OperationVindex, "", "", "", "", "")
	errorDie(err, t)

	if errObj == nil {
		fmt.Printf("Output: %s\n", output)
		errorDie(fmt.Errorf("expected error object return, got nil"), t)
	}

	if errObj.ErrorCode != QuantoError.NotImplemented { // TODO FIX-ME when implemented
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.NotImplemented, errObj.ErrorCode), t)
	}
	// endregion
	// region Operation Index
	output, errObj, err = MakeHKPLookup(HKP.OperationIndex, "", "", "", "", "")
	errorDie(err, t)

	if errObj == nil {
		fmt.Printf("Output: %s\n", output)
		errorDie(fmt.Errorf("expected error object return, got nil"), t)
	}

	if errObj.ErrorCode != QuantoError.NotImplemented { // TODO FIX-ME when implemented
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.NotImplemented, errObj.ErrorCode), t)
	}
	// endregion
}
