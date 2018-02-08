using System;
using Org.BouncyCastle.Bcpg;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.ArgumentModels;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.Endpoints {
    [REST("/gpg")]
    public class GPG {
        [Inject]
        readonly PGPManager pgpManager;

        [POST("/generateKey")]
        public string GenerateKey(GPGGenerateKeyData data) {
            try {
                var genTask = pgpManager.GenerateGPGKey(data.Identifier, data.Password, data.Bits);
                genTask.Wait();
                return genTask.Result;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Password",
                    ErrorData = e,
                    Message = "Cannot Decrypt Key"
                });
            }
        }

        [POST("/unlockKey")]
        public string UnlockKey(GPGUnlockKeyData unlockData) {
            try {
                pgpManager.UnlockKey(unlockData.FingerPrint, unlockData.Password);
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Password",
                    ErrorData = e,
                    Message = "Cannot Decrypt Key"
                });
            }

            return "OK";
        }

        [POST("/sign")]
        public string Sign(GPGSignData data) {
            try {
                byte[] signData = Convert.FromBase64String(data.Base64Data);
                var sigTask = pgpManager.SignData(data.FingerPrint, signData);
                sigTask.Wait();
                return sigTask.Result;
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Signature",
                    ErrorData = e,
                    Message = "Cannot Verify Signature"
                });
            }
        }


        [POST("/signQuanto")]
        public string SignQuanto(GPGSignData data) {
            try {
                byte[] signData = Convert.FromBase64String(data.Base64Data);
                var sigTask = pgpManager.SignData(data.FingerPrint, signData);
                sigTask.Wait();

                var key = pgpManager[data.FingerPrint];

                return Tools.GPG2Quanto(sigTask.Result, key.PublicKey.GetFingerprint().ToHexString(), HashAlgorithmTag.Sha512);
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Signature",
                    ErrorData = e,
                    Message = "Cannot Verify Signature"
                });
            }
        }

        [POST("/verifySignatureQuanto")]
        public string VerifySignatureQuanto(GPGVerifySignatureData data) {
            try {
                byte[] verifyData = Convert.FromBase64String(data.Base64Data);
                string pgpSignature = Tools.Quanto2GPG(data.Signature);

                if (pgpSignature == null) {
                    throw new ErrorObjectException(new ErrorObject {
                        ErrorCode = ErrorCodes.InvalidSignature,
                        ErrorField = "Signature",
                        Message = "The provided Signature is invalid"
                    });
                }

                if (!pgpManager.VerifySignature(verifyData, pgpSignature)) {
                    throw new ErrorObjectException(new ErrorObject {
                        ErrorCode = ErrorCodes.InvalidSignature,
                        ErrorField = "Signature",
                        Message = "The provided Signature is invalid"
                    });
                }
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Signature",
                    ErrorData = e,
                    Message = "Cannot Verify Signature"
                });
            }

            return "OK";
        }

        [POST("/verifySignature")]
        public string VerifySignature(GPGVerifySignatureData data) {
            try {
                byte[] verifyData = Convert.FromBase64String(data.Base64Data);
                if (!pgpManager.VerifySignature(verifyData, data.Signature)) {
                    throw new ErrorObjectException(new ErrorObject {
                        ErrorCode = ErrorCodes.InvalidSignature,
                        ErrorField = "Signature",
                        Message = "The provided Signature is invalid"
                    });
                }
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "Signature",
                    ErrorData = e,
                    Message = "Cannot Verify Signature"
                });
            }

            return "OK";
        }

        [POST("/encrypt")]
        public string Encrypt(GPGEncryptData data) {
            try {
                string filename = data.Filename ?? $"QuantoEncrypt-{Tools.TimestampMS()}.bin";
                byte[] encryptData = Convert.FromBase64String(data.Base64Data);
                return pgpManager.Encrypt(filename, encryptData, data.FingerPrint);
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "",
                    ErrorData = e,
                    Message = "Cannot Encrypt Data"
                });
            }
        }

        [POST("/decrypt")]
        public GPGDecryptedDataReturn Decrypt(GPGDecryptData data) {
            try {
                return pgpManager.Decrypt(data.AsciiArmoredData);
            } catch (ErrorObjectException e) {
                throw e;
            } catch (Exception e) {
                throw new ErrorObjectException(new ErrorObject {
                    ErrorCode = ErrorCodes.InvalidFieldData,
                    ErrorField = "AsciiArmoredData",
                    ErrorData = e,
                    Message = "Cannot Decrypt Data"
                });
            }
        }

    }
}
