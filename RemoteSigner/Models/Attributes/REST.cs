using System;
namespace RemoteSigner.Models.Attributes {
    [AttributeUsage(AttributeTargets.Class)]
    public class REST : Attribute {
        readonly string path;
        readonly Boolean addToRoot;

        public REST(string path) {
            this.path = path;
        }
        public REST(string path, bool addToRoot) {
            this.path = path;
            this.addToRoot = addToRoot;
        }

        public string Path {
            get { return path; }
        }
    }
}
