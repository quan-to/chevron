using System;
namespace RemoteSigner.Models.ArgumentModels {
    public class GPGUnlockKeyData {
        public String FingerPrint { get; set; }
        public String Password { get; set; }
    }
}
