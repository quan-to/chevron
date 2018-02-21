using System;
namespace RemoteSigner.Exceptions.HKP {
    public class HKPBaseException : Exception {
        public HKPBaseException() { }
        public HKPBaseException(string message) : base(message) { }
    }
}
