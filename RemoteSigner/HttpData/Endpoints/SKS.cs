using System;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/sks")]
    public class SKS {
        [Inject]
        readonly SKSManager sks;

        [GET("/getKey")]
        public string GenerateKey([QueryParam] string fingerPrint) {
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

        [POST("/addKey")]
        public string AddKey(SKSAddKeyData data) {
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
