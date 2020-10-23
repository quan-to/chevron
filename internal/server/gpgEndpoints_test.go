package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/test"
	"io/ioutil"
	"net/http"
	"testing"
)

// region GPG Endpoint Tests
func TestGenerateKey(t *testing.T) {
	InvalidPayloadTest("/gpg/generateKey", t)

	ctx := context.Background()

	// region Test Generate Key
	genKeyBody := models.GPGGenerateKeyData{
		Identifier: "Test",
		Password:   "123456",
		Bits:       gpg.MinKeyBits(),
	}

	body, err := json.Marshal(genKeyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/generateKey", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	key := string(d)

	fingerPrint, err := tools.GetFingerPrintFromKey(key)

	errorDie(err, t)

	_, err = gpg.LoadKey(ctx, key)

	errorDie(err, t)

	err = gpg.UnlockKey(ctx, fingerPrint, genKeyBody.Password)

	errorDie(err, t)
	// endregion
	// region Test MinKeyBits
	genKeyBody.Bits = gpg.MinKeyBits() - 1
	body, err = json.Marshal(genKeyBody)
	errorDie(err, t)

	r = bytes.NewReader(body)
	req, err = http.NewRequest("POST", "/gpg/generateKey", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)
	if err != nil {
		errorDie(err, t)
	}

	if errObj.ErrorField != "Bits" {
		errorDie(fmt.Errorf("expected Bits as error field. Got %s", errObj.ErrorField), t)
	}

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s as error code. Got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
	// region Test No Password

	genKeyBody.Bits = gpg.MinKeyBits()
	genKeyBody.Password = ""
	body, err = json.Marshal(genKeyBody)
	errorDie(err, t)

	r = bytes.NewReader(body)
	req, err = http.NewRequest("POST", "/gpg/generateKey", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err = ReadErrorObject(res.Body)
	if err != nil {
		errorDie(err, t)
	}

	if errObj.ErrorField != "Password" {
		errorDie(fmt.Errorf("expected Bits as error field. Got %s", errObj.ErrorField), t)
	}

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s as error code. Got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}

func TestDecryptInvalidPayload(t *testing.T) {
	InvalidPayloadTest("/gpg/decrypt", t)

	// Invalid Decrypting Data
	decryptBody := models.GPGDecryptData{
		DataOnly:         true,
		AsciiArmoredData: "huehuebrbr",
	}

	body, err := json.Marshal(decryptBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/decrypt", r)

	errorDie(err, t)

	res := executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)
	if err != nil {
		errorDie(err, t)
	}

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s in ErrorCode. Got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
}

func TestEncryptData(t *testing.T) {
	InvalidPayloadTest("/gpg/encrypt", t)

	ctx := context.Background()

	encryptBody := models.GPGEncryptData{
		DataOnly:    true,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
		Filename:    "test-encrypt",
		FingerPrint: test.TestKeyFingerprint,
	}

	body, _ := json.Marshal(encryptBody)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/encrypt", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	encryptedData := string(d)

	data, err := gpg.Decrypt(ctx, encryptedData, true)

	errorDie(err, t)

	if data.Base64Data != encryptBody.Base64Data {
		t.Errorf("expected Base64Data %s got %s", encryptBody.Base64Data, data.Base64Data)
	}

	if data.Filename != encryptBody.Filename {
		t.Errorf("expected Filename %s got %s", encryptBody.Filename, data.Filename)
	}

	// Test Invalid Base64

	encryptBody.Base64Data = "ééééaiseh - - -12= '/x. huebrbrbrbré"
	body, _ = json.Marshal(encryptBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/encrypt", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)
	if err != nil {
		errorDie(err, t)
	}

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s in ErrorCode. Got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}

	// Test Invalid Body
}

func TestDecryptDataOnly(t *testing.T) {

	decryptBody := models.GPGDecryptData{
		DataOnly:         true,
		AsciiArmoredData: test.TestDecryptDataOnly,
	}

	body, err := json.Marshal(decryptBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/decrypt", r)

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

	var decryptedData models.GPGDecryptedData

	err = json.Unmarshal(d, &decryptedData)

	errorDie(err, t)

	decryptedBytes, err := base64.StdEncoding.DecodeString(decryptedData.Base64Data)

	errorDie(err, t)

	if string(decryptedBytes) != test.TestSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", test.TestSignatureData, string(decryptedBytes))
	}
}
func TestDecrypt(t *testing.T) {
	InvalidPayloadTest("/gpg/decrypt", t)
	decryptBody := models.GPGDecryptData{
		DataOnly:         false,
		AsciiArmoredData: test.TestDecryptDataAscii,
	}

	body, err := json.Marshal(decryptBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/decrypt", r)

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

	var decryptedData models.GPGDecryptedData

	err = json.Unmarshal(d, &decryptedData)

	errorDie(err, t)

	decryptedBytes, err := base64.StdEncoding.DecodeString(decryptedData.Base64Data)

	errorDie(err, t)

	if string(decryptedBytes) != test.TestSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", test.TestSignatureData, string(decryptedBytes))
	}
}
func TestVerifySignature(t *testing.T) {
	InvalidPayloadTest("/gpg/verifySignature", t)
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
		Signature:  test.TestSignatureSignature,
	}

	body, err := json.Marshal(verifyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/verifySignature", r)

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
}
func TestVerifySignatureQuanto(t *testing.T) {
	InvalidPayloadTest("/gpg/verifySignatureQuanto", t)
	quantoSignature := tools.GPG2Quanto(test.TestSignatureSignature, test.TestKeyFingerprint, "SHA512")

	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
		Signature:  quantoSignature,
	}

	body, err := json.Marshal(verifyBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/verifySignatureQuanto", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		log.Debug(errObj.StackTrace)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
}
func TestSign(t *testing.T) {
	InvalidPayloadTest("/gpg/sign", t)
	// region Generate Signature
	signBody := models.GPGSignData{
		FingerPrint: test.TestKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
	}

	body, err := json.Marshal(signBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/sign", r)

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

	log.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
		Signature:  string(d),
	}

	body, err = json.Marshal(verifyBody)

	errorDie(err, t)

	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/verifySignature", r)

	errorDie(err, t)

	res = executeRequest(req)

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
	// endregion
	// region Test Invalid Data
	signBody = models.GPGSignData{
		FingerPrint: test.TestKeyFingerprint,
		Base64Data:  "é123 'en / .sd",
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/sign", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
	// region Test Invalid Fingerprint
	signBody = models.GPGSignData{
		FingerPrint: "ABCDEFGH",
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/sign", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err = ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}
func TestSignQuanto(t *testing.T) {
	InvalidPayloadTest("/gpg/signQuanto", t)
	// region Generate Signature
	signBody := models.GPGSignData{
		FingerPrint: test.TestKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
	}

	body, err := json.Marshal(signBody)

	errorDie(err, t)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/signQuanto", r)

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

	log.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
		Signature:  string(d),
	}

	body, err = json.Marshal(verifyBody)

	errorDie(err, t)

	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/verifySignatureQuanto", r)

	errorDie(err, t)

	res = executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		t.Errorf("%s", errObj.String())
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	if string(d) != "OK" {
		t.Errorf("Expected OK got %s", string(d))
	}
	// endregion
	// region Test Invalid Data
	signBody = models.GPGSignData{
		FingerPrint: test.TestKeyFingerprint,
		Base64Data:  "é123 'en / .sd",
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/signQuanto", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
	// region Test Invalid Fingerprint
	signBody = models.GPGSignData{
		FingerPrint: "ABCDEFGH",
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(test.TestSignatureData)),
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/signQuanto", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err = ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}

func TestUnlockKey(t *testing.T) {
	InvalidPayloadTest("/gpg/unlockKey", t)
	// region Test Unlock Key
	data := models.GPGUnlockKeyData{
		FingerPrint: test.TestKeyFingerprint,
		Password:    test.TestKeyPassword,
	}

	body, _ := json.Marshal(data)
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/unlockKey", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		t.Errorf("%s", errObj.String())
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	if string(d) != "OK" {
		errorDie(fmt.Errorf("expected response %s got %s", "OK", string(d)), t)
	}
	// endregion
	// region Test Invalid Password
	data.Password += "huebr"
	body, _ = json.Marshal(data)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/unlockKey", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err := ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected ErrorCode to be %s got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
	// endregion
}

// endregion
