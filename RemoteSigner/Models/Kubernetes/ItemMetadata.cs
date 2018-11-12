using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class ItemMetadata {
        
        [JsonProperty("name")]
        public string Name { get; set; }
        
        [JsonProperty("namespace")]
        public string Namespace { get; set; }
        
        [JsonProperty("selfLink")]
        public string SelfLink { get; set; }
        
        [JsonProperty("uid")]
        public string UID { get; set; }
        
        [JsonProperty("resourceVersion")]
        public string ResourceVersion { get; set; }
        
        [JsonProperty("creationTimestamp")]
        public DateTime CreationTimestamp { get; set; }
        
        [JsonProperty("labels")]
        public Dictionary<string, string> Labels { get; set; }
        
        [JsonProperty("annotations")]
        public Dictionary<string, string> Annotations { get; set; }
        
        [JsonProperty("ownerReferences")]
        public object OwnerReferences { get; set; } // TODO
    }
}
