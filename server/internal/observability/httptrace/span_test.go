package httptrace

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestSpanName(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register?debug=true", nil)

	got := RequestSpanName("http.server", req)
	want := "POST /api/v1/auth/register"
	if got != want {
		t.Fatalf("expected span name %q, got %q", want, got)
	}
}
