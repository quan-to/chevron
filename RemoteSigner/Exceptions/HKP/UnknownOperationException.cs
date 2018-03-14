using System;
namespace RemoteSigner.Exceptions.HKP {
    public class UnknownOperationException: HKPBaseException {
        public UnknownOperationException(string operation) : base($"Unknown operation named {operation}"){
        }
    }
}
