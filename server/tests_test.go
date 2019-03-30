package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/quan-to/remote-signer/QuantoError"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// region Generate Signature
	r := bytes.NewReader([]byte(""))

	req, err := http.NewRequest("POST", "/tests/ping", r)

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
		errorDie(fmt.Errorf("expected body to be OK got %s", string(d)), t)
	}
	// endregion
}
