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
	models.UserModelTableInit,
	models.UserTokenTableInit,
}

func init() {
	DbSetup()
	InitTables()
}

func DbSetup() {
	RthState = RethinkDbState{}
	if remote_signer.EnableRethinkSKS {
		initLock.Lock()
		defer initLock.Unlock()
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
		dbLog.Info("Connected!")
		RthState.connection = conn
	}
}

var initLock = sync.Mutex{}

func InitTables() {
	if remote_signer.EnableRethinkSKS {
		SLog.UnsetTestMode()
		initLock.Lock()
		defer initLock.Unlock()
		dbLog.Info("Running InitTables")
		var dbs = GetDatabases()
		var conn = GetConnection()
		needWait := false

		if remote_signer.StringIndexOf(remote_signer.DatabaseName, dbs) == -1 {
			dbLog.Warn("Database %s does not exists. Creating it...", remote_signer.DatabaseName)
			err := r.DBCreate(remote_signer.DatabaseName).Exec(conn)
			if err != nil {
				dbLog.Fatal(err)
			}
			time.Sleep(5 * time.Second)
		} else {
			dbLog.Debug("Database %s already exists. Skipping...", remote_signer.DatabaseName)
		}

		WaitDatabaseCreate(remote_signer.DatabaseName)

		dbLog.Warn("Waiting for database %s to be ready", remote_signer.DatabaseName)
		_ = r.DB(remote_signer.DatabaseName).Wait(r.WaitOpts{
			Timeout: 0,
		}).Exec(conn)

		dbLog.Info("Database %s is ready", remote_signer.DatabaseName)

		conn.Use(remote_signer.DatabaseName)

		tables := GetTables()
		numNodes := NumNodes()

		for _, v := range tablesToInitialize {
			if remote_signer.StringIndexOf(v.TableName, tables) == -1 {
				dbLog.Info("Table %s does not exists. Creating...", v.TableName)
				err := r.TableCreate(v.TableName, r.TableCreateOpts{
					Durability: "hard",
					Replicas:   numNodes,
				}).Exec(conn)
				if err != nil {
					dbLog.Fatal(err)
				}
				_, _ = r.Table(v.TableName).Wait(r.WaitOpts{
					Timeout: 0,
				}).Run(conn)
				time.Sleep(time.Millisecond * 500)
				needWait = true
			}

			dbLog.Info("Checking Indexes for table %s", v.TableName)
			idxs := GetTableIndexes(v.TableName)

			for _, vidx := range v.TableIndexes {
				dbLog.Debug("Checking index %s in %s", v.TableName, vidx)
				if remote_signer.StringIndexOf(vidx, idxs) == -1 {
					dbLog.Warn("Index %s not found at table %s. Creating it...", vidx, v.TableName)
					err := r.Table(v.TableName).IndexCreate(vidx).Exec(conn)
					if err != nil {
						dbLog.Fatal(err)
					}
					_ = r.Table(v.TableName).IndexWait().Exec(conn)
					needWait = true
				} else {
					dbLog.Debug("Index %s already exists in table %s. Skipping it...", vidx, v.TableName)
				}
			}
		}
		if needWait {
			time.Sleep(5 * time.Second)
		}
	}
}

func GetConnection() *r.Session {
	if RthState.connection == nil {
		DbSetup()
	}
	_ = RthState.connection.Reconnect()
	return RthState.connection
}
