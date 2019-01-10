using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Threading;
using RemoteSigner.Database.Attributes;
using RemoteSigner.Database.Models;
using RemoteSigner.Log;
using RethinkDb.Driver;
using RethinkDb.Driver.Net;

namespace RemoteSigner.Database {
    public class DatabaseManager {

        const int MAX_RETRY_COUNT = 5;
        static readonly RethinkDB R = RethinkDB.R;
        readonly List<Connection> connectionPool;

        public static DatabaseManager GlobalDm { get; private set; }

        int currentConn;
        int globalRetryCount;

        static DatabaseManager() {
            GlobalDm = new DatabaseManager();
        }

        DatabaseManager() {
            connectionPool = new List<Connection>();
            currentConn = 0;
            globalRetryCount = 0;
            InitPool();
        }

        void InitPool() {
            // DO NOT LOCK HERE
            if (Configuration.EnableRethinkSKS) {
                Logger.Log("DatabaseManager", $"Initializing RethinkDB Connection Pool for SKS at {Configuration.RethinkDBHost}:{Configuration.RethinkDBPort} with pool size {Configuration.RethinkDBPoolSize}");
                int tryCount = 0;
                while (connectionPool.Count < Configuration.RethinkDBPoolSize) {
                    try {
                        // That library has a Connection pool but does not support hostname which is used by Rancher and other orchestrations in Load Balancer
                        var c = R.Connection()
                                 .Hostname(Configuration.RethinkDBHost)
                                 .Port(Configuration.RethinkDBPort)
                                 .Timeout(60)
                                 .User(Configuration.RethinkDBUsername, Configuration.RethinkDBPassword)
                                 .Connect();
                        connectionPool.Add(c);
                        tryCount = 0;
                    } catch (Exception e) {
                        tryCount++;
                        Logger.Error("DatabaseManager", $"Error connecting to database {e}");
                    }
                    if (tryCount > MAX_RETRY_COUNT) {
                        Logger.Error("DatabaseManager", $"Max Retries of {MAX_RETRY_COUNT} trying to connect to database.");
                        throw new ApplicationException($"Max Retries of {MAX_RETRY_COUNT} trying to connect to database.");
                    }
                }
                InitData();
            }
        }

        void InitData() {
            try {
                var c = GetConnection();
                if (!R.DbList().Contains(Configuration.DatabaseName).RunAtom<bool>(c)) {
                    Logger.Warn("DatabaseManager", $"Database {Configuration.DatabaseName} not found. Creating it...");
                    var x = R.DbCreate(Configuration.DatabaseName).Run(c);
                }

                UpdateConnectionsDatabase(Configuration.DatabaseName);

                Logger.Log("DatabaseManager", "Searching for database table definitions");
                Assembly a = Assembly.GetExecutingAssembly();
                string[] namespaces = a.GetTypes().Select(x => x.Namespace).Distinct().ToArray();
                foreach (string n in namespaces) {
                    if (n != null && n.StartsWith("RemoteSigner.Database", StringComparison.InvariantCultureIgnoreCase)) {
                        Logger.Log("DatabaseManager", $"Loading DB Data for namespace {n}");
                        InitTables(c, a, n);
                    }
                }

                Migrations();
            } catch (Exception e) {
                Logger.Error("DatabaseManager", $"Error initializing Database: {e}");
            }
        }

