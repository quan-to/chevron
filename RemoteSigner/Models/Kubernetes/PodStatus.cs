using System;
using Newtonsoft.Json;

namespace RemoteSigner.Models.Kubernetes {
    public class PodStatus {
        [JsonProperty("phase")] 
        public string Phase { get; set; }
        
        [JsonProperty("hostIP")] 
        public string HostIP { get; set; }
        
        [JsonProperty("podIP")] 
        public string PodIP { get; set; }
        
        [JsonProperty("startTime")] 
        public DateTime StartTime { get; set; }
    }
}
