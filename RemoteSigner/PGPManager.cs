using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Org.BouncyCastle.Bcpg;
using Org.BouncyCastle.Bcpg.OpenPgp;
using Org.BouncyCastle.Crypto.Parameters;
using Org.BouncyCastle.Math;
using Org.BouncyCastle.Security;
using RemoteSigner.Exceptions;
using RemoteSigner.Log;

namespace RemoteSigner {
    public class PGPManager {
        public String KeyFolder { private set; get; }

        Dictionary<string, PgpSecretKey> privateKeys;
        Dictionary<string, PgpPrivateKey> decryptedKeys;

        KeyRingManager krm;

        public PGPManager() {
            KeyFolder = Configuration.PrivateKeyFolder;

            privateKeys = new Dictionary<string, PgpSecretKey>();
            decryptedKeys = new Dictionary<string, PgpPrivateKey>();
            krm = new KeyRingManager();

            LoadKeys();
        }

        void LoadKeys() {
            Logger.Log($"Loading keys from {KeyFolder}");
            try {
                var files = Directory.GetFiles(KeyFolder).ToList();
                files.ForEach((f) => {
                    try {
                        Logger.Log($"Loading key at {f}");
                        LoadPrivateKeyFromFile(f);
                    } catch (Exception e) {
                        Logger.Error($"Error loading key at {f}: {e}");
                    }
                });
            } catch (Exception e) {
                Logger.Error($"Error Loading keys from {KeyFolder}: {e}");
            }
            Logger.Log($"Done loading keys...");
        }

        public void UnlockKey(string fingerPrint, string password) {
            if (privateKeys.ContainsKey(fingerPrint)) {
                if (decryptedKeys.ContainsKey(fingerPrint)) {
                    Logger.Debug($"Key {fingerPrint} is already unlocked.");
                    return;
                }
                Logger.Debug($"Decrypting key {fingerPrint}");
                var sec = privateKeys[fingerPrint];
                try {
                    var dec = sec.ExtractPrivateKey(password.ToCharArray());
                    if (!TestPrivateKey(sec.PublicKey, dec)) {
                        throw new InvalidKeyPasswordException(fingerPrint);
                    }
                    decryptedKeys[fingerPrint] = dec;
                } catch (Exception) {
                    throw new InvalidKeyPasswordException(fingerPrint);
                }
            } else {
                Logger.Error($"Key {fingerPrint} is not loaded!");
                throw new KeyNotLoadedException(fingerPrint);
            }
        }

        public string LoadPrivateKeyFromFile(string filename) {
            using (var s = File.OpenRead(filename)) {
                return LoadPrivateKey(s);
            }
        }

        public string LoadPrivateKey(Stream s) {
            var pgpSec = ReadSecretKey(s);
            string fingerPrint = Tools.H16FP(pgpSec.PublicKey.GetFingerprint().ToHexString());
            Logger.Debug($"Loaded key {fingerPrint}");
            privateKeys[fingerPrint] = pgpSec;
            krm.AddKey(pgpSec.PublicKey, true);
            return fingerPrint;
        }

        public string LoadPrivateKey(string key) {
            using (var s = Tools.GenerateStreamFromString(key)) {
                return LoadPrivateKey(s);
            }
        }

        public void SavePrivateKey(string fingerPrint, string privateKey) {
            File.WriteAllText(Path.Combine(KeyFolder, $"{fingerPrint}.key"), privateKey);
        }

        public Task<string> SignData(string fingerPrint, byte[] data, HashAlgorithmTag hash = HashAlgorithmTag.Sha512) {
            if (!decryptedKeys.ContainsKey(fingerPrint)) {
                throw new KeyNotDecryptedException(fingerPrint);
            }

            var pgpSec = privateKeys[fingerPrint];
            var pgpPrivKey = decryptedKeys[fingerPrint];

            return Task.Run(() => {
                using (var ms = new MemoryStream()) {
                    var s = new ArmoredOutputStream(ms);
                    using (var bOut = new BcpgOutputStream(s)) {
                        var sGen = new PgpSignatureGenerator(pgpSec.PublicKey.Algorithm, hash);
                        sGen.InitSign(PgpSignature.BinaryDocument, pgpPrivKey);
                        sGen.Update(data, 0, data.Length);
                        sGen.Generate().Encode(bOut);
                        s.Close();
                        ms.Seek(0, SeekOrigin.Begin);
                        return Encoding.UTF8.GetString(ms.ToArray());
                    }
                }
            });
        }

