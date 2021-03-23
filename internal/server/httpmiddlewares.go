package server

import (
	"fmt"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/slog"
	"net"
	"net/http"
	"time"
)

// ResponseWriter is a http.ResponseWriter wrapper that provides the status code info.
type ResponseWriter struct {
	http.ResponseWriter
	status int
	contentLength int
}

// WriteHeader implements the http.ResponseWriter WriteHeader function. It makes enable to store the response status code.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

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

		startTime := time.Now()
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		requestID := r.Header.Get(config.RequestIDHeader)

		logParams := map[string]interface{}{
			"endpoint": r.URL.Path,
			"method":   r.Method,
			"host":     host,
			"requestID": requestID,
		}

		log = log.WithFields(logParams)

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
