package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	hellohttp "github.com/yorukot/netstamp/internal/transport/http/hello"
	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
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

	r.Route(dep.basePath(), func(apiRouter chi.Router) {
		api := humachi.New(apiRouter, newHumaConfig(dep))
		registerSystemRoutes(api, dep.ReadinessCheck)
		hellohttp.NewHandler(dep.HelloService).RegisterRoutes(api)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return r
}

func newHumaConfig(dep Dependencies) huma.Config {
	config := huma.DefaultConfig("Netstamp API", dep.APIVersion)
	config.Info.Description = "Controller HTTP API for Netstamp."
	config.Servers = []*huma.Server{{URL: dep.basePath()}}
	return config
}

func (d *Dependencies) basePath() string {
	return "/api/" + d.APIVersion
}
