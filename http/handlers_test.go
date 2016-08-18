package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Expected HTTP Code 200 for health check")
	}

	body := rr.Body.String()
	if `{"status": "online"}` != body {
		t.Error("Health check body was invalid.")
	}
}

func TestParseEmptyMessage(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/message", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ParseMessageHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Error("Empty requests should result in bad request.")
	}
}

func TestParseMessageBadMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/message", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ParseMessageHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Error("GET requests should result in not allowed.")
	}
}

func TestParseUnsupportedContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/message", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ParseMessageHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Error("GET requests should result in not allowed.")
	}
}

func TestParseValidMessage(t *testing.T) {
	msg := "@chris @matt http://www.google.com (coffee)(hat)"
	req, err := http.NewRequest("POST", "/v1/message", bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ParseMessageHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Expected successul request. Code: ", rr.Code)
	}

	expectedMessage := `{"mentions":["chris","matt"],"emoticons":["coffee","hat"],"links":[{"url":"http://www.google.com","title":"Google"}]}`

	if expectedMessage != rr.Body.String() {
		t.Error(fmt.Sprintf("Expected message %s, got %s", expectedMessage, rr.Body.String()))
	}
}
