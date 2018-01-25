using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Class)]
    public class REST : Attribute {
        readonly string path;
        public REST(string path) {
            this.path = path;
        }

        public string Path {
            get { return path; }
        }
    }
}
