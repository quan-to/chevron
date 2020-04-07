package database

import (
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"strings"
	"time"
)

const maxRetryCount = 5

type RethinkDbState struct {
	connection *r.Session
}

var RthState RethinkDbState
var dbLog = slog.Scope("DatabaseManager").Tag(tools.DefaultTag)

var tablesToInitialize = []models.TableInitStruct{
	models.GPGKeyTableInit,
	models.UserModelTableInit,
	models.UserTokenTableInit,
}

func init() {
	DbSetup()
	InitTables()
}

func DbSetup() {
	if config.EnableRethinkSKS && RthState.connection == nil {
		dbLog.Await("RethinkDB SKS Enabled. Starting %d connections to %s:%d", config.RethinkDBPoolSize, config.RethinkDBHost, config.RethinkDBPort)
		conn, err := r.Connect(r.ConnectOpts{
			Address:    fmt.Sprintf("%s:%d", config.RethinkDBHost, config.RethinkDBPort),
			Username:   config.RethinkDBUsername,
			Password:   config.RethinkDBPassword,
			NumRetries: maxRetryCount,
			MaxOpen:    config.RethinkDBPoolSize,
			InitialCap: config.RethinkDBPoolSize,
			Database:   config.DatabaseName,
		})

		if err != nil {
			dbLog.Fatal(err)
		}
		dbLog.Done("Connected!")
		RthState.connection = conn
	}
}

func InitTables() {
	if config.EnableRethinkSKS {
		slog.UnsetTestMode()
		dbLog.Await("Starting running InitTables")
		dbs := GetDatabases()
		conn := GetConnection()

		if tools.StringIndexOf(config.DatabaseName, dbs) == -1 {
			dbLog.Note("Database %s does not exists. Creating it...", config.DatabaseName)
			err := r.DBCreate(config.DatabaseName).Exec(conn)
			if err != nil && strings.Index(err.Error(), " already exists") == -1 {
				dbLog.Fatal(err)
			}
		} else {
			dbLog.WarnDone("Database %s already exists. Skipping...", config.DatabaseName)
		}

		WaitDatabaseCreate(config.DatabaseName)

		dbLog.Await("Waiting for database %s to be ready", config.DatabaseName)
		_ = r.DB(config.DatabaseName).Wait(r.WaitOpts{Timeout: 0}).Exec(conn)

		dbLog.Success("Database %s is ready", config.DatabaseName)

		conn.Use(config.DatabaseName)

		tables := GetTables()
		numNodes := NumNodes()

		for _, v := range tablesToInitialize {
			if tools.StringIndexOf(v.TableName, tables) == -1 {
				dbLog.Await("Table %s does not exists. Creating...", v.TableName)
				err := r.TableCreate(v.TableName, r.TableCreateOpts{
					Durability: "hard",
					Replicas:   numNodes,
				}).Exec(conn)
				if err != nil && strings.Index(err.Error(), " already exists") == -1 {
					dbLog.Fatal(err)
				}
				WaitTableCreate(v.TableName)
				dbLog.Success("Table %s created.", v.TableName)
			}

			dbLog.Await("        Checking Indexes for table %s", v.TableName)
			idxs := GetTableIndexes(v.TableName)

			for _, vidx := range v.TableIndexes {
				dbLog.Await("           Checking index %s in %s", v.TableName, vidx)
				if tools.StringIndexOf(vidx, idxs) == -1 {
					dbLog.Note("           Index %s not found at table %s. Creating it...", vidx, v.TableName)
					err := r.Table(v.TableName).IndexCreate(vidx).Exec(conn)
					if err != nil && strings.Index(err.Error(), " already exists") == -1 {
						dbLog.Fatal(err)
					}
					WaitTableIndexCreate(v.TableName, vidx)
				} else {
					dbLog.WarnDone("           Index %s already exists in table %s. Skipping it...", vidx, v.TableName)
				}
			}

			dbLog.Success("        Finished getting indexes for table %s", v.TableName)
		}
	}
}

func Cleanup() {
	if RthState.connection != nil {
		err := RthState.connection.Close(r.CloseOpts{
			NoReplyWait: false,
		})

		if err != nil {
			slog.Fatal(err)
		}
		RthState.connection = nil
	}
}

func GetConnection() *r.Session {
	if RthState.connection == nil {
		dbLog.Info("GetConnection() - Conection is nil, running DbSetup()")
		DbSetup()
	}
	return RthState.connection
}

func ResetDatabase() {
	//slog.UnsetTestMode()

	dbLog.Error("Reseting Database")
	c := GetConnection()
	dbs := GetDatabases()

	dbLog.Error("Dropping test database %s", config.DatabaseName)
	if tools.StringIndexOf(config.DatabaseName, dbs) > -1 {
		dbLog.Error("Test Database already exists, dropping.")
		_ = r.DBDrop(config.DatabaseName).Exec(c)
	}

	WaitDatabaseDrop(config.DatabaseName)
	time.Sleep(1 * time.Second)

	dbLog.Info("Database reseted")
	//slog.SetTestMode()
}