        public bool VerifySignature(byte[] data, string signature, PgpPublicKey publicKey = null) {
            PgpSignatureList p3 = null;
            using (var inputStream = PgpUtilities.GetDecoderStream(Tools.GenerateStreamFromString(signature))) {
                var pgpFact = new PgpObjectFactory(inputStream);
                var o = pgpFact.NextPgpObject();
                if (o is PgpCompressedData c1) {
                    pgpFact = new PgpObjectFactory(c1.GetDataStream());
                    p3 = (PgpSignatureList)pgpFact.NextPgpObject();
                } else {
                    p3 = (PgpSignatureList)o;
                }
            }

            var sig = p3[0];
            if (publicKey == null) {
                string fingerPrint = Tools.H16FP(sig.KeyId.ToString("X").ToUpper());
                publicKey = krm[fingerPrint];
                if (publicKey == null) {
                    throw new KeyNotLoadedException(fingerPrint);
                }
            }
            sig.InitVerify(publicKey);
            sig.Update(data);

            return sig.Verify();
        }

        public Task<string> GenerateGPGKey(string identifier, string password, int bits = 3072) {
            return Task.Run(() => {
                using (var ms = new MemoryStream()) {
                var s = new ArmoredOutputStream(ms);
                    var kpg = GeneratorUtilities.GetKeyPairGenerator("RSA");
                    kpg.Init(new RsaKeyGenerationParameters(BigInteger.ValueOf(0x10001), new SecureRandom(), bits, 25));
                    var kp = kpg.GenerateKeyPair();

                    var secretKey = new PgpSecretKey(
                        PgpSignature.DefaultCertification,
                        PublicKeyAlgorithmTag.RsaGeneral,
                        kp.Public,
                        kp.Private,
                        DateTime.UtcNow,
                        identifier,
                        SymmetricKeyAlgorithmTag.Cast5,
                        password.ToCharArray(),
                        null,
                        null,
                        new SecureRandom()
                    );

                    secretKey.Encode(s);
                    s.Close();
                    ms.Seek(0, SeekOrigin.Begin);
                    var reader = new StreamReader(ms);
                    return reader.ReadToEnd();
                }
            });
        }

        public string GetPublicKeyASCII(string fingerPrint) {
            var publicKey = krm[fingerPrint];
            if (publicKey == null) {
                throw new KeyNotLoadedException(fingerPrint);
            }
            using (var ms = new MemoryStream()) {
                var s = new ArmoredOutputStream(ms);
                publicKey.Encode(s);
                s.Close();
                ms.Seek(0, SeekOrigin.Begin);
                var reader = new StreamReader(ms);
                return reader.ReadToEnd();
            }
        }

        #region Private Methods
        /**
         * A simple routine that opens a key ring file and loads the first available key
         * suitable for signature generation.
         * 
         * @param input stream to read the secret key ring collection from.
         * @return a secret key.
         * @throws IOException on a problem with using the input stream.
         * @throws PGPException if there is an issue parsing the input stream.
         */
        internal static PgpSecretKey ReadSecretKey(Stream input) {
            var pgpSec = new PgpSecretKeyRingBundle(PgpUtilities.GetDecoderStream(input));

            foreach (PgpSecretKeyRing keyRing in pgpSec.GetKeyRings()) {
                foreach (PgpSecretKey key in keyRing.GetSecretKeys()) {
                    if (key.IsSigningKey) {
                        return key;
                    }
                }
            }

            throw new ArgumentException("Can't find signing key in key ring.");
        }

        bool TestPrivateKey(PgpPublicKey publicKey, PgpPrivateKey privateKey) {
            try {
                byte[] testData = Encoding.ASCII.GetBytes("testdata");
                var signature = "";
                using (var ms = new MemoryStream()) {
                    var s = new ArmoredOutputStream(ms);
                    using (var bOut = new BcpgOutputStream(s)) {
                        var sGen = new PgpSignatureGenerator(publicKey.Algorithm, HashAlgorithmTag.Sha512);
                        sGen.InitSign(PgpSignature.BinaryDocument, privateKey);
                        sGen.Update(testData);
                        sGen.Generate().Encode(bOut);
                        s.Close();
                        ms.Seek(0, SeekOrigin.Begin);
                        signature = Encoding.UTF8.GetString(ms.ToArray());
                    }
                }

                return VerifySignature(testData, signature, publicKey);
            } catch (Exception e) {
                Logger.Error($"Error verifing private key: {e}");
                return false;
            }
        }
        #endregion
    }
}
