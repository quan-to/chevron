using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Method)]
    public class PUT : Attribute, IHttpMethod {
        readonly string path;
        public PUT(string path) {
            this.path = path;
        }

        public string Method {
            get { return "PUT"; }
        }

        public string Path {
            get { return path; }
        }
    }
}
