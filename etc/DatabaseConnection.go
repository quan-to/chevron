package etc

import (
	"github.com/quan-to/remote-signer/database"
	"gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func GetConnection() *rethinkdb.Session {
	return database.GetConnection()
}

func GetDatabases() []string {
	return database.GetDatabases()
}

func GetTableIndexes(tableName string) []string {
	return database.GetTableIndexes(tableName)
}

func GetTables() []string {
	return database.GetTables()
}

func DbSetup() {
	database.DbSetup()
}

func InitTables() {
	database.InitTables()
}
