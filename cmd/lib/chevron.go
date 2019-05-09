package main

import "C"
import (
	"github.com/quan-to/chevron/chevronlib"
)

//export LoadKey
func LoadKey(keyData string) (loadedPrivateKeys int, err string) {
	l, e := chevronlib.LoadKey(keyData)
	loadedPrivateKeys = l
	if e != nil {
		err = e.Error()
	}
	return
}

//export UnlockKey
func UnlockKey(fingerprint, password string) (err string) {
	e := chevronlib.UnlockKey(fingerprint, password)
	if e != nil {
		err = e.Error()
	}
	return
}

//export VerifySignature
func VerifySignature(data []byte, signature string) (result bool, err string) {
	r, e := chevronlib.VerifySignature(data, signature)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

//export VerifyBase64DataSignature
func VerifyBase64DataSignature(b64data, signature string) (result bool, err string) {
	r, e := chevronlib.VerifyBase64DataSignature(b64data, signature)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

//export SignData
func SignData(data []byte, fingerprint string) (result string, err string) {
	r, e := chevronlib.SignData(data, fingerprint)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

//export SignBase64Data
func SignBase64Data(b64data, fingerprint string) (result string, err string) {
	r, e := chevronlib.SignBase64Data(b64data, fingerprint)
	result = r
	if e != nil {
		err = e.Error()
	}

	return
}

func main() {}
