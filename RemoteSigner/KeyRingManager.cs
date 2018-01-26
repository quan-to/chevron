using System;
using System.Collections.Generic;
using Org.BouncyCastle.Bcpg.OpenPgp;
using RemoteSigner.Log;

namespace RemoteSigner {
    public class KeyRingManager {
        readonly int MaxCacheKeys;
        public int Count { get { return fingerPrints.Count; } }

        Dictionary<string, PgpPublicKey> publicKeys;
        Queue<string> fingerPrints;
        SKSManager sks;

        public KeyRingManager() {
            MaxCacheKeys = Configuration.MaxKeyRingCache;
            publicKeys = new Dictionary<string, PgpPublicKey>();
            fingerPrints = new Queue<string>();
            sks = new SKSManager();
        }

        public void AddKey(string publicKey, bool nonErasable = false) {
            using (var s = Tools.GenerateStreamFromString(publicKey)) {
                var pgpPub = new PgpPublicKeyRing(PgpUtilities.GetDecoderStream(s));
                var pubKey = pgpPub.GetPublicKey();
                var fingerPrint = Tools.H16FP(pubKey.GetFingerprint().ToHexString());
                if (!publicKeys.ContainsKey(fingerPrint)) {
                    publicKeys[fingerPrint] = pubKey;
                    if (!nonErasable) {
                        fingerPrints.Enqueue(fingerPrint);
                    }
                    if (Count > MaxCacheKeys) {
                        string fpToRemove = fingerPrints.Dequeue();
                        publicKeys.Remove(fpToRemove);
                    }
                }
            }
        }


        public void AddKey(PgpPublicKey publicKey, bool nonErasable = false) {
            var fingerPrint = Tools.H16FP(publicKey.GetFingerprint().ToHexString());
            if (!publicKeys.ContainsKey(fingerPrint)) {
                publicKeys[fingerPrint] = publicKey;
                if (!nonErasable) {
                    fingerPrints.Enqueue(fingerPrint);
                }
                if (Count > MaxCacheKeys) {
                    string fpToRemove = fingerPrints.Dequeue();
                    publicKeys.Remove(fpToRemove);
                }
            }
        }

        public bool ContainsKey(string fingerPrint) {
            return publicKeys.ContainsKey(fingerPrint);
        }

        public PgpPublicKey GetKey(string fingerPrint) {
            if (!publicKeys.ContainsKey(fingerPrint)) {
                Logger.Log("KeyRingManager", $"Key {fingerPrint} not found in local keyring. Fetching from SKS...");
                var getTask = sks.GetSKSKey(fingerPrint);
                getTask.Wait();
                if (getTask.Result == null) {
                    Logger.Error("KeyRingManager", $"Key {fingerPrint} not found in SKS server.");
                    return null;
                }
                AddKey(getTask.Result);
            }

            return publicKeys[fingerPrint];
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
