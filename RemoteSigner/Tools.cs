using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Reflection;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using Org.BouncyCastle.Bcpg;
using Org.BouncyCastle.Bcpg.OpenPgp;
using RemoteSigner.Database.Models;

namespace RemoteSigner {
    public static class Tools {

        static readonly Regex NameObsEmailGPGReg = new Regex("(.*)\\s?(\\(.*\\))?\\s?<(.*)>", RegexOptions.IgnoreCase);
        static readonly Regex NameObsGPGReg = new Regex("(.*)\\s?\\((.*)\\)", RegexOptions.IgnoreCase);
        static readonly Regex NameEmailGPGReg = new Regex("(.*)\\s?<(.*)>", RegexOptions.IgnoreCase);
        static readonly HttpClient client = new HttpClient();

        public static async Task<string> Post(string url, string content) {
            var httpContent = new StringContent(content, Encoding.UTF8, "application/json");
            var response = await client.PostAsync(url, httpContent);
            return await response.Content.ReadAsStringAsync();
        }

        public static async Task<string> Get(string url) {
            var response = await client.GetAsync(url);
            return await response.Content.ReadAsStringAsync();
        }

        public static GPGKey AsciiArmored2GPGKey(string asciiArmored) {
            GPGKey key = null;
            using (var s = GenerateStreamFromString(asciiArmored)) {
                var pgpPub = new PgpPublicKeyRing(PgpUtilities.GetDecoderStream(s));
                var pubKey = pgpPub.GetPublicKey();
                key = new GPGKey {
                    Id = null,
                    FullFingerPrint = pubKey.GetFingerprint().ToHexString(),
                    AsciiArmoredPublicKey = asciiArmored,
                    AsciiArmoredPrivateKey = null,
                    Emails = new List<string>(),
                    Names = new List<string>(),
                    KeyUids = new List<GPGKeyUid>(),
                    KeyBits = pubKey.BitStrength,
                };

                foreach(string userId in pubKey.GetUserIds()) {
                    var m = NameObsEmailGPGReg.Match(userId);
                    var m2 = NameEmailGPGReg.Match(userId);
                    var email = "";
                    var name = "";
                    var obs = "";
                    if (m.Success && m.Groups.Count == 4) {
                        email = m.Groups[3].Value.Trim();
                        obs = m.Groups[2].Value.Trim();
                        key.Emails.Add(email);
                        var z = NameObsGPGReg.Match(m.Groups[1].Value);
                        if (z.Success && z.Groups.Count == 3) {
                            name = z.Groups[1].Value;
                            obs = z.Groups[2].Value;
                            key.Names.Add(name);
                        }
                    } else if (m2.Success && m2.Groups.Count == 3) {
                        name = m2.Groups[1].Value;
                        email = m2.Groups[2].Value;
                        key.Names.Add(name);
                        key.Emails.Add(email);
                    } else {
                        name = userId;
                        key.Names.Add(name);
                    }

                    key.KeyUids.Add(new GPGKeyUid {
                        Name = name,
                        Email = email,
                        Description = obs,
                    });
                }
            }
            return key;
        }

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
            string checksum = gpgSig.Substring(gpgSig.Length - 5, 5);
            for (int i = 0; i < gpgSig.Length - 5; i++) {
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

        public static String H8FP(string fingerPrint) {
            if (fingerPrint.Length < 8) {
                throw new ArgumentException("FingerPrint string has less than 8 chars!");
            }
            return fingerPrint.Substring(fingerPrint.Length - 8, 8);
        }

        public static string ToHexString(this byte[] ba) {
            StringBuilder hex = new StringBuilder(ba.Length * 2);
            foreach (byte b in ba)
                hex.AppendFormat("{0:x2}", b);
            return hex.ToString().ToUpper();
        }

        public static long TimestampMS() {
            return (Int64)(DateTime.UtcNow.Subtract(new DateTime(1970, 1, 1))).TotalMilliseconds;
        }

        public static Type[] GetTypesInNamespace(Assembly assembly, string nameSpace) {
            return assembly.GetTypes().Where(t => String.Equals(t.Namespace, nameSpace, StringComparison.Ordinal)).ToArray();
        }
    }
}
