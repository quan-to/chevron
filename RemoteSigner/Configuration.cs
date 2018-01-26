using System;
using System.IO;

namespace RemoteSigner {
    public static class Configuration {

        public static string SyslogServer { get; private set; }
        public static string SyslogFacility { get; private set; }
        public static string PrivateKeyFolder { get; private set; }
        public static string SKSServer { get; private set; }
        public static int MaxKeyRingCache { get; private set; }

        static Configuration() {
            SyslogServer = Environment.GetEnvironmentVariable("SYSLOG_IP") ?? "127.0.0.1";
            SyslogFacility = Environment.GetEnvironmentVariable("SYSLOG_FACILITY") ?? "LOG_USER";
            PrivateKeyFolder = Environment.GetEnvironmentVariable("PRIVATE_KEY_FOLDER") ?? "./keys";
            SKSServer = Environment.GetEnvironmentVariable("SKS_SERVER") ?? "http://pgp.mit.edu/";

            var mkrc = Environment.GetEnvironmentVariable("MAX_KEYRING_CACHE_SIZE") ?? "1000";
            MaxKeyRingCache = int.Parse(mkrc);

#pragma warning disable RECS0022 // A catch clause that catches System.Exception and has an empty body
            try { Directory.CreateDirectory(PrivateKeyFolder); } catch (Exception) { }
#pragma warning restore RECS0022 // A catch clause that catches System.Exception and has an empty body
        }
    }
}
