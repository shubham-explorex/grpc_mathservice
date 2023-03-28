package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	math1_v1 "example.com/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var err error

func main() {
	fmt.Println("Go gRPC Beginners Tutorial!")

	// GRPC Server
	s := math1_v1.Server{}
	grpcServer := grpc.NewServer()
	math1_v1.RegisterMathServiceServer(grpcServer, &s)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Fatalln(grpcServer.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()

	// Rest Server
	mux := runtime.NewServeMux()

	// Register handler
	err = math1_v1.RegisterMathServiceHandler(context.Background(), mux, conn)

	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":9050",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
