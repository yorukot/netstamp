package grpcserver

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/yorukot/netstamp/internal/transport/grpc/interceptors"
)

type Dependencies struct {
	Log         *zap.Logger
	ServiceName string
}

func New(dep Dependencies) *grpc.Server {
	if dep.Log == nil {
		dep.Log = zap.NewNop()
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryLogging(dep.Log),
			interceptors.UnaryRecovery(dep.Log),
		),
	)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	if dep.ServiceName != "" {
		healthServer.SetServingStatus(dep.ServiceName, healthv1.HealthCheckResponse_SERVING)
	}
	healthv1.RegisterHealthServer(server, healthServer)

	return server
}
