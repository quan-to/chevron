using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Method)]
    public class POST : Attribute, IHttpMethod {
        readonly string path;
        public POST(string path) {
            this.path = path;
        }

        public string Method {
            get { return "POST"; }
        }

        public string Path {
            get { return path; }
        }
    }
}
