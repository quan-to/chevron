package database

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func GetTableIndexes(tableName string) []string {
	var conn = GetConnection()

	c, err := r.Table(tableName).IndexList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	if err != nil {
		panic(err)
	}

	var idxI = z.([]interface{})
	var idx = make([]string, len(idxI))

	for i, v := range idxI {
		idx[i] = v.(string)
	}

	return idx
}

func GetDatabases() []string {
	var conn = GetConnection()

	c, err := r.DBList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	if err != nil {
		panic(err)
	}

	var dbsI = z.([]interface{})
	var dbs = make([]string, len(dbsI))

	for i, v := range dbsI {
		dbs[i] = v.(string)
	}

	return dbs
}

func GetTables() []string {
	var conn = GetConnection()

	c, err := r.TableList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	if err != nil {
		panic(err)
	}

	var tbsI = z.([]interface{})
	var tbs = make([]string, len(tbsI))

	for i, v := range tbsI {
		tbs[i] = v.(string)
	}

	return tbs
}
