package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"testing"

	"github.com/google/uuid"
	"github.com/quan-to/chevron/internal/agent"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/etc/magicbuilder"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/test"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

var sm interfaces.SecretsManager
var gpg interfaces.PGPManager
var log = slog.Scope("TestRemoteSigner")
var dbh DatabaseHandler

var router *mux.Router

func errorDie(err error, t *testing.T) {
	if err != nil {
		fmt.Println("----------------------------------------")
		debug.PrintStack()
		fmt.Println("----------------------------------------")
		t.Error(err)
		t.FailNow()
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestMain(m *testing.M) {
	slog.UnsetTestMode()
	var err error

	QuantoError.EnableStackTrace()

	u, _ := uuid.NewRandom()

	config.DatabaseName = "qrs_test_" + u.String()
	config.PrivateKeyFolder = "../../test/data"
	config.KeyPrefix = "testkey_"
	config.KeysBase64Encoded = false
	config.RethinkDBPoolSize = 1
	config.EnableRethinkSKS = false

	config.MasterGPGKeyBase64Encoded = false
	config.MasterGPGKeyPath = "../../test/data/testkey_privateTestKey.gpg"
	config.MasterGPGKeyPasswordPath = "../../test/data/testprivatekeyPassword.txt"

	ctx := context.Background()
	dbh, err = agent.MakeDatabaseHandler(log)

	if err != nil {
		slog.Fatal("Error initializing selected database: %s", err)
	}
	ctx = context.WithValue(ctx, tools.CtxDatabaseHandler, dbh)

	sm = magicbuilder.MakeSM(nil, dbh)
	gpg = magicbuilder.MakePGP(nil, dbh)
	gpg.LoadKeys(ctx)

	err = gpg.UnlockKey(ctx, test.TestKeyFingerprint, test.TestKeyPassword)

	if err != nil {
		slog.UnsetTestMode()
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Adding key %s to SKS Database", test.TestKeyFingerprint)
	pubKey, _ := gpg.GetPublicKeyASCII(ctx, test.TestKeyFingerprint)
	log.Info("Result: %s", keymagic.PKSAdd(ctx, pubKey))

	router = GenRemoteSignerServerMux(log, sm, gpg, dbh)

	slog.SetTestMode()
	code := m.Run()
	slog.UnsetTestMode()
	os.Exit(code)
}

func InvalidPayloadTest(endpoint string, t *testing.T) {
	r := bytes.NewReader([]byte(""))

	req, err := http.NewRequest("POST", endpoint, r)

	errorDie(err, t)

	res := executeRequest(req)

	if res.Code != 500 {
		errorDie(fmt.Errorf("expected error 500 for invalid payload"), t)
	}

	var errObj QuantoError.ErrorObject

	d, err := ioutil.ReadAll(res.Body)
	errorDie(err, t)
	err = json.Unmarshal(d, &errObj)

	if err != nil {
		errorDie(err, t)
	}

	if errObj.ErrorCode != QuantoError.InvalidFieldData {
		errorDie(fmt.Errorf("expected %s in ErrorCode. Got %s", QuantoError.InvalidFieldData, errObj.ErrorCode), t)
	}
}

func ReadErrorObject(r io.Reader) (QuantoError.ErrorObject, error) {
	var errObj QuantoError.ErrorObject
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return errObj, err
	}

	err = json.Unmarshal(data, &errObj)
	return errObj, err
}
