using System;
using RemoteSigner.Models;

namespace RemoteSigner.Exceptions {
    public class KeyNotDecryptedException: ErrorObjectException {
        public KeyNotDecryptedException(string fingerPrint) : base(new ErrorObject {
            Message = $"The key {fingerPrint} is not decrypted.",
            ErrorCode = ErrorCodes.NoDataAvailable,
            ErrorField = "key"
        }) { }
    }
}
