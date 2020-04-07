package models

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"time"
)

var UserModelTableInit = TableInitStruct{
	TableName:    "users",
	TableIndexes: []string{"Username", "FingerPrint", "CreatedAt"},
}

type UserModel struct {
	id          string
	FingerPrint string
	Username    string
	Password    string
	Fullname    string
	CreatedAt   time.Time
}

func (um *UserModel) GetUsername() string {
	return um.Username
}

func (um *UserModel) GetFullName() string {
	return um.Fullname
}

func (um *UserModel) GetUserdata() interface{} {
	return nil
}

func (um *UserModel) GetToken() string {
	return ""
}

func (um *UserModel) GetCreatedAt() time.Time {
	return um.CreatedAt
}

func (um *UserModel) GetFingerPrint() string {
	return um.FingerPrint
}

func AddUser(conn *r.Session, um *UserModel) (string, error) {
	existing, err := r.
		Table(UserModelTableInit.TableName).
		GetAllByIndex("Username", um.Username).
		Run(conn)

	if err != nil {
		return "", err
	}

	defer existing.Close()

	if !existing.IsNil() {
		return "", fmt.Errorf("already exists")
	}

	wr, err := r.Table(UserModelTableInit.TableName).
		Insert(um).
		RunWrite(conn)

	if err != nil {
		return "", err
	}

	um.id = wr.GeneratedKeys[0]

	return wr.GeneratedKeys[0], err
}

func GetUser(conn *r.Session, username string) (um *UserModel, err error) {
	var res *r.Cursor
	res, err = r.Table(UserModelTableInit.TableName).
		GetAllByIndex("Username", username).
		Limit(1).
		CoerceTo("array").
		Run(conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next(&um) {
		return um, nil
	}

	return nil, fmt.Errorf("not found")
}

func UpdateUser(conn *r.Session, um *UserModel) error {
	res, err := r.Table(UserModelTableInit.TableName).
		GetAllByIndex("Username", um.Username).
		Update(um).
		RunWrite(conn)

	if err != nil {
		return err
	}

	if res.Replaced == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
