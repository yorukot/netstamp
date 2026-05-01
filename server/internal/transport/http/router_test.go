package httpserver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	httpserver "github.com/yorukot/netstamp/internal/transport/http"
)

func TestHelloEndpoint(t *testing.T) {
	router := httpserver.NewRouter(httpserver.Dependencies{
		Log:          zap.NewNop(),
		HelloService: apphello.NewService("netstamp-api"),
		ReadinessCheck: func(context.Context) error {
			return nil
		},
		RequestTimeout: time.Second,
	})

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

	if body.Message != "Hello from NetStamp" {
		t.Fatalf("unexpected message %q", body.Message)
	}
	if body.Service != "netstamp-api" {
		t.Fatalf("unexpected service %q", body.Service)
	}
}
