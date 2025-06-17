package main

import (
	"log"
	"net"
	app_grpc "unomns/findpath/internal/grpc"
	findpathv1 "unomns/findpath/protos/gen/findpath"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	findpathv1.RegisterPathFinderServer(
		grpcServer,
		app_grpc.NewServer(),
	)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
