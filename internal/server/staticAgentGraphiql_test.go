package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/server/agent"
	"github.com/quan-to/chevron/pkg/QuantoError"
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

		if strings.Contains(v, "index.htm") {
			f := string(expected)
			f = strings.Replace(f, "{SERVER_URL}", config.AgentTargetURL, -1)
			f = strings.Replace(f, "{AGENT_URL}", config.AgentExternalURL, -1)
			f = strings.Replace(f, "{AGENT_ADMIN_URL}", config.AgentAdminExternalURL, -1)
			expected = []byte(f)
		}

		if !compareByteArray(got, expected) {
			t.Errorf("expected byte array is different than received for file %s", v)
		}
	}

	got := GetFile("/", t)
	expected, _ := agent.Asset("bundle/index.html")

	f := string(expected)
	f = strings.Replace(f, "{SERVER_URL}", config.AgentTargetURL, -1)
	f = strings.Replace(f, "{AGENT_URL}", config.AgentExternalURL, -1)
	f = strings.Replace(f, "{AGENT_ADMIN_URL}", config.AgentAdminExternalURL, -1)
	expected = []byte(f)

	if !compareByteArray(got, expected) {
		t.Errorf("expected byte array is different than received for root")
	}

	got = GetFile("", t)
	if !compareByteArray(got, expected) {
		t.Errorf("expected byte array is different than received for empty root")
	}
}
