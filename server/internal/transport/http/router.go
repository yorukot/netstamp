package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	hellohttp "github.com/yorukot/netstamp/internal/transport/http/hello"
	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
	"github.com/yorukot/netstamp/internal/transport/http/respond"
)

type Dependencies struct {
	Log            *zap.Logger
	HelloService   *apphello.Service
	ReadinessCheck func(context.Context) error
	RequestTimeout time.Duration
}

func NewRouter(dep Dependencies) http.Handler {
	if dep.Log == nil {
		dep.Log = zap.NewNop()
	}
	if dep.RequestTimeout == 0 {
		dep.RequestTimeout = 10 * time.Second
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(dep.RequestTimeout))
	r.Use(httpmiddleware.ZapRequestLogger(dep.Log))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		respond.JSON(w, http.StatusOK, map[string]string{
			"message": "NetStamp API is running",
		})
	})
	r.Get("/livez", livenessHandler)
	r.Get("/readyz", readinessHandler(dep.ReadinessCheck))

	r.Route("/v1", func(r chi.Router) {
		hellohttp.NewHandler(dep.HelloService).RegisterRoutes(r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respond.Error(w, r, http.StatusNotFound, "not_found", "route not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respond.Error(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	return r
}
