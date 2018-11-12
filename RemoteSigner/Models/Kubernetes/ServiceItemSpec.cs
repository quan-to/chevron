using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class ServiceItemSpec {
        
        [JsonProperty("ports")]
        public List<ServiceItemSpecPort> Ports { get; set; }
        
        [JsonProperty("selector")]
        public Dictionary<string, string> Selector { get; set; }
        
        [JsonProperty("clusterIP")]
        public string ClusterIP { get; set; }
        
        [JsonProperty("sessionAffinity")]
        public string SessionAffinity { get; set; }
    }
}
