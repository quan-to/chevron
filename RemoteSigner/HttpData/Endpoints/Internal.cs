using System;
using System.Collections.Generic;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/__internal")]
    public class Internal {

        #region Injection
        // Disable Warning about null. This is a runtime injection.
        #pragma warning disable CS0649
        [Inject]
        readonly SecretsManager sm;

        #pragma warning restore CS0649
        #endregion
        [GET("/__triggerKeyUnlock")]
        public string TriggerKeyUnlock() {
            sm.UnlockLocalKeys().Wait();
            return "OK";
        }

        [GET("/__getUnlockPasswords")]
        public Dictionary<string, string> GetUnlockPasswords() {
            return sm.GetKeys();
        }

        [POST("/__postEncryptedPasswords")]
        public string PostEncryptedPasswords(Dictionary<string, string> encryptedPasswords) {
            foreach (var key in encryptedPasswords.Keys) {
                sm.PutEncryptedKeyPassword(key, encryptedPasswords[key]);
            }
            return "OK";
        }
    }
}
