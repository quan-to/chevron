using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Method)]
    public class GET : Attribute, IHttpMethod {
        readonly string path;
        public GET(string path) {
            this.path = path;
        }

        public string Method {
            get { return "GET"; }
        }

        public string Path {
            get { return path; }
        }
    }
}
