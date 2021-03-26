package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/quan-to/slog"
)

// skipEndpoints represents the endpoints that must be skipped in LoggingMiddleware
var skipEndpoints = map[string]bool{
	"/test/ping": true,
}

// ResponseWriter is a http.ResponseWriter wrapper that provides the status code and content length info.
type ResponseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
}

// WriteHeader implements the http.ResponseWriter.WriteHeader function. It makes enable to store the response status code.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write implements the http.ResponseWriter.Write function. It makes enable to store the response content length.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.contentLength = len(b)
	return rw.ResponseWriter.Write(b)
}

func wrapResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

// LoggingMiddleware is a HTTP middleware that logs the entry and exit requests
func LoggingMiddleware(next http.Handler) http.Handler {
	log := slog.Scope("LoggingMiddleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if skipEndpoints[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()
		host, _, _ := net.SplitHostPort(r.RemoteAddr)

		logParams := map[string]interface{}{
			"endpoint": r.URL.Path,
			"method":   r.Method,
			"host":     host,
		}

		log = wrapLogWithRequestID(log, r).WithFields(logParams)
		log.Await("incoming request %s %s", r.Method, r.URL.Path)

		rw := wrapResponseWriter(w)
		next.ServeHTTP(rw, r)

		duration := time.Since(startTime)

		logParams["contentLength"] = rw.contentLength
		logParams["statusCode"] = fmt.Sprint(rw.status)
		logParams["elapsedTime"] = duration.Milliseconds()

		log.Done("Finished request %s - %s - %d", r.Method, r.URL.Path, rw.status)
	})
}
