using System;
using RemoteSigner.Models;

namespace RemoteSigner.Exceptions {
    public class InvalidKeyPasswordException: ErrorObjectException {
        public InvalidKeyPasswordException(string fingerPrint) : base(new ErrorObject {
            Message = $"The password for key {fingerPrint} is invalid.",
            ErrorCode = ErrorCodes.InvalidFieldData,
            ErrorField = "password"
        }) { }
    }
}
