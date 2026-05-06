package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yorukot/netstamp/internal/logger"
)

func UnaryRecovery(log *zap.Logger) grpc.UnaryServerInterceptor {
	if log == nil {
		log = zap.NewNop()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				fields := []zap.Field{
					zap.String("grpc.full_method", info.FullMethod),
					zap.Any("panic", recovered),
					zap.Stack("stacktrace"),
				}
				fields = append(fields, logger.TraceFields(ctx)...)

				log.Error("grpc_panic_recovered", fields...)
				resp = nil
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
