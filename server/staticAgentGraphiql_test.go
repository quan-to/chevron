package server

import (
	"encoding/json"
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/server/agent"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func GetFile(url string, t *testing.T) []byte {
	req, err := http.NewRequest("GET", "/graphiql"+url, nil)

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

	return d
}

func compareByteArray(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}

func TestStaticGraphiQL(t *testing.T) {
	files, err := agent.AssetDir("bundle")
	errorDie(err, t)

	// Test All Files
	for _, v := range files {
		got := GetFile("/"+v, t)
		expected, _ := agent.Asset("bundle/" + v)

		if strings.Index(v, "index.htm") > -1 {
			f := string(expected)
			f = strings.Replace(f, "{SERVER_URL}", remote_signer.AgentTargetURL, -1)
			f = strings.Replace(f, "{AGENT_URL}", remote_signer.AgentExternalURL, -1)
			f = strings.Replace(f, "{AGENT_ADMIN_URL}", remote_signer.AgentAdminExternalURL, -1)
			expected = []byte(f)
		}

		if !compareByteArray(got, expected) {
			t.Errorf("expected byte array is different than received for file %s", v)
		}
	}

	got := GetFile("/", t)
	expected, _ := agent.Asset("bundle/index.html")

	f := string(expected)
	f = strings.Replace(f, "{SERVER_URL}", remote_signer.AgentTargetURL, -1)
	f = strings.Replace(f, "{AGENT_URL}", remote_signer.AgentExternalURL, -1)
	f = strings.Replace(f, "{AGENT_ADMIN_URL}", remote_signer.AgentAdminExternalURL, -1)
	expected = []byte(f)

	if !compareByteArray(got, expected) {
		t.Errorf("expected byte array is different than received for root")
	}

	got = GetFile("", t)
	if !compareByteArray(got, expected) {
		t.Errorf("expected byte array is different than received for empty root")
	}
}
