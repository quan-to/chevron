package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	// region Test Invalid Proxy Token
	r := bytes.NewReader([]byte(""))

	req, err := http.NewRequest("POST", "/agent", r)

	errorDie(err, t)

	req.Header.Add("proxyToken", "heuehehuaheauhuae")

	res := executeRequest(req)

	d, err := ioutil.ReadAll(res.Body)

	errorDie(err, t)

	if res.Code == 200 {
		errorDie(fmt.Errorf("expected not 200, got 200"), t)
	}

	var errObj QuantoError.ErrorObject
	err = json.Unmarshal(d, &errObj)
	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.PermissionDenied {
		errorDie(fmt.Errorf("expected %s in errorCode, got %s", QuantoError.PermissionDenied, errObj.ErrorCode), t)
	}
	// endregion
	// region Test Login Bypass
	// TODO: Test without Quanto Kernel
	//remote_signer.PushVariables()
	//remote_signer.AgentBypassLogin = true
	//remote_signer.AgentTargetURL = "https://quanto-api.com.br/all"
	//
	//// The AddressPostalCode is public, so no signature needed
	//payload := map[string]interface{}{
	//    "query": "query {\n  System_GetAddressPostalCode(code: \"04348160\") {\n    state\n  }\n}",
	//    "_timestamp":    time.Now().Nanosecond() / 1000,
	//    "_timeUniqueId": fmt.Sprintf("%v-agent", time.Now().Nanosecond()/1000),
	//}
	//d, _ = json.Marshal(payload)
	//
	//r = bytes.NewReader(d)
	//
	//req, err = http.NewRequest("POST", "/agent", r)
	//errorDie(err, t)
	//
	//
	//remote_signer.PopVariables()
	// endregion
}
