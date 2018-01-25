using System;
namespace RemoteSigner.Models.ArgumentModels {
    public class GPGVerifySignatureData {
        public String Base64Data { get; set; }
        public String Signature { get; set; }
    }
}
