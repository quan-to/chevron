using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Method)]
    public class DELETE : Attribute, IHttpMethod {
        readonly string path;
        public DELETE(string path) {
            this.path = path;
        }

        public string Method {
            get { return "DELETE"; }
        }

        public string Path {
            get { return path; }
        }
    }
}
