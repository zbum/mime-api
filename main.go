package main

import (
	"log"
	"net"
	"net/http"
	"os"

	pb "mime_api/proto"

	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Start gRPC server in goroutine
	go startGRPCServer(logger)

	// Start HTTP server
	mantyMux := mux.NewMantyMux()
	mantyMux.AddMiddleware(middleware.AccessLogger(logger))
	mantyMux.HandleFunc("POST /v1/display-part", DisplayPart)
	mantyMux.HandleFunc("GET /health", Health)

	logger.Println("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mantyMux))
}

func startGRPCServer(logger *log.Logger) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDisplayPartServiceServer(grpcServer, NewDisplayPartServer())
	reflection.Register(grpcServer)

	logger.Println("gRPC server listening on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve gRPC: %v", err)
	}
}
