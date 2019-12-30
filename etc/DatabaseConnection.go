package etc

import (
	"github.com/quan-to/chevron/database"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
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

func ResetDatabase() {
	database.ResetDatabase()
}

func Cleanup() {
	database.Cleanup()
}
