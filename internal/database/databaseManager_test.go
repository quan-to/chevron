package database

import (
	"fmt"
	"github.com/google/uuid"
	config "github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
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
	InitTables()
	InitTables() // Try twice to ensure no breaks in existent adds
	Cleanup()
}

func TestGetConnection(t *testing.T) {
	s := GetConnection()
	if s == nil {
		t.Fatal("Expected to have RethinkDB Connection on GetConnection")
	}

	RthState.connection = nil

	s = GetConnection()
	if s == nil {
		t.Fatal("Expected GetConnection to setup a connection")
	}
}

func TestResetDatabase(t *testing.T) {
	u, _ := uuid.NewRandom()
	config.DatabaseName = "qrs_test_" + u.String()

	c := GetConnection()
	err := r.DBCreate(config.DatabaseName).Exec(c)
	if err != nil {
		t.Fatalf("Cannot create database %q", err)
	}

	WaitDatabaseCreate(config.DatabaseName)

	ResetDatabase()

	dbs := GetDatabases()

	if tools.StringIndexOf(config.DatabaseName, dbs) > -1 {
		t.Error("Expected ResetDatabase to drop existing database")
	}
}
