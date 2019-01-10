package remote_signer

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestPushPopVars(t *testing.T) {
	popVariables()
	pushVariables()
	popVariables()
}

func testIntVar(v *int, envName string, localName string, t *testing.T) {
	err := os.Setenv(envName, "huebr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assertPanic(t, setup, fmt.Sprintf("%s should panic with a invalid value", envName))

	val := int(rand.Int31())

	err = os.Setenv(envName, strconv.FormatInt(int64(val), 10))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	setup()

	if *v != val {
		t.Errorf("%s variable does not come from %s. Expected %d got %d", localName, envName, *v, val)
	}
}

func TestConfiguration(t *testing.T) {
	pushVariables()

	testIntVar(&MaxKeyRingCache, "MAX_KEYRING_CACHE_SIZE", "MaxKeyRingCache", t)
	testIntVar(&HttpPort, "HTTP_PORT", "HttpPort", t)
	testIntVar(&RethinkDBPoolSize, "RETHINKDB_POOL_SIZE", "RethinkDBPoolSize", t)
	testIntVar(&RethinkDBPort, "RETHINKDB_PORT", "RethinkDBPort", t)

	popVariables()
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
