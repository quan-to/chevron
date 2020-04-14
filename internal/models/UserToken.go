package models

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"time"
)

var UserTokenTableInit = TableInitStruct{
	TableName:    "tokens",
	TableIndexes: []string{"Token", "Username", "FingerPrint", "CreatedAt"},
}

type UserToken struct {
	id          string
	FingerPrint string
	Username    string
	Password    string
	Fullname    string
	Token       string
	CreatedAt   time.Time
	Expiration  time.Time
}

func (ut *UserToken) GetUsername() string {
	return ut.Username
}

func (ut *UserToken) GetFullName() string {
	return ut.Fullname
}

func (ut *UserToken) GetUserdata() interface{} {
	return nil
}

func (ut *UserToken) GetToken() string {
	return ut.Token
}

func (ut *UserToken) GetCreatedAt() time.Time {
	return ut.CreatedAt
}

func (ut *UserToken) GetFingerPrint() string {
	return ut.FingerPrint
}

func AddUserToken(conn *r.Session, ut *UserToken) (string, error) {
	wr, err := r.Table(UserTokenTableInit.TableName).
		Insert(ut).
		RunWrite(conn)

	if err != nil {
		return "", err
	}

	ut.id = wr.GeneratedKeys[0]

	return wr.GeneratedKeys[0], err
}

// RemoveUserToken removes a user token from the database
func RemoveUserToken(conn *r.Session, token string) (err error) {
	_, err = r.Table(UserTokenTableInit.TableName).
		GetAllByIndex("Token", token).
		Limit(1).
		Delete().
		RunWrite(conn)

	if err != nil {
		return err
	}

	return nil
}

func GetUserToken(conn *r.Session, token string) (ut *UserToken, err error) {
	var res *r.Cursor
	res, err = r.Table(UserTokenTableInit.TableName).
		GetAllByIndex("Token", token).
		Limit(1).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next(&ut) {
		return ut, nil
	}

	return nil, fmt.Errorf("not found")
}

func InvalidateUserTokens(conn *r.Session) (int, error) {
	wr, err := r.Table(UserTokenTableInit.TableName).
		Filter(r.Row.Field("Expiration").Lt(time.Now())).
		Delete().
		RunWrite(conn)

	if err != nil {
		return 0, err
	}

	return wr.Deleted, nil
}
