package rql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quan-to/chevron/pkg/models/testmodels"

	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func TestRethinkDBDriver_AddUserToken(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.Token)
	m2, _ := convertToRethinkDB(models.UserToken{})

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Insert(m)).
		Return(r.WriteResponse{
			GeneratedKeys: []string{testmodels.Token.ID},
		}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		Insert(m2)).
		Return(r.WriteResponse{}, fmt.Errorf("test error")))

	id, err := h.AddUserToken(testmodels.Token)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.EqualFold(id, testmodels.Token.ID) {
		t.Fatalf("expected token id to be %q but got %q", testmodels.Token.ID, id)
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

	m, _ := convertToRethinkDB(testmodels.Token)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", testmodels.Token.Token).
		Limit(1).
		CoerceTo("array")).
		Return([]map[string]interface{}{m}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", "").
		Limit(1).
		CoerceTo("array")).
		Return([]map[string]interface{}{}, nil))

	u, err := h.GetUserToken(testmodels.Token.Token)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if diff := pretty.Compare(testmodels.Token, u); diff != "" {
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
		GetAllByIndex("Token", testmodels.Token.Token).
		Limit(1).
		Delete()).
		Return(r.WriteResponse{Deleted: 1}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(userTokenTableInit.TableName).
		GetAllByIndex("Token", "").
		Limit(1).
		Delete()).
		Return(r.WriteResponse{}, fmt.Errorf("test error")))

	err := h.RemoveUserToken(testmodels.Token.Token)

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
