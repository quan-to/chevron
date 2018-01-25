using System;
namespace RemoteSigner.Exceptions {
    public class InvalidKeyPasswordException: Exception {
        public InvalidKeyPasswordException(string fingerPrint) : base($"The password for key {fingerPrint} is invalid.") { }
    }
}
