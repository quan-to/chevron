using System;
namespace RemoteSigner.Exceptions {
    public class KeyNotLoadedException: Exception {
        public KeyNotLoadedException(string fingerPrint) : base($"The key {fingerPrint} is not loaded.") {}
    }
}
