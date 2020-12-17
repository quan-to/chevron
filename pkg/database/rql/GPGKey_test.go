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

func TestRethinkDBDriver_AddGPGKey(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	toUpdate := testmodels.GpgKey
	toUpdate.FullFingerprint = "HUEBR"
	toUpdate.ID = "A123"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", testmodels.GpgKey.FullFingerprint)).
		Return(nil, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", toUpdate.FullFingerprint)).
		Return([]map[string]interface{}{
			{"id": toUpdate.ID},
		}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", "ERR")).
		Return(nil, fmt.Errorf("test error")))

	m, _ := convertToRethinkDB(testmodels.GpgKey)
	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).Insert(m)).Return(r.WriteResponse{
		GeneratedKeys: []string{testmodels.GpgKey.ID}, Inserted: 1,
	}, nil))

	m2, _ := convertToRethinkDB(toUpdate)
	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).Get(toUpdate.ID).Update(m2)).Return(r.WriteResponse{
		Updated: 1,
	}, nil))

	// Test Create
	id, added, err := h.AddGPGKey(testmodels.GpgKey)

	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if !added {
		t.Fatalf("expected item to be added")
	}

	if id != testmodels.GpgKey.ID {
		t.Fatalf("expected id to be %q but got %q", testmodels.GpgKey.ID, id)
	}

	// Test Update

	id, added, err = h.AddGPGKey(toUpdate)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if added {
		t.Fatalf("expected item to be updated not added")
	}

	if id != toUpdate.ID {
		t.Fatalf("expected id to be %q got %q", toUpdate.ID, id)
	}

	// Test Error
	_, _, err = h.AddGPGKey(models.GPGKey{FullFingerprint: "ERR"})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_DeleteGPGKey(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Get(testmodels.GpgKey.ID).
		Delete()).Return(nil, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Get("ERR").
		Delete()).Return(nil, fmt.Errorf("test error")))

	err := h.DeleteGPGKey(testmodels.GpgKey)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	err = h.DeleteGPGKey(models.GPGKey{ID: "ERR"})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FetchGPGKeyByFingerprint(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerprint").Match(fmt.Sprintf("%s$", testmodels.GpgKey.FullFingerprint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", testmodels.GpgKey.FullFingerprint))
			}).Count().Gt(0)))).
		Limit(1).
		CoerceTo("array")).Return([]map[string]interface{}{m}, nil))

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerprint").Match(fmt.Sprintf("%s$", "ERR")).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", "ERR"))
			}).Count().Gt(0)))).
		Limit(1).
		CoerceTo("array")).Return(nil, fmt.Errorf("test error")))

	gpgKey, err := h.FetchGPGKeyByFingerprint(testmodels.GpgKey.FullFingerprint)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if diff := pretty.Compare(testmodels.GpgKey, *gpgKey); diff != "" {
		t.Errorf("Expected gpgKey to be the same. (-got +want)\\n%s", diff)
	}

	_, err = h.FetchGPGKeyByFingerprint("ERR")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q but got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FetchGPGKeysWithoutSubKeys(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.HasFields("Subkeys").Not().Or(r.Row.Field("Subkeys").Count().Eq(0))).
		CoerceTo("array")).
		Return([]map[string]interface{}{m, m}, nil))

	keys, err := h.FetchGPGKeysWithoutSubKeys()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected exactly two keys got %d", len(keys))
	}

	for i, v := range keys {
		if diff := pretty.Compare(testmodels.GpgKey, v); diff != "" {
			t.Errorf("[%d] Expected gpgKey to be the same. (-got +want)\\n%s", i, diff)
		}
	}
	mock.AssertExpectations(t)

	// test error
	mock = r.NewMock()
	h.conn = mock

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.HasFields("Subkeys").Not().Or(r.Row.Field("Subkeys").Count().Eq(0))).
		CoerceTo("array")).
		Return(nil, fmt.Errorf("test error")))

	_, err = h.FetchGPGKeysWithoutSubKeys()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.EqualFold(err.Error(), "test error") {
		t.Fatalf("expected error to be %q got %q", "test error", err.Error())
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FindGPGKeyByEmail(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)

	testEmail := "a@a.com"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(func(v r.Term) interface{} {
			return v.Field("Emails").
				Filter(func(t r.Term) interface{} {
					return t.Match(testEmail)
				}).
				Count().
				Gt(0)
		}).
		Slice(models.DefaultPageStart, models.DefaultPageEnd).
		CoerceTo("array")).Return([]map[string]interface{}{m, m}, nil))

	keys, err := h.FindGPGKeyByEmail(testEmail, models.DefaultPageStart, models.DefaultPageEnd)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected exactly two keys got %d", len(keys))
	}

	for i, v := range keys {
		if diff := pretty.Compare(testmodels.GpgKey, v); diff != "" {
			t.Errorf("[%d] Expected gpgKey to be the same. (-got +want)\\n%s", i, diff)
		}
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FindGPGKeyByFingerPrint(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerprint").Match(fmt.Sprintf("%s$", testmodels.GpgKey.FullFingerprint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", testmodels.GpgKey.FullFingerprint))
			}).Count().Gt(0)))).
		Slice(models.DefaultPageStart, models.DefaultPageEnd).
		CoerceTo("array")).Return([]map[string]interface{}{m, m}, nil))

	keys, err := h.FindGPGKeyByFingerPrint(testmodels.GpgKey.FullFingerprint, models.DefaultPageStart, models.DefaultPageEnd)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected exactly two keys got %d", len(keys))
	}

	for i, v := range keys {
		if diff := pretty.Compare(testmodels.GpgKey, v); diff != "" {
			t.Errorf("[%d] Expected gpgKey to be the same. (-got +want)\\n%s", i, diff)
		}
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FindGPGKeyByName(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)
	testName := "testName"

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(func(v r.Term) interface{} {
			return v.Field("Names").
				Filter(func(t r.Term) interface{} {
					return t.Match(testName)
				}).
				Count().
				Gt(0)
		}).
		Slice(models.DefaultPageStart, models.DefaultPageEnd).
		CoerceTo("array")).Return([]map[string]interface{}{m, m}, nil))

	keys, err := h.FindGPGKeyByName(testName, models.DefaultPageStart, models.DefaultPageEnd)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected exactly two keys got %d", len(keys))
	}

	for i, v := range keys {
		if diff := pretty.Compare(testmodels.GpgKey, v); diff != "" {
			t.Errorf("[%d] Expected gpgKey to be the same. (-got +want)\\n%s", i, diff)
		}
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_FindGPGKeyByValue(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)
	testName := "testName"

	var filterEmailList = func(r r.Term) interface{} {
		return r.Match(testName)
	}

	var filterNames = func(r r.Term) interface{} {
		return r.Match(testName)
	}

	var filterSub = func(r r.Term) interface{} {
		return r.Field("Emails").Filter(filterEmailList).Count().Gt(0).
			Or(r.Field("Names").Filter(filterNames).Count().Gt(0)).
			Or(r.Field("FullFingerprint").Match(fmt.Sprintf("%s$", testName)))
	}

	mock.ExpectedQueries = append(mock.ExpectedQueries, mock.On(r.Table(gpgKeyTableInit.TableName).
		Filter(filterSub).
		Slice(models.DefaultPageStart, models.DefaultPageEnd).
		CoerceTo("array")).Return([]map[string]interface{}{m, m}, nil))

	keys, err := h.FindGPGKeyByValue(testName, models.DefaultPageStart, models.DefaultPageEnd)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected exactly two keys got %d", len(keys))
	}

	for i, v := range keys {
		if diff := pretty.Compare(testmodels.GpgKey, v); diff != "" {
			t.Errorf("[%d] Expected gpgKey to be the same. (-got +want)\\n%s", i, diff)
		}
	}

	mock.AssertExpectations(t)
}

func TestRethinkDBDriver_UpdateGPGKey(t *testing.T) {
	mock := r.NewMock()
	h := MakeRethinkDBDriver(slog.Scope("TEST"))
	h.conn = mock

	m, _ := convertToRethinkDB(testmodels.GpgKey)

	mock.ExpectedQueries = append(mock.ExpectedQueries,
		mock.On(r.Table(gpgKeyTableInit.TableName).
			Get(testmodels.GpgKey.ID).
			Update(m)).Return(nil, nil))

	err := h.UpdateGPGKey(testmodels.GpgKey)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	mock.AssertExpectations(t)
}
