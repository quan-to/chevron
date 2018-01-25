using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Class, AllowMultiple = true)]
    public class RestExceptionHandler : Attribute {
        public Type exceptionType;

        public RestExceptionHandler(Type exceptionType) {
            this.exceptionType = exceptionType;
        }

        public Type ExceptionType {
            get { return exceptionType; }
        }
    }
}
