using System;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    /// <summary>
    /// HKP Server based on https://tools.ietf.org/html/draft-shaw-openpgp-hkp-00
    /// </summary>
    [REST("/pks")]
    public class PKS {

        #region Injection
        // Disable Warning about null. This is a runtime injection.
        #pragma warning disable CS0649
        [Inject]
        readonly HKPManager hkp;
        #pragma warning restore CS0649
        #endregion
        /// <summary>
        /// HKP Standard Lookup call
        /// </summary>
        /// <returns>The lookup result (depends on operation)</returns>
        /// <param name="op">Operation</param>
        /// <param name="options">Options</param>
        /// <param name="mr">Machine Readable</param>
        /// <param name="nm">No Modified</param>
        /// <param name="fingerprint">Show Fingerprint (on / off)</param>
        /// <param name="exact">Exact Match (on / off)</param>
        /// <param name="search">Search Data</param>
        [GET("/lookup")]
        public string Lookup([QueryParam] string op, [QueryParam] string options, [QueryParam] string mr, [QueryParam] string nm, [QueryParam] string fingerprint, [QueryParam] string exact, [QueryParam] string search) {
            return hkp.Lookup(op, options, mr, nm, fingerprint, exact, search);
        }

        /// <summary>
        /// Add a key to SKS
        /// </summary>
        /// <returns>Nothing</returns>
        /// <param name="keyData">ASCII Armored Public Key</param>
        [POST("/add")]
        public void Add(string keyData) {
            var s = keyData.Split('=');
            if (s.Length != 2) {
                throw new Exception("Invalid Payload");
            }
            var key = Uri.UnescapeDataString(s[1]);
            hkp.Add(key);
        }
    }
}
