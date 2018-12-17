package remote_signer

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func getTableIndexes(tableName string) []string {
	var conn = GetConnection()

	c, err := r.Table(tableName).IndexList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	var idxI = z.([]interface{})
	var idx = make([]string, len(idxI))

	for i, v := range idxI {
		idx[i] = v.(string)
	}

	return idx
}

func getDatabases() []string {
	var conn = GetConnection()

	c, err := r.DBList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	var dbsI = z.([]interface{})
	var dbs = make([]string, len(dbsI))

	for i, v := range dbsI {
		dbs[i] = v.(string)
	}

	return dbs
}

func getTables() []string {
	var conn = GetConnection()

	c, err := r.TableList().CoerceTo("array").Run(conn)

	if err != nil {
		panic(err)
	}

	z, err := c.Interface()

	var tbsI = z.([]interface{})
	var tbs = make([]string, len(tbsI))

	for i, v := range tbsI {
		tbs[i] = v.(string)
	}

	return tbs
}
