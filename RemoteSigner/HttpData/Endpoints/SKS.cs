using System.Collections.Generic;
using RemoteSigner.Database.Models;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/sks")]
    public class SKS {

        #region Injection
        // Disable Warning about null. This is a runtime injection.
        #pragma warning disable CS0649
        [Inject]
        readonly PublicKeyStore pks;
    
        #pragma warning restore CS0649
        #endregion
        [GET("/getKey")]
        public string GetKey([QueryParam] string fingerPrint) {
            var key = pks.GetKey(fingerPrint);
            if (key == null) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.NotFound,
                    ErrorField = "FingerPrint",
                    Message = "Cannot find key in SKS Server"
                });
            }

            return key;
        }

        [GET("/searchByName")]
        public List<GPGKey> SearchByName([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            return pks.SearchByName(name, pageStart, pageEnd);
        }

        [GET("/searchByFingerPrint")]
        public List<GPGKey> SearchByFingerPrint([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            return pks.SearchByFingerPrint(name, pageStart, pageEnd);
        }

        [GET("/searchByEmail")]
        public List<GPGKey> SearchByEmail([QueryParam] string name, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            return pks.SearchByEmail(name, pageStart, pageEnd);
        }

        [GET("/search")]
        public List<GPGKey> Search([QueryParam] string valueData, [QueryParam] int? pageStart, [QueryParam] int? pageEnd) {
            return pks.Search(valueData, pageStart, pageEnd);
        }

        [POST("/addKey")]
        public string AddKey(SKSAddKeyData data) {
            return pks.AddKey(data.PublicKey);
        }
    }
}
