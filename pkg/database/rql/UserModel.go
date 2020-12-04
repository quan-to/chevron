package rql

import (
	"fmt"
	"github.com/quan-to/chevron/pkg/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var userModelTableInit = tableInitStruct{
	TableName:    "users",
	TableIndexes: []string{"Username", "Fingerprint", "CreatedAt"},
}

func (h *RethinkDBDriver) initUserTable() error {
	return h.initFromStruct(userModelTableInit)
}

func (h *RethinkDBDriver) migrateUserTable() error {
	// Legacy fields
	h.log.Info("Migrating old fields to new fields")
	result, err := r.Table(userModelTableInit.TableName).
		Filter(r.Row.HasFields("FingerPrint")).
		Update(map[string]interface{}{
			"Fingerprint": r.Row.Field("FingerPrint"),
		}).RunWrite(h.conn)

	if err != nil {
		return err
	}

	h.log.Info("Migrated %d users FingerPrint -> Fingerprint", result.Updated)

	h.log.Info("Migrating old fields to new fields")
	result, err = r.Table(userModelTableInit.TableName).
		Filter(r.Row.HasFields("Fullname")).
		Update(map[string]interface{}{
			"FullName": r.Row.Field("Fullname"),
		}).RunWrite(h.conn)

	if err != nil {
		return err
	}

	h.log.Info("Migrated %d users Fullname -> FullName", result.Updated)
	return nil
}

func (h *RethinkDBDriver) AddUser(um models.User) (string, error) {
	existing, err := r.
		Table(userModelTableInit.TableName).
		GetAllByIndex("Username", um.Username).
		Run(h.conn)

	if err != nil {
		return "", err
	}

	defer existing.Close()

	if !existing.IsNil() {
		return "", fmt.Errorf("already exists")
	}

	rum, err := convertToRethinkDB(um)
	if err != nil {
		return "", err
	}

	wr, err := r.Table(userModelTableInit.TableName).
		Insert(rum).
		RunWrite(h.conn)

	if err != nil {
		return "", err
	}

	return wr.GeneratedKeys[0], err
}

func (h *RethinkDBDriver) GetUser(username string) (um *models.User, err error) {
	var res *r.Cursor
	res, err = r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", username).
		Limit(1).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return um, err
	}

	defer res.Close()

	rdata := map[string]interface{}{}

	if res.Next(&rdata) {
		err = convertFromRethinkDB(rdata, &um)
		return um, err
	}

	return um, fmt.Errorf("not found")
}

func (h *RethinkDBDriver) UpdateUser(um models.User) error {
	rum, err := convertToRethinkDB(um)
	if err != nil {
		return err
	}

	res, err := r.Table(userModelTableInit.TableName).
		GetAllByIndex("Username", um.Username).
		Update(rum).
		RunWrite(h.conn)

	if err != nil {
		return err
	}

	if res.Replaced == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
