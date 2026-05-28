package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetJsonReturnsHTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"Message":"invalid payload"}`))
	}))
	defer server.Close()

	var target map[string]interface{}
	err := GetJson(server.URL, &target)
	if err == nil {
		t.Fatal("expected http error to be returned")
	}

	if !strings.Contains(err.Error(), "http 400") {
		t.Fatalf("expected status code in error, got %v", err)
	}
}

func TestSendJsonDecodesSuccessResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var target map[string]interface{}
	if err := SendJson(server.URL, http.MethodPost, &target, map[string]string{"demo": "value"}); err != nil {
		t.Fatalf("expected success response, got %v", err)
	}

	if ok, _ := target["ok"].(bool); !ok {
		t.Fatalf("expected decoded success body, got %+v", target)
	}
}

func TestSendJsonReturnsHTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"validation failed"}`))
	}))
	defer server.Close()

	var target map[string]interface{}
	err := SendJson(server.URL, http.MethodPost, &target, map[string]string{"demo": "value"})
	if err == nil {
		t.Fatal("expected http error to be returned")
	}

	if !strings.Contains(err.Error(), "http 422") {
		t.Fatalf("expected status code in error, got %v", err)
	}
}
