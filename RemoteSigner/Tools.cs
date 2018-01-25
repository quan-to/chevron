using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using Org.BouncyCastle.Bcpg;

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

        public static String Quanto2GPG(string signature) {
            var sig = "-----BEGIN PGP SIGNATURE-----\nVersion: Quanto\n";
            var s = signature.Split('$');
            if (s.Length != 3) {
                s = signature.Split('_');
            }
            if (s.Length != 3) {
                return null;
            }
            string gpgSig = s[2];
            string checksum = gpgSig.Substring(gpgSig.Length - 4, 4);
            for (int i = 0; i < gpgSig.Length - 4; i++) {
                if (i % 64 == 0) {
                    sig += '\n';
                }
                sig += gpgSig[i];
            }
            return $"{sig}\n{checksum}\n-----END PGP SIGNATURE-----";
        }

        public static String GPG2Quanto(string signature, string fingerPrint, HashAlgorithmTag hash) {
            string hashName = hash.ToString().ToUpper();
            string cutSig = "";

            string[] s = signature.Trim().Split('\n');

            for (int i = 2; i < s.Length - 1; i++) {
                cutSig += s[i];
            }

            return $"{fingerPrint}_{hashName}_{cutSig}";
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
