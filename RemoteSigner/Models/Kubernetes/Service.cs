using System.Collections.Generic;
using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class Service {
        
        [JsonProperty("kind")]
        public string Kind { get; set; }

        [JsonProperty("apiVersion")]
        public string APIVersion { get; set; }

        [JsonProperty("metadata")] 
        public object Metadata { get; set; }

        [JsonProperty("items")]
        public List<ServiceItem> Items { get; set; }
    }
}
