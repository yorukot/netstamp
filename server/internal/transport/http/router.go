package httpserver

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
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
	APIVersion     string
	HelloService   *apphello.Service
	ReadinessCheck func(context.Context) error
	RequestTimeout time.Duration
}

func NewRouter(dep Dependencies) http.Handler {
	if dep.Log == nil {
		dep.Log = zap.NewNop()
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(dep.RequestTimeout))
	r.Use(httpmiddleware.ZapRequestLogger(dep.Log))

	api := humachi.New(r, newHumaConfig(dep))
	registerSystemRoutes(api, dep.ReadinessCheck)
	hellohttp.NewHandler(dep.HelloService).RegisterRoutes(api)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respond.Error(w, r, http.StatusNotFound, "not_found", "route not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respond.Error(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	return r
}

func newHumaConfig(dep Dependencies) huma.Config {
	version := strings.TrimSpace(dep.APIVersion)
	if version == "" {
		version = "dev"
	}

	config := huma.DefaultConfig("Netstamp API", version)
	config.Info.Description = "Controller HTTP API for Netstamp."
	return config
}
