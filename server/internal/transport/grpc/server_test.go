package grpcserver_test

import (
	"context"
	"net"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	grpcserver "github.com/yorukot/netstamp/internal/transport/grpc"
)

func TestHealthService(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpcserver.New(grpcserver.Dependencies{
		Log:         zap.NewNop(),
		ServiceName: "netstamp-api",
	})

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("serve grpc: %v", err)
		}
	}()
	defer server.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial grpc server: %v", err)
	}
	defer conn.Close()

	client := healthv1.NewHealthClient(conn)
	resp, err := client.Check(ctx, &healthv1.HealthCheckRequest{
		Service: "netstamp-api",
	})
	if err != nil {
		t.Fatalf("check health: %v", err)
	}

	if resp.GetStatus() != healthv1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING, got %s", resp.GetStatus())
	}
}
