package server

import (
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/slog"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(test *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	contentTypeExpected := "application/json"
	log := slog.Scope("MAIN")

	testCases := []struct {
		description          string
		call				 func(w http.ResponseWriter)
		statusCodeExpected   int
		responseBodyExpected string
	}{
		{
			"Test Internal Server Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				InternalServerError("server is down!","FAILURE", w, req, log)
			},
			500,
			"{\"errorCode\":\"INTERNAL_SERVER_ERROR\",\"errorField\":\"server\",\"message\":\"server is down!\",\"errorData\":\"FAILURE\",\"stackTrace\":\"\"}",
		},
		{
			"Test Route Not Found Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				CatchAllRouter(w, req, slog.Scope("MAIN"))
			},
			404,
			"{\"errorCode\":\"NOT_FOUND\",\"errorField\":\"path\",\"message\":\"The specified URL does not exists.\",\"errorData\":null,\"stackTrace\":\"\"}",
		},
		{
			"Test Unexpected Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				CatchAllError("Unexpected Error", w, req, log)
			},
			500,
			"{\"errorCode\":\"INTERNAL_SERVER_ERROR\",\"errorField\":\"server\",\"message\":\"There was an internal server error.\",\"errorData\":\"Unexpected Error\",\"stackTrace\":\"\"}",
		},
		{
			"Test Not Implemented Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				NotImplemented(w, req, log)
			},
			400,
			"{\"errorCode\":\"NOT_IMPLEMENTED\",\"errorField\":\"server\",\"message\":\"This call is not implemented\",\"errorData\":null,\"stackTrace\":\"\"}",
		},
		{
			"Test Not Found Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				NotFound("user", "User not found!", w, req, log)
			},
			400,
			"{\"errorCode\":\"NOT_FOUND\",\"errorField\":\"user\",\"message\":\"User not found!\",\"errorData\":null,\"stackTrace\":\"\"}",
		},
		{
			"Test Permission Denied Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				PermissionDenied("proxyToken", "Please check if your proxyToken is valid", w, req, log)
			},
			400,
			"{\"errorCode\":\"PERMISSION_DENIED\",\"errorField\":\"proxyToken\",\"message\":\"Please check if your proxyToken is valid\",\"errorData\":null,\"stackTrace\":\"\"}",
		},
		{
			"Test Invalid Fields Error Response",
			func(w http.ResponseWriter) {
				QuantoError.DisableStackTrace()
				InvalidFieldData("username", "Field username is invalid", w, req, log)
			},
			400,
			"{\"errorCode\":\"INVALID_FIELD_DATA\",\"errorField\":\"username\",\"message\":\"Field username is invalid\",\"errorData\":null,\"stackTrace\":\"\"}",
		},
	}

	for _, tc := range testCases {
		test.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.call(w)
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			responseBody := string(body)

			if resp.StatusCode != tc.statusCodeExpected {
				t.Fatalf("Got %d; want %d", resp.StatusCode, tc.statusCodeExpected)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != contentTypeExpected {
				t.Fatalf("Got %s; want %s", contentType, contentTypeExpected)
			}

			if responseBody != tc.responseBodyExpected {
				t.Fatalf("Got %s; want %s", body, tc.responseBodyExpected)
			}
		})
	}


}