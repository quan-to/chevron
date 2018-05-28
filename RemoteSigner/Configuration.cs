using System;
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
        public static bool MasterGPGKeyBase64Encoded { get; private set; }
        public static bool KeysBase64Encoded { get; private set; }

        static Configuration() {
            SyslogServer = Environment.GetEnvironmentVariable("SYSLOG_IP") ?? "127.0.0.1";
            SyslogFacility = Environment.GetEnvironmentVariable("SYSLOG_FACILITY") ?? "LOG_USER";
            PrivateKeyFolder = Environment.GetEnvironmentVariable("PRIVATE_KEY_FOLDER") ?? "./keys";
            SKSServer = Environment.GetEnvironmentVariable("SKS_SERVER") ?? "http://localhost:11371";
            KeyPrefix = Environment.GetEnvironmentVariable("KEY_PREFIX") ?? "";

            var mkrc = Environment.GetEnvironmentVariable("MAX_KEYRING_CACHE_SIZE") ?? "1000";
            MaxKeyRingCache = int.Parse(mkrc);

            var hp = Environment.GetEnvironmentVariable("HTTP_PORT") ?? "5100";
            HttpPort = int.Parse(hp);

            EnableRethinkSKS = Environment.GetEnvironmentVariable("ENABLE_RETHINKDB_SKS") == "true";
            RethinkDBHost = Environment.GetEnvironmentVariable("RETHINKDB_HOST") ?? "localhost";
            var rdbport = Environment.GetEnvironmentVariable("RETHINKDB_PORT") ?? "28015";
            RethinkDBPort = int.Parse(rdbport);
            var rdbpool = Environment.GetEnvironmentVariable("RETHINKDB_POOL_SIZE") ?? "10";
            RethinkDBPoolSize = int.Parse(rdbpool);

            DatabaseName = Environment.GetEnvironmentVariable("DATABASE_NAME") ?? "remote_signer";

            MasterGPGKeyPath = Environment.GetEnvironmentVariable("MASTER_GPG_KEY_PATH") ?? null;
            MasterGPGKeyPath = MasterGPGKeyPath == null || MasterGPGKeyPath.Trim().Length > 0 ? MasterGPGKeyPath : null;
            MasterGPGKeyPasswordPath = Environment.GetEnvironmentVariable("MASTER_GPG_KEY_PASSWORD_PATH") ?? null;
            MasterGPGKeyPasswordPath = MasterGPGKeyPasswordPath == null || MasterGPGKeyPasswordPath.Trim().Length > 0 ? MasterGPGKeyPasswordPath : null;
            MasterGPGKeyBase64Encoded = Environment.GetEnvironmentVariable("MASTER_GPG_KEY_BASE64_ENCODED") == "true";
            KeysBase64Encoded = Environment.GetEnvironmentVariable("KEYS_BASE64_ENCODED") == "true";

#pragma warning disable RECS0022 // A catch clause that catches System.Exception and has an empty body
            try { Directory.CreateDirectory(PrivateKeyFolder); } catch (Exception) { }
#pragma warning restore RECS0022 // A catch clause that catches System.Exception and has an empty body
        }
    }
}
