using System;
using System.Text;
using Newtonsoft.Json;
using RemoteSigner.Exceptions;
using RemoteSigner.Models;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.ExceptionHandlers {
    [RestExceptionHandler(typeof(ErrorObjectException))]
    public class ErrorObjectExceptionHandler : IRestExceptionHandler {
        public RestResult HandleException(Exception e) {
            var ce = e as ErrorObjectException;

            RestResult result = new RestResult {
                ContentType = MimeTypes.JSON,
                StatusCode = System.Net.HttpStatusCode.InternalServerError,
                Result = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(ce.ErrorObject))
            };
            return result;
        }
    }
    [RestExceptionHandler(typeof(ErrorObjectsException))]
    public class ErrorObjectsExceptionHandler : IRestExceptionHandler {
        public RestResult HandleException(Exception e) {
            var ce = e as ErrorObjectsException;

            RestResult result = new RestResult() {
                ContentType = MimeTypes.JSON,
                StatusCode = System.Net.HttpStatusCode.InternalServerError,
                Result = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(ce.ErrorObjects))
            };
            return result;
        }
    }
}
