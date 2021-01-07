package config

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"syscall"
	"testing"

	"github.com/bouk/monkey"
	"github.com/quan-to/slog"
)

func assertPanic(t *testing.T, f func(), message string) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}

func TestPushPopVars(t *testing.T) {
	PopVariables()
	PushVariables()
	PopVariables()
}

func testIntVar(v *int, envName string, localName string, t *testing.T) {
	slog.SetTestMode()
	err := os.Setenv(envName, "huebr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assertPanic(t, Setup, fmt.Sprintf("%s should panic with a invalid value", envName))

	val := int(rand.Int31())

	err = os.Setenv(envName, strconv.FormatInt(int64(val), 10))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	Setup()

	if *v != val {
		t.Errorf("%s variable does not come from %s. Expected %d got %d", localName, envName, *v, val)
	}

	slog.UnsetTestMode()
}

func testStringVar(v *string, envName string, localName string, def string, t *testing.T) {
	slog.SetTestMode()
	err := os.Setenv(envName, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	Setup()
	if *v != def {
		t.Error(fmt.Errorf("%s: expected default %s got %s", localName, def, *v))
		t.FailNow()
	}

	val := strconv.FormatInt(int64(rand.Int31()), 32)

	err = os.Setenv(envName, val)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	Setup()

	if *v != val {
		t.Errorf("%s variable does not come from %s. Expected %s got %s", localName, envName, *v, val)
	}
	slog.UnsetTestMode()
}

func assertEqual(a interface{}, b interface{}, message string, t *testing.T) {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		t.Errorf(message)
		t.FailNow()
	}

	switch v := a.(type) {
	case bool:
		if v != b.(bool) {
			t.Errorf(message)
			t.FailNow()
		}
	case error:
		if v != b.(error) {
			t.Errorf(message)
			t.FailNow()
		}
	case string:
		if v != b.(string) {
			t.Errorf(message)
			t.FailNow()
		}
	case int:
		if v != b.(int) {
			t.Errorf(message)
			t.FailNow()
		}
	case float32:
		if v != b.(float32) {
			t.Errorf(message)
			t.FailNow()
		}
	case float64:
		if v != b.(float64) {
			t.Errorf(message)
			t.FailNow()
		}
	default:
		if v != b {
			t.Errorf(message)
			t.FailNow()
		}
	}
}

func TestConfiguration(t *testing.T) {
	slog.SetTestMode()
	PushVariables()

	testIntVar(&MaxKeyRingCache, "MAX_KEYRING_CACHE_SIZE", "MaxKeyRingCache", t)
	testIntVar(&HttpPort, "HTTP_PORT", "HttpPort", t)
	testIntVar(&AgentTokenExpiration, "AGENT_TOKEN_EXPIRATION", "AgentTokenExpiration", t)

	testStringVar(&SyslogServer, "SYSLOG_IP", "SyslogServer", "127.0.0.1", t)
	testStringVar(&SyslogFacility, "SYSLOG_FACILITY", "SyslogFacility", "LOG_USER", t)
	testStringVar(&PrivateKeyFolder, "PRIVATE_KEY_FOLDER", "PrivateKeyFolder", "./keys", t)
	testStringVar(&VaultAddress, "VAULT_ADDRESS", "VaultAddress", "http://localhost:8200", t)
	testStringVar(&VaultNamespace, "VAULT_NAMESPACE", "VaultNamespace", "remote-signer", t)
	testStringVar(&VaultBackend, "VAULT_BACKEND", "VaultBackend", "secret", t)
	testStringVar(&AgentTargetURL, "AGENT_TARGET_URL", "AgentTargetURL", "https://api.sandbox.contaquanto.com/all", t)
	testStringVar(&Environment, "Environment", "Environment", "development", t)
	testStringVar(&AgentExternalURL, "AGENT_EXTERNAL_URL", "AgentExternalURL", "/agent", t)
	testStringVar(&AgentAdminExternalURL, "AGENTADMIN_EXTERNAL_URL", "AgentAdminExternalURL", "/agentAdmin", t)

	PopVariables()

	PushVariables()
	slog.SetTestMode()

	PopVariables()
	slog.UnsetTestMode()

	_ = os.Setenv("ENABLE_RETHINKDB_SKS", "true")

	_ = os.Setenv("Environment", "development")
	Setup()
	assertEqual(slog.DebugEnabled(), true, "Debug should be enabled in development", t)

	_ = os.Setenv("Environment", "production")
	Setup()
	assertEqual(slog.DebugEnabled(), false, "Debug should be disabled in production", t)

	PopVariables()
	slog.UnsetTestMode()

	PushVariables()
	_ = syscall.Setenv("SHOW_LINES", "true")
	Setup()
	assertEqual(slog.ShowLinesEnabled(), true, "SHOW_LINES=true env should set slog.SetShowLines to true", t)
	PopVariables()

	PushVariables()
	_ = syscall.Setenv("SHOW_LINES", "false")
	Setup()
	assertEqual(slog.ShowLinesEnabled(), false, "SHOW_LINES=false env should set slog.SetShowLines to false", t)
	PopVariables()
}
