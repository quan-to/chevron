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
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;

namespace RemoteSigner {
    public class PGPManager {

        const string PGPManagerLog = "PGPManager";

        public String KeyFolder { private set; get; }

        Dictionary<string, PgpSecretKey> privateKeys;
        Dictionary<string, PgpPrivateKey> decryptedKeys;
        Dictionary<string, string> FP8TO16;

        KeyRingManager krm;

        public PGPManager() {
            KeyFolder = Configuration.PrivateKeyFolder;

            privateKeys = new Dictionary<string, PgpSecretKey>();
            decryptedKeys = new Dictionary<string, PgpPrivateKey>();
            krm = new KeyRingManager();
            FP8TO16 = new Dictionary<string, string>();
            LoadKeys();
        }

        void LoadKeys() {
            Logger.Log(PGPManagerLog, $"Loading keys from {KeyFolder}");
            try {
                var files = Directory.GetFiles(KeyFolder).ToList();
                files.ForEach((f) => {
                    if ((Configuration.KeyPrefix.Length > 0 && f.StartsWith(Configuration.KeyPrefix, StringComparison.InvariantCultureIgnoreCase)) || Configuration.KeyPrefix.Length == 0) {
                        try {
                            Logger.Log(PGPManagerLog, $"Loading key at {f}");
                            LoadPrivateKeyFromFile(f);
                        } catch (Exception e) {
                            Logger.Error(PGPManagerLog, $"Error loading key at {f}: {e}");
                        }
                    }
                });
            } catch (Exception e) {
                Logger.Error(PGPManagerLog, $"Error Loading keys from {KeyFolder}: {e}");
            }
            Logger.Log(PGPManagerLog, $"Done loading keys...");
        }

