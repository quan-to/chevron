package database

import (
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strings"
	"time"
)

const maxRetryCount = 5

type RethinkDbState struct {
	connection *r.Session
}

var RthState RethinkDbState
var dbLog = slog.Scope("DatabaseManager").Tag(remote_signer.DefaultTag)

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
	if remote_signer.EnableRethinkSKS && RthState.connection == nil {
		dbLog.Await("RethinkDB SKS Enabled. Starting %d connections to %s:%d", remote_signer.RethinkDBPoolSize, remote_signer.RethinkDBHost, remote_signer.RethinkDBPort)
		conn, err := r.Connect(r.ConnectOpts{
			Address:    fmt.Sprintf("%s:%d", remote_signer.RethinkDBHost, remote_signer.RethinkDBPort),
			Username:   remote_signer.RethinkDBUsername,
			Password:   remote_signer.RethinkDBPassword,
			NumRetries: maxRetryCount,
			MaxOpen:    remote_signer.RethinkDBPoolSize,
			InitialCap: remote_signer.RethinkDBPoolSize,
			Database:   remote_signer.DatabaseName,
		})

		if err != nil {
			dbLog.Fatal(err)
		}
		dbLog.Done("Connected!")
		RthState.connection = conn
	}
}

func InitTables() {
	if remote_signer.EnableRethinkSKS {
		slog.UnsetTestMode()
		dbLog.Await("Running InitTables")
		dbs := GetDatabases()
		conn := GetConnection()

		if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) == -1 {
			dbLog.Note("Database %s does not exists. Creating it...", remote_signer.DatabaseName)
			err := r.DBCreate(remote_signer.DatabaseName).Exec(conn)
			if err != nil && strings.Index(err.Error(), " already exists") == -1 {
				dbLog.Fatal(err)
			}
		} else {
			dbLog.WarnDone("Database %s already exists. Skipping...", remote_signer.DatabaseName)
		}

		WaitDatabaseCreate(remote_signer.DatabaseName)

		dbLog.Await("Waiting for database %s to be ready", remote_signer.DatabaseName)
		_ = r.DB(remote_signer.DatabaseName).Wait(r.WaitOpts{Timeout: 0}).Exec(conn)

		dbLog.Success("Database %s is ready", remote_signer.DatabaseName)

		conn.Use(remote_signer.DatabaseName)

		tables := GetTables()
		numNodes := NumNodes()

		for _, v := range tablesToInitialize {
			if remote_signer.StringIndexOf(v.TableName, tables) == -1 {
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
				if remote_signer.StringIndexOf(vidx, idxs) == -1 {
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
		err := RthState.connection.Close()

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

	dbLog.Error("Dropping test database %s", remote_signer.DatabaseName)
	if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) > -1 {
		dbLog.Error("Test Database already exists, dropping.")
		_ = r.DBDrop(remote_signer.DatabaseName).Exec(c)
	}

	WaitDatabaseDrop(remote_signer.DatabaseName)
	time.Sleep(1 * time.Second)

	dbLog.Info("Database reseted")
	//slog.SetTestMode()
}
