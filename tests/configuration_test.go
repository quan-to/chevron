package tests

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestPushPopVars(t *testing.T) {
	remote_signer.PopVariables()
	remote_signer.PushVariables()
	remote_signer.PopVariables()
}

func testIntVar(v *int, envName string, localName string, t *testing.T) {
	err := os.Setenv(envName, "huebr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assertPanic(t, remote_signer.Setup, fmt.Sprintf("%s should panic with a invalid value", envName))

	val := int(rand.Int31())

	err = os.Setenv(envName, strconv.FormatInt(int64(val), 10))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	remote_signer.Setup()

	if *v != val {
		t.Errorf("%s variable does not come from %s. Expected %d got %d", localName, envName, *v, val)
	}
}

func TestConfiguration(t *testing.T) {
	remote_signer.PushVariables()

	testIntVar(&remote_signer.MaxKeyRingCache, "MAX_KEYRING_CACHE_SIZE", "MaxKeyRingCache", t)
	testIntVar(&remote_signer.HttpPort, "HTTP_PORT", "HttpPort", t)
	testIntVar(&remote_signer.RethinkDBPoolSize, "RETHINKDB_POOL_SIZE", "RethinkDBPoolSize", t)
	testIntVar(&remote_signer.RethinkDBPort, "RETHINKDB_PORT", "RethinkDBPort", t)

	remote_signer.PopVariables()
}
