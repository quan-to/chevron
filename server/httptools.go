package server

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/quan-to/chevron/QuantoError"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const httpInternalTimestamp = "___HTTP_INTERNAL_TIMESTAMP___"

var httpToolsLog = slog.Scope("HTTP Tools")

func WriteJSON(data interface{}, statusCode int, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {

	if logI == nil {
		logI = httpToolsLog
	}

	b, err := json.Marshal(data)

	if err != nil {
		httpToolsLog.Error("Error serializing object: %s", err)
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

func UnmarshalBodyOrDie(outData interface{}, w http.ResponseWriter, r *http.Request, logI *slog.Instance) bool {
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

func InvalidFieldData(field string, message string, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InvalidFieldData, field, message, nil), 400, w, r, logI)
}

func PermissionDenied(field string, message string, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.PermissionDenied, field, message, nil), 400, w, r, logI)
}

func NotFound(field string, message string, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotFound, field, message, nil), 400, w, r, logI)
}

func NotImplemented(w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotImplemented, "server", "This call is not implemented", nil), 400, w, r, logI)
}

func CatchAllError(data interface{}, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InternalServerError, "server", "There was an internal server error.", data), 500, w, r, logI)
}

func CatchAllRouter(w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.NotFound, "path", "The specified URL does not exists.", nil), 404, w, r, logI)
}

func InternalServerError(message string, data interface{}, w http.ResponseWriter, r *http.Request, logI *slog.Instance) {
	WriteJSON(QuantoError.New(QuantoError.InternalServerError, "server", message, data), 500, w, r, logI)
}

func LogExit(slog *slog.Instance, r *http.Request, statusCode int, bodyLength int) {
	method := aurora.Bold(r.Method).Cyan()
	hts := r.Header.Get(httpInternalTimestamp)
	ts := float64(0)

	if hts != "" {
		v, err := strconv.ParseInt(hts, 10, 64)
		if err == nil {
			ts = time.Since(time.Unix(0, v)).Seconds() * 1000
		}
	}

	statusCodeStr := aurora.Black(fmt.Sprintf("[%d]", statusCode))
	switch statusCode {
	case 400:
		statusCodeStr = aurora.Red(statusCodeStr).Inverse().Bold()
	case 404:
		statusCodeStr = aurora.Red(statusCodeStr).Inverse().Bold()
	case 500:
		statusCodeStr = aurora.Red(statusCodeStr).Inverse().Bold()
	case 200:
		statusCodeStr = aurora.Green(statusCodeStr).Inverse().Bold()
	default:
		statusCodeStr = aurora.Gray(statusCodeStr).Bold()
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	remote := aurora.Gray(host)

	if ts != 0 {
		slog.LogNoFormat("%s (%8.2f ms) {%8d bytes} %-4s %s from %s", statusCodeStr, ts, bodyLength, method, aurora.Gray(r.URL.Path), remote)
	} else {
		slog.LogNoFormat("%s {%8d bytes}          %-4s %s from %s", statusCodeStr, bodyLength, method, aurora.Gray(r.URL.Path), remote)
	}
}

func InitHTTPTimer(r *http.Request) {
	t := time.Now().UnixNano()
	r.Header.Set(httpInternalTimestamp, fmt.Sprintf("%d", t))
}
