using System;
namespace RemoteSigner.Models {
    public class KeyInfo {
        public string FingerPrint { get; set; }
        public string Identifier { get; set; }
        public int Bits { get; set; }
        public bool ContainsPrivateKey { get; set; }
        public bool PrivateKeyDecrypted { get; set; }
    }
}
