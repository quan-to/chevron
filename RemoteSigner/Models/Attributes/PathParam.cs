using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Parameter)]
    public class PathParam : Attribute {
        readonly string paramName;
        public PathParam() {

        }

        public PathParam(string paramName) {
            this.paramName = paramName;
        }

        public string ParamName {
            get { return paramName; }
        }
    }
}
