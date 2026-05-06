package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/yorukot/netstamp/internal/logger"
)

func UnaryLogging(root *zap.Logger) grpc.UnaryServerInterceptor {
	if root == nil {
		root = zap.NewNop()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		requestID := requestIDFromMetadata(ctx)

		reqLog := root.With(
			zap.String("request_id", requestID),
			zap.String("grpc.full_method", info.FullMethod),
		)
		if traceFields := logger.TraceFields(ctx); len(traceFields) > 0 {
			reqLog = reqLog.With(traceFields...)
		}

		ctx = logger.WithContext(ctx, reqLog)
		resp, err := handler(ctx, req)

		code := status.Code(err)
		fields := []zap.Field{
			zap.String("grpc.code", code.String()),
			zap.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		switch {
		case code == codes.OK:
			reqLog.Info("grpc_request", fields...)
		case code == codes.InvalidArgument ||
			code == codes.NotFound ||
			code == codes.Unauthenticated ||
			code == codes.PermissionDenied:
			reqLog.Warn("grpc_request", fields...)
		default:
			reqLog.Error("grpc_request", fields...)
		}

		return resp, err
	}
}

func requestIDFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	for _, key := range []string{"x-request-id", "request-id"} {
		values := md.Get(key)
		if len(values) > 0 {
			return values[0]
		}
	}

	return ""
}
