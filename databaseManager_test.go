package remote_signer

import (
	"fmt"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/SLog"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
	"testing"
)

var slog *SLog.Instance

func ResetDatabase() {
	c := GetConnection()
	dbs := getDatabases()
	if stringIndexOf(DatabaseName, dbs) > -1 {
		_, err := rethinkdb.DBDrop(DatabaseName).Run(c)
		if err != nil {
			panic(err)
		}
	}
}

func TestMain(m *testing.M) {
	QuantoError.EnableStackTrace()
	SLog.SetTestMode()
	slog = SLog.Scope("TestLog")

	PrivateKeyFolder = "."
	KeyPrefix = "testkey_"
	KeysBase64Encoded = false

	MasterGPGKeyBase64Encoded = false
	MasterGPGKeyPath = "./testkey_privateTestKey.gpg"
	MasterGPGKeyPasswordPath = "./testprivatekeyPassword.txt"

	DatabaseName = "qrs_test"
	HttpPort = 40000
	SKSServer = fmt.Sprintf("http://localhost:%d/sks/", HttpPort)
	EnableRethinkSKS = true

	rthState = rethinkDbState{
		currentConn: 0,
	}

	dbSetup()
	ResetDatabase()

	sm := MakeSecretsManager()
	gpg := MakePGPManager()
	gpg.LoadKeys()

	err := gpg.UnlockKey(testKeyFingerprint, testKeyPassword)

	if err != nil {
		SLog.SetError(true)
		slog.Error(err)
		os.Exit(1)
	}

	_ = sm

	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	initTables()
	initTables() // Test if already initialized
}
