using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGGenerateKeyData {
        public String Identifier { get; set; }
        public String Password { get; set; }
        public int Bits { get; set; }
    }
}
