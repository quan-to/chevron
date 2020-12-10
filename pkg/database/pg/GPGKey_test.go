package pg

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/uuid"
	"github.com/quan-to/slog"
)

var testGPGKey = models.GPGKey{
	ID:              "abcd",
	FullFingerprint: "DEADBEEFDEADBEEFDEADBEEFDEADBEEF",
	Names:           []string{"A", "B"},
	Emails:          []string{"a@a.com", "b@a.com"},
	KeyUids: []models.GPGKeyUid{
		{
			Name:        "A",
			Email:       "a@a.com",
			Description: "desc",
		},
		{
			Name:        "B",
			Email:       "b@a.com",
			Description: "desc",
		},
	},
	KeyBits:                1234,
	Subkeys:                []string{"BABABEBE"},
	AsciiArmoredPublicKey:  "PUBKEY",
	AsciiArmoredPrivateKey: "PRIVKEY",
}

const unexpectedError = "unexpected error %q"
const expectationsDidNotMet = "expectations did not met: %s"

func init() {
	slog.SetTestMode()
}

type customConverter struct{}

func (s customConverter) ConvertValue(v interface{}) (driver.Value, error) {
	if val, ok1 := v.(driver.NamedValue); ok1 {
		return val.Value, nil
	}

	return v, nil
}

