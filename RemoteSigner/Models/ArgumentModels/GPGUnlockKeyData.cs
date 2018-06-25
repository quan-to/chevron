using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGUnlockKeyData {
        public String FingerPrint { get; set; }
        public String Password { get; set; }
    }
}