        void InitTables(Connection c, Assembly runningAssembly, string modulesAssembly) {
            Type[] typelist = Tools.GetTypesInNamespace(runningAssembly, modulesAssembly);
            for (int i = 0; i < typelist.Length; i++) {
                Type tClass = typelist[i];

                // Search for DBTable Attribute
                Attribute t = tClass.GetCustomAttribute(typeof(DBTable));
                if (t != null) {
                    Logger.Log("DatabaseManager", $"Found Table Definition {tClass.Name}");
                    DBTable dbTable = (DBTable)t;
                    var tableName = dbTable.TableName;
                    if (!R.TableList().Contains(tableName).RunAtom<bool>(c)) {
                        Logger.Warn("DatabaseManager", $"Table {tableName} does not exist. Creating...");
                        R.TableCreate(tableName).Run(c);
                    }
                    // Get Indexes
                    var properties = tClass.GetProperties().ToList();
                    List<string> existingIndexes = R.Table(tableName)
                     .IndexList()
                     .CoerceTo("array")
                     .Run<List<string>>(c);
                    properties.ForEach((prop) => {
                        var idx = prop.Name != "Id" ? prop.GetCustomAttribute(typeof(DBIndex)) : null;
                        if (idx != null) {
                            Logger.Debug("DatabaseManager", $"Checking Index {prop.Name} on table {tableName}");
                            if (!existingIndexes.Contains(prop.Name)) {
                                Logger.Warn("DatabaseManager", $"Index {prop.Name} does not exists in table {tableName}. Creating it...");
                                var propType = prop.PropertyType;
                                if (propType.IsGenericType && propType.GetGenericTypeDefinition() == typeof(List<>)) {
                                    // Multi-index 
                                    R.Table(tableName).IndexCreate(prop.Name).OptArg("multi", true).Run(c);
                                } else {
                                    // Single Index
                                    R.Table(tableName).IndexCreate(prop.Name).Run(c);
                                }
                            }
                        }
                    });
                }
            }
        }

        void UpdateConnectionsDatabase(string dbName) {
            lock (connectionPool) {
                Logger.Log("DatabaseManager", $"Changing Database Name for active connections to {dbName}");
                connectionPool.ForEach(c => c.Use(dbName));
            }
        }

        void Migrations() {
            var c = GetConnection();
            Logger.Log("DatabaseManager", "Migrating Keys without Subkey");
            var v = GPGKey.FetchKeyWithoutSubkey(c);
            if (v.Count > 0) {
                Logger.Log("DatabaseManager", $"There are {v.Count} keys without subkeys indexed.");
                foreach (var k in v) {
                    k.Subkeys = PGPManager.GetSubKeysFromPublicKey(k.AsciiArmoredPublicKey);
                    Logger.Log("DatabaseManager", $"Added {k.Subkeys.Count} subkeys to {k.FullFingerPrint}");
                    k.Save(c);
                }
            }

            Logger.Log("DatabaseManager", "Migrations Finished");
        }

        public Connection GetConnection() {
            try {
                Connection c;
                lock (connectionPool) {
                    if (connectionPool.Count == 0) {
                        Logger.Error("Empty connection pool! Running InitPool");
                        InitPool();
                    }
                    currentConn = (currentConn + 1) % connectionPool.Count;
                    c = connectionPool[currentConn];
                }
                try {
                    c.CheckOpen();
                } catch (Exception e) {
                    Logger.Warn("DatabaseManager", $"One rethinkdb connection is dead. Retrying. {e}");
                    c.Close();
                    c.Reconnect();
                }

                try {
                    c.CheckOpen();
                } catch (Exception e) {
                    Logger.Warn("DatabaseManager", $"Could not reconnect {e}");
                    globalRetryCount++;
                    if (globalRetryCount > 10) {
                        throw new ApplicationException($"Max Retries of 10 trying to connect to database.");
                    }
                    lock (connectionPool) {
                        connectionPool.Remove(c);
                    }
                    return GetConnection();
                }

                return c;
            } catch (Exception e) {
                Logger.Error("DatabaseManager", $"Fatal Error at GetConnection(): {e.Message}" +
                             $"\tCurrent Conn: {currentConn}" +
                             $"\tPool Size: {connectionPool.Count}" +
                             $"\tStack Trace: {e.StackTrace}");
                Logger.Error("DatabaseManager", "Waiting 10 seconds before continuing...");
                Thread.Sleep(10);
                throw new ApplicationException(e.Message);
            }
        }
    }
}
