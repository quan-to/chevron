using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGEncryptData {
        public String FingerPrint { get; set; }
        public String Base64Data { get; set; }
        public String Filename { get; set; }
    }
}
