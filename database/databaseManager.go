package database

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"sync"
	"time"
)

const maxRetryCount = 5

type RethinkDbState struct {
	connection *r.Session
}

var RthState RethinkDbState
var dbLog = SLog.Scope("DatabaseManager")

var tablesToInitialize = []models.TableInitStruct{
	models.GPGKeyTableInit,
}

func init() {
	DbSetup()
	InitTables()
}

func DbSetup() {
	RthState = RethinkDbState{}

	if remote_signer.EnableRethinkSKS {
		dbLog.Info("RethinkDB SKS Enabled. Starting %d connections to %s:%d", remote_signer.RethinkDBPoolSize, remote_signer.RethinkDBHost, remote_signer.RethinkDBPort)
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

		RthState.connection = conn
	}
}

var initLock = sync.Mutex{}

func InitTables() {
	if remote_signer.EnableRethinkSKS {
		initLock.Lock()
		defer initLock.Unlock()
		var dbs = GetDatabases()
		var conn = GetConnection()

		if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) == -1 {
			dbLog.Warn("Database %s does not exists. Creating it...", remote_signer.DatabaseName)
			err := r.DBCreate(remote_signer.DatabaseName).Exec(conn)
			if err != nil {
				dbLog.Fatal(err)
			}
		} else {
			dbLog.Debug("Database %s already exists. Skipping...", remote_signer.DatabaseName)
		}

		tables := GetTables()

		for _, v := range tablesToInitialize {
			if remote_signer.StringIndexOf(v.TableName, tables) == -1 {
				dbLog.Info("Table %s does not exists. Creating...", v.TableName)
				err := r.TableCreate(v.TableName).Exec(conn)
				if err != nil {
					dbLog.Fatal(err)
				}
			}

			dbLog.Info("Checking Indexes for table %s", v.TableName)
			idxs := GetTableIndexes(v.TableName)

			for _, vidx := range v.TableIndexes {
				if remote_signer.StringIndexOf(vidx, idxs) == -1 {
					dbLog.Warn("Index %s not found at table %s. Creating it...", vidx, v.TableName)
					err := r.Table(v.TableName).IndexCreate(vidx).Exec(conn)
					if err != nil {
						dbLog.Fatal(err)
					}
				} else {
					dbLog.Debug("Index %s already exists in table %s. Skipping it...", vidx, v.TableName)
				}
			}
		}
		time.Sleep(2 * time.Second) // Wait for indexes
	}
}

func GetConnection() *r.Session {
	if RthState.connection == nil {
		DbSetup()
		InitTables()
	}
	return RthState.connection
}
