package cache

import (
	"fmt"
	"testing"

	"github.com/quan-to/chevron/pkg/models/testmodels"

	"bou.ke/monkey"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/database/memory"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

func TestDriver_AddGPGKey(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	testKeyToAdd := testmodels.GpgKey
	testKeyToAdd.ID = "0000"

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return testKeyToAdd.ID
	})

	data, err := h.cache.Marshal(&testKeyToAdd)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectSet(gpgKeyByIDPrefix+testKeyToAdd.ID, data, gpgKeyExpiration).
		SetVal("")
	mock.ExpectSet(gpgKeyByFingerprintPrefix+testKeyToAdd.GetShortFingerPrint(), data, gpgKeyExpiration).
		SetVal("")

	id, added, err := h.AddGPGKey(testmodels.GpgKey)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if !added {
		t.Fatal("expected key to be added")
	}

	// Check if key has been added
	key, err := mem.FetchGPGKeyByFingerprint(testmodels.GpgKey.FullFingerprint)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if key.ID != id {
		t.Fatalf("expected added key ID to be %s got %s", id, key.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_DeleteGPGKey(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	testKeyToRemove := testmodels.GpgKey
	testKeyToRemove.ID = "0000"

	// Assume mem works, and add the gpg key
	_, _, _ = mem.AddGPGKey(testmodels.GpgKey)

	mock.ExpectDel(gpgKeyByIDPrefix + testKeyToRemove.ID).SetVal(0)
	mock.ExpectDel(gpgKeyByFingerprintPrefix + testKeyToRemove.FullFingerprint).SetVal(0)

	err := h.DeleteGPGKey(testKeyToRemove)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	key, err := mem.FetchGPGKeyByFingerprint(testKeyToRemove.FullFingerprint)
	if key != nil || err == nil {
		t.Fatalf("expected key to be removed")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_UpdateGPGKey(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	testKeyToUpdate := testmodels.GpgKey
	testKeyToUpdate.ID = "0000"

	// Assume mem works, and add the gpg key
	_, _, _ = mem.AddGPGKey(testKeyToUpdate)

	data, err := h.cache.Marshal(&testKeyToUpdate)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectSet(gpgKeyByIDPrefix+testKeyToUpdate.ID, data, gpgKeyExpiration).
		SetVal("")
	mock.ExpectSet(gpgKeyByFingerprintPrefix+testKeyToUpdate.GetShortFingerPrint(), data, gpgKeyExpiration).
		SetVal("")

	err = h.UpdateGPGKey(testKeyToUpdate)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_FetchGPGKeyByFingerprint(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	t.Log("Testing uncached key")
	testKeyToFetch := testmodels.GpgKey
	testKeyToFetch.ID = "0000"

	// Assume mem works, and add the gpg key
	_, _, _ = mem.AddGPGKey(testKeyToFetch)

	data, err := h.cache.Marshal(&testKeyToFetch)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectGet(gpgKeyByFingerprintPrefix + testKeyToFetch.GetShortFingerPrint()).
		SetErr(fmt.Errorf("not found"))
	mock.ExpectSet(gpgKeyByIDPrefix+testKeyToFetch.ID, data, gpgKeyExpiration).
		SetVal("")
	mock.ExpectSet(gpgKeyByFingerprintPrefix+testKeyToFetch.GetShortFingerPrint(), data, gpgKeyExpiration).
		SetVal("")

	key, err := h.FetchGPGKeyByFingerprint(testKeyToFetch.GetShortFingerPrint())
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if key == nil {
		t.Fatal("expected key, got nil")
	}

	if diff := pretty.Compare(testKeyToFetch, *key); diff != "" {
		t.Errorf("Expected key to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}

	// Test cached key
	t.Log("Testing cached key")
	db, mock = redismock.NewClientMock()
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	mock.ExpectGet(gpgKeyByFingerprintPrefix + testKeyToFetch.GetShortFingerPrint()).
		SetVal(string(data))

	key, err = h.FetchGPGKeyByFingerprint(testKeyToFetch.GetShortFingerPrint())
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if key == nil {
		t.Fatal("expected key, got nil")
	}

	if diff := pretty.Compare(testKeyToFetch, *key); diff != "" {
		t.Errorf("Expected key to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_FetchGPGKeysWithoutSubKeys(t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(mem, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Passthrough call, test if consistent
	gotKeys, gotErr := h.FetchGPGKeysWithoutSubKeys()
	expectedKeys, expectedErr := mem.FetchGPGKeysWithoutSubKeys()

	if gotErr != expectedErr {
		t.Fatalf("expected error to be %q got %q", expectedErr, gotErr)
	}

	if diff := pretty.Compare(expectedKeys, gotKeys); diff != "" {
		t.Errorf("Expected key list to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func testFindFunction(value, criteria string, h *Driver, f keyListFallbackFunc, t *testing.T) {
	mem := memory.MakeMemoryDBDriver(nil)
	db, mock := redismock.NewClientMock()
	h.proxy = mem
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	// Assume mem works, and add the gpg key
	testKey := testmodels.GpgKey
	id, _, _ := mem.AddGPGKey(testKey)
	testKey.ID = id

	keyList := []models.GPGKey{testKey}

	data, err := h.cache.Marshal(&keyList)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	keyString := fmt.Sprintf("%s%s%s%d%d", gpgKeyEntryList, criteria, value, 0, 10)

	mock.ExpectSet(keyString, data, gpgKeyEntriesExpiration).SetVal("")

	// Passthrough call, test if consistent
	keys, err := f(value, 0, 10)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if len(keys) == 0 || len(keys) > 1 {
		t.Fatalf("expected one key got %d", len(keys))
	}

	gotKey := keys[0]
	if diff := pretty.Compare(testKey, gotKey); diff != "" {
		t.Errorf("Expected key to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_FindGPGKeyByEmail(t *testing.T) {
	value := testmodels.GpgKey.Emails[0]
	h := MakeRedisDriver(nil, nil)
	testFindFunction(value, gpgKeysByEmailCriteria, h, h.FindGPGKeyByEmail, t)
}

func TestDriver_FindGPGKeyByFingerPrint(t *testing.T) {
	value := testmodels.GpgKey.FullFingerprint
	h := MakeRedisDriver(nil, nil)
	testFindFunction(value, gpgKeysByFingerprintCriteria, h, h.FindGPGKeyByFingerPrint, t)
}

func TestDriver_FindGPGKeyByName(t *testing.T) {
	value := testmodels.GpgKey.Names[0]
	h := MakeRedisDriver(nil, nil)
	testFindFunction(value, gpgKeysByNameCriteria, h, h.FindGPGKeyByName, t)
}

func TestDriver_FindGPGKeyByValue(t *testing.T) {
	valuesToFind := []string{
		testmodels.GpgKey.FullFingerprint,
		testmodels.GpgKey.GetShortFingerPrint(),
	}

	for _, v := range testmodels.GpgKey.Names {
		valuesToFind = append(valuesToFind, v)
	}
	for _, v := range testmodels.GpgKey.Emails {
		valuesToFind = append(valuesToFind, v)
	}

	for _, v := range valuesToFind {
		h := MakeRedisDriver(nil, nil)
		testFindFunction(v, gpgKeysByValueCriteria, h, h.FindGPGKeyByValue, t)
	}
}
