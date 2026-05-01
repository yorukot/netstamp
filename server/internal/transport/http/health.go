package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/yorukot/netstamp/internal/transport/http/respond"
)

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	respond.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func readinessHandler(check func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check == nil {
			respond.JSON(w, http.StatusOK, map[string]string{
				"status": "ready",
			})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := check(ctx); err != nil {
			respond.Error(w, r, http.StatusServiceUnavailable, "not_ready", "readiness check failed")
			return
		}

		respond.JSON(w, http.StatusOK, map[string]string{
			"status": "ready",
		})
	}
}
