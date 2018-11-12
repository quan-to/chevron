using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class ServiceItem {
        [JsonProperty("metadata")]
        public ItemMetadata Metadata { get; set; }
        [JsonProperty("spec")]
        public ServiceItemSpec Spec { get; set; }
        [JsonProperty("status")]
        public object Status { get; set; }
    }
}
