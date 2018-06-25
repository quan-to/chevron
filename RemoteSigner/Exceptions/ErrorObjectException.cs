using System;
using System.Collections.Generic;
using RemoteSigner.Models;

namespace RemoteSigner.Exceptions {
    public class ErrorObjectException : Exception {
        public ErrorObject ErrorObject { get; private set; }
        public ErrorObjectException(ErrorObject eo) : base(eo.Message) {
            ErrorObject = eo;
        }
    }
    public class ErrorObjectsException : Exception {
        public List<ErrorObject> ErrorObjects { get; private set; }
        public ErrorObjectsException(List<ErrorObject> eo) : base(eo[0].Message) {
            ErrorObjects = eo;
        }
    }
}
