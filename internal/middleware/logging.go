package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		logger.Info("grpc request started",
			zap.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("code", st.Code().String()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("grpc request failed", fields...)
		} else {
			logger.Info("grpc request completed", fields...)
		}

		return resp, err
	}
}

func StreamLoggingInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		logger.Info("grpc stream started",
			zap.String("method", info.FullMethod),
		)

		err := handler(srv, ss)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("code", st.Code().String()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("grpc stream failed", fields...)
		} else {
			logger.Info("grpc stream completed", fields...)
		}

		return err
	}
}
