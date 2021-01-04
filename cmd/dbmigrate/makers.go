package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quan-to/chevron/internal/agent"
	"github.com/quan-to/chevron/pkg/database/cache"
	"github.com/quan-to/chevron/pkg/database/pg"
	"github.com/quan-to/chevron/pkg/database/rql"
	"github.com/quan-to/slog"
)

func stringParam(config map[string]interface{}, paramName string) string {
	p, ok := config[paramName]
	if !ok || p == nil {
		return ""
	}

	ps, ok := p.(string)
	if !ok {
		return ""
	}

	return ps
}

func boolParam(config map[string]interface{}, paramName string) *bool {
	p, ok := config[paramName]
	if !ok || p == nil {
		return nil
	}

	ps, ok := p.(bool)
	if !ok {
		return nil
	}

	return &ps
}

func intParam(config map[string]interface{}, paramName string) int {
	p, ok := config[paramName]
	if !ok || p == nil {
		return -1
	}

	// JSON Numbers are float64
	ps, ok := p.(float64)
	if !ok {
		return -1
	}

	return int(ps)
}

func makeHandler(filename string, logger slog.Instance) (migrateHandler, error) {
	configData := map[string]interface{}{}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&configData)
	if err != nil {
		return nil, err
	}

	dialect := stringParam(configData, "DATABASE_DIALECT")

	if dialect == "" {
		return nil, fmt.Errorf("DATABASE_DIALECT was not specified")
	}
	logger.Info("Loading database %s", dialect)
	var dbh agent.DatabaseHandler
	switch dialect {
	case "rethinkdb":
		dbh, err = makeRethinkDBHandler(configData, logger)
		if err != nil {
			return nil, err
		}
	case "postgres":
		dbh, err = makePostgresDBHandler(configData, logger)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("database %s not supported", dialect)
	}

	redisEnable := boolParam(configData, "REDIS_ENABLE")

	if redisEnable != nil && *redisEnable {
		logger.Info("REDIS Enabled. Creating proxy layer")
		dbh, err = makeRedisLayer(configData, dbh, logger)
	}

	return dbh, err
}

func makeRedisLayer(config map[string]interface{}, dbh agent.DatabaseHandler, logger slog.Instance) (*cache.Driver, error) {
	logger.Info("Redis enabled. Wrapping cache layer")
	redisDriver := cache.MakeRedisDriver(dbh, logger)
	var tlsConfig *tls.Config
	redisEnable := boolParam(config, "REDIS_TLS_ENABLED")
	if redisEnable != nil && *redisEnable {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	redisHost := stringParam(config, "REDIS_HOST")
	redisUser := stringParam(config, "REDIS_USER")
	redisPass := stringParam(config, "REDIS_PASS")
	redisDatabaseIndex := intParam(config, "REDIS_DATABASE_INDEX")

	err := redisDriver.Setup(&redis.RingOptions{
		Addrs: map[string]string{
			"server0": redisHost,
		},
		Username:  redisUser,
		Password:  redisPass,
		DB:        redisDatabaseIndex,
		TLSConfig: tlsConfig,
	}, 1, time.Second)
	if err != nil {
		return nil, err
	}
	return redisDriver, nil
}

func makeRethinkDBHandler(config map[string]interface{}, logger slog.Instance) (*rql.RethinkDBDriver, error) {
	rqlHost := stringParam(config, "RETHINKDB_HOST")
	rqlUser := stringParam(config, "RETHINKDB_USERNAME")
	rqlPass := stringParam(config, "RETHINKDB_PASSWORD")
	rqlPort := intParam(config, "RETHINKDB_PORT")
	rqlDatabase := stringParam(config, "DATABASE_NAME")

	logger.Info("RethinkDB Database Enabled. Creating handler")
	rdb := rql.MakeRethinkDBDriver(logger)
	logger.Info("Connecting to RethinkDB")
	err := rdb.Connect(rqlHost, rqlUser, rqlPass, rqlDatabase, rqlPort, 5)
	if err != nil {
		return nil, err
	}
	logger.Info("Initializing database")
	err = rdb.InitDatabase()
	if err != nil {
		return nil, err
	}
	logger.Info("RethinkDB Handler done!")
	return rdb, nil
}

func makePostgresDBHandler(config map[string]interface{}, logger slog.Instance) (*pg.PostgreSQLDBDriver, error) {
	connectionString := stringParam(config, "CONNECTION_STRING")
	logger.Info("PostgreSQL Database Enabled. Creating handler")
	rdb := pg.MakePostgreSQLDBDriver(logger)
	logger.Info("Initializing database")
	err := rdb.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
