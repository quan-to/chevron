using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Net.Security;
using System.Reflection;
using System.Security.Cryptography.X509Certificates;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using Org.BouncyCastle.Bcpg;
using Org.BouncyCastle.Bcpg.OpenPgp;
using RemoteSigner.Database.Models;
using RemoteSigner.Log;

namespace RemoteSigner {
    public static class Tools {

        static readonly Regex NameObsEmailGPGReg = new Regex("(.*)\\s?(\\(.*\\))?\\s?<(.*)>", RegexOptions.IgnoreCase);
        static readonly Regex NameObsGPGReg = new Regex("(.*)\\s?\\((.*)\\)", RegexOptions.IgnoreCase);
        static readonly Regex NameEmailGPGReg = new Regex("(.*)\\s?<(.*)>", RegexOptions.IgnoreCase);
        static readonly Regex PGPSig = new Regex("-----BEGIN PGP SIGNATURE-----(.*)-----END PGP SIGNATURE-----", RegexOptions.IgnoreCase | RegexOptions.Singleline);
        static readonly HttpClient client = new HttpClient();
        
        public static async Task<string> Post(string url, string content) {
            var httpContent = new StringContent(content, Encoding.UTF8, "application/json");
            var response = await client.PostAsync(url, httpContent);
            return await response.Content.ReadAsStringAsync();
        }

        public static async Task<string> Get(string url) {
            var response = await client.GetAsync(url);
            if (response.StatusCode != System.Net.HttpStatusCode.OK) {
                throw new Exception("Status Code != 200");
            }
            return await response.Content.ReadAsStringAsync();
        }
        
        public static async Task<string> Get(string url, string authorization) {
            var authClient = new HttpClient();
            authClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", authorization);
            var response = await authClient.GetAsync(url);
            if (response.StatusCode != HttpStatusCode.OK) {
                var data = await response.Content.ReadAsStringAsync();
                Logger.Debug("HttpClient", data);
                throw new Exception("Status Code != 200");
            }
            return await response.Content.ReadAsStringAsync();
        }

        public static string ValidateAndTrimGPGKey(string gpgKey) {
            try {
                using (var s = GenerateStreamFromString(gpgKey)) {
                    var pgpPub = new PgpPublicKeyRing(PgpUtilities.GetDecoderStream(s));
                    var pubKey = pgpPub.GetPublicKey();
                    return gpgKey;
                }
            } catch (Exception e) {
                throw e; // TODO
            }
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

        /// <summary>
        /// Adds CRC24 to signatures that does not have it.
        /// </summary>
        /// <returns>The fix.</returns>
        /// <param name="signature">Signature.</param>
        public static string SignatureFix(string signature) {
            var retSig = signature;
            var m = PGPSig.Match(signature);
            if (m.Groups.Count > 1) {
                var sig = "";
                var data = m.Groups[1].Value.TrimStart().TrimEnd().Split('\n');
                var save = false;
                if (data.Length == 1) {
                    sig = data[0];
                } else {
                    data.ToList().ForEach((l) => {
                        if (!save) {
                            save |= l.Length == 0;
                            if (l.Length > 2 && l.Substring(0, 2) == "iQ") { // Workarround for a GPG Bug in production
                                save = true;
                                sig += l;
                            }
                        } else {
                            sig += l;
                        }
                    });
                }
                try {
                    byte[] bData = Convert.FromBase64String(sig);
                    // Append checksum
                    var crc24 = new Crc24();
                    foreach (var b in bData) {
                        crc24.Update(b);
                    }
                    var crc = crc24.Value;
                    var crcu8 = new byte[3];
                    crcu8[0] = (byte)(crc >> 16 & 0xFF);
                    crcu8[1] = (byte)(crc >> 8 & 0xFF);
                    crcu8[2] = (byte)(crc & 0xFF);

                    retSig = "-----BEGIN PGP SIGNATURE-----\n\n";
                    retSig += sig + "\n=";
                    retSig += Convert.ToBase64String(crcu8);
                    retSig += "\n-----END PGP SIGNATURE-----";
                    return retSig;
                } catch (Exception) {
                    // Signature is already with checksum
                }
            }

            return retSig;
        }

        public static Stream GenerateStreamFromByteArray(byte[] data) {
            var stream = new MemoryStream();
            var writer = new StreamWriter(stream);
            writer.Write(data);
            writer.Flush();
            stream.Seek(0, SeekOrigin.Begin);
            return stream;
        }

        public static Stream GenerateStreamFromString(string s) {
            var stream = new MemoryStream();
            var writer = new StreamWriter(stream);
            writer.Write(s);
            writer.Flush();
            stream.Seek(0, SeekOrigin.Begin);
            return stream;
        }

        public static string Raw2AsciiArmored(byte[] b) {
            var encOut = new MemoryStream();
            var s = new ArmoredOutputStream(encOut);
            s.Write(b);
            s.Close();
            encOut.Seek(0, SeekOrigin.Begin);
            var reader = new StreamReader(encOut);
            return reader.ReadToEnd();
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
