package middleware

import (
	"encoding/json"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func WriteProblem(w http.ResponseWriter, r *http.Request, status int, detail string) {
	requestID := chimw.GetReqID(r.Context())
	if requestID != "" {
		w.Header().Set("X-Request-ID", requestID)
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(&huma.ErrorModel{
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	})
}
