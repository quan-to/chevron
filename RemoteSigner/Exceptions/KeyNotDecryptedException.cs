using System;
namespace RemoteSigner.Exceptions {
    public class KeyNotDecryptedException: Exception {
        public KeyNotDecryptedException(string fingerPrint) : base($"The key {fingerPrint} is not decrypted.") { }
    }
}
