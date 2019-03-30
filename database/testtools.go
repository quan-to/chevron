package database

import (
	"github.com/quan-to/remote-signer"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"time"
)

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

	return count
}

func WaitDatabaseDrop(database string) {
	dbLog.Info("Waiting for database drop")
	nodes := NumNodes() * 4
	for i := 0; i < nodes; i++ {
		dbs := GetDatabases()
		for remote_signer.StringIndexOf(database, dbs) > -1 {
			time.Sleep(1 * time.Second)
			dbs = GetDatabases()
		}
		time.Sleep(1 * time.Second)
	}
}

func WaitDatabaseCreate(database string) {
	dbLog.Info("Waiting for database create")
	nodes := NumNodes() * 4
	for i := 0; i < nodes; i++ {
		dbs := GetDatabases()
		for remote_signer.StringIndexOf(database, dbs) == -1 {
			time.Sleep(1 * time.Second)
			dbs = GetDatabases()
		}
		time.Sleep(1 * time.Second)
	}
	_ = r.DB(database).
		Wait(r.WaitOpts{Timeout: 0}).
		Exec(GetConnection())
}

func WaitTableCreate(table string) {
	dbLog.Info("Waiting for table create")
	nodes := NumNodes() * 4
	for i := 0; i < nodes; i++ {
		tables := GetTables()
		for remote_signer.StringIndexOf(table, tables) == -1 {
			time.Sleep(1 * time.Second)
			tables = GetTables()
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = r.DB(remote_signer.DatabaseName).
		Table(table).
		Wait(r.WaitOpts{Timeout: 0}).
		Exec(GetConnection())
}

func WaitTableIndexCreate(table, index string) {
	dbLog.Info("Waiting for index create")
	nodes := NumNodes() * 4
	for i := 0; i < nodes; i++ {
		indexes := GetTableIndexes(table)
		for remote_signer.StringIndexOf(index, indexes) == -1 {
			time.Sleep(100 * time.Millisecond)
			indexes = GetTableIndexes(table)
		}
	}
}
