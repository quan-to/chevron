using System;
namespace RemoteSigner.Exceptions.HKP {
    public class OperationNotImplemented : HKPBaseException {
        public OperationNotImplemented(string operation) : base($"Operation {operation} is not implemented at this server.") {}
    }
}
