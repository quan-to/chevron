using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGSignData {
        public String FingerPrint { get; set; }
        public String Base64Data { get; set; }
    }
}
