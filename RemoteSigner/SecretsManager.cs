using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using RemoteSigner.Log;
using RemoteSigner.Models.ArgumentModels;

namespace RemoteSigner {
    public class SecretsManager {

        // fingerPrint - gpg encrypted master
        readonly Dictionary<string, string> EncryptedPasswords;
        readonly PGPManager gpg; // Isolated from global keyring
        readonly string masterKeyFingerprint;

        public SecretsManager() {
            EncryptedPasswords = new Dictionary<string, string>();
            gpg = new PGPManager();
            if (Configuration.MasterGPGKeyPath != null) {
                masterKeyFingerprint = gpg.LoadPrivateKeyFromFile(Configuration.MasterGPGKeyPath);
            } else {
                masterKeyFingerprint = null;
                Logger.Warn("SecretsManager", "No master key specified. No Secrets master available...");
            }
        }

        public void PutKeyPassword(string fingerPrint, string password) {
            if (masterKeyFingerprint == null) {
                Logger.Warn("SecretsManager", "Not saving password. Master Key not loaded.");
                return;
            }

            Logger.Log("SecretsManager", $"Saving password for key {fingerPrint}");
            EncryptedPasswords[fingerPrint] = gpg.Encrypt($"key-password-utf8-{fingerPrint}.txt", Encoding.UTF8.GetBytes(password), masterKeyFingerprint);
        }

        public void PutEncryptedKeyPassword(string fingerPrint, string password) {
            if (masterKeyFingerprint == null) {
                Logger.Warn("SecretsManager", "Not saving password. Master Key not loaded.");
                return;
            }
            EncryptedPasswords[fingerPrint] = password;
        }

        public Dictionary<string, string> GetKeys() {
            var dict = new Dictionary<string, string>(); // Copy Passwords
            foreach (var key in EncryptedPasswords.Keys) {
                dict[key] = EncryptedPasswords[key];
            }

            return dict;
        }

        public Task UnlockLocalKeys() {
            return Task.Run(async () => {
                if (masterKeyFingerprint == null) {
                    Logger.Error("SecretsManager", "Cannot unlock keys. Master key not loaded.");
                    return;
                }

                Logger.Log("SecretsManager", "Loading encrypted keys");
                var keys = GetKeys();
                Logger.Log("SecretsManager", "Decrypting master key");
                string pass = File.ReadAllText(Configuration.MasterGPGKeyPasswordPath, Encoding.UTF8);
                gpg.UnlockKey(masterKeyFingerprint, pass);
                Logger.Log("SecretsManager", "Starting key unlock");
                foreach (var key in keys.Keys) {
                    Logger.Log("SecretsManager", $"Unlocking key {key}");
                    string enc = keys[key];
                    var dec = gpg.Decrypt(enc);
                    var decPass = Encoding.UTF8.GetString(Convert.FromBase64String(dec.Base64Data));

                    var payload = new GPGUnlockKeyData {
                        FingerPrint = key,
                        Password = decPass,
                    };

                    string response = await Tools.Post("http://localhost:5100/remoteSigner/gpg/unlockKey", JsonConvert.SerializeObject(payload));
                    if (response != "OK") {
                        Logger.Error("SecretsManager", $"Error unlocking key: {response}");
                    }
                }
            });
        }
    }
}
