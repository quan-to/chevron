package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/quan-to/chevron/pkg/QuantoError"
)

func TestAdminLogin(t *testing.T) {
	// region Generate Login
	payload := map[string]interface{}{
		"query": "mutation Login($username: String!, $password: String!) { Login(username: $username, password: $password) { Value UserName ExpirationDateTimeISO UserFullName  }}",
		"variables": map[string]interface{}{
			"username": "admin",
			"password": "admin",
		},
		"operationName": "Login",
		"_timestamp":    time.Now().Nanosecond() / 1000,
		"_timeUniqueId": fmt.Sprintf("%v-agent", time.Now().Nanosecond()/1000),
	}

	d, _ := json.Marshal(payload)

	r := bytes.NewReader(d)

	req, err := http.NewRequest("POST", "/agentAdmin", r)

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

	var data map[string]interface{}

	err = json.Unmarshal(d, &data)

	errorDie(err, t)

	loginData := data["data"].(map[string]interface{})["Login"].(map[string]interface{})

	UserName := loginData["UserName"].(string)
	Value := loginData["Value"].(string)

	if UserName != "admin" {
		errorDie(fmt.Errorf("expected username to be admin got %s", UserName), t)
	}

	if Value == "" {
		errorDie(fmt.Errorf("expected Value got empty"), t)
	}

	// endregion
	// region Test Login
	payload = map[string]interface{}{
		"query":         "\n\nquery Me {\n  WhoAmI\n}\n",
		"variables":     map[string]interface{}{},
		"operationName": "Me",
		"_timestamp":    time.Now().Nanosecond() / 1000,
		"_timeUniqueId": fmt.Sprintf("%v-agent", time.Now().Nanosecond()/1000),
	}

	d, _ = json.Marshal(payload)

	r = bytes.NewReader(d)

	req, err = http.NewRequest("POST", "/agentAdmin", r)

	errorDie(err, t)

	req.Header.Add("proxyToken", Value)

	res = executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	if res.Code != 200 {
		var errObj QuantoError.ErrorObject
		err := json.Unmarshal(d, &errObj)
		errorDie(err, t)
		errorDie(fmt.Errorf(errObj.Message), t)
	}

	errorDie(err, t)

	err = json.Unmarshal(d, &data)

	errorDie(err, t)

	whoAmI := data["data"].(map[string]interface{})["WhoAmI"].(string)

	if whoAmI != "Administrator" {
		errorDie(fmt.Errorf("expected whoAmI to be Administrator got %s", whoAmI), t)
	}
	// endregion
}

func TestInvalidToken(t *testing.T) {
	payload := map[string]interface{}{
		"query":         "\n\nquery Me {\n  WhoAmI\n}\n",
		"variables":     map[string]interface{}{},
		"operationName": "Me",
		"_timestamp":    time.Now().Nanosecond() / 1000,
		"_timeUniqueId": fmt.Sprintf("%v-agent", time.Now().Nanosecond()/1000),
	}

	d, _ := json.Marshal(payload)

	r := bytes.NewReader(d)

	req, err := http.NewRequest("POST", "/agentAdmin", r)

	errorDie(err, t)

	req.Header.Add("proxyToken", "huebrbrbrbrbr")

	res := executeRequest(req)

	d, err = ioutil.ReadAll(res.Body)

	errorDie(err, t)

	if res.Code == 200 {
		errorDie(fmt.Errorf("expected not 200, got 200"), t)
	}

	var errObj QuantoError.ErrorObject
	err = json.Unmarshal(d, &errObj)
	errorDie(err, t)

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s in errorCode, got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
}
