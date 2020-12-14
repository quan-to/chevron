package pg

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
	"github.com/quan-to/chevron/pkg/models"
)

var testUser = models.User{
	ID:          "abcd",
	Username:    "johnhuebr",
	FullName:    "John HUEBR",
	Fingerprint: "DEADBEEFDEADBEEF",
	Password:    "I think you will never guess",
	CreatedAt:   time.Now().Truncate(time.Second),
}

func expectUserSelect(mock sqlmock.Sqlmock) {
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
		testUser.ID,
		testUser.Fingerprint,
		testUser.Username,
		[]byte(testUser.Password),
		testUser.FullName,
		testUser.CreatedAt,
		time.Time{},
		(*time.Time)(nil),
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user WHERE user_username = $1 LIMIT 1`)).
		WithArgs(testUser.Username).
		WillReturnRows(expectedUserRows)
}

func TestPostgreSQLDBDriver_AddUser(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	testAdd := testUser
	testAdd.ID = ""

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user WHERE user_username = $1 LIMIT 1`)).
		WithArgs(testUser.Username).
		WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chevron_user(user_id, user_fingerprint, user_username, user_password, user_full_name, user_created_at) VALUES (?, ?, ?, ?, ?, now())`)).
		WithArgs(
			sqlmock.AnyArg(),
			testAdd.Fingerprint,
			testAdd.Username,
			[]byte(testAdd.Password),
			testAdd.FullName,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	_, err := h.AddUser(testAdd)

	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_GetUser(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	// Test Existing GET
	mock.ExpectBegin()
	expectUserSelect(mock)
	mock.ExpectCommit()

	u, err := h.GetUser(testUser.Username)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if u == nil {
		t.Fatalf(unexpectedError, "user came nil")
	}

	if diff := pretty.Compare(testUser, u); diff != "" {
		t.Errorf("Expected user to be the same. (-got +want)\\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}

	// Test not found
	mockDB, mock, _ = sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM chevron_user WHERE user_username = $1 LIMIT 1`)).
		WithArgs(testUser.Username).
		WillReturnError(fmt.Errorf("sql: no rows in result set"))
	mock.ExpectRollback()

	_, err = h.GetUser(testUser.Username)
	if err == nil {
		t.Fatalf(unexpectedError, "expected error to be not nil, got nil")
	}
	if !strings.EqualFold("not found", err.Error()) {
		t.Fatalf("expected error to be %q got %q", "not found", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}

func TestPostgreSQLDBDriver_UpdateUser(t *testing.T) {
	h := MakePostgreSQLDBDriver(nil)
	converter := sqlmock.ValueConverterOption(customConverter{})

	mockDB, mock, _ := sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE chevron_user SET user_fingerprint = ?, user_password = ?, user_full_name = ?, user_updated_at = now() WHERE user_id = ?`)).
		WithArgs(
			testUser.Fingerprint,
			[]byte(testUser.Password),
			testUser.FullName,
			testUser.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := h.UpdateUser(testUser)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}

	//Test without ID
	mockDB, mock, _ = sqlmock.New(converter)
	h.conn = sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	expectUserSelect(mock)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE chevron_user SET user_fingerprint = ?, user_password = ?, user_full_name = ?, user_updated_at = now() WHERE user_id = ?`)).
		WithArgs(
			testUser.Fingerprint,
			[]byte(testUser.Password),
			testUser.FullName,
			testUser.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	upUser := testUser
	upUser.ID = ""

	err = h.UpdateUser(upUser)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(expectationsDidNotMet, err)
	}
}
