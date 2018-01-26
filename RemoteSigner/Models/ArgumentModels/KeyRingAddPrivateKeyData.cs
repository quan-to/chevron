using System;
namespace RemoteSigner.Models.ArgumentModels {
    public class KeyRingAddPrivateKeyData {
        public String EncryptedPrivateKey { get; set; }
        public bool SaveToDisk { get; set; }
    }
}
