package cache

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

const unexpectedError = "unexpected error: %s"
const expectationsWereNotMet = "expectations were not met: %s"
const userTokenExpirationTime = time.Hour

var testTime = time.Now().Truncate(time.Second)

var testToken = models.UserToken{
	Fingerprint: "DEADBEEF",
	Username:    "johnhuebr",
	Fullname:    "John HUEBR",
	Token:       uuid.EnsureUUID(nil),
	CreatedAt:   testTime.Truncate(time.Second),
	Expiration:  testTime.Add(userTokenExpirationTime).Truncate(time.Second),
}

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
		return testTime
	})

	testData := testToken
	testData.ID = "0000"

	data, err := h.cache.Marshal(&testData)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectSet(userTokenPrefix+testToken.Token, data, userTokenExpirationTime).
		SetVal("")

	entryId, err := h.AddUserToken(testToken)
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

	testData := testToken
	testData.ID = "0000"

	data, err := h.cache.Marshal(&testData)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	mock.ExpectGet(userTokenPrefix + testToken.Token).SetVal(string(data))
	mock.ExpectGet(userTokenPrefix + "HUEBR").SetErr(fmt.Errorf("test error"))

	entry, err := h.GetUserToken(testToken.Token)
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
