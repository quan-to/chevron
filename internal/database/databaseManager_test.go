package database

import (
	"fmt"
	"github.com/google/uuid"
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	slog.UnsetTestMode()

	QuantoError.EnableStackTrace()

	config.PrivateKeyFolder = ".."
	config.KeyPrefix = "testkey_"
	config.KeysBase64Encoded = false
	config.RethinkDBPoolSize = 1

	config.MasterGPGKeyBase64Encoded = false
	config.MasterGPGKeyPath = "../testkey_privateTestKey.gpg"
	config.MasterGPGKeyPasswordPath = "../testprivatekeyPassword.txt"
	u, _ := uuid.NewRandom()
	config.DatabaseName = "qrs_test_" + u.String()
	config.HttpPort = 40000
	config.SKSServer = fmt.Sprintf("http://localhost:%d/sks/", config.HttpPort)
	config.EnableRethinkSKS = true
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
