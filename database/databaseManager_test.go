package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	slog.UnsetTestMode()

	QuantoError.EnableStackTrace()

	remote_signer.PrivateKeyFolder = ".."
	remote_signer.KeyPrefix = "testkey_"
	remote_signer.KeysBase64Encoded = false
	remote_signer.RethinkDBPoolSize = 1

	remote_signer.MasterGPGKeyBase64Encoded = false
	remote_signer.MasterGPGKeyPath = "../testkey_privateTestKey.gpg"
	remote_signer.MasterGPGKeyPasswordPath = "../testprivatekeyPassword.txt"
	u, _ := uuid.NewRandom()
	remote_signer.DatabaseName = "qrs_test_" + u.String()
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
	os.Exit(code)
}

func TestInitTable(t *testing.T) {
	ResetDatabase()
	// Breaks the test due rethink non atomic operations
	InitTables()
	Cleanup()
}
