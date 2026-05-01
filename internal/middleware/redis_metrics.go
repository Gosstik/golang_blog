package middleware

import (
	"context"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// RedisMetricsHook implements redis.Hook
type RedisMetricsHook struct {
	operationDuration *prometheus.HistogramVec
}

func NewRedisMetricsHook(reg prometheus.Registerer) *RedisMetricsHook {
	h := &RedisMetricsHook{
		// Reason to use Histogram:
		// We need to analyze the *distribution* of Redis operation latencies.
		// It allows computing percentiles (p50, p90, p99) via PromQL.
		operationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "redis_likes_operation_duration_seconds",
				Help: "Histogram of Redis likes operation durations in seconds, " +
					"partitioned by command name. Used to analyze time needed for " +
					"retrieving/modifying likes from Redis.",
				// Buckets tuned for Redis: 0.1ms, 0.25ms, 0.5ms, 1ms, 2.5ms,
				// 5ms, 10ms, 25ms, 50ms, 100ms.
				Buckets: []float64{
					0.0001, 0.00025, 0.0005, 0.001, 0.0025,
					0.005, 0.01, 0.025, 0.05, 0.1,
				},
			},
			[]string{"command"},
		),
	}

	reg.MustRegister(h.operationDuration)
	return h
}

func (h *RedisMetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h *RedisMetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd)

		h.operationDuration.WithLabelValues(cmd.FullName()).Observe(time.Since(start).Seconds())

		return err
	}
}

func (h *RedisMetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)

		h.operationDuration.WithLabelValues("pipeline").Observe(time.Since(start).Seconds())

		return err
	}
}
