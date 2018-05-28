using System;
using System.Net.Http;
using System.Collections.Generic;
using System.Net;
using System.Threading.Tasks;

namespace RemoteSigner {
    public class SKSManager {

        static readonly HttpClient client = new HttpClient();

        public String SKSURL { get; set; }

        public SKSManager() {
            SKSURL = Configuration.SKSServer;
        }

        public async Task<string> GetSKSKey(string fingerPrint) {
            var response = await (new HttpClient()).GetAsync($"{SKSURL}/pks/lookup?op=get&options=mr&search=0x{fingerPrint}");
            Console.WriteLine($"{SKSURL}/pks/lookup?op=get&options=mr&search=0x{fingerPrint}");
            if (response.StatusCode == HttpStatusCode.OK) {
                return await response.Content.ReadAsStringAsync();
            }

            return null;
        }

        public async Task<bool> PutSKSKey(string publicKey) {
            var values = new Dictionary<string, string> {{"keytext", publicKey }};
            var content = new FormUrlEncodedContent(values);
            var response = await client.PostAsync(SKSURL, content);
            return response.StatusCode == HttpStatusCode.OK;
        }
    }
}
