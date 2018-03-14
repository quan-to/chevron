using System;
using System.Collections.Generic;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/keyRing")]
    public class KeyRing {

        #region Injection
        // Disable Warning about null. This is a runtime injection.
        #pragma warning disable CS0649
        [Inject]
        readonly PGPManager pgpManager;

        #pragma warning restore CS0649
        #endregion
        [GET("/getKey")]
        public string GenerateKey([QueryParam] string fingerPrint) {
            return pgpManager.GetPublicKeyASCII(fingerPrint);
        }

        [GET("/cachedKeys")]
        public List<KeyInfo> GetCachedKeys() {
            return pgpManager.GetCachedKeys();
        }

        [GET("/privateKeys")]
        public List<KeyInfo> GetPrivatedKeys() {
            return pgpManager.GetLoadedPrivateKeys();
        }

        [POST("/addPrivateKey")]
        public string AddPrivateKey(KeyRingAddPrivateKeyData data) {
            string fingerPrint = "";
            try {
                fingerPrint = pgpManager.LoadPrivateKey(data.EncryptedPrivateKey);
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "EncryptedPrivateKey",
                    ErrorData = e,
                    Message = "Invalid PGP Private Key"
                });
            }

            if (data.SaveToDisk) {
                pgpManager.SavePrivateKey(fingerPrint, data.EncryptedPrivateKey);
            }

            return "OK";
        }
    }
}
