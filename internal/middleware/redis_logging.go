package middleware

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisLoggingHook implements redis.Hook
type RedisLoggingHook struct {
	logger *zap.Logger
}

func NewRedisLoggingHook(logger *zap.Logger) *RedisLoggingHook {
	return &RedisLoggingHook{logger: logger}
}

func (h *RedisLoggingHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		h.logger.Debug("redis dial",
			zap.String("network", network),
			zap.String("addr", addr),
		)
		conn, err := next(ctx, network, addr)
		if err != nil {
			h.logger.Error("redis dial failed",
				zap.String("addr", addr),
				zap.Error(err),
			)
		}
		return conn, err
	}
}

func (h *RedisLoggingHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		h.logger.Debug("redis command started",
			zap.String("cmd", cmd.FullName()),
			zap.String("detail", cmd.String()),
		)

		err := next(ctx, cmd)

		fields := []zap.Field{
			zap.String("cmd", cmd.FullName()),
			zap.Duration("duration", time.Since(start)),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			h.logger.Error("redis command failed", fields...)
		} else {
			h.logger.Debug("redis command completed", fields...)
		}

		return err
	}
}

func (h *RedisLoggingHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		cmdNames := make([]string, len(cmds))
		for i, cmd := range cmds {
			cmdNames[i] = cmd.FullName()
		}

		h.logger.Debug("redis pipeline started",
			zap.Int("commands_count", len(cmds)),
			zap.Strings("commands", cmdNames),
		)

		err := next(ctx, cmds)

		fields := []zap.Field{
			zap.Int("commands_count", len(cmds)),
			zap.Duration("duration", time.Since(start)),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			h.logger.Error("redis pipeline failed", fields...)
		} else {
			h.logger.Debug("redis pipeline completed", fields...)
		}

		return err
	}
}
