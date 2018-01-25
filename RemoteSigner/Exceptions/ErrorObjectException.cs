using System;
using RemoteSigner.Models;

namespace RemoteSigner.Exceptions {
    public class ErrorObjectException: Exception {
        public ErrorObject ErrorObject { get; private set; }
        public ErrorObjectException(ErrorObject eo) : base(eo.Message) {
            ErrorObject = eo;
        }
    }
}
