package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/models"
	"io/ioutil"
	"net/http"
	"testing"
)

// region GPG Endpoint Tests
func TestGenerateKey(t *testing.T) {
	InvalidPayloadTest("/gpg/generateKey", t)

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

	fingerPrint, err := remote_signer.GetFingerPrintFromKey(key)

	errorDie(err, t)

	err, _ = gpg.LoadKey(key)

	errorDie(err, t)

	err = gpg.UnlockKey(fingerPrint, genKeyBody.Password)

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
	encryptBody := models.GPGEncryptData{
		DataOnly:    true,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Filename:    "test-encrypt",
		FingerPrint: testKeyFingerprint,
	}

	body, _ := json.Marshal(encryptBody)

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/encrypt", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	encryptedData := string(d)

	data, err := gpg.Decrypt(encryptedData, true)

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
		AsciiArmoredData: testDecryptDataOnly,
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

	if string(decryptedBytes) != testSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", testSignatureData, string(decryptedBytes))
	}
}
func TestDecrypt(t *testing.T) {
	InvalidPayloadTest("/gpg/decrypt", t)
	decryptBody := models.GPGDecryptData{
		DataOnly:         false,
		AsciiArmoredData: testDecryptDataAscii,
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

	if string(decryptedBytes) != testSignatureData {
		t.Errorf("Expected \"%s\" got \"%s\"", testSignatureData, string(decryptedBytes))
	}
}
func TestVerifySignature(t *testing.T) {
	InvalidPayloadTest("/gpg/verifySignature", t)
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
		Signature:  testSignatureSignature,
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
	quantoSignature := remote_signer.GPG2Quanto(testSignatureSignature, testKeyFingerprint, "SHA512")

	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
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
		slog.Debug(errObj.StackTrace)
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
		FingerPrint: testKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
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

	slog.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
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
		FingerPrint: testKeyFingerprint,
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
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/sign", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err = ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InternalServerError {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InternalServerError, errObj.ErrorCode), t)
	}
	// endregion
}
func TestSignQuanto(t *testing.T) {
	InvalidPayloadTest("/gpg/signQuanto", t)
	// region Generate Signature
	signBody := models.GPGSignData{
		FingerPrint: testKeyFingerprint,
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
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

	slog.Debug("Signature: %s", string(d))
	// endregion
	// region Verify Signature
	verifyBody := models.GPGVerifySignatureData{
		Base64Data: base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
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
		FingerPrint: testKeyFingerprint,
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
		Base64Data:  base64.StdEncoding.EncodeToString([]byte(testSignatureData)),
	}

	body, _ = json.Marshal(signBody)
	r = bytes.NewReader(body)

	req, err = http.NewRequest("POST", "/gpg/signQuanto", r)

	errorDie(err, t)

	res = executeRequest(req)

	errObj, err = ReadErrorObject(res.Body)

	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InternalServerError {
		errorDie(fmt.Errorf("expected error code %s got %s", QuantoError.InternalServerError, errObj.ErrorCode), t)
	}
	// endregion
}

func TestUnlockKey(t *testing.T) {
	InvalidPayloadTest("/gpg/unlockKey", t)
	// region Test Unlock Key
	data := models.GPGUnlockKeyData{
		FingerPrint: testKeyFingerprint,
		Password:    testKeyPassword,
	}

	body, _ := json.Marshal(data)
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "/gpg/unlockKey", r)

	errorDie(err, t)

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

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
