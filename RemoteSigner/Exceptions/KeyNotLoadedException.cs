using System;
using RemoteSigner.Models;

namespace RemoteSigner.Exceptions {
    public class KeyNotLoadedException: ErrorObjectException {
        public KeyNotLoadedException(string fingerPrint) : base(new ErrorObject {
            Message = $"The key {fingerPrint} is not loaded.",
            ErrorCode = ErrorCodes.NoDataAvailable,
            ErrorField = "key"
        }) {}
    }
}