        public void UnlockKey(string fingerPrint, string password) {
            if (fingerPrint.Length == 8 && FP8TO16.ContainsKey(fingerPrint)) {
                fingerPrint = FP8TO16[fingerPrint];
            }
            if (privateKeys.ContainsKey(fingerPrint)) {
                if (decryptedKeys.ContainsKey(fingerPrint)) {
                    Logger.Debug(PGPManagerLog, $"Key {fingerPrint} is already unlocked.");
                    return;
                }
                Logger.Debug(PGPManagerLog, $"Decrypting key {fingerPrint}");
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
                Logger.Error(PGPManagerLog, $"Key {fingerPrint} is not loaded!");
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
            Logger.Debug(PGPManagerLog, $"Loaded key {fingerPrint}");
            privateKeys[fingerPrint] = pgpSec;
            FP8TO16[Tools.H8FP(fingerPrint)] = fingerPrint;
            krm.AddKey(pgpSec.PublicKey, true);
            return fingerPrint;
        }

        public List<KeyInfo> GetCachedKeys() {
            return krm.CachedKeys;
        }

        public List<KeyInfo> GetLoadedPrivateKeys() {
            return privateKeys.Keys.Select((k) => new KeyInfo {
                FingerPrint = k,
                Identifier = privateKeys[k].UserIds.Cast<string>().First(),
                ContainsPrivateKey = true,
                PrivateKeyDecrypted = decryptedKeys.ContainsKey(k),
                Bits = privateKeys[k].PublicKey.BitStrength
            }).ToList();
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
            if (fingerPrint.Length == 8 && FP8TO16.ContainsKey(fingerPrint)) {
                fingerPrint = FP8TO16[fingerPrint];
            }
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

        public GPGDecryptedDataReturn Decrypt(string data) {
            using (var stream = PgpUtilities.GetDecoderStream(Tools.GenerateStreamFromString(data))) {
                var pgpF = new PgpObjectFactory(stream);
                var o = pgpF.NextPgpObject();
                var enc = o as PgpEncryptedDataList;
                if (enc == null) {
                    enc = (PgpEncryptedDataList)pgpF.NextPgpObject();
                }

                PgpPublicKeyEncryptedData pbe = null;
                PgpPrivateKey pgpPrivKey = null;
                PgpSecretKey pgpSec = null;
                string lastFingerPrint = "None";
                foreach (PgpPublicKeyEncryptedData pked in enc.GetEncryptedDataObjects()) {
                    string keyId = pked.KeyId.ToString("X").ToUpper();
                    string fingerPrint = keyId.Length < 16 ? FP8TO16[Tools.H8FP(keyId)] : Tools.H16FP(keyId);
                    lastFingerPrint = fingerPrint;
                    if (!decryptedKeys.ContainsKey(fingerPrint)) {
                        continue;
                    }

                    pgpSec = privateKeys[fingerPrint];
                    pgpPrivKey = decryptedKeys[fingerPrint];
                    pbe = pked;
                    break;
                }

                if (pbe == null) {
                    throw new KeyNotLoadedException(lastFingerPrint);
                }

                var clear = pbe.GetDataStream(pgpPrivKey);
                var plainFact = new PgpObjectFactory(clear);
                var message = plainFact.NextPgpObject();
                var outData = new GPGDecryptedDataReturn {
                    FingerPrint = lastFingerPrint,
                };
                if (message is PgpCompressedData cData) {
                    var pgpFact = new PgpObjectFactory(cData.GetDataStream());
                    message = pgpFact.NextPgpObject();
                }

                if (message is PgpLiteralData ld) {
                    outData.Filename = ld.FileName;
                    var iss = ld.GetInputStream();
                    byte[] buffer = new byte[16 * 1024];
                    using (var ms = new MemoryStream()) {
                        int read;
                        while ((read = iss.Read(buffer, 0, buffer.Length)) > 0) {
                            ms.Write(buffer, 0, read);
                        }
                        outData.Base64Data = Convert.ToBase64String(ms.ToArray());
                    }
                } else if (message is PgpOnePassSignatureList) {
                    throw new PgpException("Encrypted message contains a signed message - not literal data.");
                } else {
                    throw new PgpException("Message is not a simple encrypted file - type unknown.");
                }

                outData.IsIntegrityProtected = pbe.IsIntegrityProtected();

                if (outData.IsIntegrityProtected) {
                    outData.IsIntegrityOK = pbe.Verify();
                }

                return outData;
            }
        }

        public string Encrypt(string filename, byte[] data, string fingerPrint) {
            if (fingerPrint.Length == 8 && FP8TO16.ContainsKey(fingerPrint)) {
                fingerPrint = FP8TO16[fingerPrint];
            }

            var publicKey = krm[fingerPrint];
            if (publicKey == null) {
                throw new KeyNotLoadedException(fingerPrint);
            }

            return Encrypt(filename, data, publicKey);
        }

        public string Encrypt(string filename, byte[] data, PgpPublicKey publicKey) {
            using (MemoryStream encOut = new MemoryStream(), bOut = new MemoryStream()) {
                var comData = new PgpCompressedDataGenerator(CompressionAlgorithmTag.Zip);
                var cos = comData.Open(bOut); // open it with the final destination
                var lData = new PgpLiteralDataGenerator();
                var pOut = lData.Open(
                    cos,                    // the compressed output stream
                    PgpLiteralData.Binary,
                    filename,               // "filename" to store
                    data.Length,            // length of clear data
                    DateTime.UtcNow         // current time
                );
                pOut.Write(data, 0, data.Length);
                lData.Close();
                comData.Close();
                var cPk = new PgpEncryptedDataGenerator(SymmetricKeyAlgorithmTag.Cast5, true, new SecureRandom());
                cPk.AddMethod(publicKey);
                byte[] bytes = bOut.ToArray();
                var s = new ArmoredOutputStream(encOut);
                var cOut = cPk.Open(s, bytes.Length);
                cOut.Write(bytes, 0, bytes.Length);  // obtain the actual bytes from the compressed stream
                cOut.Close();
                s.Close();
                encOut.Seek(0, SeekOrigin.Begin);
                var reader = new StreamReader(encOut);
                return reader.ReadToEnd();
            }
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
                string keyId = sig.KeyId.ToString("X").ToUpper();
                string fingerPrint = keyId.Length < 16 ? Tools.H8FP(keyId) : Tools.H16FP(sig.KeyId.ToString("X").ToUpper());
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


        public PgpSecretKey GetKey(string fingerPrint) {
            if (fingerPrint.Length == 8 && FP8TO16.ContainsKey(fingerPrint)) {
                fingerPrint = FP8TO16[fingerPrint];
            }
            return privateKeys.ContainsKey(fingerPrint) ? privateKeys[fingerPrint] : null;
        }

        public PgpSecretKey this[string key] {
            get {
                return GetKey(key);
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
                Logger.Error(PGPManagerLog, $"Error verifing private key: {e}");
                return false;
            }
        }
        #endregion
    }
}
