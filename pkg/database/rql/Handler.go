package rql

import (
	"encoding/json"
	"strings"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"
	"github.com/quan-to/slog"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// tableInitStruct is used for representing a table initialization
type tableInitStruct struct {
	TableName    string
	TableIndexes []string
}

// RethinkDBDriver is a database driver for RethinkDB
type RethinkDBDriver struct {
	conn     r.QueryExecutor
	log      slog.Instance
	database string

	// Migration tools
	gpgKeysMigrationCursor *r.Cursor
	userMigrationCursor    *r.Cursor
}

// MakeRethinkDBDriver creates a new database driver for rethinkdb
func MakeRethinkDBDriver(log slog.Instance) *RethinkDBDriver {
	if log == nil {
		log = slog.Scope("RethinkDB")
	} else {
		log = log.SubScope("RethinkDB")
	}

	return &RethinkDBDriver{
		log: log,
	}
}

// HealthCheck returns nil if everything is OK with the handler
func (h *RethinkDBDriver) HealthCheck() error {
	d, err := r.Expr(1).Run(h.conn)
	_ = d.Close()
	return err
}

// InitDatabase initializes indexes and tables required to operation
func (h *RethinkDBDriver) InitDatabase() error {
	runners := []func() error{
		// Initializations
		h.initUserTable,
		h.initUserTokenTable,
		h.initGPGKeyTable,

		// Migrations
		h.migrateUserTable,
		h.migrateUserTokenTable,
		h.migrateGPGKeyTable,
	}

	for _, v := range runners {
		err := v()
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *RethinkDBDriver) initFromStruct(v tableInitStruct) error {
	tables, err := h.getTables()
	if err != nil {
		return err
	}

	if tools.StringIndexOf(v.TableName, tables) == -1 {
		h.log.Await("Table %s does not exists. Creating...", v.TableName)
		err := r.TableCreate(v.TableName, r.TableCreateOpts{Durability: "hard"}).Exec(h.conn)
		if err != nil && !strings.Contains(err.Error(), " already exists") {
			h.log.Error("Error creating table: %s", err)
			return err
		}
		err = h.waitTableCreate(v.TableName)
		if err != nil {
			return err
		}
		h.log.Success("Table %s created.", v.TableName)
	}

	h.log.Await("        Checking Indexes for table %s", v.TableName)
	idxs, err := h.getTableIndexes(v.TableName)
	if err != nil {
		return err
	}

	for _, vidx := range v.TableIndexes {
		h.log.Await("           Checking index %s in %s", v.TableName, vidx)
		if tools.StringIndexOf(vidx, idxs) == -1 {
			h.log.Note("           Index %s not found at table %s. Creating it...", vidx, v.TableName)
			err := r.Table(v.TableName).IndexCreate(vidx).Exec(h.conn)
			if err != nil && !strings.Contains(err.Error(), " already exists") {
				h.log.Error("Error creating index %s on table %s: %s", vidx, v.TableName, err)
				return err
			}
			err = h.waitTableIndexCreate(v.TableName, vidx)
			if err != nil {
				return err
			}
		} else {
			h.log.WarnDone("           Index %s already exists in table %s. Skipping it...", vidx, v.TableName)
		}
	}

	h.log.Success("        Finished getting indexes for table %s", v.TableName)
	return nil
}

// convertToRethinkDB converts to a string map changing the ID field to id
func convertToRethinkDB(data interface{}) (map[string]interface{}, error) {
	bdata, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	outData := map[string]interface{}{}
	err = json.Unmarshal(bdata, &outData)
	if err != nil {
		return nil, err
	}

	if id, ok := outData["ID"]; ok {
		outData["id"] = id
		delete(outData, "ID")
	}

	return outData, nil
}

// convertFromRethinkDB converts from a RethinkDB string map changing the id field to ID
func convertFromRethinkDB(input map[string]interface{}, output interface{}) error {
	if id, ok := input["id"]; ok {
		input["ID"] = id
		delete(input, "id")
	}

	bdata, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(bdata, output)
}

func (h *RethinkDBDriver) resultsAsArray(res *r.Cursor) ([]models.GPGKey, error) {
	results := make([]models.GPGKey, 0)
	var gpgKey map[string]interface{}

	for res.Next(&gpgKey) {
		gpgKey = h.fixGPGKey(gpgKey)
		var key models.GPGKey
		err := convertFromRethinkDB(gpgKey, &key)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	return results, nil
}
