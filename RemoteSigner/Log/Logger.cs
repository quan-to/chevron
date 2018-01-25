/**
 *  Made by Lucas Teske - Original Project: https://github.com/opensatelliteproject/
 */
using System;
using System.Net.Sockets;
using System.Threading;
using System.Text;
using System.Globalization;

namespace RemoteSigner.Log {
    /**
      * A simple logger that outputs to console and syslog
      */
    public class Logger {
        static Logger GlobalConsole;
        static readonly SyslogClient syslog = new SyslogClient();

        readonly Mutex messageMutex;
        public static bool LogConsole { get; set; }
        public delegate void ConsoleEvent(ConsoleMessage data);
        public static event ConsoleEvent MessageAvailable;

        public static Boolean GlobalEnableInfo { get { return GlobalConsole.EnableInfo; } set { GlobalConsole.EnableInfo = value; } }
        public static Boolean GlobalEnableWarn { get { return GlobalConsole.EnableWarn; } set { GlobalConsole.EnableWarn = value; } }
        public static Boolean GlobalEnableError { get { return GlobalConsole.EnableError; } set { GlobalConsole.EnableError = value; } }
        public static Boolean GlobalEnableDebug { get { return GlobalConsole.EnableDebug; } set { GlobalConsole.EnableDebug = value; } }

        public Boolean EnableInfo { get; set; }
        public Boolean EnableWarn { get; set; }
        public Boolean EnableError { get; set; }
        public Boolean EnableDebug { get; set; }

        static string RemoveDiacritics(string text) {
            var normalizedString = text.Normalize(NormalizationForm.FormD);
            var stringBuilder = new StringBuilder();

            foreach (var c in normalizedString) {
                var unicodeCategory = CharUnicodeInfo.GetUnicodeCategory(c);
                if (unicodeCategory != UnicodeCategory.NonSpacingMark) {
                    stringBuilder.Append(c);
                }
            }

            return stringBuilder.ToString().Normalize(NormalizationForm.FormC);
        }

        static Logger() {
            Init();
        }

        static void Init() {
            if (GlobalConsole == null) {
                GlobalConsole = new Logger();
            }
        }

        Logger() {
            EnableInfo = true;
            EnableError = true;
            EnableWarn = true;
            EnableDebug = false;
            LogConsole = true;
            messageMutex = new Mutex();
            if (Tools.IsLinux) {
                SyslogClient.SysLogServerIp = Configuration.SyslogServer;
            }
        }

        public static void Log(string message) {
            GlobalConsole._Log(message);
        }
        public static void Warn(string message) {
            GlobalConsole._Warn(message);
        }
        public static void Error(string message) {
            GlobalConsole._Error(message);
        }
        public static void Debug(string message) {
            GlobalConsole._Debug(message);
        }

        public static void Log(string context, string message) {
            GlobalConsole._Log($"{context} - {message}");
        }
        public static void Warn(string context, string message) {
            GlobalConsole._Warn($"{context} - {message}");
        }
        public static void Error(string context, string message) {
            GlobalConsole._Error($"{context} - {message}");
        }
        public static void Debug(string context, string message) {
            GlobalConsole._Debug($"{context} - {message}");
        }

        public void _Log(string message) {
            if (EnableInfo) {
                messageMutex.WaitOne();
                ConsoleMessage cm = new ConsoleMessage(ConsoleMessagePriority.INFO, message);
                if (LogConsole) {
                    ConsoleColor oldColor = Console.ForegroundColor;
                    Console.ForegroundColor = ConsoleColor.Blue;
                    Console.WriteLine(cm);
                    Console.ForegroundColor = oldColor;
                }

                if (syslog != null) {
                    try {
                        SyslogClient.Send(new Message(Configuration.SyslogFacility, Level.Information, cm.Message));
                    } catch (SocketException) {
                        // Syslog not configured, ignore.
                    }
                }
                messageMutex.ReleaseMutex();
                MessageAvailable?.Invoke(cm);
            }
        }

        public void _Warn(string message) {
            if (EnableWarn) {
                messageMutex.WaitOne();
                ConsoleMessage cm = new ConsoleMessage(ConsoleMessagePriority.WARN, message);
                if (LogConsole) {
                    ConsoleColor oldColor = Console.ForegroundColor;
                    Console.ForegroundColor = ConsoleColor.Yellow;
                    Console.WriteLine(cm);
                    Console.ForegroundColor = oldColor;
                }

                if (syslog != null) {
                    try {
                        SyslogClient.Send(new Message(Configuration.SyslogFacility, Level.Warning, cm.Message));
                    } catch (SocketException) {
                        // Syslog not configured, ignore.
                    }
                }
                messageMutex.ReleaseMutex();
                MessageAvailable?.Invoke(cm);
            }
        }

        public void _Error(string message) {
            if (EnableError) {
                messageMutex.WaitOne();
                ConsoleMessage cm = new ConsoleMessage(ConsoleMessagePriority.ERROR, message);
                if (LogConsole) {
                    ConsoleColor oldColor = Console.ForegroundColor;
                    Console.ForegroundColor = ConsoleColor.Red;
                    Console.WriteLine(cm);
                    Console.ForegroundColor = oldColor;
                }

                if (syslog != null) {
                    try {
                        SyslogClient.Send(new Message(Configuration.SyslogFacility, Level.Error, cm.Message));
                    } catch (SocketException) {
                        // Syslog not configured, ignore.
                    }
                }
                messageMutex.ReleaseMutex();
                MessageAvailable?.Invoke(cm);
            }
        }

        public void _Debug(string message) {
            if (EnableDebug) {
                messageMutex.WaitOne();
                ConsoleMessage cm = new ConsoleMessage(ConsoleMessagePriority.DEBUG, message);
                if (LogConsole) {
                    ConsoleColor oldColor = Console.ForegroundColor;
                    Console.ForegroundColor = ConsoleColor.DarkGray;
                    Console.WriteLine(cm);
                    Console.ForegroundColor = oldColor;
                }

                if (syslog != null) {
                    try {
                        SyslogClient.Send(new Message(Configuration.SyslogFacility, Level.Debug, cm.Message));
                    } catch (SocketException) {
                        // Syslog not configured, ignore.
                    }
                }
                messageMutex.ReleaseMutex();
                MessageAvailable?.Invoke(cm);
            }
        }
    }
}