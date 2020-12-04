package rql

import (
	"fmt"
	"github.com/quan-to/chevron/pkg/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var gpgKeyTableInit = tableInitStruct{
	TableName:    "gpgKey",
	TableIndexes: []string{"FullFingerprint", "Names", "Emails", "Subkeys"},
}

func (h *RethinkDBDriver) initGPGKeyTable() error {
	return h.initFromStruct(gpgKeyTableInit)
}

func (h *RethinkDBDriver) migrateGPGKeyTable() error {
	// Legacy fields
	h.log.Info("Migrating old fields to new fields")
	result, err := r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.HasFields("FullFingerPrint")).
		Update(map[string]interface{}{
			"FullFingerprint": r.Row.Field("FullFingerPrint"),
		}).RunWrite(h.conn)

	if err != nil {
		return err
	}

	h.log.Info("Migrated %d tokens FullFingerPrint -> FullFingerprint", result.Updated)
	return nil
}

func (h *RethinkDBDriver) fixGPGKey(k map[string]interface{}) map[string]interface{} {
	if fp, ok := k["FullFingerPrint"]; ok {
		k["FullFingerprint"] = fp
		delete(k, "FullFingerPrint")
	}
	return k
}

func (h *RethinkDBDriver) UpdateGPGKey(key models.GPGKey) error {
	rdata, err := convertToRethinkDB(key)
	if err != nil {
		return err
	}

	return r.Table(gpgKeyTableInit.TableName).
		Get(key.ID).
		Update(rdata).
		Exec(h.conn)
}

func (h *RethinkDBDriver) DeleteGPGKey(key models.GPGKey) error {
	return r.Table(gpgKeyTableInit.TableName).
		Get(key.ID).
		Delete().
		Exec(h.conn)
}

// AddGPGKey adds a GPG Key to the database or update an existing one by fingerprint
// Returns generated id / hasBeenAdded / error
func (h *RethinkDBDriver) AddGPGKey(key models.GPGKey) (string, bool, error) {
	existing, err := r.
		Table(gpgKeyTableInit.TableName).
		GetAllByIndex("FullFingerprint", key.FullFingerprint).
		Run(h.conn)

	if err != nil {
		return "", false, err
	}

	defer existing.Close()

	rdata, err := convertToRethinkDB(key)
	if err != nil {
		return "", false, err
	}

	gpgKey := map[string]interface{}{}
	exists := existing.Next(&gpgKey)
	keyId, ok := gpgKey["id"]

	if exists && ok {
		stringKeyId := keyId.(string)
		// Update
		_, err := r.Table(gpgKeyTableInit.TableName).
			Get(stringKeyId).
			Update(rdata).
			RunWrite(h.conn)

		if err != nil {
			return "", false, err
		}

		return stringKeyId, false, err
	}

	// Create
	wr, err := r.Table(gpgKeyTableInit.TableName).
		Insert(rdata).
		RunWrite(h.conn)

	if err != nil {
		return "", false, err
	}
	return wr.GeneratedKeys[0], true, err
}

func (h *RethinkDBDriver) FetchGPGKeysWithoutSubKeys() ([]models.GPGKey, error) {
	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.HasFields("Subkeys").Not().Or(r.Row.Field("Subkeys").Count().Eq(0))).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}

func (h *RethinkDBDriver) FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error) {
	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerprint").Match(fmt.Sprintf("%s$", fingerprint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", fingerprint))
			}).Count().Gt(0)))).
		Limit(1).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	var gpgKey map[string]interface{}

	if res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		return &key, nil
	}

	return nil, fmt.Errorf("not found")
}

func (h *RethinkDBDriver) FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var filterEmailList = func(r r.Term) interface{} {
		return r.Match(email)
	}
	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(func(r r.Term) interface{} {
			return r.Field("Emails").
				Filter(filterEmailList).
				Count().
				Gt(0)
		}).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()
	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}

func (h *RethinkDBDriver) FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(r.Row.Field("FullFingerprint").Match(fmt.Sprintf("%s$", fingerPrint)).
			Or(r.Row.HasFields("Subkeys").And(r.Row.Field("Subkeys").Filter(func(p r.Term) interface{} {
				return p.Match(fmt.Sprintf("%s$", fingerPrint))
			}).Count().Gt(0)))).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()
	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}

func (h *RethinkDBDriver) FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}
	var filterEmailList = func(r r.Term) interface{} {
		return r.Match(value)
	}

	var filterNames = func(r r.Term) interface{} {
		return r.Match(value)
	}

	var filterSub = func(r r.Term) interface{} {
		return r.Field("Emails").Filter(filterEmailList).Count().Gt(0).
			Or(r.Field("Names").Filter(filterNames).Count().Gt(0)).
			Or(r.Field("FullFingerprint").Match(fmt.Sprintf("%s$", value)))
	}

	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(filterSub).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}

func (h *RethinkDBDriver) FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	if pageStart < 0 {
		pageStart = models.DefaultPageStart
	}

	if pageEnd < 0 {
		pageEnd = models.DefaultPageEnd
	}

	var filterNames = func(r r.Term) interface{} {
		return r.Match(name)
	}
	res, err := r.Table(gpgKeyTableInit.TableName).
		Filter(func(r r.Term) interface{} {
			return r.Field("Names").
				Filter(filterNames).
				Count().
				Gt(0)
		}).
		Slice(pageStart, pageEnd).
		CoerceTo("array").
		Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer res.Close()

	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}
