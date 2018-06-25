using System;
using Newtonsoft.Json.Linq;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;
using ContaQuanto.FieldCipher;
using System.Linq;
using System.Collections.Generic;
using ContaQuanto.FieldCipher.Models;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/fieldCipher")]
    public class JsonFieldCipher {
        #pragma warning disable CS0649
        [Inject]
        readonly PGPManager pgpManager;

        [POST("/cipher")]
        public FieldCipherPacket Cipher(FieldCipherInput data) {
            var keys = data.Keys.Select(k => pgpManager.GetPublicKeyASCII(k)).ToList();
            var errors = new List<ErrorObject>();
            keys.ForEach((k) => {
                var idx = keys.IndexOf(k);
                var key = data.Keys[idx];
                if (k == null) {
                    errors.Add(new ErrorObject {
                        ErrorCode = ErrorCodes.NotFound,
                        ErrorField = $"Keys[{idx}]",
                        Message = $"Cannot find a key with fingerprint {key}",
                    });
                }
            });

            if (errors.Count > 0) {
                throw new ErrorObjectsException(errors);
            }

            var cipher = new Cipher(keys);
            return cipher.GenerateEncryptedPacket(data.JSON, data.SkipFields);
        }

        [POST("/decipher")]
        public FieldDecipherPacket Decipher(FieldDecipherInput data) {
            var key = pgpManager.GetPrivate(data.KeyFingerprint);
            if (key == null) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.NotFound,
                    ErrorField = "KeyFingerprint",
                    Message = $"There is no such unlocked private key with fingerprint {data.KeyFingerprint}",
                });
            }
            var decipher = new Decipher(key);
            return decipher.DecipherPacket(new FieldCipherPacket {
                EncryptedKey = data.EncryptedKey,
                EncryptedJSON = data.EncryptedJSON
            });
        }
    }
}
