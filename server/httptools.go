package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
)

const httpInternalTimestamp = "___HTTP_INTERNAL_TIMESTAMP___"

// HTTPHandleFunc is a type for a HTTP Handler Function
type HTTPHandleFunc = func(w http.ResponseWriter, r *http.Request)

// HTTPHandleFuncWithLog is a type for a HTTP Handler Function with an slog instance argument
type HTTPHandleFuncWithLog = func(log slog.Instance, w http.ResponseWriter, r *http.Request)

// WriteJSON returns a JSON Object to the specified http.ResponseWriter
func WriteJSON(data interface{}, statusCode int, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	var b []byte
	var err error

	if qErr, ok := data.(*QuantoError.ErrorObject); ok && !QuantoError.ShowStackTrace() {
		qErr.StackTrace = ""
		b, err = json.Marshal(qErr)
	} else {
		b, err = json.Marshal(data)
	}

	if err != nil {
		logI.WithFields(map[string]interface{}{
			"data": data,
		}).Error("Error serializing object: %s", err)
		w.Header().Set("Content-Type", models.MimeText)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error"))
		LogExit(logI, r, 500, len("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(statusCode)
	n, _ := w.Write(b)
	LogExit(logI, r, statusCode, n)
}

// UnmarshalBodyOrDie tries to unmarshal the request body into the specified interface and returns InvalidFieldData to the client if something is wrong
func UnmarshalBodyOrDie(outData interface{}, w http.ResponseWriter, r *http.Request, logI slog.Instance) bool {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logI.Error(err)
		WriteJSON(QuantoError.New(QuantoError.InternalServerError, "body", err.Error(), nil), 500, w, r, logI)
		return false
	}

	body = []byte(strings.Replace(string(body), "\t", "", -1))

	err = json.Unmarshal(body, outData)

	if err != nil {
		logI.Error(err)
		WriteJSON(QuantoError.New(QuantoError.InvalidFieldData, "body", err.Error(), nil), 500, w, r, logI)
		return false
	}

	return true
}

// InvalidFieldData helper method to return an invalid field data error to http client
func InvalidFieldData(field string, message string, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InvalidFieldData, field, message, nil), 400, w, r, logI)
}

// PermissionDenied helper method to return an permission denied error to http client
func PermissionDenied(field string, message string, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.PermissionDenied, field, message, nil), 400, w, r, logI)
}

// NotFound helper method to return an not found error to http client
func NotFound(field string, message string, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotFound, field, message, nil), 400, w, r, logI)
}

// NotImplemented helper method to return an not implemented error to http client
func NotImplemented(w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotImplemented, "server", "This call is not implemented", nil), 400, w, r, logI)
}

// CatchAllError helper method to return an internal server error error to http client in case of non expected errors
func CatchAllError(data interface{}, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InternalServerError, "server", "There was an internal server error.", data), 500, w, r, logI)
}

// CatchAllRouter helper method to return an not found error error to http client in case of non expected endpoints
func CatchAllRouter(w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotFound, "path", "The specified URL does not exists.", nil), 404, w, r, logI)
}

// InternalServerError helper method to return an internal server error to http client
func InternalServerError(message string, data interface{}, w http.ResponseWriter, r *http.Request, logI slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InternalServerError, "server", message, data), 500, w, r, logI)
}

// LogExit does the logging of the exit call inside a HTTP Request
func LogExit(slog slog.Instance, r *http.Request, statusCode int, bodyLength int) {
	hts := r.Header.Get(httpInternalTimestamp)
	ts := float64(0)

	if hts != "" {
		v, err := strconv.ParseInt(hts, 10, 64)
		if err == nil {
			ts = time.Since(time.Unix(0, v)).Seconds() * 1000
		}
	}

	statusCodeStr := fmt.Sprintf("[%d]", statusCode)

	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ts != 0 {
		slog.Done("%s (%.2f ms) {%d bytes} %s %s from %s", statusCodeStr, ts, bodyLength, r.Method, r.URL.Path, host)
	} else {
		slog.Done("%s {%d bytes} %s %s from %s", statusCodeStr, bodyLength, r.Method, r.URL.Path, host)
	}
}

// InitHTTPTimer initializes the HTTP Request timer and prints a log line representing a received HTTP Request
func InitHTTPTimer(log slog.Instance, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	t := time.Now().UnixNano()
	r.Header.Set(httpInternalTimestamp, fmt.Sprintf("%d", t))
	log.Await("%s %s from %s", r.Method, r.URL.Path, host)
}

func wrapWithLog(log slog.Instance, f HTTPHandleFuncWithLog) HTTPHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(log, w, r)
	}
}

func wrapContextWithRequestID(r *http.Request) context.Context {
	var requestID string

	id, ok := r.Header[remote_signer.RequestIDHeader]
	if ok && len(id) >= 1 {
		// Tag the log
		requestID = id[0]
	} else {
		requestID = remote_signer.DefaultTag
	}

	return context.WithValue(r.Context(), remote_signer.CtxRequestID, requestID)
}

func wrapLogWithRequestID(log slog.Instance, r *http.Request) slog.Instance {
	id, ok := r.Header[remote_signer.RequestIDHeader]
	if ok && len(id) >= 1 {
		// Tag the log
		return log.Tag(id[0])
	}

	return log.Tag(remote_signer.DefaultTag)
}
