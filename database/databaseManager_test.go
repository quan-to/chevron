package database

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/slog"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var rql *exec.Cmd
	var err error
	rql, err = remote_signer.RQLStart()
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}

	QuantoError.EnableStackTrace()

	remote_signer.PrivateKeyFolder = ".."
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false
	remote_signer.RethinkDBPoolSize = 1

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../testprivatekeyPassword.txt"

	remote_signer.DatabaseName = "qrs_test"
	remote_signer.HttpPort = 40000
	remote_signer.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", remote_signer.HttpPort)
	remote_signer.EnableRethinkSKS = true
	DbSetup()

	slog.SetTestMode()
	code := m.Run()
	slog.UnsetTestMode()

	ResetDatabase()
	slog.UnsetTestMode()
	Cleanup()
	slog.Warn("STOPPING RETHINKDB")
	remote_signer.RQLStop(rql)
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	ResetDatabase()
	// Breaks the test due rethink non atomic operations
	InitTables()
	Cleanup()
}
