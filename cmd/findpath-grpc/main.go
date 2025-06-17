package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	app_grpc "unomns/findpath/internal/grpc"
	findpathv1 "unomns/findpath/protos/gen/findpath"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "Port for gRPC server")
	flag.Parse()

	addr := fmt.Sprintf(":%s", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	findpathv1.RegisterPathFinderServer(
		grpcServer,
		app_grpc.NewServer(),
	)

	go func() {
		log.Printf("gRPC server listening on %s\n", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down gRPC server gracefully...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
