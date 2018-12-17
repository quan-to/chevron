package remote_signer

import (
	"fmt"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

const maxRetryCount = 5

type rethinkDbState struct {
	connection  *r.Session
	currentConn int
}

var rthState rethinkDbState
var dbLog = SLog.Scope("DatabaseManager")

var tablesToInitialize = []models.TableInitStruct{
	models.GPGKeyTableInit,
}

func init() {
	rthState = rethinkDbState{
		currentConn: 0,
	}

	if EnableRethinkSKS {
		dbLog.Info("RethinkDB SKS Enabled. Starting %d connections to %s:%d", RethinkDBPoolSize, RethinkDBHost, RethinkDBPort)
		conn, err := r.Connect(r.ConnectOpts{
			Address:    fmt.Sprintf("%s:%d", RethinkDBHost, RethinkDBPort),
			Username:   RethinkDBUsername,
			Password:   RethinkDBPassword,
			NumRetries: maxRetryCount,
			MaxOpen:    RethinkDBPoolSize,
			InitialCap: RethinkDBPoolSize,
			Database:   DatabaseName,
		})

		if err != nil {
			dbLog.Fatal(err)
		}

		rthState.connection = conn

		initTables()
	}
}

func initTables() {
	var dbs = getDatabases()
	var conn = GetConnection()

	if stringIndexOf(DatabaseName, dbs) == -1 {
		dbLog.Warn("Database %s does not exists. Creating it...", DatabaseName)
		err := r.DBCreate(DatabaseName).Exec(conn)
		if err != nil {
			dbLog.Fatal(err)
		}
	} else {
		dbLog.Debug("Database %s already exists. Skipping...", DatabaseName)
	}

	tables := getTables()

	for _, v := range tablesToInitialize {
		if stringIndexOf(v.TableName, tables) == -1 {
			dbLog.Info("Table %s does not exists. Creating...", v.TableName)
			err := r.TableCreate(v.TableName).Exec(conn)
			if err != nil {
				dbLog.Fatal(err)
			}
		}

		dbLog.Info("Checking Indexes for table %s", v.TableName)
		idxs := getTableIndexes(v.TableName)

		for _, vidx := range v.TableIndexes {
			if stringIndexOf(vidx, idxs) == -1 {
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
}

func GetConnection() *r.Session {
	return rthState.connection
}
