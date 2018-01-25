using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;

namespace RemoteSigner {
    public static class Tools {
        public static bool IsLinux {
            get {
                int p = (int)Environment.OSVersion.Platform;
                return (p == 4) || (p == 6) || (p == 128);
            }
        }

        public static Stream GenerateStreamFromString(string s) {
            MemoryStream stream = new MemoryStream();
            StreamWriter writer = new StreamWriter(stream);
            writer.Write(s);
            writer.Flush();
            stream.Seek(0, SeekOrigin.Begin);
            return stream;
        }

        public static String H16FP(string fingerPrint) {
            if (fingerPrint.Length < 16) {
                throw new ArgumentException("FingerPrint string has less than 16 chars!");
            }
            return fingerPrint.Substring(fingerPrint.Length - 16, 16);
        }

        public static string ToHexString(this byte[] ba) {
            StringBuilder hex = new StringBuilder(ba.Length * 2);
            foreach (byte b in ba)
                hex.AppendFormat("{0:x2}", b);
            return hex.ToString().ToUpper();
        }

        public static Type[] GetTypesInNamespace(Assembly assembly, string nameSpace) {
            return assembly.GetTypes().Where(t => String.Equals(t.Namespace, nameSpace, StringComparison.Ordinal)).ToArray();
        }
    }
}
