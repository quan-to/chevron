package database

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/tools"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
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

	defer c.Close()

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
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		dbs := GetDatabases()
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(database, dbs) == -1 {
			break
		}
	}
}

func WaitDatabaseCreate(database string) {
	dbLog.Await("Waiting for database create")
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		dbs := GetDatabases()
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(database, dbs) > -1 {
			break
		}
	}

	dbLog.Done("Done waiting database create")
}

func WaitTableCreate(table string) {
	dbLog.Await("Waiting for table %s create", table)
	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		tables := GetTables()
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(table, tables) > -1 {
			break
		}
	}

	_ = r.DB(config.DatabaseName).
		Table(table).
		Wait(r.WaitOpts{Timeout: 0}).
		Exec(GetConnection())
	dbLog.Done("Done waiting table %s create", table)
}

func WaitTableIndexCreate(table, index string) {
	dbLog.Await("Waiting for index %s/%s create", table, index)

	timeout := time.Now().Add(time.Second * 5)
	for time.Now().UnixNano() < timeout.UnixNano() {
		indexes := GetTableIndexes(table)
		time.Sleep(time.Millisecond * 100)
		if tools.StringIndexOf(index, indexes) > -1 {
			break
		}
	}
	dbLog.Done("Done waiting index %s/%s create", table, index)
}
