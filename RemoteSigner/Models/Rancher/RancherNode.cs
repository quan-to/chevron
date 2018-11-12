using System;
namespace RemoteSigner.Models.Rancher {
    public class RancherNode {
        public string UUID { get; set; }
        public string Name { get; set; }
        public string IPAddress { get; set; }
        public int ID { get; set; }
        public bool IsSelf { get; set; }
    }
}
