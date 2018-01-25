using System;
namespace RemoteSigner {
    public static class Configuration {

        public static string SyslogServer { get; private set; }
        public static string SyslogFacility { get; private set; }
        public static string PrivateKeyFolder { get; private set; }

        static Configuration() {
            SyslogServer = Environment.GetEnvironmentVariable("SYSLOG_IP") ?? "127.0.0.1";
            SyslogFacility = Environment.GetEnvironmentVariable("SYSLOG_FACILITY") ?? "LOG_USER";
            PrivateKeyFolder = Environment.GetEnvironmentVariable("PRIVATE_KEY_FOLDER") ?? "./keys";
        }
    }
}
