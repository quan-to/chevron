package cache

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/quan-to/chevron/pkg/models/testmodels"

	"bou.ke/monkey"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

const unexpectedError = "unexpected error: %s"
const expectationsWereNotMet = "expectations were not met: %s"
const userTokenExpirationTime = time.Hour

func TestDriver_AddUserToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(nil, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	monkey.Patch(uuid.EnsureUUID, func(log slog.Instance) string {
		return "0000"
	})
	monkey.Patch(time.Now, func() time.Time {
		return testmodels.Time
	})

	testData := testmodels.Token
	testData.ID = "0000"

	data, err := h.cache.Marshal(&testData)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectSet(userTokenPrefix+testmodels.Token.Token, data, userTokenExpirationTime).
		SetVal("")

	entryId, err := h.AddUserToken(testmodels.Token)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if entryId != testData.ID {
		t.Fatalf("expected entryId to be %s got %s", testData.ID, entryId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_GetUserToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	h := MakeRedisDriver(nil, nil)
	h.cache = cache.New(&cache.Options{
		Redis: db,
	})

	testData := testmodels.Token
	testData.ID = "0000"

	data, err := h.cache.Marshal(&testData)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectGet(userTokenPrefix + testmodels.Token.Token).SetVal(string(data))
	mock.ExpectGet(userTokenPrefix + "HUEBR").SetErr(fmt.Errorf("test error"))

	entry, err := h.GetUserToken(testmodels.Token.Token)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if diff := pretty.Compare(testData, entry); diff != "" {
		t.Errorf("Expected token to be the same. (-got +want)\\n%s", diff)
	}

	_, err = h.GetUserToken("HUEBR")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %s but got %s", "test error", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsWereNotMet, err)
	}
}

func TestDriver_InvalidateUserTokens(t *testing.T) {
	h := MakeRedisDriver(nil, nil)
	n, err := h.InvalidateUserTokens()
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if n != 0 {
		t.Fatalf("expected no invalidations, got %d", n)
	}
}
