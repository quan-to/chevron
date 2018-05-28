using System;
using System.Collections.Generic;
using RemoteSigner.Database;
using RemoteSigner.Database.Models;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;

namespace RemoteSigner {
    public class PublicKeyStore {
        readonly SKSManager sks;

        public PublicKeyStore() {
            sks = new SKSManager();
        }

        public string GetKey(string fingerPrint) {
            if (!Configuration.EnableRethinkSKS) {
                var getTask = sks.GetSKSKey(fingerPrint);
                getTask.Wait();
                return getTask.Result;
            }

            var conn = DatabaseManager.GlobalDm.GetConnection();
            var key = GPGKey.GetGPGKeyByFingerPrint(conn, fingerPrint);
            return key?.AsciiArmoredPublicKey;
        }

        public List<GPGKey> SearchByName(string name, int? pageStart, int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByName(DatabaseManager.GlobalDm.GetConnection(), name, pageStart, pageEnd);
            }
            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        public List<GPGKey> SearchByFingerPrint(string name, int? pageStart, int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByFingerPrint(DatabaseManager.GlobalDm.GetConnection(), name, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        public List<GPGKey> SearchByEmail(string name, int? pageStart, int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByEmail(DatabaseManager.GlobalDm.GetConnection(), name, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        public List<GPGKey> Search( string valueData, int? pageStart, int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByAll(DatabaseManager.GlobalDm.GetConnection(), valueData, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        public string AddKey(string publicKey) {
            if (Configuration.EnableRethinkSKS) {
                var conn = DatabaseManager.GlobalDm.GetConnection();
                var res = GPGKey.AddGPGKey(conn, Tools.AsciiArmored2GPGKey(publicKey));
                if (res.Inserted == 0 && res.Unchanged == 0 && res.Replaced == 0) {
                    return res.FirstError;
                }

                return "OK";
            }

            var addTask = sks.PutSKSKey(publicKey);
            addTask.Wait();
            if (!addTask.Result) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InternalServerError,
                    ErrorField = "PublicKey",
                    Message = "Cannot add key in SKS Server"
                });
            }

            return "OK";
        }
    }
}
