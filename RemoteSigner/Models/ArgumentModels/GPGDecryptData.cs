using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct GPGDecryptData {
        public String AsciiArmoredData { get; set; }
        public Boolean DataOnly { get; set; }
    }
}
