using System;
using System.Collections.Generic;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/__internal")]
    public class Internal {

        [Inject]
        readonly SecretsManager sm;

        [GET("/__triggerKeyUnlock")]
        public string TriggerKeyUnlock() {
            sm.UnlockLocalKeys().Wait();
            return "OK";
        }

        [GET("/__getUnlockPasswords")]
        public Dictionary<string, string> GetUnlockPasswords() {
            return sm.GetKeys();
        }
    }
}
