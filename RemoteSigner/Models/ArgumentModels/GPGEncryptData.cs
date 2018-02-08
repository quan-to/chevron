using System;
namespace RemoteSigner.Models.ArgumentModels {
    public class GPGEncryptData {
        public String FingerPrint { get; set; }
        public String Base64Data { get; set; }
        public String Filename { get; set; }
    }
}
