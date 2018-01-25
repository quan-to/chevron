using Newtonsoft.Json;

namespace RemoteSigner.Models {
    class ErrorObjectQ {
        public string errorCode { get; set; }
        public string errorField { get; set; }
        public string message { get; set; }
        public string errorData { get; set; }
    }

    public class ErrorObject {

        public string ErrorCode { get; set; }
        public string ErrorField { get; set; }
        public string Message { get; set; }
        public object ErrorData { get; set; }

        public string ToJSON() {
            var e = new ErrorObjectQ {
                errorCode = ErrorCode,
                errorField = ErrorField,
                message = Message,
                errorData = JsonConvert.SerializeObject(ErrorData),
            };
            return JsonConvert.SerializeObject(e);
        }
    }

}
