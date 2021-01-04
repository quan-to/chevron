package pg

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/chevron/pkg/models/testmodels"
	"github.com/quan-to/chevron/pkg/uuid"
)

func TestPostgreSQLDBDriver_InitCursor(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key`)).
		WillReturnRows(mock.NewRows(nil))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user`)).
		WillReturnRows(mock.NewRows(nil))

	err := h.InitCursor()

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if h.gpgKeysRows == nil {
		t.Fatal("expected h.gpgKeyRows but got nil")
	}

	if h.usersRows == nil {
		t.Fatal("expected h.usersRows but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_FinishCursor(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key`)).
		WillReturnRows(mock.NewRows(nil))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user`)).
		WillReturnRows(mock.NewRows(nil))

	_ = h.InitCursor()
	err := h.FinishCursor()

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if h.usersRows != nil {
		t.Fatal("expected h.usersRows to be cleaned up")
	}
	if h.gpgKeysRows != nil {
		t.Fatal("expected h.gpgKeysRows to be cleaned up")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_NextGPGKey(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

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
		testmodels.GpgKey.ID,
		testmodels.GpgKey.FullFingerprint,
		tools.FPto16(testmodels.GpgKey.FullFingerprint),
		testmodels.GpgKey.KeyBits,
		testmodels.GpgKey.AsciiArmoredPublicKey,
		testmodels.GpgKey.AsciiArmoredPrivateKey,
		time.Now(),
		time.Now(),
		time.Time{},
		(*string)(nil),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key`)).
		WillReturnRows(expectedGPGKeyRow)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user`)).
		WillReturnRows(mock.NewRows(nil))
	mock.ExpectBegin()

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

	for _, v := range testmodels.GpgKey.KeyUids {
		expectedGPGKeyUIDRows = expectedGPGKeyUIDRows.AddRow(
			uuid.EnsureUUID(h.log),
			testmodels.GpgKey.ID,
			v.Name,
			v.Email,
			v.Description,
			time.Now(),
			time.Now(),
			time.Time{},
		)
	}

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

	for _, v := range testmodels.GpgKey.Subkeys {
		expectedSubGPGKeyRow = expectedSubGPGKeyRow.AddRow(
			uuid.EnsureUUID(h.log),
			v,
			v,
			testmodels.GpgKey.KeyBits,
			testmodels.GpgKey.AsciiArmoredPublicKey,
			testmodels.GpgKey.AsciiArmoredPrivateKey,
			time.Now(),
			time.Now(),
			time.Time{},
			testmodels.GpgKey.ID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key_uid WHERE gpg_key_uid_parent = $1`)).
		WithArgs(testmodels.GpgKey.ID).
		WillReturnRows(expectedGPGKeyUIDRows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_parent = $1`)).
		WithArgs(testmodels.GpgKey.ID).
		WillReturnRows(expectedSubGPGKeyRow)
	mock.ExpectRollback()

	fetchKey := models.GPGKey{}
	_ = h.InitCursor()
	err := h.NextGPGKey(&fetchKey)
	if !err {
		t.Fatal("expected to fetch one key but got 0")
	}

	if diff := pretty.Compare(testmodels.GpgKey, fetchKey); diff != "" {
		t.Errorf("Expected gpgKey to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_NextUser(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	expectedUserRows := sqlmock.NewRows([]string{
		"user_id",
		"user_fingerprint",
		"user_username",
		"user_password",
		"user_full_name",
		"user_created_at",
		"user_updated_at",
		"user_deleted_at",
	}).AddRow(
		testmodels.User.ID,
		testmodels.User.Fingerprint,
		testmodels.User.Username,
		[]byte(testmodels.User.Password),
		testmodels.User.FullName,
		testmodels.User.CreatedAt,
		time.Time{},
		(*time.Time)(nil),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key`)).
		WillReturnRows(mock.NewRows(nil))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user`)).
		WillReturnRows(expectedUserRows)

	fetchUser := models.User{}
	_ = h.InitCursor()
	err := h.NextUser(&fetchUser)
	if !err {
		t.Fatal("expected to fetch one user but got 0")
	}

	if diff := pretty.Compare(testmodels.User, fetchUser); diff != "" {
		t.Errorf("Expected user to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_NumGPGKeys(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.NewRows([]string{"count"}).AddRow("100")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM chevron_gpg_key`)).
		WillReturnRows(rows)

	num, err := h.NumGPGKeys()
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if num != 100 {
		t.Fatalf("expected 100 rows got %d", num)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_AddGPGKeys(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})
	// region Test ADD
	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	// Check existance
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_gpg_key WHERE gpg_key_fingerprint16 = $1 LIMIT 1`)).
		WithArgs(tools.FPto16(testmodels.GpgKey.FullFingerprint)).
		WillReturnError(fmt.Errorf("sql: no rows in result set"))
	// Insert Key
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_gpg_key(gpg_key_id, gpg_key_full_fingerprint, gpg_key_fingerprint16, gpg_key_keybits, gpg_key_parent, gpg_key_public_key, gpg_key_private_key) VALUES (?, ?, ?, ?, ?, ?, ?)`)).
		WithArgs(sqlmock.AnyArg(), testmodels.GpgKey.FullFingerprint, tools.FPto16(testmodels.GpgKey.FullFingerprint), testmodels.GpgKey.KeyBits, (*string)(nil), testmodels.GpgKey.AsciiArmoredPublicKey, testmodels.GpgKey.AsciiArmoredPrivateKey).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Insert UIDs
	for _, uid := range testmodels.GpgKey.KeyUids {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_gpg_key_uid(gpg_key_uid_id, gpg_key_uid_name, gpg_key_uid_email, gpg_key_uid_description, gpg_key_uid_parent) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), uid.Name, uid.Email, uid.Description, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()

	_, added, err := h.AddGPGKeys([]models.GPGKey{testmodels.GpgKey})
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if len(added) == 0 || !added[0] {
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

	_, added, err = h.AddGPGKeys([]models.GPGKey{testmodels.GpgKey})
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}
	if len(added) == 0 || added[0] {
		t.Fatalf("expected updated but got added")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
	// endregion
}
