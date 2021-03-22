package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/quan-to/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestLoggingMiddleware stores the middleware logs in a buffer and validates their contents
func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	slog.UnsetTestMode()
	slog.SetDefaultOutput(&logBuffer)
	slog.SetLogFormat(slog.JSON)

	expectedScope := "LoggingMiddleware"
	expectedHTTPMethod := http.MethodPost
	someURL, _ := url.Parse("http://huehuebr.com/resource/subresource")
	expectedCode := 200
	waitMs := float64(150)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(waitMs) * time.Millisecond)
		w.WriteHeader(expectedCode)
	})

	req := httptest.NewRequest(http.MethodPost, someURL.String(), nil)

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

		currentURL := logMap["url"]
		expectURL := someURL.Path
		if currentURL != expectURL {
			t.Fatalf("[url] Got %s; want %s", currentURL, expectURL)
		}

		currentMethod := logMap["method"]
		if currentMethod != expectedHTTPMethod {
			t.Fatalf("[http method] Got %s; want %s", currentMethod, expectedHTTPMethod)
		}

		// DONE op means a finished request
		if logMap["op"] != "DONE" {
			continue
		}

		currentStatus := logMap["statusCode"]
		if fmt.Sprint(currentStatus) != fmt.Sprint(expectedCode) {
			t.Fatalf("[status code] Got %s; want %v", currentStatus, expectedCode)
		}

		currentResponseTime := logMap["responseTime"].(float64)
		if currentResponseTime < waitMs {
			t.Fatalf("[response time] Got %f; want greater than %f", currentResponseTime, waitMs)
		}

	}

	slog.SetTestMode()
	slog.SetLogFormat(slog.PIPE)
}
