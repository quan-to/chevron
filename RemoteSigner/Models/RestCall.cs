using System;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.Models {
    class RestCall {
        public string className;
        public string methodName;
        public IHttpMethod method;
        public REST baseRest;
    }
}
