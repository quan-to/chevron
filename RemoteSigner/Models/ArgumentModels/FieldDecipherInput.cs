using System;
using Newtonsoft.Json.Linq;

namespace RemoteSigner.Models.ArgumentModels {
    public struct FieldDecipherInput {
        public string KeyFingerprint { get; set; }
        public string EncryptedKey { get; set; }
        public JObject EncryptedJSON { get; set; }
    }
}
