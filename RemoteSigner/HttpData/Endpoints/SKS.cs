using System;
using System.Collections.Generic;
using RemoteSigner.Database;
using RemoteSigner.Database.Models;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/sks")]
    public class SKS {
        [Inject]
        readonly SKSManager sks;

        [Inject]
        readonly DatabaseManager dm;

        [GET("/getKey")]
        public string GetKey([QueryParam] string fingerPrint) {
            if (!Configuration.EnableRethinkSKS) {
                var getTask = sks.GetSKSKey(fingerPrint);
                getTask.Wait();
                if (getTask.Result == null) {
                    throw new ErrorObjectException(new ErrorObject {
                        ErrorCode = ErrorCodes.NotFound,
                        ErrorField = "FingerPrint",
                        Message = "Cannot find key in SKS Server"
                    });
                }

                return getTask.Result;
            }

            var conn = dm.GetConnection();
            var key = GPGKey.GetGPGKeyByFingerPrint(conn, fingerPrint);
            return key?.AsciiArmoredPublicKey;
        }

        [GET("/searchByName")]
        public List<GPGKey> SearchByName([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByName(dm.GetConnection(), name, pageStart, pageEnd);
            }
            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        [GET("/searchByFingerPrint")]
        public List<GPGKey> SearchByFingerPrint([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByFingerPrint(dm.GetConnection(), name, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        [GET("/searchByEmail")]
        public List<GPGKey> SearchByEmail([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByEmail(dm.GetConnection(), name, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        [GET("/search")]
        public List<GPGKey> Search([QueryParam] string valueData, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            if (Configuration.EnableRethinkSKS) {
                return GPGKey.SearchGPGByAll(dm.GetConnection(), valueData, pageStart, pageEnd);
            }

            throw new NotSupportedException("The server does not have RethinkDB enabled so it cannot serve search");
        }

        [POST("/addKey")]
        public string AddKey(SKSAddKeyData data) {
            if (Configuration.EnableRethinkSKS) {
                var conn = dm.GetConnection();
                var res = GPGKey.AddGPGKey(conn, Tools.AsciiArmored2GPGKey(data.PublicKey));
                if (res.Inserted == 0 && res.Unchanged == 0 && res.Replaced == 0) {
                    return res.FirstError;
                }

                return "OK";
            }

            var addTask = sks.PutSKSKey(data.PublicKey);
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
