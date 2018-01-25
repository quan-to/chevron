using System;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/tests")]
    public class Tests {
        [GET("/ping")]
        public string Ping() {
            return "pong";
        }
    }
}
