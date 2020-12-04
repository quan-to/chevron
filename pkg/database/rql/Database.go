package rql

import (
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/tools"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"time"
)

const maxRetryCount = 5

func (h *RethinkDBDriver) ResetDatabase() error {
	h.log.Error("Resetting Database")
	dbs, err := h.getDatabases()

	if err != nil {
		return err
	}

	h.log.Error("Dropping test database %s", config.DatabaseName)
	if tools.StringIndexOf(config.DatabaseName, dbs) > -1 {
		h.log.Error("Test Database already exists, dropping.")
		_ = r.DBDrop(config.DatabaseName).Exec(h.conn)
	}

	err = h.waitDatabaseDrop(h.database)
	if err != nil {
		return err
	}
	h.log.Info("Database resetted")

	return nil
}

func (h *RethinkDBDriver) Connect(host, username, password, database string, port, poolSize int) error {
	h.log.Await("RethinkDB SKS Enabled. Starting %d connections to %s:%d", poolSize, host, port)
	conn, err := r.Connect(r.ConnectOpts{
		Address:    fmt.Sprintf("%s:%d", host, port),
		Username:   username,
		Password:   password,
		NumRetries: maxRetryCount,
		MaxOpen:    poolSize,
		InitialCap: poolSize,
		Database:   database,
	})

	if err != nil {
		return err
	}

	h.log.Done("Connected!")
	h.conn = conn
	return nil
}

func (h *RethinkDBDriver) waitTableCreate(table string) error {
	h.log.Await("Waiting for table %s create", table)
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		tables, err := h.getTables()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(table, tables) > -1 {
			break
		}
	}

	_ = r.DB(config.DatabaseName).
		Table(table).
		Wait(r.WaitOpts{Timeout: 0}).
		Exec(h.conn)

	h.log.Done("Done waiting table %s create", table)

	return nil
}

func (h *RethinkDBDriver) waitTableIndexCreate(table, index string) error {
	h.log.Await("Waiting for index %s/%s create", table, index)

	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		indexes, err := h.getTableIndexes(table)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(index, indexes) > -1 {
			break
		}
	}
	h.log.Done("Done waiting index %s/%s create", table, index)

	return nil
}

func (h *RethinkDBDriver) getTables() ([]string, error) {
	c, err := r.TableList().CoerceTo("array").Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	z, err := c.Interface()

	if err != nil {
		return nil, err
	}

	var tbsI = z.([]interface{})
	var tbs = make([]string, len(tbsI))

	for i, v := range tbsI {
		tbs[i] = v.(string)
	}

	return tbs, nil
}

func (h *RethinkDBDriver) getTableIndexes(tableName string) ([]string, error) {
	c, err := r.Table(tableName).IndexList().CoerceTo("array").Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	z, err := c.Interface()

	if err != nil {
		return nil, err
	}

	var idxI = z.([]interface{})
	var idx = make([]string, len(idxI))

	for i, v := range idxI {
		idx[i] = v.(string)
	}

	return idx, nil
}

func (h *RethinkDBDriver) getDatabases() ([]string, error) {
	c, err := r.DBList().CoerceTo("array").Run(h.conn)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	z, err := c.Interface()

	if err != nil {
		return nil, err
	}

	var dbsI = z.([]interface{})
	var dbs = make([]string, len(dbsI))

	for i, v := range dbsI {
		dbs[i] = v.(string)
	}

	return dbs, nil
}

func (h *RethinkDBDriver) waitDatabaseDrop(database string) error {
	h.log.Info("Waiting for database drop")
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		dbs, err := h.getDatabases()
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(database, dbs) == -1 {
			break
		}
	}

	return nil
}

func (h *RethinkDBDriver) waitDatabaseCreate(database string) error {
	h.log.Await("Waiting for database create")
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		dbs, err := h.getDatabases()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(database, dbs) > -1 {
			break
		}
	}

	h.log.Done("Done waiting database create")

	return nil
}
