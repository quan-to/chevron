using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class Pod {
        [JsonProperty("metadata")] 
        public ItemMetadata Metadata { get; set; }
        
        [JsonProperty("status")] 
        public PodStatus Status { get; set; }
    }
}
