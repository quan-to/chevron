using System;
using System.Collections.Generic;
using Newtonsoft.Json;
using RemoteSigner.Database.Attributes;
using RethinkDb.Driver;
using RethinkDb.Driver.Model;
using RethinkDb.Driver.Net;


namespace RemoteSigner.Database.Models {
    [DBTable("gpgKey")]
    public class GPGKey {

        const int DEFAULT_PAGE_START = 0;
        const int DEFAULT_PAGE_END = 100;

        [JsonIgnore]
        static readonly RethinkDB R = RethinkDB.R;

        [JsonProperty("id", NullValueHandling = NullValueHandling.Ignore)]
        public string Id { get; set; }

        [DBIndex]
        public string FullFingerPrint { get; set; }

        [DBIndex]
        public List<string> Names { get; set; }

        [DBIndex]
        public List<string> Emails { get; set; }
        
        [DBIndex]
        public List<string> Subkeys { get; set; }

        public List<GPGKeyUid> KeyUids { get; set; }

        public int KeyBits { get; set; }

        public string AsciiArmoredPublicKey { get; set; }

        public string AsciiArmoredPrivateKey { get; set; }

        public bool ShouldSerializeAsciiArmoredPrivateKey() {
            return false;
        }

        public void Save(Connection conn) {
            R.Table("gpgKey")
                .Get(Id)
                .Update(this)
                .RunNoReply(conn);
        }

        public static Result AddGPGKey(Connection conn, GPGKey data) {
            var existing = R.Table("gpgKey")
                            .GetAll(data.FullFingerPrint)
                            .OptArg("index", "FullFingerPrint")
                            .RunCursor<GPGKey>(conn);
            if (existing.MoveNext()) {
                return R.Table("gpgKey").Get(existing.Current.Id).Update(data).RunResult(conn);
            }
            return R.Table("gpgKey").Insert(data).RunResult(conn);
        }

        public static GPGKey GetGPGKeyByFingerPrint(Connection conn, string fingerPrint) {
            List<GPGKey> s = R.Table("gpgKey")
                              .Filter(
                                (r) => r["FullFingerPrint"].Match($"{fingerPrint}$")
                                    .Or(r["Subkeys"].Filter(sk => sk.Match($"{fingerPrint}$")).Count().Gt(0)))
                              .Limit(1)
                              .CoerceTo("array")
                              .Run<List<GPGKey>>(conn);
            return s.Count > 0 ? s[0] : null;
        }

        public static List<GPGKey> SearchGPGByEmail(Connection conn, string email, int? pageStart, int? pageEnd) {
            return R.Table("gpgKey")
                    .Filter((r) => r["Email"].Filter((e) => e.Match(email)).Count().Gt(0))
                    .Slice(pageStart.GetValueOrDefault(DEFAULT_PAGE_START), pageEnd.GetValueOrDefault(DEFAULT_PAGE_END))
                    .CoerceTo("array")
                    .Run<List<GPGKey>>(conn);
        }

        public static List<GPGKey> SearchGPGByFingerPrint(Connection conn, string fingerPrint, int? pageStart, int? pageEnd) {
            return R.Table("gpgKey")
                    .Filter((r) => r["FullFingerPrint"].Match($"{fingerPrint}$"))
                    .Slice(pageStart.GetValueOrDefault(DEFAULT_PAGE_START), pageEnd.GetValueOrDefault(DEFAULT_PAGE_END))
                    .CoerceTo("array")
                    .Run<List<GPGKey>>(conn);
        }

        public static List<GPGKey> SearchGPGByName(Connection conn, string name, int? pageStart, int? pageEnd) {
            return R.Table("gpgKey")
                    .Filter((r) => r["Names"].Filter((n) => n.Match(name)).Count().Gt(0))
                    .Slice(pageStart.GetValueOrDefault(DEFAULT_PAGE_START), pageEnd.GetValueOrDefault(DEFAULT_PAGE_END))
                    .CoerceTo("array")
                    .Run<List<GPGKey>>(conn);
        }

        public static List<GPGKey> SearchGPGByAll(Connection conn, string valueData, int? pageStart, int? pageEnd) {
            return R.Table("gpgKey")
                    .Filter(
                        (r) => r["Names"].Filter((n) => n.Match(valueData)).Count().Gt(0)   // On Names
                        .Or(r["FullFingerPrint"].Match($"{valueData}$"))                    // On FingerPrint
                        .Or(r["Emails"].Filter((e) => e.Match(valueData)).Count().Gt(0)))    // On Email
                    .Slice(pageStart.GetValueOrDefault(DEFAULT_PAGE_START), pageEnd.GetValueOrDefault(DEFAULT_PAGE_END))
                    .CoerceTo("array")
                    .Run<List<GPGKey>>(conn);
        }

        internal static List<GPGKey> FetchKeyWithoutSubkey(Connection conn) {
            return R.Table("gpgKey")
                       .Filter(
                           (r) => r.HasFields("Subkeys").Not().Or(r["Subkeys"].Count().Eq(0)))
                       .CoerceTo("array")
                       .Run<List<GPGKey>>(conn);
        }
    }
}
