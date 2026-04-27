package middleware

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GRPCMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewGRPCMetrics(reg prometheus.Registerer) *GRPCMetrics {
	m := &GRPCMetrics{
		// Reason to use Counter: RPS counting
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_server_requests_total",
				Help: "Total number of gRPC requests received, partitioned by method and status code.",
			},
			[]string{"method", "code"},
		),
		// Reason to use Histogram: computing percentiles (p50, p90, p99)
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_server_request_duration_seconds",
				Help:    "Histogram of gRPC request durations in seconds, partitioned by method.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
	}

	reg.MustRegister(m.requestsTotal, m.requestDuration)
	return m
}

func (m *GRPCMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		m.requestsTotal.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
		m.requestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

func (m *GRPCMetrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		st, _ := status.FromError(err)
		m.requestsTotal.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
		m.requestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return err
	}
}
