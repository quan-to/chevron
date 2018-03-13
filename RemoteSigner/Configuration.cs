﻿using System;
using System.IO;

namespace RemoteSigner {
    public static class Configuration {

        public static string SyslogServer { get; private set; }
        public static string SyslogFacility { get; private set; }
        public static string PrivateKeyFolder { get; private set; }
        public static string KeyPrefix { get; private set; }
        public static string SKSServer { get; private set; }
        public static int HttpPort { get; private set; }
        public static int MaxKeyRingCache { get; private set; }
        public static bool EnableRethinkSKS { get; private set; }
        public static string RethinkDBHost { get; private set; }
        public static int RethinkDBPort { get; private set; }
        public static int RethinkDBPoolSize { get; private set; }
        public static string DatabaseName { get; private set; }
        public static string MasterGPGKeyPath { get; private set; }
        public static string MasterGPGKeyPasswordPath { get; private set; }

        static Configuration() {
            SyslogServer = Environment.GetEnvironmentVariable("SYSLOG_IP") ?? "127.0.0.1";
            SyslogFacility = Environment.GetEnvironmentVariable("SYSLOG_FACILITY") ?? "LOG_USER";
            PrivateKeyFolder = Environment.GetEnvironmentVariable("PRIVATE_KEY_FOLDER") ?? "./keys";
            SKSServer = Environment.GetEnvironmentVariable("SKS_SERVER") ?? "http://pgp.mit.edu/";
            KeyPrefix = Environment.GetEnvironmentVariable("KEY_PREFIX") ?? "";

            var mkrc = Environment.GetEnvironmentVariable("MAX_KEYRING_CACHE_SIZE") ?? "1000";
            MaxKeyRingCache = int.Parse(mkrc);

            var hp = Environment.GetEnvironmentVariable("HTTP_PORT") ?? "5100";
            HttpPort = int.Parse(hp);

            EnableRethinkSKS = true; // Environment.GetEnvironmentVariable("ENABLE_RETHINKDB_SKS") == "true"; // TODO: Change-me
            RethinkDBHost = Environment.GetEnvironmentVariable("RETHINKDB_HOST") ?? "localhost";
            var rdbport = Environment.GetEnvironmentVariable("RETHINKDB_PORT") ?? "28015";
            RethinkDBPort = int.Parse(rdbport);
            var rdbpool = Environment.GetEnvironmentVariable("RETHINKDB_POOL_SIZE") ?? "10";
            RethinkDBPoolSize = int.Parse(rdbpool);

            DatabaseName = Environment.GetEnvironmentVariable("DATABASE_NAME") ?? "remote_signer";

            MasterGPGKeyPath = "keys/870AFA59.key"; // Environment.GetEnvironmentVariable("MASTER_GPG_KEY_PATH") ?? null; // TODO: Change-me
            MasterGPGKeyPasswordPath = "pass.txt"; // Environment.GetEnvironmentVariable("MASTER_GPG_KEY_PASSWORD_PATH") ?? null; // TODO: Change-me

#pragma warning disable RECS0022 // A catch clause that catches System.Exception and has an empty body
            try { Directory.CreateDirectory(PrivateKeyFolder); } catch (Exception) { }
#pragma warning restore RECS0022 // A catch clause that catches System.Exception and has an empty body
        }
    }
}
