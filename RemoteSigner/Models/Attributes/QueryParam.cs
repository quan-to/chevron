using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Parameter)]
    public class QueryParam : Attribute {
        readonly string paramName;
        public QueryParam() {

        }

        public QueryParam(string paramName) {
            this.paramName = paramName;
        }

        public string ParamName {
            get { return paramName; }
        }
    }
}
