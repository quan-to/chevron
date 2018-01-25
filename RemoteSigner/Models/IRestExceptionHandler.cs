using System;
namespace RemoteSigner.Models {
    public interface IRestExceptionHandler {
        RestResult HandleException(Exception e);
    }
}
