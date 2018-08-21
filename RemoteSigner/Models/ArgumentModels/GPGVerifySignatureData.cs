using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGVerifySignatureData {
        public String Base64Data { get; set; }
        public String Signature { get; set; }
    }
}
