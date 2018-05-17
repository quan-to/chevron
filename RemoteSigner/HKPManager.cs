using System;
using RemoteSigner.Exceptions.HKP;
using RemoteSigner.Models.HKP;

namespace RemoteSigner {
    public class HKPManager {
        PublicKeyStore pks;

        public HKPManager() {
            pks = new PublicKeyStore();
        }

        public string Lookup(string operation, string options, string machineReadable, string noModification, string showFingerprint, string exactMatch, string searchData) {
            switch (operation) {
                case Operation.Get:
                    return OperationGet(options, machineReadable != null && machineReadable.ToLower() == "on", noModification != null && noModification.ToLower() == "on", searchData);
                case Operation.Index:
                    return OperationIndex(options, machineReadable != null && machineReadable.ToLower() == "on", noModification != null && noModification.ToLower() == "on", showFingerprint != null && showFingerprint.ToLower() == "on", exactMatch != null && exactMatch.ToLower() == "on", searchData);
                case Operation.Vindex:
                    return OperationVIndex(options, machineReadable != null && machineReadable.ToLower() == "on", noModification != null && noModification.ToLower() == "on", showFingerprint != null && showFingerprint.ToLower() == "on", exactMatch != null && exactMatch.ToLower() == "on", searchData);
                default:
                    throw new UnknownOperationException(operation);
            }
        }

        string OperationGet(string options, bool machineReadable, bool noModification, string searchData) {
            if (searchData.StartsWith("0x", StringComparison.InvariantCulture)) {
                return pks.GetKey(searchData.Substring(2));
            }

            var results = pks.Search(searchData, 0, 1);
            return results.Count > 0 ? results[0].AsciiArmoredPublicKey : null;
        }

        string OperationIndex(string options, bool machineReadable, bool noModification, bool showFingerPrint, bool exactMatch, string searchData) {
            throw new OperationNotImplemented("index");
        }

        string OperationVIndex(string options, bool machineReadable, bool noModification, bool showFingerPrint, bool exactMatch, string searchData) {
            throw new OperationNotImplemented("vindex");
        }

        public void Add(string key) {
            string fineKey = Tools.ValidateAndTrimGPGKey(key);
            string err = pks.AddKey(fineKey);
            if (err != "OK") {
                throw new HKPBaseException(err);
            }
        }
    }
}
