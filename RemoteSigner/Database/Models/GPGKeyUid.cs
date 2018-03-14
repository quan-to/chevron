using System;
using Newtonsoft.Json;

namespace RemoteSigner.Database.Models {
    public class GPGKeyUid {
        public string Name { get; set; }
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public string Email { get; set; }
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public string Description { get; set; }
    }
}
