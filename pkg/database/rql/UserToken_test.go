package rql

import (
	"fmt"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"strings"
	"testing"
	"time"
)

var testToken = models.UserToken{
	ID:          "abcd",
	FingerPrint: "DEADBEEF",
	Username:    "johnhuebr",
	Fullname:    "John HUEBR",
	Token:       "dummytoken",
	CreatedAt:   time.Now().Truncate(time.Second),
	Expiration:  time.Now().Add(time.Hour).Truncate(time.Second),
}

func TestRethinkDBDriver_AddUserToken(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testToken)
	m2, _ := convertToRethinkDB(models.UserToken{})

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Insert(m)).
		Return(r.WriteResponse{
			GeneratedKeys: []string{testToken.ID},
		}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Insert(m2)).
		Return(r.WriteResponse{}, fmt.Errorf("test error")))

	id, err := h.AddUserToken(testToken)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.EqualFold(id, testToken.ID) {
		t.Fatalf("expected token id to be %q but got %q", testToken.ID, id)
	}

	// Test error
	_, err = h.AddUserToken(models.UserToken{})

	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_GetUserToken(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testToken)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", testToken.Token).
		Limit(1).
		CoerceTo("array")).
		Return([]map[string]interface{}{m}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", "").
		Limit(1).
		CoerceTo("array")).
		Return([]map[string]interface{}{}, nil))

	u, err := h.GetUserToken(testToken.Token)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if diff := pretty.Compare(testToken, u); diff != "" {
		t.Errorf("Expected token to be the same. (-got +want)\\n%s", diff)
	}

	// Test not found
	_, err = h.GetUserToken("")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not found") {
		t.Fatalf("expectet error to be %q but got %q", "not found", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_InvalidateUserTokens(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Filter(r.Row.Field("Expiration").Lt(r.MockAnything())).
		Delete()).
		Return(r.WriteResponse{
			Deleted: 100,
		}, nil))

	n, err := h.InvalidateUserTokens()

	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if n != 100 {
		t.Fatalf("expected %d deletes got %d", 100, n)
	}

	mock.AssertExpectations(t)

	mock = r.NewMock()
	h.conn = mock

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Filter(r.Row.Field("Expiration").Lt(r.MockAnything())).
		Delete()).
		Return(r.WriteResponse{}, fmt.Errorf("test error")))

	_, err = h.InvalidateUserTokens()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %s got %s", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_RemoveUserToken(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", testToken.Token).
		Limit(1).
		Delete()).
		Return(r.WriteResponse{Deleted: 1}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", "").
		Limit(1).
		Delete()).
		Return(r.WriteResponse{}, fmt.Errorf("test error")))

	err := h.RemoveUserToken(testToken.Token)

	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	err = h.RemoveUserToken("")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}
