package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/slog"
)

// TestLoggingMiddleware stores the middleware logs in a buffer and validates their contents
func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	slog.UnsetTestMode()
	slog.SetDefaultOutput(&logBuffer)
	slog.SetLogFormat(slog.JSON)
	config.RequestIDHeader = "X-Request-ID"

	expectedScope := "LoggingMiddleware"
	expectedHTTPMethod := http.MethodPost
	someURL, _ := url.Parse("http://huehuebr.com/resource/subresource")
	someResponse := []byte("huehuebrbr")
	expectedCode := 200
	expectedRequestID := "01101000-01110101-01100101"
	waitMs := float64(150)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(waitMs) * time.Millisecond)
		w.WriteHeader(expectedCode)
		_, _ = w.Write(someResponse)
	})

	req := httptest.NewRequest(http.MethodPost, someURL.String(), nil)
	req.Header.Set(config.RequestIDHeader, expectedRequestID)

	LoggingMiddleware(mockHandler).ServeHTTP(httptest.NewRecorder(), req)

	if logBuffer.Len() == 0 {
		t.Fatalf("No logs found in middleware")
	}

	// each line is a log (1: incoming request. 2: finished request)
	scanner := bufio.NewScanner(&logBuffer)
	for scanner.Scan() {

		var logMap map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &logMap); err != nil {
			t.Fatalf("error on unmarshal log result: %s. captured logs: %s", err, logBuffer.String())
		}

		currentScope := logMap["scope"]
		if currentScope != expectedScope {
			t.Fatalf("[scope] Got %s; want %s", currentScope, expectedScope)
		}

		currentEndpoint := logMap["endpoint"]
		expectedEndpoint := someURL.Path
		if currentEndpoint != expectedEndpoint {
			t.Fatalf("[endpoint] Got %s; want %s", currentEndpoint, expectedEndpoint)
		}

		currentMethod := logMap["method"]
		if currentMethod != expectedHTTPMethod {
			t.Fatalf("[http method] Got %s; want %s", currentMethod, expectedHTTPMethod)
		}

		// DONE op means a finished request
		if logMap["op"] != "DONE" {
			continue
		}

		currentContentLength := logMap["contentLength"]
		expectedContentLength := len(someResponse)
		if fmt.Sprint(currentContentLength) != fmt.Sprint(expectedContentLength) {
			t.Fatalf("[content length] Got %s; want %v", currentContentLength, expectedContentLength)
		}

		currentStatus := logMap["statusCode"]
		if fmt.Sprint(currentStatus) != fmt.Sprint(expectedCode) {
			t.Fatalf("[status code] Got %s; want %v", currentStatus, expectedCode)
		}

		currentResponseTime := logMap["elapsedTime"].(float64)
		if currentResponseTime < waitMs {
			t.Fatalf("[response time] Got %f; want greater than %f", currentResponseTime, waitMs)
		}

		currentRequestID := logMap["tag"]
		if fmt.Sprint(currentRequestID) != fmt.Sprint(expectedRequestID) {
			t.Fatalf("[request ID] Got %s; want %v", currentRequestID, expectedRequestID)
		}

	}

	slog.SetTestMode()
	slog.SetLogFormat(slog.PIPE)
}
