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
                    return OperationGet(options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", searchData);
                case Operation.Index:
                    return OperationIndex(options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", showFingerprint.ToLower() == "on", exactMatch.ToLower() == "on", searchData);
                case Operation.Vindex:
                    return OperationVIndex(options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", showFingerprint.ToLower() == "on", exactMatch.ToLower() == "on", searchData);
                default:
                    throw new UnknownOperationException(operation);
            }
        }

        string OperationGet(string options, bool machineReadable, bool noModification, string searchData) {
            throw new OperationNotImplemented("index");
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
