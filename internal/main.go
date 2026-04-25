package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/clients"
	"github.com/Gosstik/golang_blog/internal/config"
	"github.com/Gosstik/golang_blog/internal/handlers/images_service"
	"github.com/Gosstik/golang_blog/internal/handlers/posts_service"
	pgRepo "github.com/Gosstik/golang_blog/internal/repositories/postgres"
	redisRepo "github.com/Gosstik/golang_blog/internal/repositories/redis"
)

func main() {
	cfg := config.Load()

	// Init clients.
	db := clients.NewPostgresDB(cfg.PostgresDSN)
	rdb := clients.NewRedisClient(cfg.RedisAddr)

	// Init repositories.
	postRepo := pgRepo.NewPostRepository(db)
	userRepo := pgRepo.NewUserRepository(db)
	likesRepo := redisRepo.NewLikesRepository(rdb)
	cacheRepo := redisRepo.NewCacheRepository(rdb)

	// Init handlers.
	postsHandler := posts_service.NewHandler(postRepo, userRepo, likesRepo, cacheRepo)
	imagesHandler := images_service.NewHandler()

	// Set up gRPC server.
	grpcServer := grpc.NewServer()
	api_proto.RegisterPostsServiceServer(grpcServer, postsHandler)
	api_proto.RegisterImagesServiceServer(grpcServer, imagesHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("gRPC server started on %s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Set up gRPC-gateway.
	conn, err := grpc.NewClient(cfg.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
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
		log.Fatalf("failed to register PostsService gateway: %v", err)
	}
	if err := api_proto.RegisterImagesServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalf("failed to register ImagesService gateway: %v", err)
	}

	// Set up HTTP mux with swagger-ui and gRPC-gateway.
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger-ui/api.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/json/api.swagger.json")
	})

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("contrib/swagger-ui"))))

	mux.Handle("/", gwmux)

	gwServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}

	log.Printf("HTTP gateway + Swagger UI started on %s\n", cfg.HTTPPort)
	log.Printf("Swagger UI available at http://localhost%s/swagger-ui/\n", cfg.HTTPPort)
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
