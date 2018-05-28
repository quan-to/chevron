using System;
using System.Collections.Generic;
using System.Linq;
using Org.BouncyCastle.Bcpg.OpenPgp;
using RemoteSigner.Log;
using RemoteSigner.Models;

namespace RemoteSigner {
    public class KeyRingManager {
        readonly int MaxCacheKeys;
        public int Count { get { return fingerPrints.Count; } }

        Dictionary<string, PgpPublicKey> publicKeys;
        Dictionary<string, KeyInfo> publicKeysInfo;
        Queue<string> fingerPrints;
        PublicKeyStore pks;

        Dictionary<string, string> FP8TO16;

        public KeyRingManager() {
            MaxCacheKeys = Configuration.MaxKeyRingCache;
            publicKeys = new Dictionary<string, PgpPublicKey>();
            publicKeysInfo = new Dictionary<string, KeyInfo>();
            fingerPrints = new Queue<string>();
            FP8TO16 = new Dictionary<string, string>();
            pks = new PublicKeyStore();
        }

        public void AddKey(string publicKey, bool nonErasable = false) {
            using (var s = Tools.GenerateStreamFromString(publicKey)) {
                var pgpPub = new PgpPublicKeyRing(PgpUtilities.GetDecoderStream(s));
                var pubKey = pgpPub.GetPublicKey();
                AddKey(pubKey, nonErasable);
            }
        }


        public void AddKey(PgpPublicKey publicKey, bool nonErasable = false) {
            var fingerPrint = Tools.H16FP(publicKey.GetFingerprint().ToHexString());
            if (!publicKeys.ContainsKey(fingerPrint)) {
                publicKeys[fingerPrint] = publicKey;
                FP8TO16[Tools.H8FP(fingerPrint)] = fingerPrint;
                publicKeysInfo[fingerPrint] = new KeyInfo {
                    FingerPrint = fingerPrint,
                    Identifier = publicKey.GetUserIds().Cast<string>().First(),
                    Bits = publicKey.BitStrength,
                    ContainsPrivateKey = false,
                    PrivateKeyDecrypted = false
                };
                if (!nonErasable) {
                    fingerPrints.Enqueue(fingerPrint);
                }
                if (Count > MaxCacheKeys) {
                    string fpToRemove = fingerPrints.Dequeue();
                    publicKeys.Remove(fpToRemove);
                    publicKeysInfo.Remove(fpToRemove);
                }
            }
        }

        public List<KeyInfo> CachedKeys {
            get { return publicKeysInfo.Values.ToList(); }
        }

        public bool ContainsKey(string fingerPrint) {
            return fingerPrint.Length == 8 ? FP8TO16.ContainsKey(fingerPrint) && publicKeys.ContainsKey(FP8TO16[fingerPrint]) : publicKeys.ContainsKey(fingerPrint);
        }

        public PgpPublicKey GetKey(string fingerPrint) {
            if (!ContainsKey(fingerPrint)) {
                Logger.Log("KeyRingManager", $"Key {fingerPrint} not found in local keyring. Fetching from KeyStore...");
                var key = pks.GetKey(fingerPrint);
                if (key == null) {
                    Logger.Error("KeyRingManager", $"Key {fingerPrint} not found in KeyStore.");
                    return null;
                }
                AddKey(key);
            }

            return fingerPrint.Length == 8 ? publicKeys[FP8TO16[fingerPrint]] : publicKeys[fingerPrint];
        }

        public PgpPublicKey this[string key] {
            get {
                return GetKey(key);
            }
            set {
                AddKey(value);
            }
        }
    }
}
