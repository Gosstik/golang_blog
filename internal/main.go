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
	"github.com/Gosstik/golang_blog/internal/handlers/images_service"
	"github.com/Gosstik/golang_blog/internal/handlers/posts_service"
)

func main() {
	postsHandler := posts_service.NewHandler()
	imagesHandler := images_service.NewHandler()

	// Set up gRPC server.
	grpcServer := grpc.NewServer()
	api_proto.RegisterPostsServiceServer(grpcServer, postsHandler)
	api_proto.RegisterImagesServiceServer(grpcServer, imagesHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Set up gRPC-gateway.
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}

	gwmux := runtime.NewServeMux(
		// Forward X-User-Uuid HTTP header to gRPC metadata.
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

	// Serve Swagger JSON.
	mux.HandleFunc("/swagger-ui/api.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/json/api.swagger.json")
	})

	// Serve Swagger UI static files.
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("contrib/swagger-ui"))))

	// All other routes go to gRPC-gateway.
	mux.Handle("/", gwmux)

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("HTTP gateway + Swagger UI started on :8090")
	log.Println("Swagger UI available at http://localhost:8090/swagger-ui/")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
