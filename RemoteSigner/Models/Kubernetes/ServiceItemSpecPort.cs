using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class ServiceItemSpecPort {
        
        [JsonProperty("name")]
        public string Name { get; set; }
        
        [JsonProperty("protocol")]
        public string Protocol { get; set; }
        
        [JsonProperty("port")]
        public int Port { get; set; }
        
        [JsonProperty("targetPort")]
        public int TargetPort { get; set; }
    }
}