func expectToUpdate(mock sqlmock.Sqlmock, withSubkeys bool) {
	mock.ExpectBegin()
	// Check existance
	testRow := sqlmock.NewRows([]string{
		"gpg_key_id",
		"gpg_key_full_fingerprint",
		"gpg_key_fingerprint16",
		"gpg_key_keybits",
		"gpg_key_public_key",
		"gpg_key_private_key",
		"gpg_key_created_at",
		"gpg_key_updated_at",
		"gpg_key_deleted_at",
		"gpg_key_parent",
	})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1`)).
		WithArgs(tools.FPto16(testGPGKey.FullFingerprint)).
		WillReturnRows(testRow.AddRow(
			testGPGKey.ID,
			testGPGKey.FullFingerprint,
			tools.FPto16(testGPGKey.FullFingerprint),
			testGPGKey.KeyBits,
			testGPGKey.AsciiArmoredPublicKey,
			testGPGKey.AsciiArmoredPrivateKey,
			time.Now(),
			time.Now(),
			time.Time{},
			(*string)(nil),
		))
	if withSubkeys {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1`)).
			WithArgs(tools.FPto16(testGPGKey.FullFingerprint)).
			WillReturnRows(testRow.AddRow(
				testGPGKey.ID,
				testGPGKey.FullFingerprint,
				tools.FPto16(testGPGKey.FullFingerprint),
				testGPGKey.KeyBits,
				testGPGKey.AsciiArmoredPublicKey,
				testGPGKey.AsciiArmoredPrivateKey,
				time.Now(),
				time.Now(),
				time.Time{},
				(*string)(nil),
			))
	}
	// Load UIDs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1`)).
		WithArgs(testGPGKey.ID).
		WillReturnRows(sqlmock.NewRows(nil))

	// Insert UIDs
	for _, uid := range testGPGKey.KeyUids {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_gpg_key_uid(gpg_key_uid_id, gpg_key_uid_name, gpg_key_uid_email, gpg_key_uid_description, gpg_key_uid_parent) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), uid.Name, uid.Email, uid.Description, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	// Update Key
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE chevron_gpg_key SET gpg_key_private_key = ?, gpg_key_public_key = ? WHERE gpg_key_id = ?`)).
		WithArgs(testGPGKey.AsciiArmoredPrivateKey, testGPGKey.AsciiArmoredPublicKey, testGPGKey.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

func TestPostgreSQLDBDriver_AddGPGKey(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	// Check existance
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1`)).
		WithArgs(tools.FPto16(testGPGKey.FullFingerprint)).
		WillReturnError(fmt.Errorf("sql: no rows in result set"))
	// Insert Key
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_gpg_key(gpg_key_id, gpg_key_full_fingerprint, gpg_key_fingerprint16, gpg_key_keybits, gpg_key_parent, gpg_key_public_key, gpg_key_private_key) VALUES (?, ?, ?, ?, ?, ?, ?)`)).
		WithArgs(sqlmock.AnyArg(), testGPGKey.FullFingerprint, tools.FPto16(testGPGKey.FullFingerprint), testGPGKey.KeyBits, (*string)(nil), testGPGKey.AsciiArmoredPublicKey, testGPGKey.AsciiArmoredPrivateKey).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Insert UIDs
	for _, uid := range testGPGKey.KeyUids {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_gpg_key_uid(gpg_key_uid_id, gpg_key_uid_name, gpg_key_uid_email, gpg_key_uid_description, gpg_key_uid_parent) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), uid.Name, uid.Email, uid.Description, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()

	_, added, err := h.AddGPGKey(testGPGKey)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if !added {
		t.Fatalf("expected added but got updated")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
	// endregion
	// region Test UPDATE
	// Test existing key
	mockDB, mock, _ = sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	expectToUpdate(mock, true)

	_, added, err = h.AddGPGKey(testGPGKey)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if added {
		t.Fatalf("expected updated but got added")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
	// endregion
}

func TestPostgreSQLDBDriver_DeleteGPGKey(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM chevron_gpg_key WHERE gpg_key_id = ?`)).
		WithArgs(testGPGKey.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := h.DeleteGPGKey(testGPGKey)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FetchGPGKeyByFingerprint(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	expectedGPGKeyRow := sqlmock.NewRows([]string{
		"gpg_key_id",
		"gpg_key_full_fingerprint",
		"gpg_key_fingerprint16",
		"gpg_key_keybits",
		"gpg_key_public_key",
		"gpg_key_private_key",
		"gpg_key_created_at",
		"gpg_key_updated_at",
		"gpg_key_deleted_at",
		"gpg_key_parent",
	}).AddRow(
		testGPGKey.ID,
		testGPGKey.FullFingerprint,
		tools.FPto16(testGPGKey.FullFingerprint),
		testGPGKey.KeyBits,
		testGPGKey.AsciiArmoredPublicKey,
		testGPGKey.AsciiArmoredPrivateKey,
		time.Now(),
		time.Now(),
		time.Time{},
		(*string)(nil),
	)

	expectedSubGPGKeyRow := sqlmock.NewRows([]string{
		"gpg_key_id",
		"gpg_key_full_fingerprint",
		"gpg_key_fingerprint16",
		"gpg_key_keybits",
		"gpg_key_public_key",
		"gpg_key_private_key",
		"gpg_key_created_at",
		"gpg_key_updated_at",
		"gpg_key_deleted_at",
		"gpg_key_parent",
	})

	for _, v := range testGPGKey.Subkeys {
		expectedSubGPGKeyRow = expectedSubGPGKeyRow.AddRow(
			uuid.EnsureUUID(h.log),
			v,
			v,
			testGPGKey.KeyBits,
			testGPGKey.AsciiArmoredPublicKey,
			testGPGKey.AsciiArmoredPrivateKey,
			time.Now(),
			time.Now(),
			time.Time{},
			testGPGKey.ID,
		)
	}

	expectedGPGKeyUIDRows := sqlmock.NewRows([]string{
		"gpg_key_uid_id",
		"gpg_key_uid_parent",
		"gpg_key_uid_name",
		"gpg_key_uid_email",
		"gpg_key_uid_description",
		"gpg_key_uid_created_at",
		"gpg_key_uid_updated_at",
		"gpg_key_uid_deleted_at",
	})

	for _, v := range testGPGKey.KeyUids {
		expectedGPGKeyUIDRows = expectedGPGKeyUIDRows.AddRow(
			uuid.EnsureUUID(h.log),
			testGPGKey.ID,
			v.Name,
			v.Email,
			v.Description,
			time.Now(),
			time.Now(),
			time.Time{},
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1`)).
		WithArgs(tools.FPto16(testGPGKey.FullFingerprint)).
		WillReturnRows(expectedGPGKeyRow)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1`)).
		WithArgs(testGPGKey.ID).
		WillReturnRows(expectedGPGKeyUIDRows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_parent = $1`)).
		WithArgs(testGPGKey.ID).
		WillReturnRows(expectedSubGPGKeyRow)

	mock.ExpectCommit()
	key, err := h.FetchGPGKeyByFingerprint(testGPGKey.FullFingerprint)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if key == nil {
		t.Fatalf("unexpected nil key")
	}
	if diff := pretty.Compare(testGPGKey, key); diff != "" {
		t.Errorf("Expected gpgKey to be the same. (-got +want)\\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FetchGPGKeysWithoutSubKeys(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	_, err := h.FetchGPGKeysWithoutSubKeys()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not supported") {
		t.Fatalf("expected error %q got %q", "not supported", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FindGPGKeyByEmail(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	_, err := h.FindGPGKeyByEmail("", 0, 0)
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not supported") {
		t.Fatalf("expected error %q got %q", "not supported", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FindGPGKeyByFingerPrint(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()

	expectedRows := sqlmock.NewRows([]string{
		"gpg_key_id",
		"gpg_key_full_fingerprint",
		"gpg_key_fingerprint16",
		"gpg_key_keybits",
		"gpg_key_public_key",
		"gpg_key_private_key",
		"gpg_key_created_at",
		"gpg_key_updated_at",
		"gpg_key_deleted_at",
		"gpg_key_parent",
	}).AddRow(
		testGPGKey.ID,
		testGPGKey.FullFingerprint,
		tools.FPto16(testGPGKey.FullFingerprint),
		testGPGKey.KeyBits,
		testGPGKey.AsciiArmoredPublicKey,
		testGPGKey.AsciiArmoredPrivateKey,
		time.Now(),
		time.Now(),
		time.Time{},
		(*string)(nil),
	).AddRow(
		testGPGKey.ID+"1234",
		testGPGKey.FullFingerprint,
		tools.FPto16(testGPGKey.FullFingerprint),
		testGPGKey.KeyBits,
		testGPGKey.AsciiArmoredPublicKey,
		testGPGKey.AsciiArmoredPrivateKey,
		time.Now(),
		time.Now(),
		time.Time{},
		(*string)(nil),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 LIKE $1 LIMIT $2 OFFSET $3`)).
		WithArgs(
			"%"+tools.FPto16(testGPGKey.FullFingerprint),
			10,
			0,
		).
		WillReturnRows(expectedRows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1`)).
		WithArgs(testGPGKey.ID).
		WillReturnRows(sqlmock.NewRows(nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_parent = $1`)).
		WithArgs(testGPGKey.ID).
		WillReturnRows(sqlmock.NewRows(nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1`)).
		WithArgs(testGPGKey.ID + "1234").
		WillReturnRows(sqlmock.NewRows(nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_parent = $1`)).
		WithArgs(testGPGKey.ID + "1234").
		WillReturnRows(sqlmock.NewRows(nil))

	mock.ExpectCommit()

	results, err := h.FindGPGKeyByFingerPrint(testGPGKey.FullFingerprint, 0, 10)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if len(results) == 0 {
		t.Fatalf("expected %d results but got %d", 2, len(results))
	}

	if results[0].ID != testGPGKey.ID {
		t.Fatalf("expected ID from row 0 to be %q but got %q", testGPGKey.ID, results[0].ID)
	}

	if results[1].ID != testGPGKey.ID+"1234" {
		t.Fatalf("expected ID from row 0 to be %q but got %q", testGPGKey.ID+"1234", results[1].ID)
	}

	for i, v := range results {
		if v.FullFingerprint != testGPGKey.FullFingerprint ||
			v.AsciiArmoredPublicKey != testGPGKey.AsciiArmoredPublicKey ||
			v.AsciiArmoredPrivateKey != testGPGKey.AsciiArmoredPrivateKey ||
			v.KeyBits != testGPGKey.KeyBits {
			t.Fatalf("expected result %d to be equal to testKey", i)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FindGPGKeyByName(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	_, err := h.FindGPGKeyByName("", 0, 0)
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not supported") {
		t.Fatalf("expected error %q got %q", "not supported", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FindGPGKeyByValue(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	_, err := h.FindGPGKeyByValue("", 0, 0)
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !strings.EqualFold(err.Error(), "not supported") {
		t.Fatalf("expected error %q got %q", "not supported", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_UpdateGPGKey(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	// Test existing key
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	expectToUpdate(mock, false)

	err := h.UpdateGPGKey(testGPGKey)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}
