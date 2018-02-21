using System;
using RemoteSigner.Exceptions.HKP;
using RemoteSigner.Models.HKP;

namespace RemoteSigner {
    public class HKPManager {
        public string Lookup(SKSManager sks, string operation, string options, string machineReadable, string noModification, string showFingerprint, string exactMatch, string searchData) {
            switch (operation) {
                case Operation.Get:
                    return OperationGet(sks, options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", searchData);
                case Operation.Index:
                    return OperationIndex(sks, options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", showFingerprint.ToLower() == "on", exactMatch.ToLower() == "on", searchData);
                case Operation.Vindex:
                    return OperationVIndex(sks, options, machineReadable.ToLower() == "on", noModification.ToLower() == "on", showFingerprint.ToLower() == "on", exactMatch.ToLower() == "on", searchData);
                default:
                    throw new UnknownOperationException(operation);
            }
        }

        string OperationGet(SKSManager sks, string options, bool machineReadable, bool noModification, string searchData) {
            throw new OperationNotImplemented("index");
        }

        string OperationIndex(SKSManager sks, string options, bool machineReadable, bool noModification, bool showFingerPrint, bool exactMatch, string searchData) {
            throw new OperationNotImplemented("index");
        }

        string OperationVIndex(SKSManager sks, string options, bool machineReadable, bool noModification, bool showFingerPrint, bool exactMatch, string searchData) {
            throw new OperationNotImplemented("vindex");
        }

        public void Add(SKSManager sks, string key) {
            throw new OperationNotImplemented("add");
        }
    }
}
