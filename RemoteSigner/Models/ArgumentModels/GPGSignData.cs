using System;
namespace RemoteSigner.Models.ArgumentModels {
    public class GPGSignData {
        public String FingerPrint { get; set; }
        public String Base64Data { get; set; }
    }
}
