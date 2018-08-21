using System;
using System.Collections.Generic;
using Newtonsoft.Json.Linq;

namespace RemoteSigner.Models.ArgumentModels {
    public struct FieldCipherInput {
        public JObject JSON { get; set; }
        public List<string> Keys { get; set; }
        public List<string> SkipFields { get; set; }
    }
}
