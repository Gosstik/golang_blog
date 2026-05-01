package main

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/clients"
	"github.com/Gosstik/golang_blog/internal/config"
	"github.com/Gosstik/golang_blog/internal/handlers/images_service"
	"github.com/Gosstik/golang_blog/internal/handlers/posts_service"
	"github.com/Gosstik/golang_blog/internal/logger"
	"github.com/Gosstik/golang_blog/internal/middleware"
	pgRepo "github.com/Gosstik/golang_blog/internal/repositories/postgres"
	redisRepo "github.com/Gosstik/golang_blog/internal/repositories/redis"
)

func main() {
	cfg := config.Load()

	// Init logger.
	zapLogger := logger.New()
	defer zapLogger.Sync()

	// Init Prometheus registry and metrics collectors.
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	promRegistry.MustRegister(prometheus.NewGoCollector())

	// Init gRPC metrics.
	grpcMetrics := middleware.NewGRPCMetrics(promRegistry)

	// Init Redis metrics hook.
	redisMetricsHook := middleware.NewRedisMetricsHook(promRegistry)

	// Init clients.
	db := clients.NewPostgresDB(cfg.PostgresDSN)
	rdb := clients.NewRedisClient(cfg.RedisAddr)

	// Add Redis hooks for logging and metrics collection.
	rdb.AddHook(middleware.NewRedisLoggingHook(zapLogger.Named("redis")))
	rdb.AddHook(redisMetricsHook)

	// Init repositories.
	postRepo := pgRepo.NewPostRepository(db)
	userRepo := pgRepo.NewUserRepository(db)
	likesRepo := redisRepo.NewLikesRepository(rdb)
	cacheRepo := redisRepo.NewCacheRepository(rdb)

	// Init handlers.
	postsHandler := posts_service.NewHandler(postRepo, userRepo, likesRepo, cacheRepo)
	imagesHandler := images_service.NewHandler()

	// Set up gRPC server.
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryLoggingInterceptor(zapLogger.Named("grpc")),
			grpcMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamLoggingInterceptor(zapLogger.Named("grpc")),
			grpcMetrics.StreamServerInterceptor(),
		),
	)
	api_proto.RegisterPostsServiceServer(grpcServer, postsHandler)
	api_proto.RegisterImagesServiceServer(grpcServer, imagesHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		zapLogger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		zapLogger.Info("gRPC server started", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			zapLogger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// Set up gRPC-gateway.
	conn, err := grpc.NewClient(cfg.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zapLogger.Fatal("failed to create gRPC client", zap.Error(err))
	}

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.EqualFold(key, "X-User-Uuid") {
				return "x-user-uuid", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	if err := api_proto.RegisterPostsServiceHandler(context.Background(), gwmux, conn); err != nil {
		zapLogger.Fatal("failed to register PostsService gateway", zap.Error(err))
	}
	if err := api_proto.RegisterImagesServiceHandler(context.Background(), gwmux, conn); err != nil {
		zapLogger.Fatal("failed to register ImagesService gateway", zap.Error(err))
	}

	// Set up HTTP mux with swagger-ui and gRPC-gateway.
	mux := http.NewServeMux()

	// Admin endpoint.
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/swagger-ui/api.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/json/api.swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("contrib/swagger-ui"))))

	// Other endpoints.
	mux.Handle("/", gwmux)

	gwServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}

	zapLogger.Info("HTTP gateway started",
		zap.String("port", cfg.HTTPPort),
		zap.String("metrics", "http://localhost"+cfg.HTTPPort+"/metrics"),
		zap.String("swagger", "http://localhost"+cfg.HTTPPort+"/swagger-ui/"),
	)
	if err := gwServer.ListenAndServe(); err != nil {
		zapLogger.Fatal("failed to serve HTTP", zap.Error(err))
	}
}
