package httpserver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	httpserver "github.com/yorukot/netstamp/internal/transport/http"
)

func TestHelloEndpoint(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/hello", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body struct {
		Message string `json:"message"`
		Service string `json:"service"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Message != "Hello from Netstamp" {
		t.Fatalf("unexpected message %q", body.Message)
	}
	if body.Service != "netstamp-api" {
		t.Fatalf("unexpected service %q", body.Service)
	}
}

func TestOpenAPISpecIncludesHelloEndpoint(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body struct {
		OpenAPI string                     `json:"openapi"`
		Info    struct{ Title string }     `json:"info"`
		Paths   map[string]json.RawMessage `json:"paths"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode openapi response: %v", err)
	}

	if body.OpenAPI == "" {
		t.Fatal("expected openapi version to be set")
	}
	if body.Info.Title != "Netstamp API" {
		t.Fatalf("unexpected title %q", body.Info.Title)
	}
	if _, ok := body.Paths["/v1/hello"]; !ok {
		t.Fatal("expected /v1/hello to be documented")
	}
}

func TestDocsEndpoint(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected html docs response, got content type %q", contentType)
	}
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	return httpserver.NewRouter(httpserver.Dependencies{
		Log:          zap.NewNop(),
		APIVersion:   "test",
		HelloService: apphello.NewService("netstamp-api"),
		ReadinessCheck: func(context.Context) error {
			return nil
		},
		RequestTimeout: time.Second,
	})
}
