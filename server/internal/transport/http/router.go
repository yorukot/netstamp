package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	appteam "github.com/yorukot/netstamp/internal/application/team"
	"github.com/yorukot/netstamp/internal/observability/httptrace"
	authhttp "github.com/yorukot/netstamp/internal/transport/http/auth"
	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
	teamhttp "github.com/yorukot/netstamp/internal/transport/http/team"
)

type Dependencies struct {
	Log            *zap.Logger
	APIVersion     string
	AuthService    *appauth.Service
	AuthVerifier   appauth.TokenVerifier
	TeamService    *appteam.Service
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
	r.Use(otelhttp.NewMiddleware("http.server",
		otelhttp.WithSpanNameFormatter(httptrace.RequestSpanName),
	))
	r.Use(httpmiddleware.ZapRecoverer(dep.Log))
	r.Use(chimw.Timeout(dep.RequestTimeout))
	r.Use(httpmiddleware.ZapRequestLogger(dep.Log))

	r.Route(dep.basePath(), func(apiRouter chi.Router) {
		api := humachi.New(apiRouter, newHumaConfig(dep))
		registerSystemRoutes(api, dep.ReadinessCheck)

		// Auth handler
		if dep.AuthService != nil {
			authhttp.NewHandler(dep.AuthService, dep.AuthVerifier).RegisterRoutes(api)
		}
		// Team handler
		if dep.TeamService != nil {
			teamhttp.NewHandler(dep.TeamService, dep.AuthVerifier).RegisterRoutes(api)
		}
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
	if config.Components.SecuritySchemes == nil {
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	config.Components.SecuritySchemes["bearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
	return config
}

func (d *Dependencies) basePath() string {
	return "/api/" + d.APIVersion
}
