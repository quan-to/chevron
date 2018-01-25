/**
 *  Made by Lucas Teske - Original Project: https://github.com/opensatelliteproject/
 */
using System;
using System.Net;
using System.Net.Sockets;
using System.Collections.Generic;

namespace RemoteSigner.Log {
    public class SyslogClient {
        static IPHostEntry ipHostInfo;
        static IPAddress ipAddress;
        static IPEndPoint ipLocalEndPoint;
        static UdpClient udpClient;
        static Dictionary<string, Facility> FacilityMap;
        public static int Port { get; set; }
        public static string SysLogServerIp { get; set; }
        public static bool IsActive { get; set; }

        static SyslogClient() {
            if (!Tools.IsLinux) {
                Console.WriteLine("Syslog only works on Linux");
                return;
            }

            ipHostInfo = Dns.GetHostEntry(Dns.GetHostName());
            ipAddress = ipHostInfo.AddressList[0];
            ipLocalEndPoint = new IPEndPoint(ipAddress, 0);
            udpClient = new UdpClient(ipLocalEndPoint);
            Port = 514;
            InitFacilityMap();
        }

        public void Close() {
            if (IsActive) {
                udpClient.Close();
                IsActive = false;
            }
        }

        public static void Send(string syslogFacility, Level level, string message) {
            Send(new Message(syslogFacility, level, message));
        }

        public static void Send(Message message) {
            if (Tools.IsLinux) {
                if (!IsActive) {
                    udpClient.Connect(SysLogServerIp, Port);
                    IsActive = true;
                }

                if (IsActive) {
                    int priority = (int)FacilityMap[message.Facility] * 8 + message.Level;
                    string msg = System.String.Format("<{0}>{1} {2} {3}", priority, DateTime.Now.ToString("MMM dd HH:mm:ss"), "XRIT", message.Text);
                    byte[] bytes = System.Text.Encoding.ASCII.GetBytes(msg);
                    udpClient.Send(bytes, bytes.Length);
                }
            }
        }

        static void InitFacilityMap() {
            FacilityMap = new Dictionary<string, Facility> {
                ["LOG_KERNEL"] = Facility.Kernel,
                ["LOG_USER"] = Facility.User,
                ["LOG_MAIL"] = Facility.Mail,
                ["LOG_DAEMON"] = Facility.Daemon,
                ["LOG_AUTH"] = Facility.Auth,
                ["LOG_SYSLOG"] = Facility.Syslog,
                ["LOG_LPR"] = Facility.Lpr,
                ["LOG_NEWS"] = Facility.News,
                ["LOG_UUCP"] = Facility.UUCP,
                ["LOG_CRON"] = Facility.Cron,
                ["LOG_LOCAL0"] = Facility.Local0,
                ["LOG_LOCAL1"] = Facility.Local1,
                ["LOG_LOCAL2"] = Facility.Local2,
                ["LOG_LOCAL3"] = Facility.Local3,
                ["LOG_LOCAL4"] = Facility.Local4,
                ["LOG_LOCAL5"] = Facility.Local5,
                ["LOG_LOCAL6"] = Facility.Local6,
                ["LOG_LOCAL7"] = Facility.Local7
            };
        }

    }
}