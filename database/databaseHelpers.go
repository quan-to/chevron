package database

import (
	"github.com/quan-to/remote-signer"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"time"
)

func GetTableIndexes(tableName string) []string {
	dbLog.Debug("Listing table indexes")
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
	dbLog.Debug("Listing databases")
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
	dbLog.Debug("Listing tables")
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

func NumNodes() int {
	var conn = GetConnection()
	c, err := r.DB("rethinkdb").
		Table("stats").
		Filter(
			r.Row.HasFields("server").
				And(
					r.Row.HasFields("table").Not(),
				)).
		Count().
		Run(conn)
	if err != nil {
		dbLog.Error("Error fetching number of nodes: %v", err)
		return 1
	}

	z, err := c.Interface()

	if err != nil {
		dbLog.Error("Error fetching number of nodes: %v", err)
		return 1
	}

	count := int(z.(float64))

	dbLog.Debug("Number of RethinkDB Nodes: %d", count)

	return count
}

func WaitDatabaseDrop(database string) {
	nodes := NumNodes()
	for i := 0; i < nodes; i++ {
		dbs := GetDatabases()
		for remote_signer.StringIndexOf(database, dbs) > -1 {
			time.Sleep(1 * time.Second)
			dbs = GetDatabases()
		}
	}
}

func WaitDatabaseCreate(database string) {
	nodes := NumNodes()
	for i := 0; i < nodes; i++ {
		dbs := GetDatabases()
		for remote_signer.StringIndexOf(database, dbs) == -1 {
			time.Sleep(1 * time.Second)
			dbs = GetDatabases()
		}
	}
}
