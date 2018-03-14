using System;
using System.Text;
using RemoteSigner.Exceptions.HKP;
using RemoteSigner.Models;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.HttpData.ExceptionHandlers {
    [RestExceptionHandler(typeof(HKPBaseException))]
    public class HKPExceptionHandler: IRestExceptionHandler {
        public RestResult HandleException(Exception e) {

            string exceptionType = e.GetType().FullName;

            if (exceptionType == typeof(OperationNotImplemented).FullName) {
                return new RestResult() {
                    ContentType = MimeTypes.Text,
                    StatusCode = System.Net.HttpStatusCode.NotImplemented,
                    Result = Encoding.UTF8.GetBytes(e.Message)
                };
            }

            if (exceptionType == typeof(UnknownOperationException).FullName) {
                return new RestResult() {
                    ContentType = MimeTypes.Text,
                    StatusCode = System.Net.HttpStatusCode.NotFound,
                    Result = Encoding.UTF8.GetBytes(e.Message)
                };
            }

            return new RestResult() {
                ContentType = MimeTypes.Text,
                StatusCode = System.Net.HttpStatusCode.InternalServerError,
                Result = Encoding.UTF8.GetBytes($"Internal Server Error - {e.Message}")
            };
        }
    }
}
